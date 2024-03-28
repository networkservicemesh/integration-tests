// Copyright (c) 2024 Pragmagic Inc. and/or its affiliates.
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

package parallel_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/networkservicemesh/integration-tests/extensions/parallel"
)

type negativeSuite struct {
	suite.Suite
}

func (s *negativeSuite) TestParallel1() {}
func (s *negativeSuite) TestParallel2() {}
func (s *negativeSuite) TestParallel3() {}

func (s *negativeSuite) TestSynchronously() {
	s.Error(goleak.Find())
}

type positiveSuite struct {
	suite.Suite
}

func (s *positiveSuite) TestParallel1() {}
func (s *positiveSuite) TestParallel2() {}
func (s *positiveSuite) TestParallel3() {}

func (s *positiveSuite) TestSynchronously() {
	s.NoError(goleak.Find())
}

func Test_OptionWithRunningTestsSynchronously_ShouldExcludeTestsFromParallelExecution(t *testing.T) {
	var s1 = new(negativeSuite)
	parallel.Run(t, s1)
	var s2 = new(positiveSuite)
	parallel.Run(t, s2, parallel.WithRunningTestsSynchronously(s2.TestSynchronously))
}
