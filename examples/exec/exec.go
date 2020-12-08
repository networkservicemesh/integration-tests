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

package exec

import (
	"strings"
	"sync"

	"github.com/edwarnicke/exechelper"
	"github.com/kelseyhightower/envconfig"
	"github.com/networkservicemesh/integration-tests/tools/versioning"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var config Config
var once sync.Once

type Exec struct {
	canFail func(string) bool
	config  *Config
	opts    []*exechelper.Option
	assert  *require.Assertions
}

func New(assert *require.Assertions) *Exec {
	writer := logrus.StandardLogger().Writer()

	result := &Exec{
		config: &config,
		assert: assert,
		opts: []*exechelper.Option{
			exechelper.WithStderr(writer),
			exechelper.WithStdout(writer),
			exechelper.WithDir(versioning.Dir()),
		},
	}

	once.Do(func() {
		err := envconfig.Usage("exec", &config)
		assert.NoError(err)

		err = envconfig.Process("exec", &config)
		assert.NoError(err)

	})

	result.registerTimeoutCommand(func(s string) bool {
		return strings.HasPrefix(s, "kubectl wait")
	})

	return result
}

func (e *Exec) Run(cmd string) {
	cond := func() bool {
		err := exechelper.Run(cmd, e.opts...)
		if !e.canFail(cmd) {
			e.assert.NoError(err)
		}
		return err == nil
	}
	e.assert.Eventually(cond, e.config.Timeout, e.config.Timeout/10)

}

// Sometimes some of the running commands can fail under exec on ci, but can not fail on the user manually run.
// These command should be registered here.
func (e *Exec) registerTimeoutCommand(filter func(string) bool) {
	if filter == nil {
		panic("filter cannot be nil")
	}
	old := e.canFail
	e.canFail = func(s string) bool {
		if filter(s) {
			return true
		}

		if old != nil {
			return old(s)
		}

		return false
	}
}
