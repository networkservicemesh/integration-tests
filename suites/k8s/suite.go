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

// Package k8s provides namespace handling logic for the tests
package k8s

import (
	"path"

	"github.com/networkservicemesh/integration-tests/suites/shell"
)

type Suite struct {
	shell.Suite
}

func (s *Suite) Namespace() string {
	dir := s.Dir()
	_, name := path.Split(dir)
	return name
}

func (s *Suite) SetupTest() {
	s.Suite.SetupTest()
	_, err := s.Run("kubectl create ns " + s.Namespace())
	s.Require().NoError(err)

	s.T().Cleanup(func() {
		_, _ = s.Run("kubectl delete ns " + s.Namespace())
	})
}
