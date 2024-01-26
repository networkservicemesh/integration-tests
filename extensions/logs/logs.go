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
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/networkservicemesh/gotestmd/pkg/bash"
)

const (
	defaultQPS        = 500 // this is default value for QPS of kubeconfig. See at documentation.
	fromAllNamespaces = ""
)

var (
	once                       sync.Once
	config                     Config
	ctx                        context.Context
	kubeConfigs                []string
	kubeClients                []kubernetes.Interface
	matchRegex                 *regexp.Regexp
	runner                     *bash.Bash
	clusterDumpSingleOperation *singleOperation
)

// Config is env config to setup log collecting.
type Config struct {
	ArtifactsDir         string        `default:"logs" desc:"Directory for storing container logs" envconfig:"ARTIFACTS_DIR"`
	Timeout              time.Duration `default:"10s" desc:"Context timeout for kubernetes queries" split_words:"true"`
	WorkerCount          int           `default:"8" desc:"Number of log collector workers" split_words:"true"`
	MaxKubeConfigs       int           `default:"3" desc:"Number of used kubeconfigs" split_words:"true"`
	AllowedNamespaces    string        `default:"(ns-.*)|(spire)|(observability)" desc:"Regex of allowed namespaces" split_words:"true"`
	LogCollectionEnabled bool          `default:"true" desc:"Boolean variable which enables log collection" split_words:"true"`
}

// nolint: gocyclo
func initialize(suiteName string) {
	if err := envconfig.Usage("logs", &config); err != nil {
		logrus.Fatal(err.Error())
	}

	if err := envconfig.Process("logs", &config); err != nil {
		logrus.Fatal(err.Error())
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

	runner, _ = bash.New()

	ctx, _ = signal.NotifyContext(context.Background(),
		os.Interrupt,
		os.Kill,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	clusterDumpSingleOperation = newSingleOperation(func() {
		if ctx.Err() != nil {
			return
		}
		for i := range kubeConfigs {
			suitedir := filepath.Join(config.ArtifactsDir, fmt.Sprintf("cluster%v", i), suiteName)

			nsString, _, _, _ := runner.Run(fmt.Sprintf(`kubectl --kubeconfig %v get ns -o go-template='{{range .items}}{{ .metadata.name }} {{end}}'`, kubeConfigs[i]))
			nsList := strings.Split(nsString, " ")

			_, _, exitCode, err := runner.Run(fmt.Sprintf("kubectl --kubeconfig %v cluster-info dump --output-directory=%s --namespaces %s",
				kubeConfigs[i],
				suitedir,
				strings.Join(filterNamespaces(nsList), ",")))

			if exitCode != 0 {
				logrus.Errorf("An error while getting cluster dump. Exit Code: %v", exitCode)
			}
			if err != nil {
				logrus.Errorf("An error while getting cluster dump. Error: %s", err.Error())
			}
		}
	})
}

// ClusterDump saves logs from all pods in specified namespaces
func ClusterDump(suiteName string) {
	once.Do(func() { initialize(suiteName) })
	clusterDumpSingleOperation.Run()
}

func filterNamespaces(nsList []string) []string {
	result := make([]string, 0)

	for i := range nsList {
		if matchRegex.MatchString(nsList[i]) {
			result = append(result, nsList[i])
		}
	}

	return result
}

func MonitorNSMSystem(ctx context.Context, suiteName string) {
	once.Do(func() { initialize(suiteName) })

	for i := range kubeClients {
		suitedir := filepath.Join(config.ArtifactsDir, fmt.Sprintf("cluster%v", i), suiteName)
		go monitorNamespaces(ctx, kubeClients[i], suitedir)
	}
}

type logCollector struct {
	kubeClient kubernetes.Interface
	suiteName  string
}

func monitorNamespaces(ctx context.Context, kubeClient kubernetes.Interface, suiteName string) {
	podMap := make(map[string]func())
	watcher, _ := kubeClient.CoreV1().Pods("nsm-system").Watch(ctx, v1.ListOptions{})
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
				fmt.Printf("POD ADDED: %s\n", pod.Name)
				podKey := fmt.Sprintf("%v-%v", pod.Namespace, pod.Name)
				if _, ok := podMap[podKey]; ok {
					return
				}

				collectCtx, collectCancel := context.WithCancel(ctx)
				podMap[podKey] = collectCancel
				go collector.collectLogs(collectCtx, pod, uuid.NewString())
			}

			if event.Type == watch.Deleted {
				podKey := fmt.Sprintf("%v-%v", pod.Namespace, pod.Name)
				fmt.Printf("POD DELETED: %s\n", pod.Name)
				if cancel, ok := podMap[podKey]; ok {
					cancel()
					delete(podMap, podKey)
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
	buf           [65536]byte
}

func (l *logReader) Save() {
	err := os.MkdirAll(filepath.Dir(l.outputFile), os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating dir: %v\n", err.Error())
	}
	file, err := os.Create(l.outputFile)
	if err != nil {
		fmt.Printf("Error: %v occured when saving logs to file\n", err.Error())
	}
	defer file.Close()

	file.Write(l.logBuffer.Bytes())
	l.logBuffer.Reset()
}

// func (l *logReader) Read() (int, error) {
// 	if l.stream == nil {
// 		return 0, errors.New("stream is nil")
// 	}

// 	n, err := l.stream.Read(l.buf[:])
// 	if n != 0 {
// 		l.LastTimeRead = time.Now()
// 	}

// 	return n, err
// }

func (l *logCollector) collectLogs(collectCtx context.Context, pod *corev1.Pod, id string) {
	readers := make([]*logReader, len(pod.Spec.Containers))

	for i := range readers {
		readers[i] = &logReader{}
		container := pod.Spec.Containers[i].Name
		readers[i].outputFile = filepath.Join(l.suiteName, pod.Namespace, pod.Name) + "-" + container + ".log"
		readers[i].podLogOptions = &corev1.PodLogOptions{}
		// readers[i].podLogOptions.Follow = true
		readers[i].podLogOptions.SinceTime = &v1.Time{Time: time.Now()}
		readers[i].podLogOptions.Container = container

		err := os.MkdirAll(filepath.Dir(readers[i].outputFile), os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating dir: %v\n", err.Error())
		}
		file, err := os.Create(readers[i].outputFile)
		readers[i].fileWriter = file
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err.Error())
		}

		readers[i].doStream = func(opts *corev1.PodLogOptions) {
			stream, _ := l.kubeClient.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, opts).Stream(collectCtx)
			for stream == nil {
				select {
				case <-collectCtx.Done():
					return
				default:
					time.Sleep(time.Second)
					stream, _ = l.kubeClient.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, opts).Stream(collectCtx)
				}
			}
			readers[i].stream = stream
		}
		// readers[i].doStream(readers[i].podLogOptions)
	}

	// for i := range readers {
	// 	stream, err := l.kubeClient.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, readers[i].podLogOptions).Stream(collectCtx)
	// 	fmt.Printf("OPENING STREAM FOR POD: %s-%s\n", pod.Name, readers[i].podLogOptions.Container)
	// 	if err != nil {
	// 		fmt.Printf("ERROR WHILE OPENING STREAM: %s\n", err.Error())
	// 	}

	// 	for stream == nil {
	// 		select {
	// 		case <-collectCtx.Done():
	// 			return
	// 		default:
	// 			time.Sleep(time.Second)
	// 			fmt.Printf("OPENING STREAM\n")
	// 			stream, err = l.kubeClient.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, readers[i].podLogOptions).Stream(collectCtx)
	// 			if err != nil {
	// 				fmt.Printf("ERROR WHILE OPENING STREAM: %s-s\n", err.Error())
	// 			}
	// 		}
	// 	}

	// 	readers[i].stream = stream
	// }

	for i := range readers {
		index := i
		go func() {
			// for {
			// 	if readers[index].stream == nil {
			// 		readers[index].doStream(readers[index].podLogOptions)
			// 	}
			// 	n, err := readers[index].stream.Read(readers[index].buf[:])
			// 	if err != nil {
			// 		if collectCtx.Err() != nil {
			// 			fmt.Printf("SAVING LOGS FROM: %s-%s", pod.Name, readers[index].podLogOptions.Container)
			// 			readers[index].Save()
			// 			if readers[index].stream != nil {
			// 				readers[index].stream.Close()
			// 			}
			// 			return
			// 		}
			// 	}
			// 	readers[index].logBuffer.Write(readers[index].buf[:n])
			// 	time.Sleep(time.Millisecond * 500)
			// }

			defer func() {
				fmt.Printf("SAVING LOGS FROM: %s-%s\n", pod.Name, readers[index].podLogOptions.Container)
				readers[index].Save()
			}()

			for {
				buf, err := l.kubeClient.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, readers[index].podLogOptions).DoRaw(collectCtx)
				if err != nil {
					fmt.Printf("ERROR WHILE READING: %s\n", err.Error())
					if collectCtx.Err() != nil {
						return
					}
					time.Sleep(time.Second)
					continue
				}

				readers[index].logBuffer.Write(buf)
				if collectCtx.Err() != nil {
					return
				}

				time.Sleep(time.Millisecond * 500)
				readers[index].podLogOptions.SinceTime = &v1.Time{Time: time.Now()}
			}
		}()
	}

	// for {
	// 	for _, reader := range readers {
	// 		buf, err := l.kubeClient.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, reader.podLogOptions).DoRaw(collectCtx)
	// 		if err == nil {
	// 			reader.podLogOptions.SinceTime = &v1.Time{time.Now()}
	// 			reader.fileWriter.Write(buf)
	// 		}
	// 		// if n, err := reader.Read(buffer); err != nil {
	// 		// 	if collectCtx.Err() == nil {
	// 		// 		if err == io.EOF {
	// 		// 			reader.podLogOptions.SinceTime = &v1.Time{Time: reader.LastTimeRead}
	// 		// 		}
	// 		// 		if reader.stream != nil {
	// 		// 			reader.stream.Close()
	// 		// 		}
	// 		// 		reader.doStream(reader.podLogOptions)
	// 		// 	}
	// 		// } else {
	// 		// 	reader.logBuffer.Write(buffer[:n])
	// 		// }
	// 	}

	// 	if collectCtx.Err() != nil {
	// 		break
	// 	}

	// 	time.Sleep(time.Second)
	// }

	// for _, reader := range readers {
	// 	reader.fileWriter.Close()
	// 	//reader.Save()
	// 	if reader.stream != nil {
	// 		reader.stream.Close()
	// 	}
	// }
}
