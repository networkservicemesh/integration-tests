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
	r := s.Runner("/home/nikita/repos/NSM/deployments-k8s/examples/multicluster/spiffe_federation")
	r.Run(`bundle1=$(kubectl --kubeconfig=$KUBECONFIG1 exec spire-server-0 -n spire -- bin/spire-server bundle show -format spiffe)` + "\n" + `bundle2=$(kubectl --kubeconfig=$KUBECONFIG2 exec spire-server-0 -n spire -- bin/spire-server bundle show -format spiffe)` + "\n" + `bundle3=$(kubectl --kubeconfig=$KUBECONFIG3 exec spire-server-0 -n spire -- bin/spire-server bundle show -format spiffe)`)
	r.Run(`echo $bundle2 | kubectl --kubeconfig=$KUBECONFIG1 exec -i spire-server-0 -n spire -- bin/spire-server bundle set -format spiffe -id "spiffe://nsm.cluster2"` + "\n" + `echo $bundle3 | kubectl --kubeconfig=$KUBECONFIG1 exec -i spire-server-0 -n spire -- bin/spire-server bundle set -format spiffe -id "spiffe://nsm.cluster3"`)
	r.Run(`echo $bundle1 | kubectl --kubeconfig=$KUBECONFIG2 exec -i spire-server-0 -n spire -- bin/spire-server bundle set -format spiffe -id "spiffe://nsm.cluster1"` + "\n" + `echo $bundle3 | kubectl --kubeconfig=$KUBECONFIG2 exec -i spire-server-0 -n spire -- bin/spire-server bundle set -format spiffe -id "spiffe://nsm.cluster3"`)
	r.Run(`echo $bundle1 | kubectl --kubeconfig=$KUBECONFIG3 exec -i spire-server-0 -n spire -- bin/spire-server bundle set -format spiffe -id "spiffe://nsm.cluster1"` + "\n" + `echo $bundle2 | kubectl --kubeconfig=$KUBECONFIG3 exec -i spire-server-0 -n spire -- bin/spire-server bundle set -format spiffe -id "spiffe://nsm.cluster2"`)
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
