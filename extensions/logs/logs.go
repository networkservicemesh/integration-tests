// Copyright (c) 2021-2022 Doc.ai and/or its affiliates.
//
// Copyright (c) 2023 Cisco and/or its affiliates.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package logs exports helper functions for storing logs from containers.
package logs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/edwarnicke/genericsync"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"github.com/networkservicemesh/gotestmd/pkg/bash"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	defaultQPS        = 500 // this is default value for QPS of kubeconfig. See at documentation.
	fromAllNamespaces = ""
)

var (
	once        sync.Once
	config      Config
	kubeClients []kubernetes.Interface
	kubeConfigs []string
	matchRegex  *regexp.Regexp
	suiteMap    genericsync.Map[string, struct{}]
)

// Config is env config to setup log collecting.
type Config struct {
	ArtifactsDir         string        `default:"logs" desc:"Directory for storing container logs" envconfig:"ARTIFACTS_DIR"`
	Timeout              time.Duration `default:"10s" desc:"Context timeout for kubernetes queries" split_words:"true"`
	WorkerCount          int           `default:"8" desc:"Number of log collector workers" split_words:"true"`
	MaxKubeConfigs       int           `default:"3" desc:"Number of used kubeconfigs" split_words:"true"`
	AllowedNamespaces    string        `default:"(ns-.*)|(nsm-system)|(spire)|(observability)" desc:"Regex of allowed namespaces" split_words:"true"`
	LogCollectionEnabled bool          `default:"true" desc:"Boolean variable which enables log collection" split_words:"true"`
}

func initialize() {
	const prefix = "logs"
	if err := envconfig.Usage(prefix, &config); err != nil {
		logrus.Fatal(err.Error())
	}

	if err := envconfig.Process(prefix, &config); err != nil {
		logrus.Fatal(err.Error())
	}

	if !config.LogCollectionEnabled {
		return
	}

	matchRegex = regexp.MustCompile(config.AllowedNamespaces)

	var singleClusterKubeConfig = os.Getenv("KUBECONFIG")

	if singleClusterKubeConfig == "" {
		singleClusterKubeConfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	kubeConfigs = []string{}

	for i := 1; i <= config.MaxKubeConfigs; i++ {
		kubeConfig := os.Getenv("KUBECONFIG" + fmt.Sprint(i))
		if kubeConfig != "" {
			kubeConfigs = append(kubeConfigs, kubeConfig)
		}
	}

	if len(kubeConfigs) == 0 {
		kubeConfigs = append(kubeConfigs, singleClusterKubeConfig)
	}

	for _, cfg := range kubeConfigs {
		kubeconfig, err := clientcmd.BuildConfigFromFlags("", cfg)
		if err != nil {
			logrus.Fatal(err.Error())
		}

		kubeconfig.QPS = float32(config.WorkerCount) * defaultQPS
		kubeconfig.Burst = int(kubeconfig.QPS) * 2

		kubeClient, err := kubernetes.NewForConfig(kubeconfig)
		if err != nil {
			logrus.Fatal(err.Error())
		}

		kubeClients = append(kubeClients, kubeClient)
	}
}

func ClusterDump(ctx context.Context, name string) {
	if _, ok := suiteMap.LoadOrStore(name, struct{}{}); ok {
		return
	}

	runner, err := bash.New()
	if err != nil {
		logrus.Errorf("An error while getting cluster dump")
		return
	}

	matchRegex = regexp.MustCompile(config.AllowedNamespaces)

	var singleClusterKubeConfig = os.Getenv("KUBECONFIG")

	if singleClusterKubeConfig == "" {
		singleClusterKubeConfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", singleClusterKubeConfig)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	kubeconfig.QPS = 500
	kubeconfig.Burst = int(kubeconfig.QPS) * 2

	kubeClient, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	suitedir := filepath.Join(config.ArtifactsDir, fmt.Sprintf("cluster%v", 0), name)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			nsList, _ := kubeClient.CoreV1().Namespaces().List(ctx, v1.ListOptions{})

			filtered := filterNamespaces(nsList)
			_, _, exitCode, err := runner.Run(fmt.Sprintf("kubectl cluster-info dump --output-directory=%s --namespaces %s", suitedir, strings.Join(filtered, ",")))
			if exitCode != 0 || err != nil {
				logrus.Errorf("An error while getting cluster dump. Exit Code: %v, Error: %s", exitCode, err)
			}

			time.Sleep(2 * time.Second)
		}
	}
}

func filterNamespaces(nsList *corev1.NamespaceList) []string {
	result := make([]string, 0)

	for _, ns := range nsList.Items {
		if matchRegex.MatchString(ns.Name) {
			result = append(result, ns.Name)
		}
	}

	return result
}

func MonitorNamespaces(ctx context.Context, name string) {
	if _, ok := suiteMap.LoadOrStore(name, struct{}{}); ok {
		return
	}
	once.Do(initialize)
	for i := range kubeClients {
		suitedir := filepath.Join(config.ArtifactsDir, fmt.Sprintf("cluster%v", i), name)
		go monitorNamespaces(ctx, kubeClients[i], suitedir)
	}
}

type logCollector struct {
	kubeClient kubernetes.Interface
	suiteName  string
}

func monitorNamespaces(ctx context.Context, kubeClient kubernetes.Interface, suiteName string) {
	podMap := make(map[string]func())
	watcher, _ := kubeClient.CoreV1().Pods(v1.NamespaceAll).Watch(ctx, v1.ListOptions{})
	eventCh := watcher.ResultChan()

	collector := &logCollector{
		kubeClient: kubeClient,
		suiteName:  suiteName,
	}
	for {
		select {
		case event := <-eventCh:
			pod, ok := event.Object.(*corev1.Pod)
			if !ok {
				return
			}

			if event.Type == watch.Added {
				if matchRegex.MatchString(pod.Namespace) {
					podKey := fmt.Sprintf("%v-%v", pod.Namespace, pod.Name)
					collectCtx, collectCancel := context.WithCancel(ctx)

					if _, ok := podMap[podKey]; ok {
						collectCancel()
						return
					}

					podMap[podKey] = collectCancel
					go collector.collectLogs(collectCtx, pod, uuid.NewString())
				}
			}

			if event.Type == watch.Deleted {
				if matchRegex.MatchString(pod.Namespace) {
					cancel := podMap[fmt.Sprintf("%v-%v", pod.Namespace, pod.Name)]
					cancel()
				}
			}
		case <-ctx.Done():
			watcher.Stop()
		}

	}
}

type logReader struct {
	podLogOptions *corev1.PodLogOptions
	LastTimeRead  time.Time
	stream        io.ReadCloser
	logBuffer     bytes.Buffer
	outputFile    string
	doStream      func(opts *corev1.PodLogOptions)
	fileWriter    *os.File
}

func (l *logReader) Save() {
	err := os.MkdirAll(filepath.Dir(l.outputFile), os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating dir: %v", err.Error())
	}
	file, err := os.Create(l.outputFile)
	if err != nil {
		fmt.Printf("Error: %v occured when saving logs to file\n", err.Error())
	}
	defer file.Close()

	file.Write(l.logBuffer.Bytes())
	l.logBuffer.Reset()
}

func (l *logReader) Read(buffer []byte) (int, error) {
	if l.stream == nil {
		return 0, errors.New("stream is nil")
	}
	n, err := l.stream.Read(buffer)
	if n != 0 {
		l.LastTimeRead = time.Now()
	}

	return n, err
}

func (l *logCollector) collectLogs(collectCtx context.Context, pod *corev1.Pod, id string) {
	readers := make([]*logReader, len(pod.Spec.Containers))
	//bufferSize := int64(32 * 1024)
	//buffer := make([]byte, bufferSize)

	for i := range readers {
		readers[i] = &logReader{}
		container := pod.Spec.Containers[i].Name
		readers[i].outputFile = filepath.Join(l.suiteName, pod.Namespace, pod.Name) + "-" + container + ".log"
		readers[i].podLogOptions = &corev1.PodLogOptions{}
		//readers[i].podLogOptions.Follow = true
		readers[i].podLogOptions.Container = container

		err := os.MkdirAll(filepath.Dir(readers[i].outputFile), os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating dir: %v", err.Error())
		}
		file, err := os.Create(readers[i].outputFile)
		readers[i].fileWriter = file
		if err != nil {
			fmt.Printf("Error opening file: %v", err.Error())
		}

		// readers[i].doStream = func(opts *corev1.PodLogOptions) {
		// 	stream, _ := l.kubeClient.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, opts).Stream(collectCtx)
		// 	for stream == nil {
		// 		select {
		// 		case <-collectCtx.Done():
		// 			return
		// 		default:
		// 			time.Sleep(time.Second)
		// 			stream, _ = l.kubeClient.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, opts).Stream(collectCtx)
		// 		}
		// 	}
		// 	readers[i].stream = stream
		// }
		// readers[i].doStream(readers[i].podLogOptions)
	}

	for {
		for _, reader := range readers {
			buf, err := l.kubeClient.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, reader.podLogOptions).DoRaw(collectCtx)
			if err == nil {
				reader.podLogOptions.SinceTime = &v1.Time{time.Now()}
				reader.fileWriter.Write(buf)
			}
			// if n, err := reader.Read(buffer); err != nil {
			// 	if collectCtx.Err() == nil {
			// 		if err == io.EOF {
			// 			reader.podLogOptions.SinceTime = &v1.Time{Time: reader.LastTimeRead}
			// 		}
			// 		if reader.stream != nil {
			// 			reader.stream.Close()
			// 		}
			// 		reader.doStream(reader.podLogOptions)
			// 	}
			// } else {
			// 	reader.logBuffer.Write(buffer[:n])
			// }
		}

		if collectCtx.Err() != nil {
			break
		}

		time.Sleep(time.Second)
	}

	for _, reader := range readers {
		reader.fileWriter.Close()
		//reader.Save()
		if reader.stream != nil {
			reader.stream.Close()
		}
	}
}
