// Code generated by gotestmd DO NOT EDIT.
package external_nse

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/k8s_monolith/configuration/loadbalancer"
	"github.com/networkservicemesh/integration-tests/suites/k8s_monolith/external_nse/dns"
	"github.com/networkservicemesh/integration-tests/suites/k8s_monolith/external_nse/docker"
	"github.com/networkservicemesh/integration-tests/suites/k8s_monolith/external_nse/spiffe_federation"
	"github.com/networkservicemesh/integration-tests/suites/spire/single_cluster"
)

type Suite struct {
	base.Suite
	loadbalancerSuite      loadbalancer.Suite
	dockerSuite            docker.Suite
	dnsSuite               dns.Suite
	single_clusterSuite    single_cluster.Suite
	spiffe_federationSuite spiffe_federation.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite, &s.loadbalancerSuite, &s.dockerSuite, &s.dnsSuite, &s.single_clusterSuite, &s.spiffe_federationSuite}
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
		r.Run(`WH=$(kubectl get pods -l app=admission-webhook-k8s -n nsm-system --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `kubectl delete mutatingwebhookconfiguration ${WH}` + "\n" + `kubectl delete ns nsm-system`)
	})
	r.Run(`kubectl create ns nsm-system`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/k8s_monolith/configuration/cluster?ref=8684eeab1df8c5be333664c1ba8cd1680a4e7451`)
	r.Run(`WH=$(kubectl get pods -l app=admission-webhook-k8s -n nsm-system --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `kubectl wait --for=condition=ready --timeout=1m pod ${WH} -n nsm-system`)
	r.Run(`kubectl get services registry -n nsm-system -o go-template='{{index (index (index (index .status "loadBalancer") "ingress") 0) "ip"}}'`)
}
func (s *Suite) TestKernel2Wireguard2Kernel() {
	r := s.Runner("../deployments-k8s/examples/k8s_monolith/external_nse/usecases/Kernel2Wireguard2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-kernel2wireguard2kernel-monolith-nse`)
	})
	r.Run(`kubectl create ns ns-kernel2wireguard2kernel-monolith-nse`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/k8s_monolith/external_nse/usecases/Kernel2Wireguard2Kernel?ref=8684eeab1df8c5be333664c1ba8cd1680a4e7451`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-kernel2wireguard2kernel-monolith-nse`)
	r.Run(`nscs=$(kubectl  get pods -l app=alpine -o go-template --template="{{range .items}}{{.metadata.name}} {{end}}" -n ns-kernel2wireguard2kernel-monolith-nse)` + "\n" + `[[ ! -z $nscs ]]`)
	r.Run(`(` + "\n" + `for nsc in $nscs` + "\n" + `do` + "\n" + `    ipAddr=$(kubectl exec -n ns-kernel2wireguard2kernel-monolith-nse $nsc -- ifconfig nsm-1) || exit` + "\n" + `    ipAddr=$(echo $ipAddr | grep -Eo 'inet addr:[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}'| cut -c 11-)` + "\n" + `    for pinger in $nscs` + "\n" + `    do` + "\n" + `        echo $pinger pings $ipAddr` + "\n" + `        kubectl exec $pinger -n ns-kernel2wireguard2kernel-monolith-nse -- ping -c2 -i 0.5 $ipAddr || exit` + "\n" + `    done` + "\n" + `done` + "\n" + `)`)
	r.Run(`(` + "\n" + `for nsc in $nscs` + "\n" + `do` + "\n" + `    echo $nsc pings docker-nse` + "\n" + `    kubectl exec -n ns-kernel2wireguard2kernel-monolith-nse $nsc -- ping 169.254.0.1 -c2 -i 0.5  || exit` + "\n" + `done` + "\n" + `)`)
	r.Run(`(` + "\n" + `for nsc in $nscs` + "\n" + `do` + "\n" + `    ipAddr=$(kubectl exec -n ns-kernel2wireguard2kernel-monolith-nse $nsc -- ifconfig nsm-1) || exit` + "\n" + `    ipAddr=$(echo $ipAddr | grep -Eo 'inet addr:[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}'| cut -c 11-)` + "\n" + `    docker exec nse-simple-vl3-docker ping -c2 -i 0.5 $ipAddr || exit` + "\n" + `done` + "\n" + `)`)
}
