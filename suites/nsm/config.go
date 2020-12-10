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

// Package nsm provides a suit that setups NSM infrastructure for a suite tests. NSM infrastructure setups once on setup suite.
package nsm

// Config is nsm env configuration. Should be used only for debugging goals.
type Config struct {
	Namespace string `default:"nsm-system" desc:"k8s namespace for nsm deployments." split_words:"true."`
	Cleanup   bool   `default:"true" desc:"if true deletes nsm namespace on the suite tear down." split_words:"true."`
}
