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

// Package k8s provides kubernetes helper functions
package k8s

import (
	corev1 "k8s.io/api/core/v1"
)

// SetNode sets NodeSelector for the pod based on passed nodeName
func SetNode(nodeName string) func(*corev1.Pod) {
	return func(pod *corev1.Pod) {
		pod.Spec.NodeSelector = map[string]string{
			"kubernetes.io/hostname": nodeName,
		}
	}
}
