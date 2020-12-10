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

// Package versioning provides helper functions for working with deployments versioning
package versioning

import (
	"os"
	"path"
	"sync"

	"github.com/edwarnicke/exechelper"
	"github.com/sirupsen/logrus"
)

const (
	commit = "c88987aceb80b32b3be33d0346b8b4014fe71252"
	url    = "http://github.com/networkservicemesh/deployments-k8s.git"
)

var once sync.Once
var dir string

// Dir returns dir of repositorty deployments-k8s
func Dir() string {
	once.Do(func() {
		dir = path.Join(os.Getenv("GOPATH"), "src", "github.com", "networkservicemesh", "deployments-k8s")
		var err error

		if _, err = os.Open(path.Clean(dir)); err == nil {
			return
		}

		err = exechelper.Run("git clone " + url + " " + dir)
		if err != nil {
			logrus.Fatal(err.Error())
		}

		err = exechelper.Run("git checkout "+commit, exechelper.WithDir(dir))

		if err != nil {
			logrus.Fatal(err.Error())
		}
	})
	return dir
}
