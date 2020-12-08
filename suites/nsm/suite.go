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

package nsm

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/networkservicemesh/integration-tests/suites/spire"
	"github.com/networkservicemesh/integration-tests/tools/k8s"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Suite is testify.Suite that setups NSM infrastructure for all nodes once.
type Suite struct {
	spire.Suite
	config Config
}

// SetupSuite setups NSM
func (s *Suite) SetupSuite() {
	s.Suite.SetupSuite()

	s.Require().NoError(envconfig.Usage("nsm", &s.config))
	s.Require().NoError(envconfig.Process("nsm", &s.config))

	s.T().Cleanup(func() {
		if !s.config.Cleanup {
			return
		}
		err := s.Client().CoreV1().Namespaces().Delete(s.config.Namespace, &metav1.DeleteOptions{})
		s.Require().NoError(err)
	})

	t := k8s.NewK8sTesting(s.T())
	t.SetNamespace(s.config.Namespace)

	_, err := s.Client().CoreV1().Namespaces().Get(s.config.Namespace, metav1.GetOptions{})
	if err == nil {
		return
	}

	_, err = s.Client().CoreV1().Namespaces().Create(&v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: s.config.Namespace,
		},
	})
	s.Require().NoError(err)

	s.RegisterNamespace(s.config.Namespace)

	t.ApplyService("registry-service.yaml")
	t.ApplyDeployment("registry-memory.yaml")
	t.ApplyDaemonSet("nsmgr.yaml")
	t.ApplyDaemonSet("fake-cross-nse.yaml")
}
