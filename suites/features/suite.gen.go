// Code generated by gotestmd DO NOT EDIT.
package features

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/basic"
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
func (s *Suite) TestAnnotated_namespace() {
	r := s.Runner("../deployments-k8s/examples/features/annotated-namespace")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-annotated-namespace`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/annotated-namespace?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-annotated-namespace`)
	r.Run(`kubectl annotate ns ns-annotated-namespace networkservicemesh.io=kernel://annotated-namespace/nsm-1`)
	r.Run(`kubectl apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/a084f5035978c429cb30137b6400bdceb82bfe90/examples/features/annotated-namespace/client.yaml`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-annotated-namespace`)
	r.Run(`kubectl logs deployments/alpine -n ns-annotated-namespace -c cmd-nsc-init | grep -c '\[id:alpine-.*-0\]'`)
	r.Run(`kubectl exec deployments/alpine -n ns-annotated-namespace -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl exec deployments/nse-kernel -n ns-annotated-namespace -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestDns() {
	r := s.Runner("../deployments-k8s/examples/features/dns")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-dns`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/dns?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait --for=condition=ready --timeout=5m pod dnsutils -n ns-dns`)
	r.Run(`kubectl wait --for=condition=ready --timeout=5m pod -l app=nse-kernel -n ns-dns`)
	r.Run(`kubectl exec pods/dnsutils -c dnsutils -n ns-dns -- nslookup -norec -nodef my.coredns.service`)
	r.Run(`kubectl exec pods/dnsutils -c dnsutils -n ns-dns -- ping -c 4 my.coredns.service`)
	r.Run(`kubectl exec pods/dnsutils -c dnsutils -n ns-dns -- dig kubernetes.default A kubernetes.default AAAA | grep "kubernetes.default.svc.cluster.local"`)
}
func (s *Suite) TestKernel2IP2Kernel_dual_stack() {
	r := s.Runner("../deployments-k8s/examples/features/dual-stack/Kernel2IP2Kernel_dual_stack")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-kernel2ip2kernel-dual-stack`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/dual-stack/Kernel2IP2Kernel_dual_stack?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-kernel2ip2kernel-dual-stack`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-kernel2ip2kernel-dual-stack`)
	r.Run(`kubectl exec pods/alpine -n ns-kernel2ip2kernel-dual-stack -- ping -c 4 2001:db8::`)
	r.Run(`kubectl exec deployments/nse-kernel -n ns-kernel2ip2kernel-dual-stack -- ping -c 4 2001:db8::1`)
	r.Run(`kubectl exec pods/alpine -n ns-kernel2ip2kernel-dual-stack -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl exec deployments/nse-kernel -n ns-kernel2ip2kernel-dual-stack -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestKernel2Kernel_dual_stack() {
	r := s.Runner("../deployments-k8s/examples/features/dual-stack/Kernel2Kernel_dual_stack")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-kernel2kernel-dual-stack`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/dual-stack/Kernel2Kernel_dual_stack?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-kernel2kernel-dual-stack`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-kernel2kernel-dual-stack`)
	r.Run(`kubectl exec pods/alpine -n ns-kernel2kernel-dual-stack -- ping -c 4 2001:db8::`)
	r.Run(`kubectl exec deployments/nse-kernel -n ns-kernel2kernel-dual-stack -- ping -c 4 2001:db8::1`)
	r.Run(`kubectl exec pods/alpine -n ns-kernel2kernel-dual-stack -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl exec deployments/nse-kernel -n ns-kernel2kernel-dual-stack -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestExclude_prefixes() {
	r := s.Runner("../deployments-k8s/examples/features/exclude-prefixes")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete configmap excluded-prefixes-config` + "\n" + `kubectl delete ns ns-exclude-prefixes`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/exclude-prefixes/configmap?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/exclude-prefixes?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-exclude-prefixes`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-exclude-prefixes`)
	r.Run(`kubectl exec pods/alpine -n ns-exclude-prefixes -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl exec deployments/nse-kernel -n ns-exclude-prefixes -- ping -c 4 172.16.1.103`)
}
func (s *Suite) TestExclude_prefixes_client() {
	r := s.Runner("../deployments-k8s/examples/features/exclude-prefixes-client")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-exclude-prefixes-client`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/exclude-prefixes-client?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-exclude-prefixes-client`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel-1 -n ns-exclude-prefixes-client`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel-2 -n ns-exclude-prefixes-client`)
	r.Run(`kubectl exec pods/alpine -n ns-exclude-prefixes-client -- ping -c 4 172.16.1.96`)
	r.Run(`kubectl exec pods/alpine -n ns-exclude-prefixes-client -- ping -c 4 172.16.1.98`)
	r.Run(`kubectl exec deployments/nse-kernel-1 -n ns-exclude-prefixes-client -- ping -c 4 172.16.1.97`)
	r.Run(`kubectl exec deployments/nse-kernel-2 -n ns-exclude-prefixes-client -- ping -c 4 172.16.1.99`)
}
func (s *Suite) TestKernel2IP2Kernel_ipv6() {
	r := s.Runner("../deployments-k8s/examples/features/ipv6/Kernel2IP2Kernel_ipv6")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-kernel2ip2kernel-ipv6`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/ipv6/Kernel2IP2Kernel_ipv6?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-kernel2ip2kernel-ipv6`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-kernel2ip2kernel-ipv6`)
	r.Run(`kubectl exec pods/alpine -n ns-kernel2ip2kernel-ipv6 -- ping -c 4 2001:db8::`)
	r.Run(`kubectl exec deployments/nse-kernel -n ns-kernel2ip2kernel-ipv6 -- ping -c 4 2001:db8::1`)
}
func (s *Suite) TestKernel2IP2Memif_ipv6() {
	r := s.Runner("../deployments-k8s/examples/features/ipv6/Kernel2IP2Memif_ipv6")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-kernel2ip2memif-ipv6`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/ipv6/Kernel2IP2Memif_ipv6?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-kernel2ip2memif-ipv6`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-memif -n ns-kernel2ip2memif-ipv6`)
	r.Run(`kubectl exec pods/alpine -n ns-kernel2ip2memif-ipv6 -- ping -c 4 2001:db8::`)
	r.Run(`result=$(kubectl exec deployments/nse-memif -n "ns-kernel2ip2memif-ipv6" -- vppctl ping 2001:db8::1 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}
func (s *Suite) TestKernel2Kernel_ipv6() {
	r := s.Runner("../deployments-k8s/examples/features/ipv6/Kernel2Kernel_ipv6")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-kernel2kernel-ipv6`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/ipv6/Kernel2Kernel_ipv6?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-kernel2kernel-ipv6`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-kernel2kernel-ipv6`)
	r.Run(`kubectl exec pods/alpine -n ns-kernel2kernel-ipv6 -- ping -c 4 2001:db8::`)
	r.Run(`kubectl exec deployments/nse-kernel -n ns-kernel2kernel-ipv6 -- ping -c 4 2001:db8::1`)
}
func (s *Suite) TestMemif2IP2Kernel_ipv6() {
	r := s.Runner("../deployments-k8s/examples/features/ipv6/Memif2IP2Kernel_ipv6")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-memif2ip2kernel-ipv6`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/ipv6/Memif2IP2Kernel_ipv6?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-memif -n ns-memif2ip2kernel-ipv6`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-memif2ip2kernel-ipv6`)
	r.Run(`result=$(kubectl exec deployments/nsc-memif -n "ns-memif2ip2kernel-ipv6" -- vppctl ping 2001:db8:: repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`kubectl exec deployments/nse-kernel -n ns-memif2ip2kernel-ipv6 -- ping -c 4 2001:db8::1`)
}
func (s *Suite) TestMemif2IP2Memif_ipv6() {
	r := s.Runner("../deployments-k8s/examples/features/ipv6/Memif2IP2Memif_ipv6")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-memif2ip2memif-ipv6`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/ipv6/Memif2IP2Memif_ipv6?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-memif -n ns-memif2ip2memif-ipv6`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-memif -n ns-memif2ip2memif-ipv6`)
	r.Run(`result=$(kubectl exec deployments/nsc-memif -n "ns-memif2ip2memif-ipv6" -- vppctl ping 2001:db8:: repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`result=$(kubectl exec deployments/nse-memif -n "ns-memif2ip2memif-ipv6" -- vppctl ping 2001:db8::1 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}
func (s *Suite) TestMemif2Memif_ipv6() {
	r := s.Runner("../deployments-k8s/examples/features/ipv6/Memif2Memif_ipv6")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-memif2memif-ipv6`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/ipv6/Memif2Memif_ipv6?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-memif -n ns-memif2memif-ipv6`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-memif -n ns-memif2memif-ipv6`)
	r.Run(`result=$(kubectl exec deployments/nsc-memif -n "ns-memif2memif-ipv6" -- vppctl ping ipv6 2001:db8:: repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`result=$(kubectl exec deployments/nse-memif -n "ns-memif2memif-ipv6" -- vppctl ping ipv6 2001:db8::1 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}
func (s *Suite) TestMutually_aware_nses() {
	r := s.Runner("../deployments-k8s/examples/features/mutually-aware-nses")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-mutually-aware-nses`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/mutually-aware-nses?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-kernel -n ns-mutually-aware-nses`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel-1 -n ns-mutually-aware-nses`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel-2 -n ns-mutually-aware-nses`)
	r.Run(`kubectl exec deployments/nsc-kernel -n ns-mutually-aware-nses -- apk update` + "\n" + `kubectl exec deployments/nsc-kernel -n ns-mutually-aware-nses -- apk add iproute2`)
	r.Run(`result=$(kubectl exec deployments/nsc-kernel -n ns-mutually-aware-nses -- ip r get 172.16.1.100 from 172.16.1.101 ipproto tcp dport 6666)` + "\n" + `echo ${result}` + "\n" + `echo ${result} | grep -E -q "172.16.1.100 from 172.16.1.101 dev nsm-1"`)
	r.Run(`result=$(kubectl exec deployments/nsc-kernel -n ns-mutually-aware-nses -- ip r get 172.16.1.100 from 172.16.1.101 ipproto udp dport 5555)` + "\n" + `echo ${result}` + "\n" + `echo ${result} | grep -E -q "172.16.1.100 from 172.16.1.101 dev nsm-2"`)
}
func (s *Suite) TestNse_composition() {
	r := s.Runner("../deployments-k8s/examples/features/nse-composition")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-nse-composition`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/nse-composition?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait --for=condition=ready --timeout=5m pod -l app=alpine -n ns-nse-composition`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-nse-composition`)
	r.Run(`kubectl exec pods/alpine -n ns-nse-composition -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl exec pods/alpine -n ns-nse-composition -- wget -O /dev/null --timeout 5 "172.16.1.100:8080"`)
	r.Run(`kubectl exec pods/alpine -n ns-nse-composition -- wget -O /dev/null --timeout 5 "172.16.1.100:80"` + "\n" + `if [ 0 -eq $? ]; then` + "\n" + `  echo "error: port :80 is available" >&2` + "\n" + `  false` + "\n" + `else` + "\n" + `  echo "success: port :80 is unavailable"` + "\n" + `fi`)
	r.Run(`kubectl exec deployments/nse-kernel -n ns-nse-composition -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestOpa() {
	r := s.Runner("../deployments-k8s/examples/features/opa")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-opa`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/opa?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-kernel -n ns-opa`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-opa`)
	r.Run(`kubectl logs deployments/nsc-kernel -n ns-opa | grep "PermissionDenied desc = no sufficient privileges"`)
}
func (s *Suite) TestPolicy_based_routing() {
	r := s.Runner("../deployments-k8s/examples/features/policy-based-routing")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-policy-based-routing`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/policy-based-routing?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nettools -n ns-policy-based-routing`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-policy-based-routing`)
	r.Run(`kubectl exec pods/nettools -n ns-policy-based-routing -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl exec deployments/nse-kernel -n ns-policy-based-routing -- ping -c 4 172.16.1.101`)
	r.Run(`result=$(kubectl exec pods/nettools -n ns-policy-based-routing -- ip r get 172.16.3.1 from 172.16.2.201 ipproto tcp dport 6666)` + "\n" + `echo ${result}` + "\n" + `echo ${result} | grep -E -q "172.16.3.1 from 172.16.2.201 via 172.16.2.200 dev nsm-1 table 1"`)
	r.Run(`result=$(kubectl exec pods/nettools -n ns-policy-based-routing -- ip r get 172.16.3.1 from 172.16.2.201 ipproto tcp sport 5555)` + "\n" + `echo ${result}` + "\n" + `echo ${result} | grep -E -q "172.16.3.1 from 172.16.2.201 dev nsm-1 table 2"`)
	r.Run(`result=$(kubectl exec pods/nettools -n ns-policy-based-routing -- ip r get 172.16.4.1 ipproto udp dport 6666)` + "\n" + `echo ${result}` + "\n" + `echo ${result} | grep -E -q "172.16.4.1 dev nsm-1 table 3 src 172.16.1.101"`)
	r.Run(`result=$(kubectl exec pods/nettools -n ns-policy-based-routing -- ip r get 172.16.4.1 ipproto udp dport 6668)` + "\n" + `echo ${result}` + "\n" + `echo ${result} | grep -E -q "172.16.4.1 dev nsm-1 table 4 src 172.16.1.101"`)
	r.Run(`result=$(kubectl exec pods/nettools -n ns-policy-based-routing -- ip -6 route get 2004::5 from 2004::3 ipproto udp dport 5555)` + "\n" + `echo ${result}` + "\n" + `echo ${result} | grep -E -q "via 2004::6 dev nsm-1 table 5 src 2004::3"`)
}
func (s *Suite) TestScale_from_zero() {
	r := s.Runner("../deployments-k8s/examples/features/scale-from-zero")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-scale-from-zero`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/scale-from-zero?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait -n ns-scale-from-zero --for=condition=ready --timeout=1m pod -l app=nse-supplier-k8s`)
	r.Run(`kubectl wait -n ns-scale-from-zero --for=condition=ready --timeout=1m pod -l app=alpine`)
	r.Run(`kubectl wait -n ns-scale-from-zero --for=condition=ready --timeout=1m pod -l app=nse-icmp-responder`)
	r.Run(`NSE=$(kubectl get pod -n ns-scale-from-zero --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' -l app=nse-icmp-responder)`)
	r.Run(`kubectl exec pods/alpine -n ns-scale-from-zero -- ping -c 4 169.254.0.0`)
	r.Run(`kubectl exec $NSE -n ns-scale-from-zero -- ping -c 4 169.254.0.1`)
	r.Run(`NSE_NODE=$(kubectl get pod -n ns-scale-from-zero --template '{{range .items}}{{.spec.nodeName}}{{"\n"}}{{end}}' -l app=nse-icmp-responder)` + "\n" + `NSC_NODE=$(kubectl get pod -n ns-scale-from-zero --template '{{range .items}}{{.spec.nodeName}}{{"\n"}}{{end}}' -l app=alpine)`)
	r.Run(`if [ $NSC_NODE == $NSE_NODE ]; then echo "OK"; else echo "different nodes"; false; fi`)
	r.Run(`kubectl delete pod -n ns-scale-from-zero alpine`)
	r.Run(`kubectl wait -n ns-scale-from-zero --for=delete --timeout=1m pod -l app=nse-icmp-responder`)
}
func (s *Suite) TestSelect_forwarder() {
	r := s.Runner("../deployments-k8s/examples/features/select-forwarder")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-select-forwarder`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/select-forwarder?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-select-forwarder`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-select-forwarder`)
	r.Run(`kubectl exec pods/alpine -n ns-select-forwarder -- ping -c 4 169.254.0.0`)
	r.Run(`kubectl exec deployments/nse-kernel -n ns-select-forwarder -- ping -c 4 169.254.0.1`)
	r.Run(`kubectl logs pods/alpine -c cmd-nsc -n ns-select-forwarder | grep "my-forwarder-vpp"`)
}
func (s *Suite) TestVl3_basic() {
	r := s.Runner("../deployments-k8s/examples/features/vl3-basic")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-vl3`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/vl3-basic?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait --for=condition=ready --timeout=2m pod -l app=alpine -n ns-vl3`)
	r.Run(`nscs=$(kubectl  get pods -l app=alpine -o go-template --template="{{range .items}}{{.metadata.name}} {{end}}" -n ns-vl3)` + "\n" + `[[ ! -z $nscs ]]`)
	r.Run(`(` + "\n" + `for nsc in $nscs ` + "\n" + `do` + "\n" + `    ipAddr=$(kubectl exec -n ns-vl3 $nsc -- ifconfig nsm-1) || exit` + "\n" + `    ipAddr=$(echo $ipAddr | grep -Eo 'inet addr:[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}'| cut -c 11-)` + "\n" + `    for pinger in $nscs` + "\n" + `    do` + "\n" + `        echo $pinger pings $ipAddr` + "\n" + `        kubectl exec $pinger -n ns-vl3 -- ping -c2 -i 0.5 $ipAddr || exit` + "\n" + `    done` + "\n" + `done` + "\n" + `)`)
	r.Run(`(` + "\n" + `for nsc in $nscs ` + "\n" + `do` + "\n" + `    echo $nsc pings nses` + "\n" + `    kubectl exec -n ns-vl3 $nsc -- ping 172.16.0.0 -c2 -i 0.5 || exit` + "\n" + `    kubectl exec -n ns-vl3 $nsc -- ping 172.16.1.0 -c2 -i 0.5 || exit` + "\n" + `done` + "\n" + `)`)
}
func (s *Suite) TestVl3_dns() {
	r := s.Runner("../deployments-k8s/examples/features/vl3-dns")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-vl3-dns`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/vl3-dns?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait --for=condition=ready --timeout=2m pod -l app=alpine -n ns-vl3-dns`)
	r.Run(`nscs=$(kubectl  get pods -l app=alpine -o go-template --template="{{range .items}}{{.metadata.name}} {{end}}" -n ns-vl3-dns)` + "\n" + `[[ ! -z $nscs ]]`)
	r.Run(`(` + "\n" + `for nsc in $nscs` + "\n" + `do` + "\n" + `    for pinger in $nscs` + "\n" + `    do` + "\n" + `        kubectl exec $pinger -n ns-vl3-dns -- ping -c2 -i 0.5 $nsc.vl3-dns -4 || exit` + "\n" + `    done` + "\n" + `done` + "\n" + `)`)
	r.Run(`(` + "\n" + `for nsc in $nscs` + "\n" + `do` + "\n" + `    for pinger in $nscs` + "\n" + `    do` + "\n" + `        # Get IP address for PTR request` + "\n" + `        nscAddr=$(kubectl exec $pinger -n ns-vl3-dns -- nslookup -type=a $nsc.vl3-dns | grep -A1 Name | tail -n1 | sed 's/Address: //')` + "\n" + `        kubectl exec $pinger -n ns-vl3-dns -- nslookup $nscAddr || exit` + "\n" + `    done` + "\n" + `done` + "\n" + `)`)
	r.Run(`nses=$(kubectl get pods -l app=nse-vl3-vpp -o go-template --template="{{range .items}}{{.metadata.name}} {{end}}" -n ns-vl3-dns)` + "\n" + `[[ ! -z nses ]]`)
	r.Run(`(` + "\n" + `for nse in $nses` + "\n" + `do` + "\n" + `    for pinger in $nscs` + "\n" + `    do` + "\n" + `        kubectl exec $pinger -n ns-vl3-dns -- ping -c2 -i 0.5 $nse.vl3-dns -4 || exit` + "\n" + `    done` + "\n" + `done` + "\n" + `)`)
	r.Run(`(` + "\n" + `for nse in $nses` + "\n" + `do` + "\n" + `    for pinger in $nscs` + "\n" + `    do` + "\n" + `        # Get IP address for PTR request` + "\n" + `        nseAddr=$(kubectl exec $pinger -n ns-vl3-dns -- nslookup -type=a $nse.vl3-dns | grep -A1 Name | tail -n1 | sed 's/Address: //')` + "\n" + `        kubectl exec $pinger -n ns-vl3-dns -- nslookup $nseAddr || exit` + "\n" + `    done` + "\n" + `done` + "\n" + `)`)
}
func (s *Suite) TestVl3_scale_from_zero() {
	r := s.Runner("../deployments-k8s/examples/features/vl3-scale-from-zero")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-vl3-scale-from-zero`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/vl3-scale-from-zero?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait -n ns-vl3-scale-from-zero --for=condition=ready --timeout=1m pod -l app=nse-supplier-k8s`)
	r.Run(`kubectl wait -n ns-vl3-scale-from-zero --for=condition=ready --timeout=1m pod -l app=alpine`)
	r.Run(`kubectl wait -n ns-vl3-scale-from-zero --for=condition=ready --timeout=1m pod -l app=nse-vl3-vpp`)
	r.Run(`nscs=$(kubectl  get pods -l app=alpine -o go-template --template="{{range .items}}{{.metadata.name}} {{end}}" -n ns-vl3-scale-from-zero)` + "\n" + `[[ ! -z $nscs ]]`)
	r.Run(`(` + "\n" + `for nsc in $nscs ` + "\n" + `do` + "\n" + `    ipAddr=$(kubectl exec -n ns-vl3-scale-from-zero $nsc -- ifconfig nsm-1) || exit` + "\n" + `    ipAddr=$(echo $ipAddr | grep -Eo 'inet addr:[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}'| cut -c 11-)` + "\n" + `    for pinger in $nscs` + "\n" + `    do` + "\n" + `        echo $pinger pings $ipAddr` + "\n" + `        kubectl exec $pinger -n ns-vl3-scale-from-zero -- ping -c2 -i 0.5 $ipAddr || exit` + "\n" + `    done` + "\n" + `done` + "\n" + `)`)
	r.Run(`(` + "\n" + `for nsc in $nscs ` + "\n" + `do` + "\n" + `    echo $nsc pings nses` + "\n" + `    kubectl exec -n ns-vl3-scale-from-zero $nsc -- ping 172.16.0.0 -c2 -i 0.5 || exit` + "\n" + `    kubectl exec -n ns-vl3-scale-from-zero $nsc -- ping 172.16.1.0 -c2 -i 0.5 || exit` + "\n" + `done` + "\n" + `)`)
}
func (s *Suite) TestWebhook() {
	r := s.Runner("../deployments-k8s/examples/features/webhook")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-webhook`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/features/webhook?ref=a084f5035978c429cb30137b6400bdceb82bfe90`)
	r.Run(`kubectl wait --for=condition=ready --timeout=5m pod -l app=nse-kernel -n ns-webhook`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nettools -n ns-webhook`)
	r.Run(`kubectl exec pods/nettools -n ns-webhook -- curl 172.16.1.100:80 | grep -o "<title>Welcome to nginx!</title>"`)
}
