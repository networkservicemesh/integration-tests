// Code generated by gotestmd DO NOT EDIT.
package pss

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/spire/single_cluster_csi"
)

type Suite struct {
	base.Suite
	single_cluster_csiSuite single_cluster_csi.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite, &s.single_cluster_csiSuite}
	for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
	r := s.Runner("../deployments-k8s/examples/pss")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ds/forwarder-vpp -n nsm-system`)
		r.Run(`kubectl delete mutatingwebhookconfiguration nsm-mutating-webhook` + "\n" + `kubectl delete ns nsm-system`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/pss/nsm-system?ref=48f2cd6fcc26d9a1f98368342f98a611abf23e82`)
	r.Run(`WH=$(kubectl get pods -l app=admission-webhook-k8s -n nsm-system --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `kubectl wait --for=condition=ready --timeout=1m pod ${WH} -n nsm-system`)
}
func (s *Suite) TestNginx() {
	r := s.Runner("../deployments-k8s/examples/pss/use-cases/nginx")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-nginx`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/pss/use-cases/nginx?ref=48f2cd6fcc26d9a1f98368342f98a611abf23e82`)
	r.Run(`kubectl wait --for=condition=ready --timeout=5m pod -l app=nse-kernel -n ns-nginx`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nettools -n ns-nginx`)
	r.Run(`kubectl exec pods/nettools -n ns-nginx -- curl 172.16.1.100:8080 | grep -o "<title>Welcome to nginx"`)
}
