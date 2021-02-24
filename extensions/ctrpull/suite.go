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

const (
	defaultDomain = "docker.io/"
	defaultTag    = ":latest"
)

// Suite creates `ctr-pull` daemonset which pulls all test images for all cluster nodes.
type Suite struct {
	shell.Suite
	Dir string
}

var once sync.Once

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

func (s *Suite) findTestImages() ([]string, error) {
	imagePattern := regexp.MustCompile(".*image: (?P<image>.*)")
	imageSubexpIndex := imagePattern.SubexpIndex("image")

	var testImages []string
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if ok, skipErr := s.shouldSkipWithError(info, err); ok {
			return skipErr
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if imagePattern.MatchString(scanner.Text()) {
				image := imagePattern.FindAllStringSubmatch(scanner.Text(), -1)[0][imageSubexpIndex]
				testImages = append(testImages, s.fullImageName(image))
			}
		}

		return nil
	}

	if err := filepath.Walk(filepath.Join(s.Dir, "apps"), walkFunc); err != nil {
		return nil, err
	}
	if err := filepath.Walk(filepath.Join(s.Dir, "examples", "spire"), walkFunc); err != nil {
		return nil, err
	}

	return testImages, nil
}

func (s *Suite) shouldSkipWithError(info os.FileInfo, err error) (bool, error) {
	if err != nil {
		return true, err
	}

	if info.IsDir() {
		if _, ok := ignored[info.Name()]; ok {
			return true, filepath.SkipDir
		}
		return true, nil
	}

	if !strings.HasSuffix(info.Name(), ".yaml") {
		return true, nil
	}

	return false, nil
}

func (s *Suite) fullImageName(image string) string {
	// domain/library/name:tag

	split := strings.Split(image, "/")
	switch len(split) {
	case 3:
		// nothing to do
	case 2:
		image = defaultDomain + image
	default:
		return ""
	}

	split = strings.Split(image, ":")
	switch len(split) {
	case 2:
		// nothing to do
	case 1:
		image += defaultTag
	default:
		return ""
	}

	return image
}
