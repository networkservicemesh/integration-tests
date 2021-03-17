// Copyright (c) 2021 Doc.ai and/or its affiliates.
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

package multisuite

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/gotestmd/pkg/suites/shell"
)

type Suite struct {
	shell.Suite

	suits []suite.TestingSuite
}

func (s *Suite) WithSuits(suits ...suite.TestingSuite) {
	s.suits = suits
}

func (s *Suite) SetT(t *testing.T) {
	s.Suite.SetT(t)
	for _, p := range s.suits {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(t)
		}
	}
}

func (s *Suite) SetupSuite() {
	for _, p := range s.suits {
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
}

func (s *Suite) SetupTest() {
	for _, p := range s.suits {
		if v, ok := p.(suite.SetupTestSuite); ok {
			v.SetupTest()
		}
	}
}

func (s *Suite) TearDownSuite() {
	for _, p := range s.suits {
		if v, ok := p.(suite.TearDownAllSuite); ok {
			v.TearDownSuite()
		}
	}
}

func (s *Suite) TearDownTest() {
	for _, p := range s.suits {
		if v, ok := p.(suite.TearDownTestSuite); ok {
			v.TearDownTest()
		}
	}
}

func (s *Suite) BeforeTest(suiteName, testName string) {
	for _, p := range s.suits {
		if v, ok := p.(suite.BeforeTest); ok {
			v.BeforeTest(suiteName, testName)
		}
	}
}

func (s *Suite) AfterTest(suiteName, testName string) {
	for _, p := range s.suits {
		if v, ok := p.(suite.AfterTest); ok {
			v.AfterTest(suiteName, testName)
		}
	}
}
