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

package examples

import (
	"github.com/networkservicemesh/integration-tests/examples/exec"
	_ "github.com/networkservicemesh/integration-tests/tools/versioning"
	"github.com/stretchr/testify/suite"
)

type Examples struct {
	suite.Suite
}

func (s *Examples) SetupSuite() {
	s.T().Cleanup(s.Cleanup)
	s.SetupSpire()
}

func (s *Examples) SetupSpire() {
	runner := exec.New(s.Require())

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
	runner.Run(`kubectl exec -n spire spire-server-0 --
					/opt/spire/bin/spire-server entry create
					-spiffeID spiffe://example.org/ns/default/sa/default
					-parentID spiffe://example.org/ns/spire/sa/spire-agent
					-selector k8s:ns:default
					-selector k8s:sa:default`)
	runner.Run(`kubectl exec -n spire spire-server-0 --
					/opt/spire/bin/spire-server entry create
					-spiffeID spiffe://example.org/ns/nsm-system/sa/default
					-parentID spiffe://example.org/ns/spire/sa/spire-agent
					-selector k8s:ns:nsm-system
					-selector k8s:sa:default`)
}

func (s *Examples) Cleanup() {
	runner := exec.New(s.Require())

	runner.Run("kubectl delete ns spire")
	runner.Run("kubectl delete ns nsm-system")
}
