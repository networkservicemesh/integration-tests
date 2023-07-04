// Code generated by gotestmd DO NOT EDIT.
package nsm_istio_booking

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/interdomain/dns"
	"github.com/networkservicemesh/integration-tests/suites/interdomain/loadbalancer"
	"github.com/networkservicemesh/integration-tests/suites/interdomain/nsm"
	"github.com/networkservicemesh/integration-tests/suites/interdomain/spiffe_federation"
	"github.com/networkservicemesh/integration-tests/suites/spire/cluster1"
	"github.com/networkservicemesh/integration-tests/suites/spire/cluster2"
)

type Suite struct {
	base.Suite
	loadbalancerSuite      loadbalancer.Suite
	dnsSuite               dns.Suite
	cluster1Suite          cluster1.Suite
	cluster2Suite          cluster2.Suite
	spiffe_federationSuite spiffe_federation.Suite
	nsmSuite               nsm.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite, &s.loadbalancerSuite, &s.dnsSuite, &s.cluster1Suite, &s.cluster2Suite, &s.spiffe_federationSuite, &s.nsmSuite}
	for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
	r := s.Runner("../deployments-k8s/examples/interdomain/nsm_istio_booking")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete -f https://raw.githubusercontent.com/istio/istio/release-1.13/samples/bookinfo/platform/kube/bookinfo.yaml` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 delete -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/nsm_istio_booking/nse-auto-scale?ref=a521fec15f2eb7f15a6278b9236cc02336361a45 ` + "\n" + `kubectl --kubeconfig=$KUBECONFIG1 delete -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/a521fec15f2eb7f15a6278b9236cc02336361a45/examples/interdomain/nsm_istio_booking/productpage/productpage.yaml` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 delete -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/a521fec15f2eb7f15a6278b9236cc02336361a45/examples/interdomain/nsm_istio_booking/netsvc.yaml` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 delete ns istio-system` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 label namespace default istio-injection-` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 delete pods --all`)
	})
	r.Run(`curl -sL https://istio.io/downloadIstioctl | sh -` + "\n" + `export PATH=$PATH:$HOME/.istioctl/bin` + "\n" + `istioctl install --readiness-timeout 10m0s --set profile=minimal -y --kubeconfig=$KUBECONFIG2` + "\n" + `istioctl --kubeconfig=$KUBECONFIG2 proxy-status`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/a521fec15f2eb7f15a6278b9236cc02336361a45/examples/interdomain/nsm_istio_booking/netsvc.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/a521fec15f2eb7f15a6278b9236cc02336361a45/examples/interdomain/nsm_istio_booking/productpage/productpage.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/nsm_istio_booking/nse-auto-scale?ref=a521fec15f2eb7f15a6278b9236cc02336361a45`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 label namespace default istio-injection=enabled` + "\n" + `` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 apply -f https://raw.githubusercontent.com/istio/istio/release-1.13/samples/bookinfo/platform/kube/bookinfo.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --timeout=5m --for=condition=ready pod -l app=productpage`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec deploy/productpage-v1 -c cmd-nsc -- apk add curl`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec deploy/productpage-v1 -c cmd-nsc -- curl -s productpage.default:9080/productpage | grep -o "<title>Simple Bookstore App</title>"`)
}
func (s *Suite) Test() {}
