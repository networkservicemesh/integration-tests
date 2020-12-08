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

import "github.com/networkservicemesh/integration-tests/examples/exec"

func (s *Examples) TestBasicConnection() {
	runner := exec.New(s.Require())

	s.T().Cleanup(func() {
		runner.Run("kubectl delete -f nsc.yaml")
		runner.Run("kubectl delete -f nse.yaml")
	})

	// Setup NSM

	runner.Run("kubectl apply -f namespace.yaml")
	runner.Run("kubectl apply -f registry-service.yaml")
	runner.Run("kubectl apply -f registry-memory.yaml")
	runner.Run("kubectl wait --for=condition=ready pod -l app=nsm-registry --namespace nsm-system")
	runner.Run("kubectl apply -f nsmgr.yaml")
	runner.Run("kubectl wait --for=condition=ready pod -l app=nsmgr --namespace nsm-system")
	runner.Run("kubectl apply -f fake-cross-nse.yaml")
	runner.Run("kubectl wait --for=condition=ready pod -l app=fake-cross-nse  --namespace nsm-system")

	// Setup NSC, NSE

	runner.Run("kubectl apply -f nse.yaml")
	runner.Run("kubectl wait --for=condition=ready pod/nse")
	runner.Run("kubectl apply -f nsc.yaml")
	runner.Run("kubectl wait --for=condition=ready pod/nsc")

	//TODO: Add ping here
}
