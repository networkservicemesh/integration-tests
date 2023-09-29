// Code generated by gotestmd DO NOT EDIT.
package external_nsc

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/k8s_monolith/configuration/loadbalancer"
	"github.com/networkservicemesh/integration-tests/suites/k8s_monolith/external_nsc/dns"
	"github.com/networkservicemesh/integration-tests/suites/k8s_monolith/external_nsc/docker"
	"github.com/networkservicemesh/integration-tests/suites/k8s_monolith/external_nsc/spiffe_federation"
	"github.com/networkservicemesh/integration-tests/suites/spire/single_cluster"
)

type Suite struct {
	base.Suite
	loadbalancerSuite      loadbalancer.Suite
	dockerSuite            docker.Suite
	dnsSuite               dns.Suite
	single_clusterSuite    single_cluster.Suite
	spiffe_federationSuite spiffe_federation.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite, &s.loadbalancerSuite, &s.dockerSuite, &s.dnsSuite, &s.single_clusterSuite, &s.spiffe_federationSuite}
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
		r.Run(`WH=$(kubectl get pods -l app=admission-webhook-k8s -n nsm-system --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `kubectl delete mutatingwebhookconfiguration ${WH}` + "\n" + `kubectl delete ns nsm-system`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/k8s_monolith/configuration/cluster?ref=45a9d9d819424900a4499764ebaf568d1af6dd29`)
	r.Run(`kubectl get services registry -n nsm-system -o go-template='{{index (index (index (index .status "loadBalancer") "ingress") 0) "ip"}}'`)
}
func (s *Suite) TestKernel2IP2Kernel() {
	r := s.Runner("../deployments-k8s/examples/k8s_monolith/external_nsc/usecases/Kernel2IP2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-kernel2ip2kernel-monolith-nsc`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/k8s_monolith/external_nsc/usecases/Kernel2IP2Kernel?ref=45a9d9d819424900a4499764ebaf568d1af6dd29`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-kernel2ip2kernel-monolith-nsc`)
	r.Run(`docker exec nsc-simple-docker ping -c4 172.16.1.100`)
	r.Run(`kubectl exec deployments/nse-kernel -n ns-kernel2ip2kernel-monolith-nsc -- ping -c 4 172.16.1.101`)
}
