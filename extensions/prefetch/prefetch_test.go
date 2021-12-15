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

// +build linux

package prefetch_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/networkservicemesh/gotestmd/pkg/bash"
	suites "github.com/networkservicemesh/integration-tests"
	"github.com/networkservicemesh/integration-tests/extensions/prefetch/images"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func Test_ImagesCanBePrefetch(t *testing.T) {
	var repo = "networkservicemesh/deployments-k8s"
	var b, err = bash.New(bash.WithEnv(os.Environ()))
	require.NoError(t, err)

	var list = images.ReteriveList(
		[]string{
			fmt.Sprintf("https://raw.githubusercontent.com/%v/%v/external-images.yaml", repo, suites.SHA),
			fmt.Sprintf("https://api.github.com/repos/%v/contents/apps?ref=%v", repo, suites.SHA),
		},
		func(s string) bool {
			return strings.HasSuffix(s, ".yaml")
		},
	)

	for _, img := range list.Images {
		var cmd = "docker pull " + img
		var logger = logrus.WithFields(map[string]interface{}{
			"run": cmd,
		})
		var out, errOut, exitCode, cmdErr = b.Run(cmd)
		if out != "" {
			logger.Info(out)
		}
		if errOut != "" {
			logrus.Error(errOut)
		}
		require.Zero(t, exitCode)
		require.NoError(t, cmdErr)
	}

}
