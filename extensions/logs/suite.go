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
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var workerCount int = runtime.NumCPU()
var excludeNamespaceFlag = flag.String("log-exclude-k8s-ns", "kube-system,local-path-storage",
	"comma separated list of excluded kubernetes namespaces")
var contextTimeoutFlag = flag.String("log-context-timeout", "15s",
	"log context timeout")
var logDirFlag = flag.String("log-dir", envOrDefault("ARTIFACTS_DIR", "logs"),
	"container logs directory")

type Suite struct {
	suite.Suite

	testStartTime time.Time
	ctxTimeout    time.Duration
	kubeClient    kubernetes.Interface
}

type logSource struct {
	namespace  string
	pod        string
	logOptions *corev1.PodLogOptions
}

func (s *Suite) SetupSuite() {
	var err error

	s.ctxTimeout, err = time.ParseDuration(*contextTimeoutFlag)
	require.NoError(s.T(), err)

	s.kubeClient, err = newKubeClient()
	require.NoError(s.T(), err)
}

func (s *Suite) SetupTest() {
	s.testStartTime = time.Now()
}

func (s *Suite) AfterTest(suiteName, testName string) {
	logPath := fmt.Sprintf("%s/%s/%s", *logDirFlag, suiteName, testName)
	require.NoError(s.T(), os.MkdirAll(logPath, os.ModePerm))

	logOptions := corev1.PodLogOptions{
		Timestamps: true,
		SinceTime:  &metav1.Time{Time: s.testStartTime},
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	list, err := s.kubeClient.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	require.NoError(s.T(), err)

	var waitGroup sync.WaitGroup
	var logSources = make(chan logSource)
	for i := 0; i < workerCount; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			s.worker(ctx, logPath, logSources)
		}()
	}

	for nsIdx := range list.Items {
		ns := &list.Items[nsIdx]

		if isExcludedNamespace(ns.Name) {
			continue
		}

		pods, err := s.kubeClient.CoreV1().Pods(ns.Name).List(ctx, metav1.ListOptions{})
		require.NoError(s.T(), err)

		for podIdx := range pods.Items {
			pod := &pods.Items[podIdx]
			logSources <- logSource{
				namespace:  ns.Name,
				pod:        pod.Name,
				logOptions: &logOptions,
			}
		}
	}

	close(logSources)
	waitGroup.Wait()
}

// newKubeClient creates new k8s client
func newKubeClient() (kubernetes.Interface, error) {
	defaultPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	path := envOrDefault("KUBECONFIG", defaultPath)

	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		return nil, err
	}

	config.Burst = workerCount * 100
	config.QPS = float32(workerCount) * 100

	return kubernetes.NewForConfig(config)
}

func isExcludedNamespace(ns string) bool {
	excludeList := strings.Split(*excludeNamespaceFlag, `,`)
	for _, excluded := range excludeList {
		if ns == excluded {
			return true
		}
	}

	return false
}

func (s *Suite) worker(ctx context.Context, logPath string, sources <-chan logSource) {
	for src := range sources {
		data, err := s.kubeClient.CoreV1().
			Pods(src.namespace).
			GetLogs(src.pod, src.logOptions).
			DoRaw(ctx)

		require.NoError(s.T(), err)

		if len(data) > 0 {
			logFile := fmt.Sprintf("%s/%s.log", logPath, src.pod)
			require.NoError(s.T(), ioutil.WriteFile(logFile, data, os.ModePerm))
		}
	}
}

func envOrDefault(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		value = defaultValue
	}

	return value
}
