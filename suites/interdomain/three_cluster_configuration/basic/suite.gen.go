// Code generated by gotestmd DO NOT EDIT.
package basic

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/interdomain/three_cluster_configuration/dns"
	"github.com/networkservicemesh/integration-tests/suites/interdomain/three_cluster_configuration/loadbalancer"
	"github.com/networkservicemesh/integration-tests/suites/interdomain/three_cluster_configuration/spiffe_federation"
	"github.com/networkservicemesh/integration-tests/suites/spire/cluster1"
	"github.com/networkservicemesh/integration-tests/suites/spire/cluster2"
	"github.com/networkservicemesh/integration-tests/suites/spire/cluster3"
)

type Suite struct {
	base.Suite
	loadbalancerSuite      loadbalancer.Suite
	dnsSuite               dns.Suite
	cluster1Suite          cluster1.Suite
	cluster2Suite          cluster2.Suite
	cluster3Suite          cluster3.Suite
	spiffe_federationSuite spiffe_federation.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite, &s.loadbalancerSuite, &s.dnsSuite, &s.cluster1Suite, &s.cluster2Suite, &s.cluster3Suite, &s.spiffe_federationSuite}
	for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
	r := s.Runner("../deployments-k8s/examples/interdomain/three_cluster_configuration/basic")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete mutatingwebhookconfiguration nsm-mutating-webhook` + "\n" + `kubectl --kubeconfig=$KUBECONFIG1 delete ns nsm-system`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete mutatingwebhookconfiguration nsm-mutating-webhook` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 delete ns nsm-system`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG3 delete ns nsm-system`)
	})
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/three_cluster_configuration/basic/cluster1?ref=1c78ec010b85613f4e0961c03a881a15a6484f35`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/three_cluster_configuration/basic/cluster2?ref=1c78ec010b85613f4e0961c03a881a15a6484f35`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/three_cluster_configuration/basic/cluster3?ref=1c78ec010b85613f4e0961c03a881a15a6484f35`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 get services nsmgr-proxy -n nsm-system -o go-template='{{index (index (index (index .status "loadBalancer") "ingress") 0) "ip"}}'`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 get services nsmgr-proxy -n nsm-system -o go-template='{{index (index (index (index .status "loadBalancer") "ingress") 0) "ip"}}'`)
	r.Run(`WH=$(kubectl --kubeconfig=$KUBECONFIG1 get pods -l app=admission-webhook-k8s -n nsm-system --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `kubectl --kubeconfig=$KUBECONFIG1 wait --for=condition=ready --timeout=1m pod ${WH} -n nsm-system`)
	r.Run(`WH=$(kubectl --kubeconfig=$KUBECONFIG2 get pods -l app=admission-webhook-k8s -n nsm-system --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 wait --for=condition=ready --timeout=1m pod ${WH} -n nsm-system`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 get services registry -n nsm-system -o go-template='{{index (index (index (index .status "loadBalancer") "ingress") 0) "ip"}}'`)
}
func (s *Suite) Test() {}
