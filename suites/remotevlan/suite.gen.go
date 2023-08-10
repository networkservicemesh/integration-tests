// Code generated by gotestmd DO NOT EDIT.
package remotevlan

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/remotevlan/rvlanovs"
	"github.com/networkservicemesh/integration-tests/suites/remotevlan/rvlanvpp"
	"github.com/networkservicemesh/integration-tests/suites/spire/single_cluster"
)

type Suite struct {
	base.Suite
	single_clusterSuite single_cluster.Suite
	rvlanovsSuite       rvlanovs.Suite
	rvlanvppSuite       rvlanvpp.Suite
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
	r := s.Runner("../deployments-k8s/examples/remotevlan")
	s.T().Cleanup(func() {
		r.Run(`WH=$(kubectl get pods -l app=admission-webhook-k8s -n nsm-system --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `kubectl delete mutatingwebhookconfiguration ${WH}` + "\n" + `kubectl delete ns nsm-system`)
		r.Run(`docker network disconnect bridge-2 kind-worker` + "\n" + `docker network disconnect bridge-2 kind-worker2` + "\n" + `docker network rm bridge-2` + "\n" + `docker exec kind-worker ip link del ext_net1` + "\n" + `docker exec kind-worker2 ip link del ext_net1` + "\n" + `true`)
	})
	r.Run(`docker network create bridge-2` + "\n" + `docker network connect bridge-2 kind-worker` + "\n" + `docker network connect bridge-2 kind-worker2`)
	r.Run(`MACS=($(docker inspect --format='{{range .Containers}}{{.MacAddress}}{{"\n"}}{{end}}' bridge-2))` + "\n" + `ifw1=$(docker exec kind-worker ip -o link | grep ${MACS[@]/#/-e } | cut -f1 -d"@" | cut -f2 -d" ")` + "\n" + `ifw2=$(docker exec kind-worker2 ip -o link | grep ${MACS[@]/#/-e } | cut -f1 -d"@" | cut -f2 -d" ")` + "\n" + `` + "\n" + `(docker exec kind-worker ip link set $ifw1 down &&` + "\n" + `docker exec kind-worker ip link set $ifw1 name ext_net1 &&` + "\n" + `docker exec kind-worker ip link set dev ext_net1 mtu 1450 &&` + "\n" + `docker exec kind-worker ip link set ext_net1 up &&` + "\n" + `docker exec kind-worker2 ip link set $ifw2 down &&` + "\n" + `docker exec kind-worker2 ip link set $ifw2 name ext_net1 &&` + "\n" + `docker exec kind-worker2 ip link set dev ext_net1 mtu 1450 &&` + "\n" + `docker exec kind-worker2 ip link set ext_net1 up)`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/remotevlan?ref=fa884d7fc02955ed77798d4b95bd12caa5a366b5`)
	r.Run(`kubectl -n nsm-system wait --for=condition=ready --timeout=2m pod -l app=nse-remote-vlan`)
	r.Run(`WH=$(kubectl get pods -l app=admission-webhook-k8s -n nsm-system --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `kubectl wait --for=condition=ready --timeout=1m pod ${WH} -n nsm-system`)
	s.RunIncludedSuites()
}
func (s *Suite) RunIncludedSuites() {
	s.Run("Rvlanovs", func() {
		suite.Run(s.T(), &s.rvlanovsSuite)
	})
	s.Run("Rvlanvpp", func() {
		suite.Run(s.T(), &s.rvlanvppSuite)
	})
}
func (s *Suite) Test() {}
