// Code generated by gotestmd DO NOT EDIT.
package dashboard

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
	r := s.Runner("../deployments-k8s/examples/observability/dashboard")
	s.T().Cleanup(func() {
		r.Run(`pkill -f "kubectl port-forward -n nsm-system service/dashboard-backend 3001:3001"` + "\n" + `pkill -f "kubectl port-forward -n nsm-system service/dashboard-ui 3000:3000"` + "\n" + `kubectl delete service/dashboard-ui service/dashboard-backend pod/dashboard -n=nsm-system`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/apps/dashboard?ref=aef6419c496cd5330a58df15e7eb69eee41db5dc`)
	r.Run(`kubectl wait --for=condition=ready pod -l app=dashboard --timeout=5m -n nsm-system`)
	r.Run(`nohup kubectl port-forward -n nsm-system service/dashboard-backend 3001:3001 &`)
	r.Run(`nohup kubectl port-forward -n nsm-system service/dashboard-ui 3000:3000 &`)
}
func (s *Suite) Test() {}
