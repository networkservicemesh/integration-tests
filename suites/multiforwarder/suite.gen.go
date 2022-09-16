// Code generated by gotestmd DO NOT EDIT.
package multiforwarder

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/spire"
)

type Suite struct {
	base.Suite
	spireSuite spire.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite, &s.spireSuite}
	for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
	r := s.Runner("../deployments-k8s/examples/multiforwarder")
	s.T().Cleanup(func() {
		r.Run(`WH=$(kubectl get pods -l app=admission-webhook-k8s -n nsm-system --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `kubectl delete mutatingwebhookconfiguration ${WH}` + "\n" + `kubectl delete ns nsm-system`)
	})
	r.Run(`kubectl create ns nsm-system`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/multiforwarder?ref=d9c9ce7b315b179188c887b89ad1843af03a2dfd`)
	r.Run(`WH=$(kubectl get pods -l app=admission-webhook-k8s -n nsm-system --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `kubectl wait --for=condition=ready --timeout=1m pod ${WH} -n nsm-system`)
}
func (s *Suite) TestKernel2Kernel() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-kernel2kernel`)
	})
	r.Run(`kubectl create ns ns-kernel2kernel`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Kernel2Kernel?ref=d9c9ce7b315b179188c887b89ad1843af03a2dfd`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-kernel2kernel`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-kernel2kernel`)
	r.Run(`NSC=$(kubectl get pods -l app=alpine -n ns-kernel2kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-kernel -n ns-kernel2kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl exec ${NSC} -n ns-kernel2kernel -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl exec ${NSE} -n ns-kernel2kernel -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestKernel2Kernel_Vfio2Noop() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2Kernel&Vfio2Noop")
	s.T().Cleanup(func() {
		r.Run(`NSE_VFIO=$(kubectl get pods -l app=nse-vfio -n ns-kernel2kernel-vfio2noop --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
		r.Run(`kubectl -n ns-kernel2kernel-vfio2noop exec ${NSE_VFIO} --container ponger -- /bin/bash -c '\` + "\n" + `  sleep 10 && kill $(pgrep "pingpong") 1>/dev/null 2>&1 &                    \` + "\n" + `'`)
		r.Run(`kubectl delete ns ns-kernel2kernel-vfio2noop`)
	})
	r.Run(`kubectl create ns ns-kernel2kernel-vfio2noop`)
	r.Run(`function mac_create(){` + "\n" + `    echo -n 00` + "\n" + `    dd bs=1 count=5 if=/dev/random 2>/dev/null | hexdump -v -e '/1 ":%02x"'` + "\n" + `}`)
	r.Run(`CLIENT_MAC=$(mac_create)` + "\n" + `echo Client MAC: ${CLIENT_MAC}`)
	r.Run(`SERVER_MAC=$(mac_create)` + "\n" + `echo Server MAC: ${SERVER_MAC}`)
	r.Run(`cat > patch-nse-vfio.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse-vfio` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: sidecar` + "\n" + `          env:` + "\n" + `            - name: NSM_SERVICES` + "\n" + `              value: "pingpong@worker.domain: { addr: ${SERVER_MAC} }"` + "\n" + `        - name: ponger` + "\n" + `          command: ["/bin/bash", "/root/scripts/pong.sh", "ens6f3", "31", ${SERVER_MAC}]` + "\n" + `EOF`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Kernel2Kernel&Vfio2Noop?ref=d9c9ce7b315b179188c887b89ad1843af03a2dfd`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-kernel -n ns-kernel2kernel-vfio2noop`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-kernel2kernel-vfio2noop`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-vfio -n ns-kernel2kernel-vfio2noop`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-vfio -n ns-kernel2kernel-vfio2noop`)
	r.Run(`NSC_KERNEL=$(kubectl get pods -l app=nsc-kernel -n ns-kernel2kernel-vfio2noop --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE_KERNEL=$(kubectl get pods -l app=nse-kernel -n ns-kernel2kernel-vfio2noop --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSC_VFIO=$(kubectl get pods -l app=nsc-vfio -n ns-kernel2kernel-vfio2noop --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`function dpdk_ping() {` + "\n" + `  err_file="$(mktemp)"` + "\n" + `  trap 'rm -f "${err_file}"' RETURN` + "\n" + `` + "\n" + `  client_mac="$1"` + "\n" + `  server_mac="$2"` + "\n" + `` + "\n" + `  command="/root/dpdk-pingpong/build/app/pingpong \` + "\n" + `      --no-huge                                   \` + "\n" + `      --                                          \` + "\n" + `      -n 500                                      \` + "\n" + `      -c                                          \` + "\n" + `      -C ${client_mac}                            \` + "\n" + `      -S ${server_mac}` + "\n" + `      "` + "\n" + `  out="$(kubectl -n ns-kernel2kernel-vfio2noop exec ${NSC_VFIO} --container pinger -- /bin/bash -c "${command}" 2>"${err_file}")"` + "\n" + `` + "\n" + `  if [[ "$?" != 0 ]]; then` + "\n" + `    echo "${out}"` + "\n" + `    cat "${err_file}" 1>&2` + "\n" + `    return 1` + "\n" + `  fi` + "\n" + `` + "\n" + `  if ! pong_packets="$(echo "${out}" | grep "rx .* pong packets" | sed -E 's/rx ([0-9]*) pong packets/\1/g')"; then` + "\n" + `    echo "${out}"` + "\n" + `    cat "${err_file}" 1>&2` + "\n" + `    return 1` + "\n" + `  fi` + "\n" + `` + "\n" + `  if [[ "${pong_packets}" == 0 ]]; then` + "\n" + `    echo "${out}"` + "\n" + `    cat "${err_file}" 1>&2` + "\n" + `    return 1` + "\n" + `  fi` + "\n" + `` + "\n" + `  echo "${out}"` + "\n" + `  return 0` + "\n" + `}`)
	r.Run(`kubectl exec ${NSC_KERNEL} -n ns-kernel2kernel-vfio2noop -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl exec ${NSE_KERNEL} -n ns-kernel2kernel-vfio2noop -- ping -c 4 172.16.1.101`)
	r.Run(`dpdk_ping ${CLIENT_MAC} ${SERVER_MAC}`)
}
func (s *Suite) TestKernel2Vxlan2Kernel() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2Vxlan2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-kernel2vxlan2kernel`)
	})
	r.Run(`kubectl create ns ns-kernel2vxlan2kernel`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Kernel2Vxlan2Kernel?ref=d9c9ce7b315b179188c887b89ad1843af03a2dfd`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-kernel2vxlan2kernel`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-kernel2vxlan2kernel`)
	r.Run(`NSC=$(kubectl get pods -l app=alpine -n ns-kernel2vxlan2kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-kernel -n ns-kernel2vxlan2kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl exec ${NSC} -n ns-kernel2vxlan2kernel -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl exec ${NSE} -n ns-kernel2vxlan2kernel -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestKernel2Vxlan2Kernel_Vfio2Noop() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2Vxlan2Kernel&Vfio2Noop")
	s.T().Cleanup(func() {
		r.Run(`NSE_VFIO=$(kubectl get pods -l app=nse-vfio -n ns-kernel2vxlan2kernel-vfio2noop --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
		r.Run(`kubectl -n ns-kernel2vxlan2kernel-vfio2noop exec ${NSE_VFIO} --container ponger -- /bin/bash -c '\` + "\n" + `  sleep 10 && kill $(pgrep "pingpong") 1>/dev/null 2>&1 &                    \` + "\n" + `'`)
		r.Run(`kubectl delete ns ns-kernel2vxlan2kernel-vfio2noop`)
	})
	r.Run(`kubectl create ns ns-kernel2vxlan2kernel-vfio2noop`)
	r.Run(`function mac_create(){` + "\n" + `    echo -n 00` + "\n" + `    dd bs=1 count=5 if=/dev/random 2>/dev/null | hexdump -v -e '/1 ":%02x"'` + "\n" + `}`)
	r.Run(`CLIENT_MAC=$(mac_create)` + "\n" + `echo Client MAC: ${CLIENT_MAC}`)
	r.Run(`SERVER_MAC=$(mac_create)` + "\n" + `echo Server MAC: ${SERVER_MAC}`)
	r.Run(`cat > patch-nse-vfio.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse-vfio` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: sidecar` + "\n" + `          env:` + "\n" + `            - name: NSM_SERVICES` + "\n" + `              value: "pingpong@worker.domain: { addr: ${SERVER_MAC} }"` + "\n" + `        - name: ponger` + "\n" + `          command: ["/bin/bash", "/root/scripts/pong.sh", "ens6f3", "31", ${SERVER_MAC}]` + "\n" + `EOF`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Kernel2Vxlan2Kernel&Vfio2Noop?ref=d9c9ce7b315b179188c887b89ad1843af03a2dfd`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-kernel -n ns-kernel2vxlan2kernel-vfio2noop`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-kernel2vxlan2kernel-vfio2noop`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-vfio -n ns-kernel2vxlan2kernel-vfio2noop`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-vfio -n ns-kernel2vxlan2kernel-vfio2noop`)
	r.Run(`NSC_KERNEL=$(kubectl get pods -l app=nsc-kernel -n ns-kernel2vxlan2kernel-vfio2noop --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE_KERNEL=$(kubectl get pods -l app=nse-kernel -n ns-kernel2vxlan2kernel-vfio2noop --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSC_VFIO=$(kubectl get pods -l app=nsc-vfio -n ns-kernel2vxlan2kernel-vfio2noop --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`function dpdk_ping() {` + "\n" + `  err_file="$(mktemp)"` + "\n" + `  trap 'rm -f "${err_file}"' RETURN` + "\n" + `` + "\n" + `  client_mac="$1"` + "\n" + `  server_mac="$2"` + "\n" + `` + "\n" + `  command="/root/dpdk-pingpong/build/app/pingpong \` + "\n" + `      --no-huge                                   \` + "\n" + `      --                                          \` + "\n" + `      -n 500                                      \` + "\n" + `      -c                                          \` + "\n" + `      -C ${client_mac}                            \` + "\n" + `      -S ${server_mac}` + "\n" + `      "` + "\n" + `  out="$(kubectl -n ns-kernel2vxlan2kernel-vfio2noop exec ${NSC_VFIO} --container pinger -- /bin/bash -c "${command}" 2>"${err_file}")"` + "\n" + `` + "\n" + `  if [[ "$?" != 0 ]]; then` + "\n" + `    echo "${out}"` + "\n" + `    cat "${err_file}" 1>&2` + "\n" + `    return 1` + "\n" + `  fi` + "\n" + `` + "\n" + `  if ! pong_packets="$(echo "${out}" | grep "rx .* pong packets" | sed -E 's/rx ([0-9]*) pong packets/\1/g')"; then` + "\n" + `    echo "${out}"` + "\n" + `    cat "${err_file}" 1>&2` + "\n" + `    return 1` + "\n" + `  fi` + "\n" + `` + "\n" + `  if [[ "${pong_packets}" == 0 ]]; then` + "\n" + `    echo "${out}"` + "\n" + `    cat "${err_file}" 1>&2` + "\n" + `    return 1` + "\n" + `  fi` + "\n" + `` + "\n" + `  echo "${out}"` + "\n" + `  return 0` + "\n" + `}`)
	r.Run(`kubectl exec ${NSC_KERNEL} -n ns-kernel2vxlan2kernel-vfio2noop -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl exec ${NSE_KERNEL} -n ns-kernel2vxlan2kernel-vfio2noop -- ping -c 4 172.16.1.101`)
	r.Run(`dpdk_ping ${CLIENT_MAC} ${SERVER_MAC}`)
}
func (s *Suite) TestMemif2Memif() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Memif2Memif")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-memif2memif`)
	})
	r.Run(`kubectl create ns ns-memif2memif`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Memif2Memif?ref=d9c9ce7b315b179188c887b89ad1843af03a2dfd`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-memif -n ns-memif2memif`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-memif -n ns-memif2memif`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc-memif -n ns-memif2memif --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-memif -n ns-memif2memif --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`result=$(kubectl exec "${NSC}" -n "ns-memif2memif" -- vppctl ping 172.16.1.100 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`result=$(kubectl exec "${NSE}" -n "ns-memif2memif" -- vppctl ping 172.16.1.101 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}
func (s *Suite) TestSriovKernel2Noop() {
	r := s.Runner("../deployments-k8s/examples/use-cases/SriovKernel2Noop")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-sriov-kernel2noop`)
	})
	r.Run(`kubectl create ns ns-sriov-kernel2noop`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/SriovKernel2Noop?ref=d9c9ce7b315b179188c887b89ad1843af03a2dfd`)
	r.Run(`kubectl -n ns-sriov-kernel2noop wait --for=condition=ready --timeout=1m pod -l app=nsc-kernel`)
	r.Run(`kubectl -n ns-sriov-kernel2noop wait --for=condition=ready --timeout=1m pod -l app=nse-kernel`)
	r.Run(`kubectl -n ns-sriov-kernel2noop wait --for=condition=ready --timeout=1m pod -l app=ponger`)
	r.Run(`NSC=$(kubectl -n ns-sriov-kernel2noop get pods -l app=nsc-kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl -n ns-sriov-kernel2noop exec ${NSC} -- ping -c 4 172.16.1.100`)
}
func (s *Suite) TestVfio2Noop() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Vfio2Noop")
	s.T().Cleanup(func() {
		r.Run(`NSE=$(kubectl -n ns-vfio2noop get pods -l app=nse-vfio --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
		r.Run(`kubectl -n ns-vfio2noop exec ${NSE} --container ponger -- /bin/bash -c '\` + "\n" + `  sleep 10 && kill $(pgrep "pingpong") 1>/dev/null 2>&1 &               \` + "\n" + `'`)
		r.Run(`kubectl delete ns ns-vfio2noop`)
	})
	r.Run(`kubectl create ns ns-vfio2noop`)
	r.Run(`function mac_create(){` + "\n" + `    echo -n 00` + "\n" + `    dd bs=1 count=5 if=/dev/random 2>/dev/null | hexdump -v -e '/1 ":%02x"'` + "\n" + `}`)
	r.Run(`CLIENT_MAC=$(mac_create)` + "\n" + `echo Client MAC: ${CLIENT_MAC}`)
	r.Run(`SERVER_MAC=$(mac_create)` + "\n" + `echo Server MAC: ${SERVER_MAC}`)
	r.Run(`cat > patch-nse-vfio.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse-vfio` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: sidecar` + "\n" + `          env:` + "\n" + `            - name: NSM_SERVICES` + "\n" + `              value: "pingpong@worker.domain: { addr: ${SERVER_MAC} }"` + "\n" + `        - name: ponger` + "\n" + `          command: ["/bin/bash", "/root/scripts/pong.sh", "ens6f3", "31", ${SERVER_MAC}]` + "\n" + `EOF`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Vfio2Noop?ref=d9c9ce7b315b179188c887b89ad1843af03a2dfd`)
	r.Run(`kubectl -n ns-vfio2noop wait --for=condition=ready --timeout=1m pod -l app=nsc-vfio`)
	r.Run(`kubectl -n ns-vfio2noop wait --for=condition=ready --timeout=1m pod -l app=nse-vfio`)
	r.Run(`NSC_VFIO=$(kubectl -n ns-vfio2noop get pods -l app=nsc-vfio --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`function dpdk_ping() {` + "\n" + `  err_file="$(mktemp)"` + "\n" + `  trap 'rm -f "${err_file}"' RETURN` + "\n" + `` + "\n" + `  client_mac="$1"` + "\n" + `  server_mac="$2"` + "\n" + `` + "\n" + `  command="/root/dpdk-pingpong/build/app/pingpong \` + "\n" + `      --no-huge                                   \` + "\n" + `      --                                          \` + "\n" + `      -n 500                                      \` + "\n" + `      -c                                          \` + "\n" + `      -C ${client_mac}                            \` + "\n" + `      -S ${server_mac}` + "\n" + `      "` + "\n" + `  out="$(kubectl -n ns-vfio2noop exec ${NSC_VFIO} --container pinger -- /bin/bash -c "${command}" 2>"${err_file}")"` + "\n" + `` + "\n" + `  if [[ "$?" != 0 ]]; then` + "\n" + `    echo "${out}"` + "\n" + `    cat "${err_file}" 1>&2` + "\n" + `    return 1` + "\n" + `  fi` + "\n" + `` + "\n" + `  if ! pong_packets="$(echo "${out}" | grep "rx .* pong packets" | sed -E 's/rx ([0-9]*) pong packets/\1/g')"; then` + "\n" + `    echo "${out}"` + "\n" + `    cat "${err_file}" 1>&2` + "\n" + `    return 1` + "\n" + `  fi` + "\n" + `` + "\n" + `  if [[ "${pong_packets}" == 0 ]]; then` + "\n" + `    echo "${out}"` + "\n" + `    cat "${err_file}" 1>&2` + "\n" + `    return 1` + "\n" + `  fi` + "\n" + `` + "\n" + `  echo "${out}"` + "\n" + `  return 0` + "\n" + `}`)
	r.Run(`dpdk_ping ${CLIENT_MAC} ${SERVER_MAC}`)
}
