// Code generated by gotestmd DO NOT EDIT.
package external_nsc

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/k8s_monolith/external_nsc/dns"
	"github.com/networkservicemesh/integration-tests/suites/k8s_monolith/external_nsc/docker"
	"github.com/networkservicemesh/integration-tests/suites/k8s_monolith/external_nsc/spire"
)

type Suite struct {
	base.Suite
	dockerSuite docker.Suite
	dnsSuite    dns.Suite
	spireSuite  spire.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite, &s.dockerSuite, &s.dnsSuite, &s.spireSuite}
	for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
	r := s.Runner("../deployments-k8s/examples/k8s_monolith/external_nsc")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns nsm-system`)
	})
	r.Run(`kubectl create ns nsm-system`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/k8s_monolith/configuration/cluster?ref=ffeaa33de70ebb1bd0ab3a822c5244b6afca8bf0`)
	r.Run(`kubectl get services registry -n nsm-system -o go-template='{{index (index (index (index .status "loadBalancer") "ingress") 0) "ip"}}'`)
}
func (s *Suite) TestKernel2Wireguard2Kernel() {
	r := s.Runner("../deployments-k8s/examples/k8s_monolith/external_nsc/usecases/Kernel2Wireguard2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/ffeaa33de70ebb1bd0ab3a822c5244b6afca8bf0/examples/k8s_monolith/external_nsc/usecases/namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- https://github.com/networkservicemesh/deployments-k8s/apps/nse-kernel?ref=ffeaa33de70ebb1bd0ab3a822c5244b6afca8bf0` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse-kernel` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nse` + "\n" + `          env:` + "\n" + `            - name: NSM_CIDR_PREFIX` + "\n" + `              value: 172.16.1.100/31` + "\n" + `            - name: NSM_PAYLOAD` + "\n" + `              value: IP` + "\n" + `            - name: NSM_SERVICE_NAMES` + "\n" + `              value: icmp-responder-ip` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ${NAMESPACE}`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-kernel -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`docker exec nsc-simple-docker ping -c4 172.16.1.100`)
	r.Run(`kubectl exec ${NSE} -n ${NAMESPACE} -- ping -c 4 172.16.1.101`)
}
