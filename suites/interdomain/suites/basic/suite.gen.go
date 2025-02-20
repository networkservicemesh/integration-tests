// Code generated by gotestmd DO NOT EDIT.
package basic

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/interdomain/three_cluster_configuration/basic"
)

type Suite struct {
	base.Suite
	basicSuite basic.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite, &s.basicSuite}
	for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
}
func (s *Suite) TestFloating_Kernel2Ethernet2Kernel() {
	r := s.Runner("../deployments-k8s/examples/interdomain/usecases/floating_Kernel2Ethernet2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete ns ns-floating-kernel2ethernet2kernel`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete ns ns-floating-kernel2ethernet2kernel`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG3 delete ns ns-floating-kernel2ethernet2kernel`)
	})
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Kernel2Ethernet2Kernel/cluster3?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Kernel2Ethernet2Kernel/cluster2?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-floating-kernel2ethernet2kernel`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Kernel2Ethernet2Kernel/cluster1?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --for=condition=ready --timeout=5m pod -l app=alpine -n ns-floating-kernel2ethernet2kernel`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pods/alpine -n ns-floating-kernel2ethernet2kernel -- ping -c 4 172.16.1.2`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 exec deployments/nse-kernel -n ns-floating-kernel2ethernet2kernel -- ping -c 4 172.16.1.3`)
}
func (s *Suite) TestFloating_Kernel2Ethernet2Memif() {
	r := s.Runner("../deployments-k8s/examples/interdomain/usecases/floating_Kernel2Ethernet2Memif")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete ns ns-floating-kernel2ethernet2memif`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete ns ns-floating-kernel2ethernet2memif`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG3 delete ns ns-floating-kernel2ethernet2memif`)
	})
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Kernel2Ethernet2Memif/cluster3?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Kernel2Ethernet2Memif/cluster2?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 wait --for=condition=ready --timeout=2m pod -l app=nse-memif -n ns-floating-kernel2ethernet2memif`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Kernel2Ethernet2Memif/cluster1?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --for=condition=ready --timeout=5m pod -l app=alpine -n ns-floating-kernel2ethernet2memif`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pods/alpine -n ns-floating-kernel2ethernet2memif -- ping -c 4 172.16.1.2`)
	r.Run(`result=$(kubectl --kubeconfig=$KUBECONFIG2 exec deployments/nse-memif -n "ns-floating-kernel2ethernet2memif" -- vppctl ping 172.16.1.3 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}
func (s *Suite) TestFloating_Kernel2IP2Kernel() {
	r := s.Runner("../deployments-k8s/examples/interdomain/usecases/floating_Kernel2IP2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete ns ns-floating-kernel2ip2kernel`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete ns ns-floating-kernel2ip2kernel`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG3 delete ns ns-floating-kernel2ip2kernel`)
	})
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Kernel2IP2Kernel/cluster3?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Kernel2IP2Kernel/cluster2?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-floating-kernel2ip2kernel`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Kernel2IP2Kernel/cluster1?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --for=condition=ready --timeout=5m pod -l app=alpine -n ns-floating-kernel2ip2kernel`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pods/alpine -n ns-floating-kernel2ip2kernel -- ping -c 4 172.16.1.2`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 exec deployments/nse-kernel -n ns-floating-kernel2ip2kernel -- ping -c 4 172.16.1.3`)
}
func (s *Suite) TestFloating_Kernel2IP2Memif() {
	r := s.Runner("../deployments-k8s/examples/interdomain/usecases/floating_Kernel2IP2Memif")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete ns ns-floating-kernel2ip2memif`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete ns ns-floating-kernel2ip2memif`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG3 delete ns ns-floating-kernel2ip2memif`)
	})
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Kernel2IP2Memif/cluster3?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Kernel2IP2Memif/cluster2?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 wait --for=condition=ready --timeout=2m pod -l app=nse-memif -n ns-floating-kernel2ip2memif`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Kernel2IP2Memif/cluster1?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --for=condition=ready --timeout=5m pod -l app=alpine -n ns-floating-kernel2ip2memif`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pods/alpine -n ns-floating-kernel2ip2memif -- ping -c 4 172.16.1.2`)
	r.Run(`result=$(kubectl --kubeconfig=$KUBECONFIG2 exec deployments/nse-memif -n "ns-floating-kernel2ip2memif" -- vppctl ping 172.16.1.3 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}
func (s *Suite) TestFloating_Memif2Ethernet2Kernel() {
	r := s.Runner("../deployments-k8s/examples/interdomain/usecases/floating_Memif2Ethernet2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete ns ns-floating-memif2ethernet2kernel`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete ns ns-floating-memif2ethernet2kernel`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG3 delete ns ns-floating-memif2ethernet2kernel`)
	})
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Memif2Ethernet2Kernel/cluster3?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Memif2Ethernet2Kernel/cluster2?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-floating-memif2ethernet2kernel`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Memif2Ethernet2Kernel/cluster1?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --for=condition=ready --timeout=2m pod -l app=nsc-memif -n ns-floating-memif2ethernet2kernel`)
	r.Run(`result=$(kubectl --kubeconfig=$KUBECONFIG1 exec deployments/nsc-memif -n "ns-floating-memif2ethernet2kernel" -- vppctl ping 172.16.1.2 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 exec deployments/nse-kernel -n ns-floating-memif2ethernet2kernel -- ping -c 4 172.16.1.3`)
}
func (s *Suite) TestFloating_Memif2Ethernet2Memif() {
	r := s.Runner("../deployments-k8s/examples/interdomain/usecases/floating_Memif2Ethernet2Memif")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete ns ns-floating-memif2ethernet2memif`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete ns ns-floating-memif2ethernet2memif`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG3 delete ns ns-floating-memif2ethernet2memif`)
	})
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Memif2Ethernet2Memif/cluster3?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Memif2Ethernet2Memif/cluster2?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 wait --for=condition=ready --timeout=2m pod -l app=nse-memif -n ns-floating-memif2ethernet2memif`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Memif2Ethernet2Memif/cluster1?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --for=condition=ready --timeout=2m pod -l app=nsc-memif -n ns-floating-memif2ethernet2memif`)
	r.Run(`result=$(kubectl --kubeconfig=$KUBECONFIG1 exec deployments/nsc-memif -n "ns-floating-memif2ethernet2memif" -- vppctl ping 172.16.1.2 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`result=$(kubectl --kubeconfig=$KUBECONFIG2 exec deployments/nse-memif -n "ns-floating-memif2ethernet2memif" -- vppctl ping 172.16.1.3 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}
func (s *Suite) TestFloating_Memif2IP2Kernel() {
	r := s.Runner("../deployments-k8s/examples/interdomain/usecases/floating_Memif2IP2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete ns ns-floating-memif2ip2kernel`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete ns ns-floating-memif2ip2kernel`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG3 delete ns ns-floating-memif2ip2kernel`)
	})
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Memif2IP2Kernel/cluster3?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Memif2IP2Kernel/cluster2?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-floating-memif2ip2kernel`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Memif2IP2Kernel/cluster1?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --for=condition=ready --timeout=2m pod -l app=nsc-memif -n ns-floating-memif2ip2kernel`)
	r.Run(`result=$(kubectl --kubeconfig=$KUBECONFIG1 exec deployments/nsc-memif -n "ns-floating-memif2ip2kernel" -- vppctl ping 172.16.1.2 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 exec deployments/nse-kernel -n ns-floating-memif2ip2kernel -- ping -c 4 172.16.1.3`)
}
func (s *Suite) TestFloating_Memif2IP2Memif() {
	r := s.Runner("../deployments-k8s/examples/interdomain/usecases/floating_Memif2IP2Memif")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete ns ns-floating-memif2ip2memif`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete ns ns-floating-memif2ip2memif`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG3 delete ns ns-floating-memif2ip2memif`)
	})
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Memif2IP2Memif/cluster3?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Memif2IP2Memif/cluster2?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 wait --for=condition=ready --timeout=2m pod -l app=nse-memif -n ns-floating-memif2ip2memif`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_Memif2IP2Memif/cluster1?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --for=condition=ready --timeout=2m pod -l app=nsc-memif -n ns-floating-memif2ip2memif`)
	r.Run(`result=$(kubectl --kubeconfig=$KUBECONFIG1 exec deployments/nsc-memif -n "ns-floating-memif2ip2memif" -- vppctl ping 172.16.1.2 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`result=$(kubectl --kubeconfig=$KUBECONFIG2 exec deployments/nse-memif -n "ns-floating-memif2ip2memif" -- vppctl ping 172.16.1.3 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}
func (s *Suite) TestFloating_dns() {
	r := s.Runner("../deployments-k8s/examples/interdomain/usecases/floating_dns")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete ns ns-floating-dns`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete ns ns-floating-dns`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG3 delete ns ns-floating-dns`)
	})
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_dns/cluster3?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_dns/cluster2?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_dns/cluster1?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --for=condition=ready --timeout=5m pod -l app=dnsutils -n ns-floating-dns`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pods/dnsutils -c dnsutils -n ns-floating-dns -- nslookup -norec -nodef my.coredns.service`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pods/dnsutils -c dnsutils -n ns-floating-dns -- ping -c 4 my.coredns.service`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pods/dnsutils -c dnsutils -n ns-floating-dns -- dig kubernetes.default A kubernetes.default AAAA | grep "kubernetes.default.svc.cluster.local"`)
}
func (s *Suite) TestFloating_nse_composition() {
	r := s.Runner("../deployments-k8s/examples/interdomain/usecases/floating_nse_composition")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete ns ns-interdomain-nse-composition`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete ns ns-interdomain-nse-composition`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG3 delete ns ns-interdomain-nse-composition`)
	})
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_nse_composition/cluster3?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_nse_composition/cluster2?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-interdomain-nse-composition`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_nse_composition/cluster1?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --for=condition=ready --timeout=5m pod -l app=alpine -n ns-interdomain-nse-composition`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pods/alpine -n ns-interdomain-nse-composition -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pods/alpine -n ns-interdomain-nse-composition -- wget -O /dev/null --timeout 5 "172.16.1.100:8080"`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pods/alpine -n ns-interdomain-nse-composition -- wget -O /dev/null --timeout 5 "172.16.1.100:80"` + "\n" + `if [ 0 -eq $? ]; then` + "\n" + `  echo "error: port :80 is available" >&2` + "\n" + `  false` + "\n" + `else` + "\n" + `  echo "success: port :80 is unavailable"` + "\n" + `fi`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 exec deployments/nse-kernel -n ns-interdomain-nse-composition -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestFloating_vl3_basic() {
	r := s.Runner("../deployments-k8s/examples/interdomain/usecases/floating_vl3-basic")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG3 delete -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_vl3-basic/cluster3?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_vl3-basic/cluster2?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_vl3-basic/cluster1?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	})
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_vl3-basic/cluster3?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_vl3-basic/cluster1?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_vl3-basic/cluster2?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-floating-vl3-basic`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-floating-vl3-basic`)
	r.Run(`ipAddr2=$(kubectl --kubeconfig=$KUBECONFIG2 exec -n ns-floating-vl3-basic pods/alpine -- ifconfig nsm-1)` + "\n" + `ipAddr2=$(echo $ipAddr2 | grep -Eo 'inet addr:[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}'| cut -c 11-)` + "\n" + `kubectl --kubeconfig=$KUBECONFIG1 exec pods/alpine -n ns-floating-vl3-basic -- ping -c 4 $ipAddr2`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pods/alpine -n ns-floating-vl3-basic -- ping -c 4 172.16.0.0` + "\n" + `kubectl --kubeconfig=$KUBECONFIG1 exec pods/alpine -n ns-floating-vl3-basic -- ping -c 4 172.16.1.0`)
	r.Run(`ipAddr1=$(kubectl --kubeconfig=$KUBECONFIG1 exec -n ns-floating-vl3-basic pods/alpine -- ifconfig nsm-1)` + "\n" + `ipAddr1=$(echo $ipAddr1 | grep -Eo 'inet addr:[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}'| cut -c 11-)` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 exec pods/alpine -n ns-floating-vl3-basic -- ping -c 4 $ipAddr1`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 exec pods/alpine -n ns-floating-vl3-basic -- ping -c 4 172.16.0.0` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 exec pods/alpine -n ns-floating-vl3-basic -- ping -c 4 172.16.1.0`)
}
func (s *Suite) TestFloating_vl3_dns() {
	r := s.Runner("../deployments-k8s/examples/interdomain/usecases/floating_vl3-dns")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG3 delete -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_vl3-dns/cluster3?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_vl3-dns/cluster2?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_vl3-dns/cluster1?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	})
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_vl3-dns/cluster3?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_vl3-dns/cluster1?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --for=condition=ready --timeout=5m pod -l app=alpine -n ns-floating-vl3-dns`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_vl3-dns/cluster2?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 wait --for=condition=ready --timeout=5m pod -l app=alpine -n ns-floating-vl3-dns`)
	r.Run(`nsc1=$(kubectl --kubeconfig=$KUBECONFIG1 get pods -l app=alpine -n ns-floating-vl3-dns --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`nse1=$(kubectl --kubeconfig=$KUBECONFIG1 get pods -l app=nse-vl3-vpp -n ns-floating-vl3-dns --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`nsc2=$(kubectl --kubeconfig=$KUBECONFIG2 get pods -l app=alpine -n ns-floating-vl3-dns --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`nse2=$(kubectl --kubeconfig=$KUBECONFIG2 get pods -l app=nse-vl3-vpp -n ns-floating-vl3-dns --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pods/alpine-1 -n ns-floating-vl3-dns -- ping -c2 -i 0.5 $nsc2.floating-vl3-dns.my.cluster3. -4`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pods/alpine-1 -n ns-floating-vl3-dns -- ping -c2 -i 0.5 $nse2.floating-vl3-dns.my.cluster3. -4`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pods/alpine-1 -n ns-floating-vl3-dns -- ping -c2 -i 0.5 $nse1.floating-vl3-dns.my.cluster3. -4`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 exec pods/alpine-2 -n ns-floating-vl3-dns -- ping -c2 -i 0.5 $nsc1.floating-vl3-dns.my.cluster3. -4`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 exec pods/alpine-2 -n ns-floating-vl3-dns -- ping -c2 -i 0.5 $nse1.floating-vl3-dns.my.cluster3. -4`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 exec pods/alpine-2 -n ns-floating-vl3-dns -- ping -c2 -i 0.5 $nse2.floating-vl3-dns.my.cluster3. -4`)
}
func (s *Suite) TestFloating_vl3_scale_from_zero() {
	r := s.Runner("../deployments-k8s/examples/interdomain/usecases/floating_vl3-scale-from-zero")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG3 delete -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_vl3-scale-from-zero/cluster3?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_vl3-scale-from-zero/cluster2?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_vl3-scale-from-zero/cluster1?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	})
	r.Run(`kubectl --kubeconfig=$KUBECONFIG3 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_vl3-scale-from-zero/cluster3?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_vl3-scale-from-zero/cluster1?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/floating_vl3-scale-from-zero/cluster2?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-floating-vl3-scale-from-zero`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-floating-vl3-scale-from-zero`)
	r.Run(`ipAddr2=$(kubectl --kubeconfig=$KUBECONFIG2 exec -n ns-floating-vl3-scale-from-zero pods/alpine -- ifconfig nsm-1)` + "\n" + `ipAddr2=$(echo $ipAddr2 | grep -Eo 'inet addr:[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}'| cut -c 11-)` + "\n" + `kubectl --kubeconfig=$KUBECONFIG1 exec pods/alpine -n ns-floating-vl3-scale-from-zero -- ping -c 4 $ipAddr2`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pods/alpine -n ns-floating-vl3-scale-from-zero -- ping -c 4 172.16.0.0` + "\n" + `kubectl --kubeconfig=$KUBECONFIG1 exec pods/alpine -n ns-floating-vl3-scale-from-zero -- ping -c 4 172.16.1.0`)
	r.Run(`ipAddr1=$(kubectl --kubeconfig=$KUBECONFIG1 exec -n ns-floating-vl3-scale-from-zero pods/alpine -- ifconfig nsm-1)` + "\n" + `ipAddr1=$(echo $ipAddr1 | grep -Eo 'inet addr:[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}'| cut -c 11-)` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 exec pods/alpine -n ns-floating-vl3-scale-from-zero -- ping -c 4 $ipAddr1`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 exec pods/alpine -n ns-floating-vl3-scale-from-zero -- ping -c 4 172.16.0.0` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 exec pods/alpine -n ns-floating-vl3-scale-from-zero -- ping -c 4 172.16.1.0`)
}
func (s *Suite) TestInterdomain_Kernel2Ethernet2Kernel() {
	r := s.Runner("../deployments-k8s/examples/interdomain/usecases/interdomain_Kernel2Ethernet2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete ns ns-interdomain-kernel2ethernet2kernel`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete ns ns-interdomain-kernel2ethernet2kernel`)
	})
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/interdomain_Kernel2Ethernet2Kernel/cluster2?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-interdomain-kernel2ethernet2kernel`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/interdomain_Kernel2Ethernet2Kernel/cluster1?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --for=condition=ready --timeout=5m pod -l app=alpine -n ns-interdomain-kernel2ethernet2kernel`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pods/alpine -n ns-interdomain-kernel2ethernet2kernel -- ping -c 4 172.16.1.2`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 exec deployments/nse-kernel -n ns-interdomain-kernel2ethernet2kernel -- ping -c 4 172.16.1.3`)
}
func (s *Suite) TestInterdomain_Kernel2IP2Kernel() {
	r := s.Runner("../deployments-k8s/examples/interdomain/usecases/interdomain_Kernel2IP2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete ns ns-interdomain-kernel2ip2kernel`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete ns ns-interdomain-kernel2ip2kernel`)
	})
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/interdomain_Kernel2IP2Kernel/cluster2?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-interdomain-kernel2ip2kernel`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/interdomain_Kernel2IP2Kernel/cluster1?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --for=condition=ready --timeout=5m pod -l app=alpine -n ns-interdomain-kernel2ip2kernel`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pods/alpine -n ns-interdomain-kernel2ip2kernel -- ping -c 4 172.16.1.2`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 exec deployments/nse-kernel -n ns-interdomain-kernel2ip2kernel -- ping -c 4 172.16.1.3`)
}
func (s *Suite) TestInterdomain_dns() {
	r := s.Runner("../deployments-k8s/examples/interdomain/usecases/interdomain_dns")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete ns ns-interdomain-dns`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete ns ns-interdomain-dns`)
	})
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/interdomain_dns/cluster2?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/usecases/interdomain_dns/cluster1?ref=c6d9d31f167b61e02474747cfc92ac0ef9d3dd90`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --for=condition=ready --timeout=5m pod -l app=dnsutils -n ns-interdomain-dns`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pods/dnsutils -c dnsutils -n ns-interdomain-dns -- nslookup -norec -nodef my.coredns.service`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pods/dnsutils -c dnsutils -n ns-interdomain-dns -- ping -c 4 my.coredns.service`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pods/dnsutils -c dnsutils -n ns-interdomain-dns -- dig kubernetes.default A kubernetes.default AAAA | grep "kubernetes.default.svc.cluster.local"`)
}
