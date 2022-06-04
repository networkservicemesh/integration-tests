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
	r := s.Runner("../deployments-k8s/examples/k8s_monolith/docker")
	s.T().Cleanup(func() {
		r.Run(`docker compose -f docker-compose.yaml -f docker-compose.override.yaml down`)
	})
	r.Run(`cat > docker-compose.override.yaml <<EOF` + "\n" + `---` + "\n" + `networks:` + "\n" + `  kind:` + "\n" + `    external: true` + "\n" + `` + "\n" + `services:` + "\n" + `  nse-simple-vl3-docker:` + "\n" + `    networks:` + "\n" + `      - kind` + "\n" + `EOF`)
	r.Run(`curl https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/8d2c08bc176cc71b6939d570e0a25cc366bafba1/apps/nse-simple-vl3-docker/docker-compose.yaml -o docker-compose.yaml`)
	r.Run(`docker compose -f docker-compose.yaml -f docker-compose.override.yaml up -d`)
}
func (s *Suite) Test() {}
