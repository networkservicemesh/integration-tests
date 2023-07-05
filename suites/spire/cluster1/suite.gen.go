// Code generated by gotestmd DO NOT EDIT.
package cluster1

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
	r := s.Runner("../deployments-k8s/examples/spire/cluster1")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete crd clusterspiffeids.spire.spiffe.io` + "\n" + `kubectl --kubeconfig=$KUBECONFIG1 delete crd clusterfederatedtrustdomains.spire.spiffe.io` + "\n" + `kubectl --kubeconfig=$KUBECONFIG1 delete validatingwebhookconfiguration.admissionregistration.k8s.io/spire-controller-manager-webhook` + "\n" + `kubectl --kubeconfig=$KUBECONFIG1 delete ns spire`)
	})
	r.Run(`[[ ! -z $KUBECONFIG1 ]]`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/spire/cluster1?ref=4ccb4e6afca027ba6822b5dc849d716aaeae74d2`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait -n spire --timeout=1m --for=condition=ready pod -l app=spire-server`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait -n spire --timeout=1m --for=condition=ready pod -l app=spire-agent`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/4ccb4e6afca027ba6822b5dc849d716aaeae74d2/examples/spire/cluster1/clusterspiffeid-template.yaml`)
}
func (s *Suite) Test() {}
