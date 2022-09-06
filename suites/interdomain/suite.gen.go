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
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/nsm_consul/nse-auto-scale?ref=358779b8e18dca94b9789d7bb66094f8bb28a9a5`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/358779b8e18dca94b9789d7bb66094f8bb28a9a5/examples/interdomain/nsm_consul/client/dashboard.yaml`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/358779b8e18dca94b9789d7bb66094f8bb28a9a5/examples/interdomain/nsm_consul/networkservice.yaml`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete pods --all`)
		r.Run(`consul-k8s uninstall --kubeconfig=$KUBECONFIG2 -auto-approve=true -wipe-data=true`)
	})
	r.Run(`brew tap hashicorp/tap` + "\n" + `brew install hashicorp/tap/consul-k8s`)
	r.Run(`consul-k8s install -config-file=helm-consul-values.yaml -set global.image=hashicorp/consul:1.12.0 -auto-approve --kubeconfig=$KUBECONFIG2`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/358779b8e18dca94b9789d7bb66094f8bb28a9a5/examples/interdomain/nsm_consul/networkservice.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/358779b8e18dca94b9789d7bb66094f8bb28a9a5/examples/interdomain/nsm_consul/client/dashboard.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/358779b8e18dca94b9789d7bb66094f8bb28a9a5/examples/interdomain/nsm_consul/service.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/nsm_consul/nse-auto-scale?ref=358779b8e18dca94b9789d7bb66094f8bb28a9a5`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/358779b8e18dca94b9789d7bb66094f8bb28a9a5/examples/interdomain/nsm_consul/server/counting.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --timeout=5m --for=condition=ready pod -l app=dashboard-nsc`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pod/dashboard-nsc -c cmd-nsc -- apk add curl`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec pod/dashboard-nsc -c cmd-nsc -- curl counting:9001`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 port-forward pod/dashboard-nsc 9002:9002 &`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete deploy counting`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/358779b8e18dca94b9789d7bb66094f8bb28a9a5/examples/interdomain/nsm_consul/server/counting_nsm.yaml`)
}
func (s *Suite) TestNsm_istio() {
	r := s.Runner("../deployments-k8s/examples/interdomain/nsm_istio")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/358779b8e18dca94b9789d7bb66094f8bb28a9a5/examples/interdomain/nsm_istio/greeting/server.yaml` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 delete -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/nsm_istio/nse-auto-scale?ref=358779b8e18dca94b9789d7bb66094f8bb28a9a5` + "\n" + `kubectl --kubeconfig=$KUBECONFIG1 delete -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/358779b8e18dca94b9789d7bb66094f8bb28a9a5/examples/interdomain/nsm_istio/greeting/client.yaml` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 delete -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/358779b8e18dca94b9789d7bb66094f8bb28a9a5/examples/interdomain/nsm_istio/networkservice.yaml` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 delete ns istio-system` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 label namespace default istio-injection-` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 delete pods --all`)
	})
	r.Run(`curl -sL https://istio.io/downloadIstioctl | sh -` + "\n" + `export PATH=$PATH:$HOME/.istioctl/bin` + "\n" + `istioctl  install --set profile=minimal -y --kubeconfig=$KUBECONFIG2` + "\n" + `istioctl --kubeconfig=$KUBECONFIG2 proxy-status`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/358779b8e18dca94b9789d7bb66094f8bb28a9a5/examples/interdomain/nsm_istio/networkservice.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/358779b8e18dca94b9789d7bb66094f8bb28a9a5/examples/interdomain/nsm_istio/greeting/client.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k https://github.com/networkservicemesh/deployments-k8s/examples/interdomain/nsm_istio/nse-auto-scale?ref=358779b8e18dca94b9789d7bb66094f8bb28a9a5`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 label namespace default istio-injection=enabled` + "\n" + `` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/358779b8e18dca94b9789d7bb66094f8bb28a9a5/examples/interdomain/nsm_istio/greeting/server.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --timeout=2m --for=condition=ready pod -l app=alpine`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec deploy/alpine -c cmd-nsc -- apk add curl`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec deploy/alpine -c cmd-nsc -- curl -s greeting.default:9080 | grep -o "hello world from istio"`)
}
func (s *Suite) TestNsm_kuma_universal_vl3() {
	r := s.Runner("../deployments-k8s/examples/interdomain/nsm_kuma_universal_vl3")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete ns kuma-system kuma-demo ns-dns-vl3` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 delete ns kuma-demo` + "\n" + `rm tls.crt tls.key ca.crt` + "\n" + `rm -rf kuma-1.7.0`)
	})
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -k ./vl3-dns` + "\n" + `kubectl --kubeconfig=$KUBECONFIG1 -n ns-dns-vl3 wait --for=condition=ready --timeout=2m pod -l app=vl3-ipam`)
	r.Run(`curl -L https://kuma.io/installer.sh | VERSION=1.7.0 ARCH=amd64 bash -` + "\n" + `export PATH=$PWD/kuma-1.7.0/bin:$PATH`)
	r.Run(`kumactl generate tls-certificate --hostname=control-plane-kuma.my-vl3-network --hostname=kuma-control-plane.kuma-system.svc --type=server --key-file=./tls.key --cert-file=./tls.crt`)
	r.Run(`cp ./tls.crt ./ca.crt`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -f namespace.yaml` + "\n" + `kubectl --kubeconfig=$KUBECONFIG1 create secret generic general-tls-certs --namespace=kuma-system --from-file=./tls.key --from-file=./tls.crt --from-file=./ca.crt`)
	r.Run(`kumactl install control-plane --tls-general-secret=general-tls-certs --tls-general-ca-bundle=$(cat ./ca.crt | base64) > ./control-plane/control-plane.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -k ./control-plane`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -f demo-redis.yaml` + "\n" + `kubectl --kubeconfig=$KUBECONFIG1 -n kuma-demo wait --for=condition=ready --timeout=3m pod -l app=redis`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -f demo-app.yaml` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 -n kuma-demo wait --for=condition=ready --timeout=3m pod -l app=demo-app`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 port-forward svc/demo-app -n kuma-demo 5000:5000 &`)
	r.Run(`response=$(curl -X POST localhost:5000/increment)`)
	r.Run(`echo $response | grep '"err":null'`)
}
