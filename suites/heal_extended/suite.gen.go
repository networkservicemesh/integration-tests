// Code generated by gotestmd DO NOT EDIT.
package heal_extended

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
func (s *Suite) TestComponent_restart() {
	r := s.Runner("../deployments-k8s/examples/heal_extended/component-restart")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-component-restart`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/heal_extended/component-restart?ref=v1.13.1-rc.4`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=client -n ns-component-restart`)
	r.Run(`# N_RESTARTS - number of restarts` + "\n" + `# TEST_TIME - determines how long the test will take (sec)` + "\n" + `# DELAY - delay between restarts (sec)` + "\n" + `# INTERFACE_READY_WAIT - how long do we wait for the interface to be ready (sec). Equals to NSM_REQUEST_TIMEOUT * 2 (for Close and Request)` + "\n" + `N_RESTARTS=15` + "\n" + `TEST_TIME=900` + "\n" + `DELAY=$(($TEST_TIME/$N_RESTARTS))` + "\n" + `INTERFACE_READY_WAIT=10`)
	r.Run(`# Iterates over NSCs and checks connectivity to NSE (sends pings)` + "\n" + `function connectivity_check() {` + "\n" + `echo -e "\n-- Connectivity check --"` + "\n" + `nscs=$(kubectl  get pods -l app=client -o go-template --template="{{range .items}}{{.metadata.name}} {{end}}" -n ns-component-restart)` + "\n" + `for nsc in $nscs` + "\n" + `do` + "\n" + `    echo -e "\nNSC: $nsc"` + "\n" + `    echo "Wait for NSM interface to be ready"` + "\n" + `    for i in $(seq 1 $INTERFACE_READY_WAIT)` + "\n" + `    do` + "\n" + `        if [ $i -eq $INTERFACE_READY_WAIT ] ; then` + "\n" + `          echo "NSM interface is not ready after $INTERFACE_READY_WAIT s"` + "\n" + `          return 1` + "\n" + `        fi` + "\n" + `        sleep 1` + "\n" + `        routes=$(kubectl exec -n ns-component-restart $nsc -- ip route)` + "\n" + `        nseAddr=$(echo $routes | grep -Eo '172\.16\.1\.[0-9]{1,3}')` + "\n" + `        test $? -ne 0 || break` + "\n" + `    done` + "\n" + `    echo "NSM interface is ready"` + "\n" + `    kubectl exec $nsc -n ns-component-restart -- ping -c2 -i 0.5 $nseAddr || return 2` + "\n" + `done` + "\n" + `return 0` + "\n" + `}` + "\n" + `` + "\n" + `# Restarts NSM components and checks connectivity.` + "\n" + `# $1 is used to define NSM-component type (e.g. forwarder or nsmgr)` + "\n" + `# -a defines the restart method.` + "\n" + `#   if specified - all NSM-pods of this type will be restarted at the same time.` + "\n" + `#   else - they will be restarted one by one.` + "\n" + `function restart_nsm_component() {` + "\n" + `nsm_component=$1` + "\n" + `shift` + "\n" + `` + "\n" + `a_flag=0` + "\n" + `while getopts 'a' flag; do` + "\n" + `  case "${flag}" in` + "\n" + `    a) a_flag=1 ;;` + "\n" + `  esac` + "\n" + `done` + "\n" + `` + "\n" + `for i in $(seq 1 $N_RESTARTS)` + "\n" + `do` + "\n" + `    echo -e "\n-------- $nsm_component restart $i of $N_RESTARTS --------"` + "\n" + `    echo "Wait $DELAY sec before restart..."` + "\n" + `    sleep $DELAY` + "\n" + `    if [ $a_flag -eq 1 ]; then` + "\n" + `        kubectl delete pod -n nsm-system -l app=${nsm_component}` + "\n" + `        kubectl wait --for=condition=ready --timeout=1m pod -l app=${nsm_component} -n nsm-system || return 1` + "\n" + `        connectivity_check || return 2` + "\n" + `    else` + "\n" + `        nodes=$(kubectl get pods -l app=${nsm_component} -n nsm-system --template '{{range .items}}{{.spec.nodeName}}{{"\n"}}{{end}}')` + "\n" + `        for node in $nodes` + "\n" + `        do` + "\n" + `            kubectl delete pod -n nsm-system -l app=${nsm_component} --field-selector spec.nodeName==${node}` + "\n" + `            kubectl wait --for=condition=ready --timeout=1m pod -l app=${nsm_component} --field-selector spec.nodeName==${node} -n nsm-system || return 1` + "\n" + `            connectivity_check || return 2` + "\n" + `        done` + "\n" + `    fi` + "\n" + `done` + "\n" + `return 0` + "\n" + `}`)
	r.Run(`connectivity_check`)
	r.Run(`restart_nsm_component forwarder-vpp`)
	r.Run(`restart_nsm_component forwarder-vpp -a`)
}
