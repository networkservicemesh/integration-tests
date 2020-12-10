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

package shell

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/edwarnicke/exechelper"
	"github.com/kelseyhightower/envconfig"
	"github.com/networkservicemesh/integration-tests/tools/versioning"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

// Suite is testify suite that provides a shell helper functions for each test.
// For each test generates a unique folder.
// Shell for each test located in the unique test folder.
type Suite struct {
	suite.Suite
	config  Config
	options []*exechelper.Option
	dir     string
}

// SetupSuite setups shell suite
func (s *Suite) SetupSuite() {
	err := envconfig.Usage("shell", &s.config)
	s.Require().NoError(err)
	err = envconfig.Process("shell", &s.config)
	s.Require().NoError(err)
	s.options = []*exechelper.Option{exechelper.WithDir(versioning.Dir())}
}

// Dir returns current test directory
func (s *Suite) Dir() string {
	return s.dir
}

// SetupTest creates a unique dir for the test
func (s *Suite) SetupTest() {
	dir, err := ioutil.TempDir(versioning.Dir(), "nsm-test-*")
	s.Require().NoError(err)
	s.options = []*exechelper.Option{exechelper.WithDir(dir)}
	s.dir = dir
	s.T().Cleanup(func() {
		_ = os.RemoveAll(dir)
		_ = os.Remove(dir)
	})
}

// Runf runs shell command with specific format
func (s *Suite) Runf(cmd string, args ...interface{}) (string, error) {
	return s.Run(fmt.Sprintf(cmd, args...))
}

// Run runs shell command. If this function is using in SetupTest then the location is the deployments repository.
// Otherwise the location is a unique test folder.
func (s *Suite) Run(cmd string, additionalOptions ...*exechelper.Option) (output string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.Timeout)
	defer cancel()

	result := new(strings.Builder)
	out := &groupWriter{writers: []io.Writer{result, logrus.StandardLogger().Out}}

	runOptions := append(s.options,
		exechelper.WithStderr(logrus.StandardLogger().Out),
		exechelper.WithStdout(out))

	runOptions = append(runOptions, additionalOptions...)

	cmd = fmt.Sprintf("sh -c \"%v\"", cmd)

	if strings.IndexAny(cmd, "\n") > 0 {
		cmd = fmt.Sprintf("sh -c '''%v'''", cmd)
	}

	for ctx.Err() == nil {
		err = exechelper.Run(cmd, runOptions...)
		if err == nil {
			return result.String(), nil
		}
		time.Sleep(s.config.Timeout / 10)
	}
	return "", err
}

type groupWriter struct {
	writers []io.Writer
}

func (w *groupWriter) Write(p []byte) (n int, err error) {
	for _, writer := range w.writers {
		n, err = writer.Write(p)
	}
	return n, err
}
