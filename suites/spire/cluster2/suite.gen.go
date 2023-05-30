// Code generated by gotestmd DO NOT EDIT.
package cluster2

import (
	"fmt"
	"sync"
	"testing"

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
	r := s.Runner("../deployments-k8s/examples/spire/cluster2")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete crd clusterspiffeids.spire.spiffe.io` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 delete crd clusterfederatedtrustdomains.spire.spiffe.io` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 delete validatingwebhookconfiguration.admissionregistration.k8s.io/spire-controller-manager-webhook` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 delete ns spire`)
	})
	r.Run(`[[ ! -z $KUBECONFIG2 ]]`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/spire/cluster2?ref=5a9bdf42902474b17fea95ab459ce98d7b5aa3d0`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 wait -n spire --timeout=1m --for=condition=ready pod -l app=spire-server`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 wait -n spire --timeout=1m --for=condition=ready pod -l app=spire-agent`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/5a9bdf42902474b17fea95ab459ce98d7b5aa3d0/examples/spire/cluster2/clusterspiffeid-template.yaml`)
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
	tests := []func(t *testing.T){}
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
func (s *Suite) Test() {}
