// Code generated by gotestmd DO NOT EDIT.
package spiffe_federation

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
	r := s.Runner("../deployments-k8s/examples/interdomain/spiffe_federation")
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/5a9bdf42902474b17fea95ab459ce98d7b5aa3d0/examples/interdomain/spiffe_federation/cluster1-spiffeid-template.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/5a9bdf42902474b17fea95ab459ce98d7b5aa3d0/examples/interdomain/spiffe_federation/cluster2-spiffeid-template.yaml`)
	r.Run(`bundle1=$(kubectl --kubeconfig=$KUBECONFIG1 exec spire-server-0 -n spire -- bin/spire-server bundle show -format spiffe)` + "\n" + `bundle2=$(kubectl --kubeconfig=$KUBECONFIG2 exec spire-server-0 -n spire -- bin/spire-server bundle show -format spiffe)`)
	r.Run(`echo $bundle2 | kubectl --kubeconfig=$KUBECONFIG1 exec -i spire-server-0 -n spire -- bin/spire-server bundle set -format spiffe -id "spiffe://nsm.cluster2"` + "\n" + `echo $bundle1 | kubectl --kubeconfig=$KUBECONFIG2 exec -i spire-server-0 -n spire -- bin/spire-server bundle set -format spiffe -id "spiffe://nsm.cluster1"`)
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
