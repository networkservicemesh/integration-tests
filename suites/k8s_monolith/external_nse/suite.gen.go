// Code generated by gotestmd DO NOT EDIT.
package external_nse

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/k8s_monolith/external_nse/dns"
	"github.com/networkservicemesh/integration-tests/suites/k8s_monolith/external_nse/docker"
	"github.com/networkservicemesh/integration-tests/suites/k8s_monolith/external_nse/spire"
)

type Suite struct {
	base.Suite
	dockerSuite docker.Suite
	dnsSuite    dns.Suite
	spireSuite  spire.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite, &s.dockerSuite, &s.dnsSuite, &s.spireSuite}
	for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
	r := s.Runner("../deployments-k8s/examples/k8s_monolith/external_nse")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns nsm-system`)
	})
	r.Run(`kubectl create ns nsm-system`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/k8s_monolith/configuration/cluster?ref=52b14acc34af63586e6665f509f73745cf1e85d9`)
	r.Run(`kubectl get services registry -n nsm-system -o go-template='{{index (index (index (index .status "loadBalancer") "ingress") 0) "ip"}}'`)
}
func (s *Suite) TestKernel2Wireguard2Kernel() {
	r := s.Runner("../deployments-k8s/examples/k8s_monolith/external_nse/usecases/Kernel2Wireguard2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-kernel2wireguard2kernel-monolith-nse`)
	})
	r.Run(`kubectl create ns ns-kernel2wireguard2kernel-monolith-nse`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/k8s_monolith/external_nse/usecases/Kernel2Wireguard2Kernel?ref=52b14acc34af63586e6665f509f73745cf1e85d9`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-kernel -n ns-kernel2wireguard2kernel-monolith-nse`)
	r.Run(`nscs=$(kubectl  get pods -l app=nsc-kernel -o go-template --template="{{range .items}}{{.metadata.name}} {{end}}" -n ns-kernel2wireguard2kernel-monolith-nse)` + "\n" + `[[ ! -z $nscs ]]`)
	r.Run(`for nsc in $nscs` + "\n" + `do` + "\n" + `    ipAddr=$(kubectl exec -n ns-kernel2wireguard2kernel-monolith-nse $nsc -- ifconfig nsm-1)` + "\n" + `    ipAddr=$(echo $ipAddr | grep -Eo 'inet addr:[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}'| cut -c 11-)` + "\n" + `    for pinger in $nscs` + "\n" + `    do` + "\n" + `        echo $pinger pings $ipAddr` + "\n" + `        kubectl exec $pinger -n ns-kernel2wireguard2kernel-monolith-nse -- ping -c4 $ipAddr` + "\n" + `    done` + "\n" + `done`)
	r.Run(`for nsc in $nscs` + "\n" + `do` + "\n" + `    echo $nsc pings docker-nse` + "\n" + `    kubectl exec -n ns-kernel2wireguard2kernel-monolith-nse $nsc -- ping 169.254.0.1 -c4` + "\n" + `done`)
	r.Run(`for nsc in $nscs` + "\n" + `do` + "\n" + `    ipAddr=$(kubectl exec -n ns-kernel2wireguard2kernel-monolith-nse $nsc -- ifconfig nsm-1)` + "\n" + `    ipAddr=$(echo $ipAddr | grep -Eo 'inet addr:[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}'| cut -c 11-)` + "\n" + `    docker exec nse-simple-vl3-docker ping -c4 $ipAddr` + "\n" + `done`)
}
