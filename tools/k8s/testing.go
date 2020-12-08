// Copyright (c) 2020 Doc.ai and/or its affiliates.
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

package k8s

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/networkservicemesh/integration-tests/tools/versioning"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

const timeout = time.Minute / 2

// Testing provides assert and helper functions for k8s testing
type Testing struct {
	assert *require.Assertions
	client kubernetes.Interface
	ns     string
}

// NewK8sTesting creates new k8s testing
func NewK8sTesting(t *testing.T) *Testing {
	assert := require.New(t)
	client, err := NewClient()
	assert.NoError(err)
	return &Testing{
		client: client,
		assert: assert,
		ns:     "default",
	}
}

// SetNamespace sets default namespace
func (t *Testing) SetNamespace(ns string) {
	t.ns = ns
}

// Namespace gets default namespace
func (t *Testing) Namespace() string {
	return t.ns
}

// ApplyService is analogy of 'kubeclt apply -f path'
func (t *Testing) ApplyService(name string) *corev1.Service {
	loc := path.Join(versioning.Dir(), name)
	b, err := ioutil.ReadFile(path.Clean(loc))
	t.assert.NoError(err)
	var s corev1.Service
	err = yaml.Unmarshal(b, &s)
	t.assert.NoError(err)
	result, err := t.client.CoreV1().Services(t.ns).Create(&s)
	t.assert.NoError(err)
	return result
}

// Nodes returns a slice of Nodes where can be deployed deployment
func (t *Testing) Nodes() []*corev1.Node {
	response, err := t.client.CoreV1().Nodes().List(metav1.ListOptions{})
	t.assert.NoError(err)

	var result []*corev1.Node
	for i := 0; i < len(response.Items); i++ {
		node := &response.Items[i]
		name := node.Labels["kubernetes.io/hostname"]
		if !strings.HasSuffix(name, "control-plane") {
			result = append(result, node)
		}
	}

	return result
}

// ApplyServiceAccount is analogy of 'kubeclt apply -f path'
func (t *Testing) ApplyServiceAccount(name string) *corev1.ServiceAccount {
	loc := path.Join(versioning.Dir(), name)
	b, err := ioutil.ReadFile(path.Clean(loc))
	t.assert.NoError(err)
	var p corev1.ServiceAccount
	err = yaml.Unmarshal(b, &p)
	t.assert.NoError(err)
	if p.Namespace == "" {
		p.Namespace = t.ns
	}
	result, err := t.client.CoreV1().ServiceAccounts(p.Namespace).Create(&p)
	t.assert.NoError(err)
	return result
}

// ApplyPod is analogy of 'kubeclt apply -f path' but with mutating pod before apply
func (t *Testing) ApplyPod(name string, mutators ...func(pod *corev1.Pod)) *corev1.Pod {
	loc := path.Join(versioning.Dir(), name)
	b, err := ioutil.ReadFile(path.Clean(loc))
	t.assert.NoError(err)
	var p corev1.Pod
	err = yaml.Unmarshal(b, &p)
	t.assert.NoError(err)
	for _, m := range mutators {
		m(&p)
	}
	if p.Namespace == "" {
		p.Namespace = t.ns
	}
	result, err := t.client.CoreV1().Pods(p.Namespace).Create(&p)
	t.assert.NoError(err)
	return result
}

// ApplyDeployment is analogy of 'kubeclt apply -f path' but with mutating deployment before apply
func (t *Testing) ApplyDeployment(name string, mutators ...func(deployment *v1.Deployment)) *v1.Deployment {
	loc := path.Join(versioning.Dir(), name)
	b, err := ioutil.ReadFile(path.Clean(loc))
	t.assert.NoError(err)
	var d v1.Deployment
	err = yaml.Unmarshal(b, &d)
	t.assert.NoError(err)
	for _, m := range mutators {
		m(&d)
	}
	if d.Namespace == "" {
		d.Namespace = t.ns
	}
	result, err := t.client.AppsV1().Deployments(d.Namespace).Create(&d)
	t.assert.NoError(err)
	return result
}

// ApplyDaemonSet is analogy of 'kubeclt apply -f path' but with mutating DaemonSet before apply
func (t *Testing) ApplyDaemonSet(name string, mutators ...func(*v1.DaemonSet)) *v1.DaemonSet {
	loc := path.Join(versioning.Dir(), name)
	b, err := ioutil.ReadFile(path.Clean(loc))
	t.assert.NoError(err)
	var d v1.DaemonSet
	err = yaml.Unmarshal(b, &d)
	t.assert.NoError(err)
	for _, m := range mutators {
		m(&d)
	}
	if d.Namespace == "" {
		d.Namespace = t.ns
	}
	result, err := t.client.AppsV1().DaemonSets(d.Namespace).Create(&d)
	t.assert.NoError(err)
	return result
}

// WaitLogsMatch waits regex pattern in logs. Note: this function will be deleted when forwarders will be available for testing.
func (t *Testing) WaitLogsMatch(podName, containerName, pattern string) {
	r, err := regexp.Compile(pattern)
	t.assert.NoError(err)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for ctx.Err() == nil {
		result, err := t.client.CoreV1().Pods(t.ns).GetLogs(podName, &corev1.PodLogOptions{Container: containerName}).DoRaw()
		if err == nil && r.MatchString(string(result)) {
			return
		}
		time.Sleep(timeout / 10)
	}
	t.assert.FailNowf(ctx.Err().Error(), "cannot match pattern: %v in logs of %v:%v ", pattern, podName, containerName)
}

// NoRestarts check that pods have not restarts in namepace
func (t *Testing) NoRestarts(ns string) {
	list, err := t.client.CoreV1().Pods(ns).List(metav1.ListOptions{})
	t.assert.NoError(err)
	for i := 0; i < len(list.Items); i++ {
		pod := &list.Items[i]
		for j := 0; j < len(pod.Status.ContainerStatuses); j++ {
			status := &pod.Status.ContainerStatuses[j]
			reason := ""
			if status.LastTerminationState.Terminated != nil {
				reason = status.LastTerminationState.Terminated.Reason
			}

			t.assert.Zero(status.RestartCount, fmt.Sprintf("Container %v of Pod %v has restart count more then zero. Reason: %v", status.Name, pod.Name, reason))
		}
	}
}
