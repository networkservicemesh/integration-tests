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

package spire

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/networkservicemesh/integration-tests/examples/exec"
	k8s "github.com/networkservicemesh/integration-tests/suites/k8s"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Suite setups spire and exports spire helper functions
type Suite struct {
	k8s.Suite
	config Config
}

// SetupTest registers current test namespace
func (s *Suite) SetupTest() {
	s.Suite.SetupTest()
	s.RegisterNamespace(s.K8sT().Namespace())
}

// SetupSuite setups spire
func (s *Suite) SetupSuite() {
	s.Suite.SetupSuite()

	s.Require().NoError(envconfig.Usage("spire", &s.config))
	s.Require().NoError(envconfig.Process("spire", &s.config))

	s.T().Cleanup(func() {
		if !s.config.Cleanup {
			logrus.WithField("spire.Suite", "TearDownSuite").Warn("cleanup skipped due to config.Cleanup = false")
			return
		}
		runner := exec.New(s.Require())
		runner.Run("kubectl delete ns " + s.config.Namespace)
	})

	runner := exec.New(s.Require())

	_, err := s.Client().CoreV1().Namespaces().Get(s.config.Namespace, v1.GetOptions{})

	if err == nil {
		return
	}

	// Setup spire
	runner.Run("kubectl apply -f spire/spire-namespace.yaml")
	runner.Run("kubectl apply -f spire")
	runner.Run("kubectl wait -n spire --timeout=60s --for=condition=ready pod -l app=spire-agent")
	runner.Run("kubectl wait -n spire --timeout=60s --for=condition=ready pod -l app=spire-server")
	runner.Run(`kubectl exec -n spire spire-server-0 --
					/opt/spire/bin/spire-server entry create
					-spiffeID spiffe://example.org/ns/spire/sa/spire-agent
					-selector k8s_sat:cluster:nsm-cluster
					-selector k8s_sat:agent_ns:spire
					-selector k8s_sat:agent_sa:spire-agent
					-node`)
	s.RegisterNamespace("default")
}

// RegisterNamespace registers namespace entry
func (s *Suite) RegisterNamespace(namespace string) {
	runner := exec.New(s.Require())
	runner.Run(fmt.Sprintf(`kubectl exec -n spire spire-server-0 --
								/opt/spire/bin/spire-server entry create
								-spiffeID spiffe://example.org/ns/%v/sa/default
								-parentID spiffe://example.org/ns/spire/sa/spire-agent
								-selector k8s:ns:%v
								-selector k8s:sa:default`, namespace, namespace))
}

// DeleteNamespace deletes namespace entry
func (s *Suite) DeleteNamespace(namespace string) {
	runner := exec.New(s.Require())
	runner.Run(fmt.Sprintf(`kubectl exec -n spire spire-server-0 --
								/opt/spire/bin/spire-server entry delete
								-spiffeID spiffe://example.org/ns/%v/sa/default
								-parentID spiffe://example.org/ns/spire/sa/spire-agent
								-selector k8s:ns:%v
								-selector k8s:sa:default`, namespace, namespace))
}
