// Code generated by gotestmd DO NOT EDIT.
package nsm_system

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/spire/single_cluster"
)

type Suite struct {
	base.Suite
	single_clusterSuite single_cluster.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite, &s.single_clusterSuite}
	for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
	r := s.Runner("../deployments-k8s/examples/observability/nsm_system")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete mutatingwebhookconfiguration nsm-mutating-webhook` + "\n" + `kubectl delete ns nsm-system`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/observability/nsm_system?ref=11869f163c260a8492c8d25e7abd4902542ad0d3`)
	r.Run(`WH=$(kubectl get pods -l app=admission-webhook-k8s -n nsm-system --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `kubectl wait --for=condition=ready --timeout=1m pod ${WH} -n nsm-system`)
}
func (s *Suite) Test() {}
