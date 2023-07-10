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
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/networkservicemesh/gotestmd/pkg/bash"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	defaultQPS        = 500 // this is default value for QPS of kubeconfig. See at documentation.
	fromAllNamespaces = ""
)

var (
	m           sync.Mutex
	once        sync.Once
	config      Config
	kubeClients []kubernetes.Interface
	kubeConfigs []string
	matchRegex  *regexp.Regexp
	runner      *bash.Bash
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
}

func ClusterDump(ctx context.Context, name string) {
	once.Do(initialize)

	m.Lock()
	defer m.Unlock()
	suitedir := filepath.Join(config.ArtifactsDir, fmt.Sprintf("cluster%v", 0), name)

	nsList, _ := kubeClients[0].CoreV1().Namespaces().List(ctx, v1.ListOptions{})

	filtered := filterNamespaces(nsList)
	_, _, exitCode, err := runner.Run(fmt.Sprintf("kubectl cluster-info dump --output-directory=%s --namespaces %s", suitedir, strings.Join(filtered, ",")))
	if exitCode != 0 || err != nil {
		logrus.Errorf("An error while getting cluster dump. Exit Code: %v, Error: %s", exitCode, err)
	}

	time.Sleep(2 * time.Second)

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
