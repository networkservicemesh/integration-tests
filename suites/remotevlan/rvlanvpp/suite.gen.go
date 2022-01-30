// Code generated by gotestmd DO NOT EDIT.
package rvlanvpp

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
)

type Suite struct {
	base.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite}
	for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
	r := s.Runner("../deployments-k8s/examples/remotevlan/rvlanvpp")
	r.Run(`kubectl apply -k .`)
}
func (s *Suite) TestKernel2RVlan() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2RVlan")
	s.T().Cleanup(func() {
		r.Run(`docker stop rvm-tester` + "\n" + `docker image rm rvm-tester:latest` + "\n" + `true`)
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/afac44e3497dd0389d1355703f992e2f73e0be47/examples/use-cases/namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`cat > first-iperf-s.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: iperf1-s` + "\n" + `  labels:` + "\n" + `    app: iperf1-s` + "\n" + `spec:` + "\n" + `  replicas: 2` + "\n" + `  selector:` + "\n" + `    matchLabels:` + "\n" + `      app: iperf1-s` + "\n" + `  template:` + "\n" + `    metadata:` + "\n" + `      labels:` + "\n" + `        app: iperf1-s` + "\n" + `      annotations:` + "\n" + `        networkservicemesh.io: kernel://finance-bridge/nsm-1` + "\n" + `    spec:` + "\n" + `      affinity:` + "\n" + `        podAntiAffinity:` + "\n" + `          requiredDuringSchedulingIgnoredDuringExecution:` + "\n" + `          - labelSelector:` + "\n" + `              matchExpressions:` + "\n" + `              - key: app` + "\n" + `                operator: In` + "\n" + `                values:` + "\n" + `                - iperf1-s` + "\n" + `            topologyKey: "kubernetes.io/hostname"` + "\n" + `      containers:` + "\n" + `      - name: iperf-server` + "\n" + `        image: networkstatic/iperf3:latest` + "\n" + `        imagePullPolicy: IfNotPresent` + "\n" + `        command: ["tail", "-f", "/dev/null"]` + "\n" + `EOF`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `resources:` + "\n" + `- first-iperf-s.yaml` + "\n" + `` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl -n ${NAMESPACE} wait --for=condition=ready --timeout=1m pod -l app=iperf1-s`)
	r.Run(`NSCS=($(kubectl get pods -l app=iperf1-s -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'))`)
	r.Run(`IS_FIRST=$(kubectl exec ${NSCS[0]} -c iperf-server -n ${NAMESPACE} -- ip a s nsm-1 | grep 172.10.0.1)` + "\n" + `if [ -n "$IS_FIRST" ]; then` + "\n" + `  kubectl exec ${NSCS[0]} -c iperf-server -n ${NAMESPACE} -- iperf3 -sD -B 172.10.0.1 -1` + "\n" + `  kubectl exec ${NSCS[1]} -c iperf-server -n ${NAMESPACE} -- iperf3 -i0 t 5 -c 172.10.0.1 -B 172.10.0.2` + "\n" + `else` + "\n" + `  kubectl exec ${NSCS[1]} -c iperf-server -n ${NAMESPACE} -- iperf3 -sD -B 172.10.0.1 -1` + "\n" + `  kubectl exec ${NSCS[0]} -c iperf-server -n ${NAMESPACE} -- iperf3 -i0 t 5 -c 172.10.0.1 -B 172.10.0.2` + "\n" + `fi`)
	r.Run(`cat > Dockerfile <<EOF` + "\n" + `FROM networkstatic/iperf3` + "\n" + `` + "\n" + `RUN apt-get update \` + "\n" + `    && apt-get install -y ethtool tcpdump \` + "\n" + `    && rm -rf /var/lib/apt/lists/*` + "\n" + `` + "\n" + `ENTRYPOINT [ "tail", "-f", "/dev/null" ]` + "\n" + `EOF` + "\n" + `docker build . -t rvm-tester`)
	r.Run(`docker run --cap-add=NET_ADMIN --rm -d --network bridge-2 --name rvm-tester rvm-tester tail -f /dev/null` + "\n" + `docker exec rvm-tester ip link set eth0 down` + "\n" + `docker exec rvm-tester ip link add link eth0 name eth0.100 type vlan id 100` + "\n" + `docker exec rvm-tester ip link set eth0 up` + "\n" + `docker exec rvm-tester ip addr add 172.10.0.254/24 dev eth0.100` + "\n" + `docker exec rvm-tester ethtool -K eth0 tx off`)
	r.Run(`docker exec rvm-tester ping -c 1 172.10.0.1`)
	r.Run(`IS_FIRST=$(kubectl exec ${NSCS[0]} -c iperf-server -n ${NAMESPACE} -- ip a s nsm-1 | grep 172.10.0.1)` + "\n" + `if [ -n "$IS_FIRST" ]; then` + "\n" + `  kubectl exec ${NSCS[0]} -c iperf-server -n ${NAMESPACE} -- iperf3 -sD -B 172.10.0.1 -1` + "\n" + `else` + "\n" + `  kubectl exec ${NSCS[1]} -c iperf-server -n ${NAMESPACE} -- iperf3 -sD -B 172.10.0.1 -1` + "\n" + `fi` + "\n" + `docker exec rvm-tester iperf3 -i0 t 5 -c 172.10.0.1`)
	r.Run(`docker exec rvm-tester iperf3 -sD -B 172.10.0.254 -1` + "\n" + `kubectl exec ${NSCS[0]} -c iperf-server -n ${NAMESPACE} -- iperf3 -i0 t 5 -c 172.10.0.254`)
}
