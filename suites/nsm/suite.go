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
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/networkservicemesh/integration-tests/suites/spire"
)

// Suite adds to the suite NSM deploy logic
type Suite struct {
	spire.Suite
	config Config
}

// SetupSuite setups nsm once for the suite
func (s *Suite) SetupSuite() {
	s.Suite.SetupSuite()

	err := envconfig.Usage("nsm", &s.config)
	s.Require().NoError(err)
	err = envconfig.Process("nsm", &s.config)
	s.Require().NoError(err)

	nss, err := s.Run("kubectl get ns")
	s.Require().NoError(err)

	if strings.Contains(nss, s.config.Namespace) {
		return
	}

	_, err = s.Run("kubectl apply -k nsm/ ")
	s.Require().NoError(err)

	id := s.RegisterSpireEntry(s.config.Namespace, "default")

	s.T().Cleanup(func() {
		s.DeleteSpireEntry(id)
		_, _ = s.Run("kubectl delete ns " + s.config.Namespace)
	})
}
