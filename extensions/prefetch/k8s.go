// Copyright (c) 2021 Doc.ai and/or its affiliates.
//
// Copyright (c) 2024 Cisco and/or its affiliates.
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
	"bytes"
	"text/template"
)

func createDaemonSet(number int, containers string) string {
	const text = `
cat > prefetch-{{.Number}}.yaml <<EOF
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: prefetch-{{.Number}}
  labels:
    app: prefetch-{{.Number}}
spec:
  selector:
    matchLabels:
      app: prefetch-{{.Number}}
  template:
    metadata:
      labels:
        app: prefetch-{{.Number}}
    spec:
      initContainers:
        # return container is used to pass the return application to other containers.
        # containers use it to exit immediately after loading the image.
        - name: return
          image: ghcr.io/networkservicemesh/cmd-return
          imagePullPolicy: IfNotPresent
          command: ["cp", "/bin/return", "/out/return"]
          volumeMounts:
            - name: bin
              mountPath: /out
{{.Containers}}
      containers:
        - name: pause
          image: registry.k8s.io/pause:3.9
      volumes:
        - name: bin
          emptyDir: { }
EOF
`
	return substitute(text, &struct {
		Number     int
		Containers string
	}{
		Number:     number,
		Containers: containers,
	})
}

func container(name, image string) string {
	const text = `
        - name: {{.Name}}
          image: {{.Image}}
          imagePullPolicy: IfNotPresent
          command: ["/bin/return"]
          volumeMounts:
            - name: bin
              mountPath: /bin
`
	return substitute(text, &struct {
		Name, Image string
	}{
		Name:  name,
		Image: image,
	})
}

func substitute(text string, data interface{}) string {
	t, _ := template.New("").Parse(text)

	buf := new(bytes.Buffer)
	_ = t.Execute(buf, data)

	return buf.String()
}
