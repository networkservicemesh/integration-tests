// Code generated by gotestmd DO NOT EDIT.
package ipsec_mechanism

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/spire/single_cluster"
)

type Suite struct {
	base.Suite
	single_clusterSuite single_cluster.Suite
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
	r := s.Runner("../deployments-k8s/examples/ipsec_mechanism")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete mutatingwebhookconfiguration nsm-mutating-webhook` + "\n" + `kubectl delete ns nsm-system`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/ipsec_mechanism?ref=09d7c7a298f58007f2a76f3cfaa8cc08a36fc147`)
	r.Run(`WH=$(kubectl get pods -l app=admission-webhook-k8s -n nsm-system --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `kubectl wait --for=condition=ready --timeout=1m pod ${WH} -n nsm-system`)
}
func (s *Suite) TestKernel2IP2Kernel() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2IP2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-kernel2ip2kernel`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Kernel2IP2Kernel?ref=09d7c7a298f58007f2a76f3cfaa8cc08a36fc147`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-kernel2ip2kernel`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-kernel2ip2kernel`)
	r.Run(`kubectl exec pods/alpine -n ns-kernel2ip2kernel -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl exec deployments/nse-kernel -n ns-kernel2ip2kernel -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestKernel2IP2Memif() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2IP2Memif")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-kernel2ip2memif`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Kernel2IP2Memif?ref=09d7c7a298f58007f2a76f3cfaa8cc08a36fc147`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-kernel2ip2memif`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-memif -n ns-kernel2ip2memif`)
	r.Run(`kubectl exec pods/alpine -n ns-kernel2ip2memif -- ping -c 4 172.16.1.100`)
	r.Run(`result=$(kubectl exec deployments/nse-memif -n "ns-kernel2ip2memif" -- vppctl ping 172.16.1.101 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}
func (s *Suite) TestMemif2IP2Kernel() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Memif2IP2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-memif2ip2kernel`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Memif2IP2Kernel?ref=09d7c7a298f58007f2a76f3cfaa8cc08a36fc147`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-memif -n ns-memif2ip2kernel`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-memif2ip2kernel`)
	r.Run(`result=$(kubectl exec deployments/nsc-memif -n "ns-memif2ip2kernel" -- vppctl ping 172.16.1.100 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`kubectl exec deployments/nse-kernel -n ns-memif2ip2kernel -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestMemif2IP2Memif() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Memif2IP2Memif")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-memif2ip2memif`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Memif2IP2Memif?ref=09d7c7a298f58007f2a76f3cfaa8cc08a36fc147`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-memif -n ns-memif2ip2memif`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-memif -n ns-memif2ip2memif`)
	r.Run(`result=$(kubectl exec deployments/nsc-memif -n "ns-memif2ip2memif" -- vppctl ping 172.16.1.100 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`result=$(kubectl exec deployments/nse-memif -n "ns-memif2ip2memif" -- vppctl ping 172.16.1.101 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}
func (s *Suite) TestVl3_basic() {
	r := s.Runner("../deployments-k8s/examples/use-cases/vl3-basic")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-vl3`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/vl3-basic?ref=09d7c7a298f58007f2a76f3cfaa8cc08a36fc147`)
	r.Run(`kubectl wait --for=condition=ready --timeout=2m pod -l app=alpine -n ns-vl3`)
	r.Run(`nscs=$(kubectl  get pods -l app=alpine -o go-template --template="{{range .items}}{{.metadata.name}} {{end}}" -n ns-vl3)` + "\n" + `[[ ! -z $nscs ]]`)
	r.Run(`(` + "\n" + `for nsc in $nscs ` + "\n" + `do` + "\n" + `    ipAddr=$(kubectl exec -n ns-vl3 $nsc -- ifconfig nsm-1) || exit` + "\n" + `    ipAddr=$(echo $ipAddr | grep -Eo 'inet addr:[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}'| cut -c 11-)` + "\n" + `    for pinger in $nscs` + "\n" + `    do` + "\n" + `        echo $pinger pings $ipAddr` + "\n" + `        kubectl exec $pinger -n ns-vl3 -- ping -c2 -i 0.5 $ipAddr || exit` + "\n" + `    done` + "\n" + `done` + "\n" + `)`)
	r.Run(`(` + "\n" + `for nsc in $nscs ` + "\n" + `do` + "\n" + `    echo $nsc pings nses` + "\n" + `    kubectl exec -n ns-vl3 $nsc -- ping 172.16.0.0 -c2 -i 0.5 || exit` + "\n" + `    kubectl exec -n ns-vl3 $nsc -- ping 172.16.1.0 -c2 -i 0.5 || exit` + "\n" + `done` + "\n" + `)`)
}
