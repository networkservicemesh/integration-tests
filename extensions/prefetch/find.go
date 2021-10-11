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
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/go-getter"
	"github.com/sirupsen/logrus"
	"mvdan.cc/xurls/v2"

	"github.com/networkservicemesh/integration-tests/internal/git"
)

const kustomization = "kustomization.yaml"

var (
	imagePattern     = regexp.MustCompile(".*image: (?P<image>.*)")
	imageSubexpIndex = imagePattern.SubexpIndex("image")
	urlPattern       = xurls.Relaxed()
)

type finder struct {
	tmpRoot string
	images  map[string]struct{}
	urls    map[string]string
}

func find(sources ...string) (map[string]struct{}, error) {
	f := &finder{
		images: make(map[string]struct{}),
		urls:   make(map[string]string),
	}

	tmpRoot, err := os.MkdirTemp("", "prefetch-")
	if err != nil {
		return nil, err
	}
	defer func() { _ = os.RemoveAll(tmpRoot) }()

	return f.find(sources)
}

func (f *finder) find(baseSources []string) (map[string]struct{}, error) {
	var files []string
	for _, source := range baseSources {
		newFiles, err := f.processSource(source)
		if err != nil {
			return nil, err
		}
		files = append(files, newFiles...)
	}

	for ; len(files) > 0; files = files[1:] {
		newFiles, err := f.processFile(files[0])
		if err != nil {
			return nil, err
		}
		files = append(files, newFiles...)
	}

	return f.images, nil
}

func (f *finder) processSource(source string) (files []string, err error) {
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

		if strings.HasSuffix(d.Name(), ".yaml") || strings.HasSuffix(d.Name(), ".md") {
			newFiles, err := f.processFile(path)
			if err != nil {
				return err
			}
			files = append(files, newFiles...)
		}

		return nil
	})
	return files, err
}

func (f *finder) processFile(file string) (files []string, err error) {
	// #nosec
	in, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() { _ = in.Close() }()

	var base, resource, skip bool
	var newFiles []string
	for scanner := bufio.NewScanner(in); err == nil && scanner.Scan(); files = append(files, newFiles...) {
		newFiles = nil

		raw := scanner.Text()
		trim := strings.TrimSpace(raw)
		if trim == "" {
			continue
		}

		if base, resource, skip = f.baseResourceSkip(base, resource, raw, trim); skip {
			continue
		}

		newFiles, err = f.processLine(file, trim, base, resource)
	}
	return files, err
}

func (f *finder) baseResourceSkip(base, resource bool, raw, trim string) (newBase, newResource, skip bool) {
	if !base && strings.HasPrefix(raw, "bases:") {
		return true, false, true
	}
	if !resource && strings.HasPrefix(raw, "resources:") {
		return false, true, true
	}
	if (base || resource) && !strings.HasPrefix(trim, "-") {
		return false, false, false
	}
	return base, resource, false
}

func (f *finder) processLine(file, trim string, base, resource bool) ([]string, error) {
	if imagePattern.MatchString(trim) {
		f.images[imagePattern.FindStringSubmatch(trim)[imageSubexpIndex]] = struct{}{}
		return nil, nil
	}

	if base || resource || f.isKubectl(trim) {
		if urlPattern.MatchString(trim) {
			return f.processURL(urlPattern.FindString(trim))
		}
	}

	if base || resource {
		dir, _ := filepath.Split(file)

		var newFile string
		if base {
			newFile = filepath.Clean(filepath.Join(dir, strings.TrimSpace(trim[1:]), kustomization))
		}
		if resource {
			newFile = filepath.Clean(filepath.Join(dir, strings.TrimSpace(trim[1:])))
		}
		if validateErr := f.validateFile(newFile, logrus.WithField("file", file)); validateErr != nil {
			return nil, nil
		}

		return []string{newFile}, nil
	}

	return nil, nil
}

func (f *finder) isKubectl(trim string) bool {
	return strings.Contains(trim, "kubectl create") ||
		strings.Contains(trim, "kubectl apply") ||
		strings.Contains(trim, "kubectl replace")
}

func (f *finder) processURL(u string) (files []string, err error) {
	logger := logrus.WithField("url", u)

	var filePath string
	var spec *git.RepoSpec
	if strings.HasSuffix(u, ".yaml") {
		var ur *url.URL
		if ur, err = url.Parse(u); err != nil {
			logger.Warn(err.Error())
			return nil, nil
		}
		_, filePath = path.Split(ur.Path)
	} else {
		if spec, err = git.NewRepoSpecFromURL(u); err != nil {
			logger.Warn(err.Error())
			return nil, nil
		}
		u = spec.CloneSpec() + "?ref=" + spec.Ref
		filePath = filepath.Join(filepath.Join(strings.Split(spec.Path, "/")...), kustomization)
	}

	var dir string
	if _, ok := f.urls[u]; !ok {
		dir, err = os.MkdirTemp(f.tmpRoot, "")
		if err != nil {
			return nil, err
		}

		if strings.HasSuffix(u, ".yaml") {
			err = getter.GetFile(filepath.Join(dir, filePath), u)
		} else {
			err = f.loadGitRepo(spec, dir)
		}
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}

		f.urls[u] = dir
	} else {
		dir = f.urls[u]
	}

	files = append(files, filepath.Join(dir, filePath))

	if validateErr := f.validateFile(files[0], logger); validateErr != nil {
		return nil, nil
	}

	return files, nil
}

func (f *finder) loadGitRepo(spec *git.RepoSpec, dir string) error {
	if spec == nil {
		return nil
	}

	cloneURL := "git::" + spec.CloneSpec() + "//."
	if err := getter.Get(dir, cloneURL); err != nil {
		return err
	}

	checkoutURL := "git::" + spec.CloneSpec() + "?ref=" + spec.Ref
	return getter.Get(dir, checkoutURL)
}

func (f *finder) validateFile(file string, logger *logrus.Entry) error {
	if _, err := os.Lstat(file); err != nil {
		logger.Warn(err.Error())
		return err
	}
	return nil
}
