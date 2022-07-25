// Code generated by gotestmd DO NOT EDIT.
package spire

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
	r := s.Runner("../deployments-k8s/examples/k8s_monolith/external_nsc/spire")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete crd spiffeids.spiffeid.spiffe.io` + "\n" + `kubectl delete ns spire`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/k8s_monolith/configuration/spire?ref=afc6528948a37cb7ee0ec978743e919550506102`)
	r.Run(`kubectl wait -n spire --timeout=1m --for=condition=ready pod -l app=spire-agent`)
	r.Run(`kubectl wait -n spire --timeout=1m --for=condition=ready pod -l app=spire-server`)
	r.Run(`bundlek8s=$(kubectl exec spire-server-0 -n spire -- bin/spire-server bundle show -format spiffe)` + "\n" + `bundledock=$(docker exec nsc-simple-docker bin/spire-server bundle show -format spiffe)` + "\n" + `echo $bundledock | kubectl exec -i spire-server-0 -n spire -- bin/spire-server bundle set -format spiffe -id "spiffe://docker.nsm/cmd-nsc-simple-docker"` + "\n" + `echo $bundlek8s | docker exec -i nsc-simple-docker bin/spire-server bundle set -format spiffe -id "spiffe://k8s.nsm"`)
}
func (s *Suite) Test() {}
