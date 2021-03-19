// Copyright (c) 2021 Doc.ai and/or its affiliates.
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

package logs

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	allNamespaces = ""
	configPrefix  = "LOG"
)

var initConfig sync.Once
var config Config

type Config struct {
	KubeConfig   string `default:"~/.kube/config" desc:"Kubernetes configuration file" envconfig:"KUBECONFIG"`
	ArtifactsDir string `default:"logs" desc:"Directory for storing container logs" envconfig:"ARTIFACTS_DIR"`

	ExcludeK8sNs   []string      `default:"kube-system,local-path-storage" desc:"Comma-separated list of excluded kubernetes namespaces" split_words:"true"`
	ContextTimeout time.Duration `default:"15s" desc:"Context timeout for kubernetes queries" split_words:"true"`
	WorkerCount    int           `default:"4" desc:"Number of log collector workers" split_words:"true"`
}

type Suite struct {
	suite.Suite

	kubeClient kubernetes.Interface

	testStartTime time.Time
	nsmContainers map[types.UID]bool
	logQueue      chan logItem
	waitGroup     sync.WaitGroup
}

type logItem struct {
	namespace  string
	pod        string
	logDir     string
	logOptions *corev1.PodLogOptions
}

func (s *Suite) SetupSuite() {
	initConfig.Do(func() {
		if err := envconfig.Usage(configPrefix, &config); err != nil {
			panic(err)
		}

		if err := envconfig.Process(configPrefix, &config); err != nil {
			panic(err)
		}
	})

	var err error

	s.kubeClient, err = newKubeClient()
	require.NoError(s.T(), err)

	s.logQueue = make(chan logItem)
	for i := 0; i < config.WorkerCount; i++ {
		s.waitGroup.Add(1)
		go func() {
			defer s.waitGroup.Done()
			for src := range s.logQueue {
				s.saveLog(src)
			}
		}()
	}
}

func (s *Suite) TearDownSuite() {
	close(s.logQueue)
	s.waitGroup.Wait()
}

func (s *Suite) SetupTest() {
	s.testStartTime = time.Now()
	s.nsmContainers = make(map[types.UID]bool)

	ctx, cancel := context.WithTimeout(context.Background(), config.ContextTimeout)
	defer cancel()

	pods, err := s.kubeClient.CoreV1().Pods(allNamespaces).List(ctx, metav1.ListOptions{})
	require.NoError(s.T(), err)

	for podIdx := range pods.Items {
		pod := &pods.Items[podIdx]

		if !isExcludedNamespace(pod.Namespace) {
			s.nsmContainers[pod.UID] = true
		}
	}
}

func (s *Suite) AfterTest(suiteName, testName string) {
	logDir := filepath.Join(config.ArtifactsDir, suiteName, testName)
	require.NoError(s.T(), os.MkdirAll(logDir, os.ModePerm))

	logOptions := corev1.PodLogOptions{
		Timestamps: true,
		SinceTime:  &metav1.Time{Time: s.testStartTime},
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.ContextTimeout)
	defer cancel()

	pods, err := s.kubeClient.CoreV1().Pods(allNamespaces).List(ctx, metav1.ListOptions{})
	require.NoError(s.T(), err)

	var waitGroup sync.WaitGroup
	for podIdx := range pods.Items {
		pod := &pods.Items[podIdx]

		if isExcludedNamespace(pod.Namespace) {
			continue
		}

		nextLogItem := logItem{
			namespace:  pod.Namespace,
			pod:        pod.Name,
			logDir:     logDir,
			logOptions: &logOptions,
		}

		if _, ok := s.nsmContainers[pod.UID]; ok {
			s.logQueue <- nextLogItem
			continue
		}

		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			s.saveLog(nextLogItem)
		}()
	}

	waitGroup.Wait()
}

func (s *Suite) saveLog(src logItem) {
	ctx, cancel := context.WithTimeout(context.Background(), config.ContextTimeout)
	defer cancel()

	data, err := s.kubeClient.CoreV1().
		Pods(src.namespace).
		GetLogs(src.pod, src.logOptions).
		DoRaw(ctx)

	require.NoError(s.T(), err)

	if len(data) > 0 {
		logFile := filepath.Join(src.logDir, src.pod+".log")
		require.NoError(s.T(), ioutil.WriteFile(logFile, data, os.ModePerm))
	}
}

// newKubeClient creates new k8s client
func newKubeClient() (kubernetes.Interface, error) {
	if strings.HasPrefix(config.KubeConfig, "~") {
		config.KubeConfig = strings.Replace(config.KubeConfig, "~", os.Getenv("HOME"), 1)
	}

	kubeconfig, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil {
		return nil, err
	}

	kubeconfig.Burst = config.WorkerCount * 100
	kubeconfig.QPS = float32(config.WorkerCount) * 100

	return kubernetes.NewForConfig(kubeconfig)
}

func isExcludedNamespace(ns string) bool {
	for _, excluded := range config.ExcludeK8sNs {
		if ns == excluded {
			return true
		}
	}

	return false
}
