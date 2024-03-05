// Code generated by gotestmd DO NOT EDIT.
package docker

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
	r := s.Runner("../deployments-k8s/examples/k8s_monolith/external_nsc/docker")
	s.T().Cleanup(func() {
		r.Run(`docker compose -f docker-compose.yaml -f docker-compose.override.yaml down`)
		r.Run(`rm docker-compose.yaml`)
	})
	r.Run(`cat > docker-compose.override.yaml <<EOF` + "\n" + `---` + "\n" + `networks:` + "\n" + `  kind:` + "\n" + `    external: true` + "\n" + `` + "\n" + `services:` + "\n" + `  nsc-simple-docker:` + "\n" + `    networks:` + "\n" + `      - kind` + "\n" + `    environment:` + "\n" + `      NSM_NETWORK_SERVICES: kernel://kernel2ip2kernel-monolith-nsc@k8s.nsm/nsm-1` + "\n" + `EOF`)
	r.Run(`curl https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/b2e0469c921138e8e2e527b3299d1ce98222b44d/apps/nsc-simple-docker/docker-compose.yaml -o docker-compose.yaml`)
	r.Run(`docker compose -f docker-compose.yaml -f docker-compose.override.yaml up -d`)
}
func (s *Suite) Test() {}
