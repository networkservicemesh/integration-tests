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

import "regexp"

var (
	ignoreSRIOV = true

	ignoreSRIOVPattern = regexp.MustCompile(".*-sriov")
	ignoreVFIOPattern  = regexp.MustCompile(".*-vfio")
)

func ignored() (ignoreList []*regexp.Regexp) {
	if ignoreSRIOV {
		ignoreList = append(ignoreList, ignoreSRIOVPattern, ignoreVFIOPattern)
	}
	return ignoreList
}

// WithSRIOV enables prefetching SR-IOV test applications
func WithSRIOV() {
	ignoreSRIOV = false
}
