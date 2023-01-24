// Code generated by gotestmd DO NOT EDIT.
package cluster3

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
	r := s.Runner("../deployments-k8s/examples/spire/cluster3")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG3 delete crd spiffeids.spiffeid.spiffe.io` + "\n" + `kubectl --kubeconfig=$KUBECONFIG3 delete ns spire`)
	})
	r.Run(`[[ ! -z $KUBECONFIG3 ]]`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/spire/cluster3?ref=47e94a59121886f9f1cbf7d382c6b438e0c28487`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 wait -n spire --timeout=1m --for=condition=ready pod -l app=spire-server`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 wait -n spire --timeout=1m --for=condition=ready pod -l app=spire-agent`)
}
func (s *Suite) Test() {}
