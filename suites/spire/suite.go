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
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/networkservicemesh/integration-tests/suites/k8s"
)

// Suite is testify suite that setups spire for the tests.
type Suite struct {
	k8s.Suite
	config Config
}

// SetupTest registers unique namespace of the test in the spire
func (s *Suite) SetupTest() {
	s.Suite.SetupTest()
	entryuID := s.RegisterSpireEntry(s.Namespace(), "default")
	s.T().Cleanup(func() {
		s.DeleteSpireEntry(entryuID)
	})
}

// SetupSuite setups spire once for the suite.
func (s *Suite) SetupSuite() {
	s.Suite.SetupSuite()

	err := envconfig.Usage("spire", &s.config)
	s.Require().NoError(err)
	err = envconfig.Process("spire", &s.config)
	s.Require().NoError(err)

	nss, err := s.Run("kubectl get ns")
	s.Require().NoError(err)

	if strings.Contains(nss, s.config.Namespace) {
		return
	}

	_, err = s.Run("kubectl apply -k spire")
	s.Require().NoError(err)

	s.T().Cleanup(func() {
		if !s.config.Cleanup {
			return
		}
		_, _ = s.Run("kubectl delete ns " + s.config.Namespace)
	})

	_, err = s.Run("kubectl wait -n spire --timeout=15s --for=condition=ready pod -l app=spire-agent")
	s.Require().NoError(err)
	_, err = s.Run("kubectl wait -n spire --timeout=15s --for=condition=ready pod -l app=spire-server")
	s.Require().NoError(err)
	_, err = s.Run(`kubectl exec -n spire spire-server-0 -- \
					/opt/spire/bin/spire-server entry create \
					-spiffeID spiffe://example.org/ns/spire/sa/spire-agent \
					-selector k8s_sat:cluster:nsm-cluster \
					-selector k8s_sat:agent_ns:spire \
					-selector k8s_sat:agent_sa:spire-agent \
					-node`)
	s.Require().NoError(err)

	_ = s.RegisterSpireEntry("default", "default")
}

// DeleteSpireEntry deletes spire entry.
func (s *Suite) DeleteSpireEntry(entryID string) {
	_, err := s.Runf(`kubectl exec -n spire spire-server-0 -- \
						/opt/spire/bin/spire-server entry delete \
						-entryID %v`, entryID)
	s.Require().NoError(err)
}

// RegisterSpireEntry register spire entry by specific namespace and service account.
func (s *Suite) RegisterSpireEntry(ns, sa string) (entryID string) {
	out, err := s.Runf(`kubectl exec -n spire spire-server-0 -- \
						/opt/spire/bin/spire-server entry create \
						-spiffeID spiffe://example.org/ns/%v/sa/%v \
						-parentID spiffe://example.org/ns/spire/sa/spire-agent \
						-selector k8s:ns:%v \
						-selector k8s:sa:%v`, ns, sa, ns, sa)
	s.Require().NoError(err)
	line := strings.Split(out, "\n")[0]
	words := strings.Split(line, " ")
	entryID = words[len(words)-1]
	return entryID
}
