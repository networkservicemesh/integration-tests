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

const createNamespace = `
cat >prefetch-namespace.yaml <<EOF
---
apiVersion: v1
kind: Namespace
metadata:
  name:prefetch
EOF
`

const createConfigMap = `
cat >prefetch-configmap.yaml <<EOF
---
apiVersion: v1
kind: ConfigMap
metadata:
  name:prefetch
data:
 prefetch.sh: |
    #!/bin/sh

    for image in {{.TestImages}}; do
      if ! ctr -n=k8s.io image ls -q | grep "\${image}"; then
        ctr -n=k8s.io image pull "\${image}"
      fi
    done
EOF
`

const createDaemonSet = `
cat >prefetch.yaml <<EOF
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name:prefetch
  labels:
    app:prefetch
spec:
  selector:
    matchLabels:
      app:prefetch
  template:
    metadata:
      labels:
        app:prefetch
    spec:
      initContainers:
        - name:prefetch
          image: docker:latest
          imagePullPolicy: IfNotPresent
          command: ["/bin/sh", "/root/scripts/ctr-pull.sh"]
          volumeMounts:
            - name: containerd
              mountPath: /run/containerd/containerd.sock
            - name: scripts
              mountPath: /root/scripts
      containers:
        - name: pause
          image: google/pause:latest
      volumes:
        - name: containerd
          hostPath:
            path: /run/containerd/containerd.sock
        - name: scripts
          configMap:
            name:prefetch
EOF
`

const createKustomization = `
cat > kustomization.yaml <<EOF
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace:prefetch

resources:
-prefetch-namespace.yaml
-prefetch-configmap.yaml
-prefetch.yaml
EOF
`
