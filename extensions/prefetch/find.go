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

package prefetch

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/go-getter"
	"mvdan.cc/xurls/v2"

	"github.com/networkservicemesh/integration-tests/internal/git"
)

type finder struct {
	tmpRoot      string
	images, urls map[string]struct{}
}

func find(sources ...string) (map[string]struct{}, error) {
	f := &finder{
		images: make(map[string]struct{}),
		urls:   make(map[string]struct{}),
	}

	tmpRoot, err := os.MkdirTemp("", "prefetch-")
	if err != nil {
		return nil, err
	}
	defer func() { _ = os.RemoveAll(tmpRoot) }()

	return f.find(sources)
}

func (f *finder) find(baseSources []string) (map[string]struct{}, error) {
	var sources []string
	for _, source := range baseSources {
		newSources, err := f.processSource(source, true)
		if err != nil {
			return nil, err
		}
		sources = append(sources, newSources...)
	}

	for ; len(sources) > 0; sources = sources[1:] {
		newSources, err := f.processSource(sources[0], false)
		if err != nil {
			return nil, err
		}
		sources = append(sources, newSources...)
	}

	return f.images, nil
}

func (f *finder) processSource(source string, isBaseSource bool) (sources []string, err error) {
	info, err := os.Lstat(source)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return f.processFile(source)
	}

	err = filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if IsExcluded(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if strings.HasSuffix(d.Name(), ".yaml") || isBaseSource && strings.HasSuffix(d.Name(), ".md") {
			sources = append(sources, path)
		}

		return nil
	})
	return sources, err
}

var (
	imagePattern     = regexp.MustCompile(".*image: (?P<image>.*)")
	imageSubexpIndex = imagePattern.SubexpIndex("image")
	urlPattern       = xurls.Strict()
)

func (f *finder) processFile(path string) (sources []string, err error) {
	// #nosec
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var urlAllowed, bashSection = strings.HasSuffix(path, ".yaml"), false
	for scanner := bufio.NewScanner(file); scanner.Scan(); {
		line := scanner.Text()

		if !urlAllowed {
			if !bashSection && strings.Contains(line, "```bash") {
				bashSection = true
			} else if bashSection && strings.Contains(line, "```") {
				bashSection = false
			}
		}

		if imagePattern.MatchString(line) {
			f.images[imagePattern.FindStringSubmatch(line)[imageSubexpIndex]] = struct{}{}
		}

		if (urlAllowed || bashSection) && urlPattern.MatchString(line) {
			if u, ok := f.formatURL(urlPattern.FindString(line)); ok {
				if _, ok = f.urls[u]; !ok {
					println(u)
					f.urls[u] = struct{}{}

					tmpDir, err := os.MkdirTemp(f.tmpRoot, "")
					if err != nil {
						return nil, err
					}

					var processURL func(dst, u string) error
					if strings.HasSuffix(u, ".yaml") {
						processURL = f.processFileURL
					} else {
						processURL = f.processGitURL
					}

					uErr := processURL(tmpDir, u)
					if uErr == nil {
						sources = append(sources, tmpDir)
					} else {
						println(u, ":", uErr.Error())
					}
				}
			}
		}
	}

	return sources, nil
}

func (f *finder) formatURL(u string) (string, bool) {
	if strings.HasSuffix(u, ".yaml") {
		return u, true
	}

	spec, err := git.NewRepoSpecFromURL(u)
	if err == nil {
		return spec.CloneSpec() + "?ref=" + spec.Ref, true
	}

	return "", false
}

func (f *finder) processFileURL(dst, u string) error {
	return getter.GetFile(filepath.Join(dst, "file.yaml"), u)
}

func (f *finder) processGitURL(dst, u string) error {
	spec, err := git.NewRepoSpecFromURL(u)
	if err != nil {
		return err
	}

	cloneURL := "git::" + spec.CloneSpec() + "//."
	if err := getter.Get(dst, cloneURL); err != nil {
		return err
	}

	checkoutURL := "git::" + spec.CloneSpec() + "?ref=" + spec.Ref
	return getter.Get(dst, checkoutURL)
}
