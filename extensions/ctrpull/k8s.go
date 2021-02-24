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

const createNamespace = `
cat > ctr-pull-namespace.yaml <<EOF
---
apiVersion: v1
kind: Namespace
metadata:
  name: ctr-pull
EOF
`

const createConfigMap = `
cat > ctr-pull-configmap.yaml <<EOF
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: ctr-pull
data:
  ctr-pull.sh: |
    #!/bin/sh

    for image in {{.TestImages}}; do
      if ! ctr -n=k8s.io image ls -q | grep "\${image}"; then
        ctr -n=k8s.io image pull "docker.io/\${image}"
      fi
    done
EOF
`

const createDaemonSet = `
cat > ctr-pull.yaml <<EOF
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: ctr-pull
  labels:
    app: ctr-pull
spec:
  selector:
    matchLabels:
      app: ctr-pull
  template:
    metadata:
      labels:
        app: ctr-pull
    spec:
      initContainers:
        - name: ctr-pull
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
            name: ctr-pull
EOF
`

const createKustomization = `
cat > kustomization.yaml <<EOF
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: ctr-pull

resources:
- ctr-pull-namespace.yaml
- ctr-pull-configmap.yaml
- ctr-pull.yaml
EOF
`
