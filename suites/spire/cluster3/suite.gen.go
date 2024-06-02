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
		r.Run(`kubectl --kubeconfig=$KUBECONFIG3 delete crd clusterspiffeids.spire.spiffe.io` + "\n" + `kubectl --kubeconfig=$KUBECONFIG3 delete crd clusterfederatedtrustdomains.spire.spiffe.io` + "\n" + `kubectl --kubeconfig=$KUBECONFIG3 delete validatingwebhookconfiguration.admissionregistration.k8s.io/spire-controller-manager-webhook` + "\n" + `kubectl --kubeconfig=$KUBECONFIG3 delete ns spire`)
	})
	r.Run(`[[ ! -z $KUBECONFIG3 ]]`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/spire/cluster3?ref=b6d93d249d50f4e8c0c38b92570536ec53e26f2b`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 wait -n spire --timeout=3m --for=condition=ready pod -l app=spire-server`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 wait -n spire --timeout=1m --for=condition=ready pod -l app=spire-agent`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/b6d93d249d50f4e8c0c38b92570536ec53e26f2b/examples/spire/cluster3/clusterspiffeid-template.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/b6d93d249d50f4e8c0c38b92570536ec53e26f2b/examples/spire/base/clusterspiffeid-webhook-template.yaml`)
}
func (s *Suite) Test() {}
