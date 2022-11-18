// Copyright (c) 2022 Cisco and/or its affiliates.
//
// Copyright (c) 2021-2022 Doc.ai and/or its affiliates.
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
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/networkservicemesh/gotestmd/pkg/bash"
)

const (
	defaultQPS        = 5 // this is default value for QPS of kubeconfig. See at documentation.
	fromAllNamespaces = ""
	kubeconfigEnv     = "KUBECONFIG"
)

var (
	once       sync.Once
	config     Config
	jobsCh     chan func()
	ctx        context.Context
	kubeClient kubernetes.Interface
	matchRegex *regexp.Regexp
)

// Config is env config to setup log collecting.
type Config struct {
	KubeConfig        string        `default:"" desc:".kube config file path" envconfig:"KUBECONFIG"`
	ArtifactsDir      string        `default:"logs" desc:"Directory for storing container logs" envconfig:"ARTIFACTS_DIR"`
	Timeout           time.Duration `default:"5s" desc:"Context timeout for kubernetes queries" split_words:"true"`
	WorkerCount       int           `default:"8" desc:"Number of log collector workers" split_words:"true"`
	AllowedNamespaces string        `default:"(ns-.*)|(nsm-system)|(spire)|(observability)" desc:"Regex of allowed namespaces" split_words:"true"`
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func captureLogs(initialNsList []string, name string) {
	dir := filepath.Join(config.ArtifactsDir, name)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		logrus.Errorf("captureLogs: MkdirAll failed: %v", err.Error())
		return
	}

	nsList, err := listNamespaces()
	if err != nil {
		logrus.Errorf("captureLogs: can't list namespaces: %v", err.Error())
		return
	}

	var newNamespaces []string
	for _, ns := range nsList {
		if contains(initialNsList, ns) {
			continue
		}
		newNamespaces = append(newNamespaces, ns)
	}

	runner, err := bash.New()
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}
	dumpCommand := fmt.Sprintf(
		"kubectl cluster-info dump"+
			" --output yaml"+
			" --output-directory \"%v\""+
			" --namespaces %v",
		dir,
		strings.Join(newNamespaces, ","),
	)
	_, _, exitCode, err := runner.Run(dumpCommand)
	if exitCode != 0 || err != nil {
		logrus.Errorf("An error while retrieving cluster-info dump")
		return
	}
	logrus.Infof("Successfully retrieved cluster-info dump")
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

	jobsCh = make(chan func(), config.WorkerCount)

	if config.KubeConfig == "" {
		config.KubeConfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	kubeconfig, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	kubeconfig.QPS = float32(config.WorkerCount) * defaultQPS
	kubeconfig.Burst = int(kubeconfig.QPS) * 2

	kubeClient, err = kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	var cancel context.CancelFunc
	ctx, cancel = signal.NotifyContext(context.Background(),
		os.Interrupt,
		os.Kill,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	for i := 0; i < config.WorkerCount; i++ {
		go func() {
			for j := range jobsCh {
				j()
			}
		}()
	}

	go func() {
		defer cancel()
		<-ctx.Done()
		close(jobsCh)
	}()
}

func listNamespaces() ([]string, error) {
	operationCtx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	resp, err := kubeClient.CoreV1().Namespaces().List(operationCtx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nsList []string
	for _, ns := range resp.Items {
		if !matchRegex.MatchString(ns.Name) {
			continue
		}
		nsList = append(nsList, ns.Name)
	}

	return nsList, nil
}

func capture(name string) context.CancelFunc {
	once.Do(initialize)

	nsList, err := listNamespaces()
	if err != nil {
		logrus.Errorf("log saver init for %v: can't list namespaces: %v", name, err.Error())
		return func() {}
	}

	return func() {
		captureLogs(nsList, name)
	}
}

func describePods(name string) {
	getCtx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	nsList, err := kubeClient.CoreV1().Namespaces().List(getCtx, metav1.ListOptions{})
	if err != nil {
		return
	}

	runner, err := bash.New()
	if err != nil {
		return
	}

	for _, ns := range filterNamespaces(nsList) {
		_, _, exitCode, err := runner.Run("kubectl describe pods -n " + ns + ">" + filepath.Join(config.ArtifactsDir, name, "describe-"+ns+".log"))
		if exitCode != 0 || err != nil {
			logrus.Errorf("An error while retrieving describe for namespace: %v", ns)
		}
	}
}

func filterNamespaces(nsList *corev1.NamespaceList) []string {
	var rv []string

	for i := 0; i < len(nsList.Items); i++ {
		if matchRegex.MatchString(nsList.Items[i].Name) && nsList.Items[i].Status.Phase == corev1.NamespaceActive {
			rv = append(rv, nsList.Items[i].Name)
		}
	}

	return rv
}

// Capture returns a function that saves logs since Capture function has been called.
func Capture(name string) context.CancelFunc {
	c := capture(name)

	return func() {
		describePods(name)

		kubeconfigValue := os.Getenv(kubeconfigEnv)
		c()
		for i := 0; ; i++ {
			val := os.Getenv(kubeconfigEnv + fmt.Sprint(i))

			if val == "" {
				break
			}

			_ = os.Setenv(kubeconfigEnv, val)
			c()
		}
		_ = os.Setenv(kubeconfigEnv, kubeconfigValue)
	}
}
