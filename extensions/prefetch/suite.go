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

// Package prefetch exports suite that can do prefetch of required images once per suite.
package prefetch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/networkservicemesh/gotestmd/pkg/suites/shell"
)

const (
	defaultDomain = "docker.io"
	officialLib   = "library"
	defaultTag    = ":latest"
)

// Suite creates `prefetch` daemonset which pulls all test images for all cluster nodes.
type Suite struct {
	shell.Suite
	Dir string
}

var once sync.Once

// SetupSuite prefetches docker images for each k8s node.
func (s *Suite) SetupSuite() {
	once.Do(func() {
		images := s.findImages()

		tmpDir := uuid.NewString()
		require.NoError(s.T(), os.MkdirAll(tmpDir, 0750))
		s.T().Cleanup(func() { _ = os.RemoveAll(tmpDir) })

		r := s.Runner(tmpDir)

		var daemonSets []string
		for d := 0; d*10 < len(images); d++ {
			var containers string
			for c := 0; c < 10 && d*10+c < len(images); c++ {
				containers += container(uuid.NewString(), images[d*10+c])
			}

			r.Run(createDaemonSet(d, containers))

			daemonSets = append(daemonSets, fmt.Sprintf("prefetch-%d", d))
		}

		r.Run("kubectl create ns prefetch")
		s.T().Cleanup(func() { r.Run("kubectl delete ns prefetch") })

		var wg sync.WaitGroup
		for _, daemonSet := range daemonSets {
			wg.Add(1)
			go func(daemonSet string) {
				defer wg.Done()

				dr := s.Runner(tmpDir)
				dr.Run(fmt.Sprintf("kubectl -n prefetch apply -f %s.yaml", daemonSet))
				dr.Run(fmt.Sprintf("kubectl -n prefetch rollout status daemonset/%s --timeout=5m", daemonSet))
				dr.Run(fmt.Sprintf("kubectl -n prefetch delete -f %s.yaml", daemonSet))
			}(daemonSet)
		}
		wg.Wait()
	})
}

func (s *Suite) findImages() []string {
	rawImages, err := find(filepath.Join(s.Dir, "examples"), filepath.Join(s.Dir, "apps"))
	require.NoError(s.T(), err)

	preparedImages := make(map[string]struct{})
	for image := range rawImages {
		preparedImages[s.fullImageName(image)] = struct{}{}
	}

	var images []string
	for image := range preparedImages {
		if image != "" {
			images = append(images, image)
		}
	}

	return images
}

func (s *Suite) fullImageName(image string) string {
	var domain, remainder string
	i := strings.IndexRune(image, '/')
	if i == -1 || (!strings.ContainsAny(image[:i], ".:")) {
		domain, remainder = defaultDomain, image
	} else {
		domain, remainder = image[:i], image[i+1:]
	}
	if domain == defaultDomain && !strings.ContainsRune(remainder, '/') {
		remainder = officialLib + "/" + remainder
	}

	switch len(strings.Split(remainder, ":")) {
	case 2:
		// nothing to do
	case 1:
		remainder += defaultTag
	default:
		return ""
	}

	return domain + "/" + remainder
}
