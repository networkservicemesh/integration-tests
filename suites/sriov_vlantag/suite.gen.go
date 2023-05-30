// Code generated by gotestmd DO NOT EDIT.
package sriov_vlantag

import (
	"fmt"
	"sync"
	"testing"

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
	r := s.Runner("/home/nikita/repos/NSM/deployments-k8s/examples/sriov_vlantag")
	s.T().Cleanup(func() {
		r.Run(`WH=$(kubectl get pods -l app=admission-webhook-k8s -n nsm-system --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `kubectl delete mutatingwebhookconfiguration ${WH}` + "\n" + `kubectl delete ns nsm-system`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/sriov?ref=5a9bdf42902474b17fea95ab459ce98d7b5aa3d0`)
}

const workerCount = 5

func worker(jobsCh <-chan func(), wg *sync.WaitGroup) {
	for j := range jobsCh {
		fmt.Println("Executing a job...")
		j()
	}
	fmt.Println("Worker is finishing...")
	wg.Done()
}
func (s *Suite) TestAll() {
	tests := []func(t *testing.T){
		s.SriovKernel2NoopVlanTag,
		s.Vfio2NoopVlanTag,
	}
	jobCh := make(chan func(), len(tests))
	wg := new(sync.WaitGroup)
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker(jobCh, wg)
	}
	for i := range tests {
		test := tests[i]
		jobCh <- func() {
			s.T().Run("TestName", test)
		}
	}
	wg.Wait()
}
func (s *Suite) SriovKernel2NoopVlanTag(t *testing.T) {
	r := s.Runner("/home/nikita/repos/NSM/deployments-k8s/examples/use-cases/SriovKernel2NoopVlanTag")
	t.Cleanup(func() {
		r.Run(`kubectl delete ns ns-sriov-kernel2noop-vlantag`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/SriovKernel2NoopVlanTag/ponger?ref=5a9bdf42902474b17fea95ab459ce98d7b5aa3d0`)
	r.Run(`kubectl -n ns-sriov-kernel2noop-vlantag wait --for=condition=ready --timeout=1m pod -l app=ponger`)
	r.Run(`kubectl -n ns-sriov-kernel2noop-vlantag exec deploy/ponger -- ip a | grep "172.16.1.100"`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/SriovKernel2NoopVlanTag?ref=5a9bdf42902474b17fea95ab459ce98d7b5aa3d0`)
	r.Run(`kubectl -n ns-sriov-kernel2noop-vlantag wait --for=condition=ready --timeout=1m pod -l app=nsc-kernel`)
	r.Run(`kubectl -n ns-sriov-kernel2noop-vlantag wait --for=condition=ready --timeout=1m pod -l app=nse-noop`)
	r.Run(`kubectl -n ns-sriov-kernel2noop-vlantag exec deployments/nsc-kernel -- ping -c 4 172.16.1.100`)
}
func (s *Suite) Vfio2NoopVlanTag(t *testing.T) {
	r := s.Runner("/home/nikita/repos/NSM/deployments-k8s/examples/use-cases/Vfio2NoopVlanTag")
	t.Cleanup(func() {
		r.Run(`kubectl -n ns-vfio2noop-vlantag exec deployments/nse-vfio --container ponger -- /bin/bash -c '\` + "\n" + `  (sleep 10 && kill $(pgrep "pingpong")) 1>/dev/null 2>&1 &             \` + "\n" + `'`)
		r.Run(`kubectl delete ns ns-vfio2noop-vlantag`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Vfio2NoopVlanTag?ref=5a9bdf42902474b17fea95ab459ce98d7b5aa3d0`)
	r.Run(`kubectl -n ns-vfio2noop-vlantag wait --for=condition=ready --timeout=1m pod -l app=nsc-vfio`)
	r.Run(`kubectl -n ns-vfio2noop-vlantag wait --for=condition=ready --timeout=1m pod -l app=nse-vfio`)
	r.Run(`function dpdk_ping() {` + "\n" + `  err_file="$(mktemp)"` + "\n" + `  trap 'rm -f "${err_file}"' RETURN` + "\n" + `` + "\n" + `  client_mac="$1"` + "\n" + `  server_mac="$2"` + "\n" + `` + "\n" + `  command="/root/dpdk-pingpong/build/pingpong \` + "\n" + `      --no-huge                               \` + "\n" + `      --                                      \` + "\n" + `      -n 500                                  \` + "\n" + `      -c                                      \` + "\n" + `      -C ${client_mac}                        \` + "\n" + `      -S ${server_mac}` + "\n" + `      "` + "\n" + `  out="$(kubectl -n ns-vfio2noop-vlantag exec deployments/nsc-vfio --container pinger -- /bin/bash -c "${command}" 2>"${err_file}")"` + "\n" + `` + "\n" + `  if [[ "$?" != 0 ]]; then` + "\n" + `    echo "${out}"` + "\n" + `    cat "${err_file}" 1>&2` + "\n" + `    return 1` + "\n" + `  fi` + "\n" + `` + "\n" + `  if ! pong_packets="$(echo "${out}" | grep "rx .* pong packets" | sed -E 's/rx ([0-9]*) pong packets/\1/g')"; then` + "\n" + `    echo "${out}"` + "\n" + `    cat "${err_file}" 1>&2` + "\n" + `    return 1` + "\n" + `  fi` + "\n" + `` + "\n" + `  if [[ "${pong_packets}" == 0 ]]; then` + "\n" + `    echo "${out}"` + "\n" + `    cat "${err_file}" 1>&2` + "\n" + `    return 1` + "\n" + `  fi` + "\n" + `` + "\n" + `  echo "${out}"` + "\n" + `  return 0` + "\n" + `}`)
	r.Run(`dpdk_ping "0a:55:44:33:22:00" "0a:55:44:33:22:11"`)
}
