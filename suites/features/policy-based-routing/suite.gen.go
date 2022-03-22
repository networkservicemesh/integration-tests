// Code generated by gotestmd DO NOT EDIT.
package policy_based_routing

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/basic"
)

type Suite struct {
	base.Suite
	basicSuite basic.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite, &s.basicSuite}
	for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
	r := s.Runner("../deployments-k8s/examples/features/policy-based-routing")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/41cd9995434986adcb18e4202be3c552d21485a8/examples/features/namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`NODE=($(kubectl get nodes -o go-template='{{range .items}}{{ if not .spec.taints  }}{{index .metadata.labels "kubernetes.io/hostname"}} {{end}}{{end}}')[0])`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `resources:` + "\n" + `- client.yaml` + "\n" + `- config-file-nse.yaml` + "\n" + `bases:` + "\n" + `- https://github.com/networkservicemesh/deployments-k8s/apps/nse-kernel?ref=41cd9995434986adcb18e4202be3c552d21485a8` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > client.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: v1` + "\n" + `kind: Pod` + "\n" + `metadata:` + "\n" + `  name: alpine` + "\n" + `  labels:` + "\n" + `    app: alpine` + "\n" + `  annotations:` + "\n" + `    networkservicemesh.io: kernel://icmp-responder/nsm-1` + "\n" + `spec:` + "\n" + `  containers:` + "\n" + `  - name: alpine` + "\n" + `    image: artgl/alpine_iproute2:3.15` + "\n" + `    imagePullPolicy: IfNotPresent` + "\n" + `    stdin: true` + "\n" + `    tty: true` + "\n" + `  nodeSelector:` + "\n" + `    kubernetes.io/hostname: ${NODE}` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse-kernel` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nse` + "\n" + `          env:` + "\n" + `            - name: NSM_CIDR_PREFIX` + "\n" + `              value: 172.16.1.100/31` + "\n" + `          volumeMounts:` + "\n" + `            - mountPath: /etc/policy-based-routing/config.yaml` + "\n" + `              subPath: config.yaml` + "\n" + `              name: policies-config-volume` + "\n" + `      volumes:` + "\n" + `        - name: policies-config-volume` + "\n" + `          configMap:` + "\n" + `            name: policies-config-file` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODE}` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ${NAMESPACE}`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ${NAMESPACE}`)
	r.Run(`NSC=$(kubectl get pods -l app=alpine -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-kernel -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl exec ${NSC} -n ${NAMESPACE} -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl exec ${NSE} -n ${NAMESPACE} -- ping -c 4 172.16.1.101`)
	r.Run(`result=$(kubectl exec ${NSC} -n ${NAMESPACE} -- ip r get 172.16.3.1 from 172.16.2.201 ipproto tcp dport 6666)` + "\n" + `echo ${result}` + "\n" + `echo ${result} | grep -E -q "172.16.3.1 from 172.16.2.201 via 172.16.2.200 dev nsm-1 table 1"`)
	r.Run(`result=$(kubectl exec ${NSC} -n ${NAMESPACE} -- ip r get 172.16.4.1 ipproto udp dport 6666)` + "\n" + `echo ${result}` + "\n" + `echo ${result} | grep -E -q "172.16.4.1 dev nsm-1 table 2 src 172.16.1.101"`)
	r.Run(`result=$(kubectl exec ${NSC} -n ${NAMESPACE} -- ip -6 route get 2004::5 from 2004::3 ipproto udp dport 5555)` + "\n" + `echo ${result}` + "\n" + `echo ${result} | grep -E -q "via 2004::6 dev nsm-1 table 3 src 2004::3"`)
}
func (s *Suite) Test() {}
