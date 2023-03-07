// Code generated by gotestmd DO NOT EDIT.
package memory

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
	r := s.Runner("../deployments-k8s/examples/memory")
	s.T().Cleanup(func() {
		r.Run(`WH=$(kubectl get pods -l app=admission-webhook-k8s -n nsm-system --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `kubectl delete mutatingwebhookconfiguration ${WH}` + "\n" + `kubectl delete ns nsm-system`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/memory?ref=db1124735917c9347b3286bccee58a4a42401a93`)
	r.Run(`WH=$(kubectl get pods -l app=admission-webhook-k8s -n nsm-system --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `kubectl wait --for=condition=ready --timeout=1m pod ${WH} -n nsm-system`)
}
func (s *Suite) TestKernel2Kernel() {
	r := s.Runner("../deployments-k8s/examples/memory/Kernel2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-kernel2kernel`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/memory/Kernel2Kernel?ref=db1124735917c9347b3286bccee58a4a42401a93`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-kernel2kernel`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-kernel2kernel`)
	r.Run(`NSC=$(kubectl get pods -l app=alpine -n ns-kernel2kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-kernel -n ns-kernel2kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl exec ${NSC} -n ns-kernel2kernel -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl exec ${NSE} -n ns-kernel2kernel -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestKernel2Vxlan2Kernel() {
	r := s.Runner("../deployments-k8s/examples/memory/Kernel2Vxlan2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-kernel2vxlan2kernel`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/memory/Kernel2Vxlan2Kernel?ref=db1124735917c9347b3286bccee58a4a42401a93`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-kernel2vxlan2kernel`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-kernel2vxlan2kernel`)
	r.Run(`NSC=$(kubectl get pods -l app=alpine -n ns-kernel2vxlan2kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-kernel -n ns-kernel2vxlan2kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl exec ${NSC} -n ns-kernel2vxlan2kernel -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl exec ${NSE} -n ns-kernel2vxlan2kernel -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestMemif2Memif() {
	r := s.Runner("../deployments-k8s/examples/memory/Memif2Memif")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-memif2memif`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/memory/Memif2Memif?ref=db1124735917c9347b3286bccee58a4a42401a93`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-memif -n ns-memif2memif`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-memif -n ns-memif2memif`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc-memif -n ns-memif2memif --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-memif -n ns-memif2memif --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`result=$(kubectl exec "${NSC}" -n "ns-memif2memif" -- vppctl ping 172.16.1.100 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`result=$(kubectl exec "${NSE}" -n "ns-memif2memif" -- vppctl ping 172.16.1.101 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}
