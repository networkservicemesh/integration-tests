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

// Package provides basic tests for nsm examples
package basic

import (
	"strings"

	"github.com/networkservicemesh/integration-tests/suites/nsm"
)

// Suite provides kustomization based remote/local connection NSM tests
type Suite struct {
	nsm.Suite
}

// TestLocalConnection deploys nsc, nse on the same node and checks that ping is working
func (s *Suite) TestLocalConnection() {
	out, err := s.Run(`kubectl get nodes -o go-template='{{range .items}}{{ if not .spec.taints  }}{{index .metadata.labels \"kubernetes.io/hostname\"}} {{end}}{{end}}'`)
	s.Require().NoError(err)
	nodes := strings.Split(strings.TrimSpace(out), " ")
	s.Require().Greater(len(nodes), 1)

	// Add kustomization

	_, err = s.Runf(`
cat << EOF > kustomization.yaml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

bases:
- ../examples/kernel-nsc
- ../examples/kernel-nse

namespace: %v

patchesStrategicMerge:
- patch-nsc.yaml
- patch-nse.yaml
EOF`, s.Namespace())

	s.Require().NoError(err)

	// Add patches
	_, err = s.Runf(`
cat << EOF > patch-nsc.yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: nsc
spec:
  nodeSelector: 
    kubernetes.io/hostname: %v
EOF`, nodes[0])

	s.Require().NoError(err)

	_, err = s.Runf(`
cat << EOF > patch-nse.yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: nse
spec:
  nodeSelector: 
    kubernetes.io/hostname: %v
EOF`, nodes[0])

	s.Require().NoError(err)

	_, err = s.Run(`kubectl apply -k .`)
	s.Require().NoError(err)

	_, err = s.Run("kubectl wait --for=condition=ready --timeout=15s pod -l app=nse -n " + s.Namespace())
	s.Require().NoError(err)
	_, err = s.Run("kubectl wait --for=condition=ready --timeout=15s pod -l app=nsc -n " + s.Namespace())
	s.Require().NoError(err)

	// TODO: Replace this to ping
	_, err = s.Runf("kubectl logs nsc -n %v | grep \"All client init operations are done.\"", s.Namespace())
	s.Require().NoError(err)
}

// TestRemoteConnection deploys nsc, nse on the different nodes and checks that ping is working
func (s *Suite) TestRemoteConnection() {
	out, err := s.Run(`kubectl get nodes -o go-template='{{range .items}}{{ if not .spec.taints  }}{{index .metadata.labels \"kubernetes.io/hostname\"}} {{end}}{{end}}'`)
	s.Require().NoError(err)
	nodes := strings.Split(strings.TrimSpace(out), " ")
	s.Require().Greater(len(nodes), 1)

	// Add kustomization

	_, err = s.Runf(`
cat << EOF > kustomization.yaml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

bases:
- ../examples/kernel-nsc
- ../examples/kernel-nse

namespace: %v

patchesStrategicMerge:
- patch-nsc.yaml
- patch-nse.yaml
EOF`, s.Namespace())

	s.Require().NoError(err)
	// Add patches
	_, err = s.Runf(`
cat << EOF > patch-nsc.yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: nsc
spec:
  nodeSelector:
    kubernetes.io/hostname: %v
EOF`, nodes[0])

	s.Require().NoError(err)

	_, err = s.Runf(`
cat << EOF > patch-nse.yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: nse
spec:
  nodeSelector:
    kubernetes.io/hostname: %v
EOF`, nodes[1])

	s.Require().NoError(err)

	_, err = s.Run(`kubectl apply -k .`)
	s.Require().NoError(err)

	_, err = s.Run("kubectl wait --for=condition=ready --timeout=15s pod -l app=nse -n " + s.Namespace())
	s.Require().NoError(err)
	_, err = s.Run("kubectl wait --for=condition=ready --timeout=15s pod -l app=nsc -n " + s.Namespace())
	s.Require().NoError(err)

	// TODO: Replace this to ping
	_, err = s.Runf("kubectl logs nsc -n %v | grep \"All client init operations are done.\"", s.Namespace())
	s.Require().NoError(err)
}
