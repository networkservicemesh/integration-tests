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
	"runtime"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/networkservicemesh/gotestmd/pkg/bash"
	suites "github.com/networkservicemesh/integration-tests"
	"github.com/networkservicemesh/integration-tests/extensions/prefetch/images"
)

func Test_ImagesCanBePrefetched(t *testing.T) {
	t.Cleanup(func() {
		goleak.VerifyNone(t)
	})
	var repo = "networkservicemesh/deployments-k8s"

	var list = images.ReteriveList(
		[]string{
			fmt.Sprintf("https://raw.githubusercontent.com/%v/%v/external-images.yaml", repo, suites.SHA),
			fmt.Sprintf("https://api.github.com/repos/%v/contents/apps?ref=%v", repo, suites.SHA),
		},
		func(s string) bool {
			return strings.HasSuffix(s, ".yaml")
		},
	)

	var errCh = make(chan error, len(list.Images))
	var imageCh = make(chan string, len(list.Images))

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			var b, err = bash.New(bash.WithEnv(os.Environ()))
			if err != nil {
				errCh <- err
				return
			}
			for img := range imageCh {
				var cmd = "docker pull " + img
				var logger = logrus.WithFields(map[string]interface{}{
					"run": cmd,
				})
				var out, errOut, _, cmdErr = b.Run(cmd)
				if out != "" {
					logger.Info(out)
				}
				if errOut != "" {
					logrus.Error(errOut)
				}
				if cmdErr != nil {
					errCh <- cmdErr
					return
				}
			}
			errCh <- nil
		}()
	}

	for _, image := range list.Images {
		imageCh <- image
	}

	close(imageCh)

	for i := 0; i < len(list.Images); i++ {
		require.NoError(t, <-errCh)
	}
}
