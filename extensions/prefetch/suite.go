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
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/require"

	"github.com/networkservicemesh/gotestmd/pkg/bash"
	"github.com/networkservicemesh/gotestmd/pkg/suites/shell"
	"github.com/networkservicemesh/integration-tests/extensions/prefetch/images"
)

// Config is env config to setup images prefetching.
type Config struct {
	ImagesPerDaemonset int    `default:"10" desc:"Number of images created per DaemonSet" split_words:"true"`
	Timeout            string `default:"10m" desc:"Kubectl rollout status timeout for the DaemonSet" split_words:"true"`
}

// Suite creates `prefetch` daemonset which pulls all test images for all cluster nodes.
type Suite struct {
	shell.Suite
	SourcesURLs []string
}

var once sync.Once

// SetupSuite prefetches docker images for each k8s node.
func (s *Suite) SetupSuite() {
	once.Do(s.initialize)
}

func (s *Suite) initialize() {
	var config Config
	require.NoError(s.T(), envconfig.Usage("prefetch", &config))
	require.NoError(s.T(), envconfig.Process("prefetch", &config))

	prefetchImages := images.ReteriveList(s.SourcesURLs, func(s string) bool {
		return strings.HasSuffix(s, ".yaml") && !IsExcluded(s)
	}).Images

	prefetchImages = removeDuplicates(prefetchImages)

	tmpDir := uuid.NewString()
	require.NoError(s.T(), os.MkdirAll(tmpDir, 0750))
	s.T().Cleanup(func() { _ = os.RemoveAll(tmpDir) })

	r := s.Runner(tmpDir)

	var daemonSets []string
	for d := 0; d*config.ImagesPerDaemonset < len(prefetchImages); d++ {
		var containers string
		for c := 0; c < config.ImagesPerDaemonset && d*config.ImagesPerDaemonset+c < len(prefetchImages); c++ {
			containers += container(uuid.NewString(), prefetchImages[d*config.ImagesPerDaemonset+c])
		}

		r.Run(createDaemonSet(d, containers))

		daemonSets = append(daemonSets, fmt.Sprintf("prefetch-%d", d))
	}

	r.Run("kubectl create ns prefetch")
	s.T().Cleanup(func() {
		r.Run("kubectl describe pods -n prefetch")
	})
	var b, err = bash.New(bash.WithDir(tmpDir))
	s.Require().NoError(err)
	for _, daemonSet := range daemonSets {
		b.Run(fmt.Sprintf("kubectl -n prefetch apply -f %s.yaml", daemonSet))
		b.Run(fmt.Sprintf("kubectl -n prefetch rollout status daemonset/%s --timeout=%s", daemonSet, config.Timeout))
		b.Run("kubectl -n prefetch describe pods")
		b.Run(fmt.Sprintf("kubectl -n prefetch delete -f %s.yaml", daemonSet))
	}
}

func removeDuplicates(source []string) []string {
	var allKeys = make(map[string]bool)
	var result []string
	for _, item := range source {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			result = append(result, item)
		}
	}
	return result
}
