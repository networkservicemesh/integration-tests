// Code generated by gotestmd DO NOT EDIT.
package spiffe_federation

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
	r := s.Runner("../deployments-k8s/examples/k8s_monolith/external_nsc/spiffe_federation")
	r.Run(`kubectl apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/c97ffd635e654ee03099ac44602321e7c76a5ded/examples/k8s_monolith/external_nsc/spiffe_federation/clusterspiffeid-template.yaml`)
	r.Run(`kubectl apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/c97ffd635e654ee03099ac44602321e7c76a5ded/examples/spire/base/clusterspiffeid-webhook-template.yaml`)
	r.Run(`bundlek8s=$(kubectl exec spire-server-0 -n spire -- bin/spire-server bundle show -format spiffe)` + "\n" + `bundledock=$(docker exec nsc-simple-docker bin/spire-server bundle show -format spiffe)`)
	r.Run(`echo $bundledock | kubectl exec -i spire-server-0 -n spire -- bin/spire-server bundle set -format spiffe -id "spiffe://docker.nsm/cmd-nsc-simple-docker"` + "\n" + `echo $bundlek8s | docker exec -i nsc-simple-docker bin/spire-server bundle set -format spiffe -id "spiffe://k8s.nsm"`)
}
func (s *Suite) Test() {}
