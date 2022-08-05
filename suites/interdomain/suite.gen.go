// Code generated by gotestmd DO NOT EDIT.
package interdomain

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/interdomain/dns"
	"github.com/networkservicemesh/integration-tests/suites/interdomain/loadbalancer"
	"github.com/networkservicemesh/integration-tests/suites/interdomain/nsm"
	"github.com/networkservicemesh/integration-tests/suites/interdomain/spire"
)

type Suite struct {
	base.Suite
	loadbalancerSuite loadbalancer.Suite
	dnsSuite          dns.Suite
	spireSuite        spire.Suite
	nsmSuite          nsm.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite, &s.loadbalancerSuite, &s.dnsSuite, &s.spireSuite, &s.nsmSuite}
	for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
}
func (s *Suite) TestNsm_consul() {
	r := s.Runner("../deployments-k8s/examples/interdomain/nsm_consul")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete deployment counting`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/nsm_consul/nse-auto-scale?ref=3d1dcfe1de90681213c7f0006f25279bb4699966`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/3d1dcfe1de90681213c7f0006f25279bb4699966/examples/interdomain/nsm_consul/client/dashboard.yaml`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/3d1dcfe1de90681213c7f0006f25279bb4699966/examples/interdomain/nsm_consul/networkservice.yaml`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete pods --all`)
		r.Run(`consul-k8s uninstall --kubeconfig=$KUBECONFIG2 -auto-approve=true -wipe-data=true`)
	})
	r.Run(`brew tap hashicorp/tap` + "\n" + `brew install hashicorp/tap/consul-k8s`)
	r.Run(`consul-k8s install -config-file=helm-consul-values.yaml -set global.image=hashicorp/consul:1.12.0 -auto-approve --kubeconfig=$KUBECONFIG2`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/3d1dcfe1de90681213c7f0006f25279bb4699966/examples/interdomain/nsm_consul/networkservice.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/3d1dcfe1de90681213c7f0006f25279bb4699966/examples/interdomain/nsm_consul/client/dashboard.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/3d1dcfe1de90681213c7f0006f25279bb4699966/examples/interdomain/nsm_consul/service.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/nsm_consul/nse-auto-scale?ref=3d1dcfe1de90681213c7f0006f25279bb4699966`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/3d1dcfe1de90681213c7f0006f25279bb4699966/examples/interdomain/nsm_consul/server/counting.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --timeout=5m --for=condition=ready pod -l app=dashboard-nsc`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pod/dashboard-nsc -c cmd-nsc -- apk add curl`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pod/dashboard-nsc -c cmd-nsc -- curl counting:9001`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 port-forward pod/dashboard-nsc 9002:9002 &`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete deploy counting`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/3d1dcfe1de90681213c7f0006f25279bb4699966/examples/interdomain/nsm_consul/server/counting_nsm.yaml`)
}
func (s *Suite) TestNsm_istio() {
	r := s.Runner("../deployments-k8s/examples/interdomain/nsm_istio")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/3d1dcfe1de90681213c7f0006f25279bb4699966/examples/interdomain/nsm_istio/greeting/server.yaml` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 delete -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/nsm_istio/nse-auto-scale?ref=3d1dcfe1de90681213c7f0006f25279bb4699966` + "\n" + `kubectl --kubeconfig=$KUBECONFIG1 delete -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/3d1dcfe1de90681213c7f0006f25279bb4699966/examples/interdomain/nsm_istio/greeting/client.yaml` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 delete -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/3d1dcfe1de90681213c7f0006f25279bb4699966/examples/interdomain/nsm_istio/networkservice.yaml` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 delete ns istio-system` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 label namespace default istio-injection-` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 delete pods --all`)
	})
	r.Run(`curl -sL https://istio.io/downloadIstioctl | sh -` + "\n" + `export PATH=$PATH:$HOME/.istioctl/bin` + "\n" + `istioctl  install --set profile=minimal -y --kubeconfig=$KUBECONFIG2` + "\n" + `istioctl --kubeconfig=$KUBECONFIG2 proxy-status`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/3d1dcfe1de90681213c7f0006f25279bb4699966/examples/interdomain/nsm_istio/networkservice.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/3d1dcfe1de90681213c7f0006f25279bb4699966/examples/interdomain/nsm_istio/greeting/client.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/nsm_istio/nse-auto-scale?ref=3d1dcfe1de90681213c7f0006f25279bb4699966`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 label namespace default istio-injection=enabled` + "\n" + `` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/3d1dcfe1de90681213c7f0006f25279bb4699966/examples/interdomain/nsm_istio/greeting/server.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --timeout=2m --for=condition=ready pod -l app=alpine`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec deploy/alpine -c cmd-nsc -- apk add curl`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec deploy/alpine -c cmd-nsc -- curl -s greeting.default:9080 | grep -o "hello world from istio"`)
}
