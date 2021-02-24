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

package ctrpull

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/networkservicemesh/gotestmd/pkg/suites/shell"
)

var once sync.Once

// Suite creates `ctr-pull` daemonset which pulls all test images for all cluster nodes.
type Suite struct {
	shell.Suite
	Dir string
}

func (s *Suite) SetupSuite() {
	once.Do(func() {
		testImages, err := s.findTestImages()
		require.NoError(s.T(), err)

		tmpDir := uuid.NewString()
		require.NoError(s.T(), os.MkdirAll(tmpDir, 0750))

		r := s.Runner(tmpDir)

		r.Run(createNamespace)
		r.Run(strings.ReplaceAll(createConfigMap, "{{.TestImages}}", strings.Join(testImages, " ")))
		r.Run(createDaemonSet)
		r.Run(createKustomization)

		r.Run("kubectl apply -k .")
		r.Run("kubectl -n ctr-pull wait --timeout=10m --for=condition=ready pod -l app=ctr-pull")

		r.Run("kubectl delete ns ctr-pull")
		_ = os.RemoveAll(tmpDir)
	})
}

func (s *Suite) findTestImages() (testImages []string, err error) {
	imagePattern := regexp.MustCompile(".*image: (?P<image>.*)")
	imageSubexpIndex := imagePattern.SubexpIndex("image")

	err = filepath.Walk(filepath.Join(s.Dir, "apps"), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".yaml") {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if imagePattern.MatchString(scanner.Text()) {
				image := imagePattern.FindAllStringSubmatch(scanner.Text(), -1)[0][imageSubexpIndex]
				testImages = append(testImages, image)
			}
		}

		return nil
	})

	return testImages, err
}
