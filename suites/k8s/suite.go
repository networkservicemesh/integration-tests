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

// Package k8s provides k8s specific suite that provides k8s api into each test
package k8s

import (
	"github.com/networkservicemesh/integration-tests/tools/k8s"
	"github.com/stretchr/testify/suite"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Suite is testify.Suite that provides k8s helpers and setups k8s ns for each test.
type Suite struct {
	suite.Suite
	kt     *k8s.Testing
	client kubernetes.Interface
}

// SetupSuite initializes k8s api client and helpers
func (s *Suite) SetupSuite() {
	s.kt = k8s.NewK8sTesting(s.T())
	client, err := k8s.NewClient()
	s.Require().NoError(err)
	s.client = client
}

// Client returns raw k8s client
func (s *Suite) Client() kubernetes.Interface {
	return s.client
}

// K8sT returns wrapped k8s client with helper and assertion functions
func (s *Suite) K8sT() *k8s.Testing {
	return s.kt
}

// SetupTest setups ns for the current test
func (s *Suite) SetupTest() {
	c := s.Client()
	ns, err := c.CoreV1().Namespaces().Create(&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{GenerateName: "ns-"}})
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = c.CoreV1().Namespaces().Delete(ns.Name, &metav1.DeleteOptions{})
		s.kt.SetNamespace("default")
	})
	s.kt.SetNamespace(ns.Name)
}
