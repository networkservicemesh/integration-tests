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

// Package basic provides a suit that tests basic NSM scenarios.
package basic

import (
	"github.com/networkservicemesh/integration-tests/suites/nsm"
	"github.com/networkservicemesh/integration-tests/tools/k8s"
)

// Suite provides basic NSM scenarios
type Suite struct {
	nsm.Suite
}

// TestLocalUsecase checks that nsc can connect to local nse
func (s *Suite) TestLocalUsecase() {
	kt := s.K8sT()

	nodes := kt.Nodes()

	s.Require().Greater(len(nodes), 0)

	_ = kt.ApplyPod("nse.yaml", k8s.SetNode(nodes[0].Labels["kubernetes.io/hostname"]))
	nsc := kt.ApplyPod("nsc.yaml", k8s.SetNode(nodes[0].Labels["kubernetes.io/hostname"]))

	kt.WaitLogsMatch(nsc.Name, nsc.Spec.Containers[0].Name, "All client init operations are done")
}

// TestRemoteUsecase checks that nsc can connect to remote nse
func (s *Suite) TestRemoteUsecase() {
	kt := s.K8sT()

	nodes := kt.Nodes()

	s.Require().Greater(len(nodes), 1)

	_ = kt.ApplyPod("nse.yaml", k8s.SetNode(nodes[1].Labels["kubernetes.io/hostname"]))
	nsc := kt.ApplyPod("nsc.yaml", k8s.SetNode(nodes[0].Labels["kubernetes.io/hostname"]))

	kt.WaitLogsMatch(nsc.Name, nsc.Spec.Containers[0].Name, "All client init operations are done")
}
