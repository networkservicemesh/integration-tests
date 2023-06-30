// Copyright (c) 2021 Doc.ai and/or its affiliates.
//
// Copyright (c) 2023 Cisco and/or its affiliates.
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

// Package images provides helpful utilities for searching images for prefetching.
package images

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	fileScheme = "file"
)

// ImageList represents list of open container images
type ImageList struct {
	Images []string
}

// ReteriveList gets list of all images from the source.
// sources can be in format
// 1. Local files: file://..
// 2. Remote gettable content: https://raw.githubusercontent.com/...
// 3. Remote files and dirs via github api: https://api.github.com/repos/...
func ReteriveList(sources []string, match func(string) bool) *ImageList {
	var result = new(ImageList)
	var filesURls []string

	for _, source := range sources {
		filesURls = append(filesURls, reteriveFileList(source, match)...)
	}

	for _, fileURL := range filesURls {
		var content = readContent(fileURL)
		var images = getImages(content)
		result.Images = append(result.Images, images...)
	}

	return result
}

func getImages(content []byte) []string {
	var imagesList ImageList

	_ = yaml.Unmarshal(content, &imagesList)

	if len(imagesList.Images) > 0 {
		return imagesList.Images
	}

	var imagePattern = regexp.MustCompile(".*image: (?P<image>.*)")
	var imageSubexpIndex = imagePattern.SubexpIndex("image")
	var result []string

	if imagePattern.Match(content) {
		for _, match := range imagePattern.FindAllStringSubmatch(string(content), -1) {
			result = append(result, match[imageSubexpIndex])
		}
	}
	return result
}

func readContent(rawurl string) []byte {
	var u, _ = url.Parse(rawurl)
	if u.Scheme == fileScheme {
		var p = filepath.Join(u.Hostname(), u.Path)
		b, err := os.ReadFile(filepath.Clean(p))
		if err == nil {
			return b
		}
		return nil
	}
	if u.Scheme == "http" || u.Scheme == "https" {
		// #nosec
		resp, err := http.Get(rawurl)
		if err != nil {
			return nil
		}
		defer func() {
			_ = resp.Body.Close()
		}()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil
		}
		return b
	}

	return nil
}

func reteriveLocalFileList(rawurl string, match func(string) bool) []string {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil
	}

	basePath := filepath.Join(u.Hostname(), u.Path)

	root, err := os.Open(filepath.Clean(basePath))
	if err != nil {
		return nil
	}

	stat, err := root.Stat()
	if err != nil {
		return nil
	}

	if !stat.IsDir() {
		return []string{fmt.Sprintf("%v://%v", fileScheme, basePath)}
	}

	var result []string

	files, err := os.ReadDir(basePath)
	if err != nil {
		return nil
	}

	for _, f := range files {
		var p = fmt.Sprintf("%v://%v", fileScheme, filepath.Join(basePath, f.Name()))
		if f.IsDir() {
			result = append(result, reteriveFileList(p, match)...)
		} else if match(f.Name()) {
			result = append(result, p)
		}
	}
	return result
}

func reteriveGithubFileList(rawurl string, match func(string) bool) []string {
	var b = readContent(rawurl)

	if b == nil {
		return nil
	}

	var objects []map[string]interface{}
	var result []string

	if err := json.Unmarshal(b, &objects); err != nil {
		var object map[string]interface{}
		if err = json.Unmarshal(b, &object); err != nil {
			return nil
		}
		objects = append(objects, object)
	}

	for _, obj := range objects {
		if _, ok := obj["path"]; !ok {
			continue
		}
		if _, ok := obj["type"]; !ok {
			continue
		}
		var p = obj["path"].(string)

		if obj["type"] == "file" {
			if match(obj["name"].(string)) {
				result = append(result, obj["download_url"].(string))
			}
		}

		if obj["type"] == "dir" {
			var nextContentsURL = apiContentsURL(rawurl, p)
			result = append(result, reteriveFileList(nextContentsURL, match)...)
		}
	}

	return result
}

func reteriveFileList(u string, match func(string) bool) []string {
	if strings.HasPrefix(u, "https://raw.githubusercontent.com") {
		return []string{u}
	}
	if strings.HasPrefix(u, "file://") {
		return reteriveLocalFileList(u, match)
	}
	if strings.HasPrefix(u, "https://api.github.com/repos/") {
		return reteriveGithubFileList(u, match)
	}
	return nil
}

func apiContentsURL(contentsURL, newPath string) string {
	var u, _ = url.Parse(contentsURL)
	var segments = strings.Split(u.Path, string(filepath.Separator))

	segments = append(segments[:5], newPath)
	u.Path = strings.Join(segments, string(filepath.Separator))

	return u.String()
}
