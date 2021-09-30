// Code generated by gotestmd DO NOT EDIT.
package sriov

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
	r := s.Runner("../deployments-k8s/examples/sriov")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete mutatingwebhookconfiguration --all` + "\n" + `kubectl delete ns nsm-system`)
	})
	r.Run(`kubectl create ns nsm-system`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/nsm-system/sa/default \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:nsm-system \` + "\n" + `-selector k8s:sa:default`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/nsm-system/sa/registry-k8s-sa \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:nsm-system \` + "\n" + `-selector k8s:sa:registry-k8s-sa`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/sriov?ref=05197ded333d78d1c2073f4152ab3e5f57a3d858`)
}
func (s *Suite) TestSriovKernel2Noop() {
	r := s.Runner("../deployments-k8s/examples/use-cases/SriovKernel2Noop")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/05197ded333d78d1c2073f4152ab3e5f57a3d858/examples/use-cases/namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/${NAMESPACE}/sa/default \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:${NAMESPACE} \` + "\n" + `-selector k8s:sa:default`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- github.com/networkservicemesh/deployments-k8s/apps/nsc-kernel?ref=05197ded333d78d1c2073f4152ab3e5f57a3d858` + "\n" + `- github.com/networkservicemesh/deployments-k8s/apps/nse-kernel?ref=05197ded333d78d1c2073f4152ab3e5f57a3d858` + "\n" + `- github.com/networkservicemesh/deployments-k8s/apps/nsc-kernel-ponger?ref=05197ded333d78d1c2073f4152ab3e5f57a3d858` + "\n" + `` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nsc.yaml` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nsc.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nsc-kernel` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nsc` + "\n" + `          env:` + "\n" + `            - name: NSM_NETWORK_SERVICES` + "\n" + `              value: kernel://icmp-responder/nsm-1?sriovToken=worker.domain/10G` + "\n" + `          resources:` + "\n" + `            limits:` + "\n" + `              worker.domain/10G: 1` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse-kernel` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nse` + "\n" + `          env:` + "\n" + `            - name: NSM_LABELS` + "\n" + `              value: serviceDomain:worker.domain` + "\n" + `            - name: NSM_CIDR_PREFIX` + "\n" + `              value: 172.16.1.100/31` + "\n" + `          resources:` + "\n" + `            limits:` + "\n" + `              master.domain/10G: 1` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl -n ${NAMESPACE} wait --for=condition=ready --timeout=1m pod -l app=nsc-kernel`)
	r.Run(`kubectl -n ${NAMESPACE} wait --for=condition=ready --timeout=1m pod -l app=nse-kernel`)
	r.Run(`kubectl -n ${NAMESPACE} wait --for=condition=ready --timeout=1m pod -l app=ponger`)
	r.Run(`NSC=$(kubectl -n ${NAMESPACE} get pods -l app=nsc-kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl -n ${NAMESPACE} exec ${NSC} -- ping -c 4 172.16.1.100`)
}
func (s *Suite) TestVfio2Noop() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Vfio2Noop")
	s.T().Cleanup(func() {
		r.Run(`NSE=$(kubectl -n ${NAMESPACE} get pods -l app=nse-vfio --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
		r.Run(`kubectl -n ${NAMESPACE} exec ${NSE} --container ponger -- /bin/bash -c '\` + "\n" + `  sleep 10 && kill $(pgrep "pingpong") 1>/dev/null 2>&1 &               \` + "\n" + `'`)
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/05197ded333d78d1c2073f4152ab3e5f57a3d858/examples/use-cases/namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/${NAMESPACE}/sa/default \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:${NAMESPACE} \` + "\n" + `-selector k8s:sa:default`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- github.com/networkservicemesh/deployments-k8s/apps/nsc-vfio?ref=05197ded333d78d1c2073f4152ab3e5f57a3d858` + "\n" + `- github.com/networkservicemesh/deployments-k8s/apps/nse-vfio?ref=05197ded333d78d1c2073f4152ab3e5f57a3d858` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl -n ${NAMESPACE} wait --for=condition=ready --timeout=1m pod -l app=nsc-vfio`)
	r.Run(`kubectl -n ${NAMESPACE} wait --for=condition=ready --timeout=1m pod -l app=nse-vfio`)
	r.Run(`NSC=$(kubectl -n ${NAMESPACE} get pods -l app=nsc-vfio --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`function dpdk_ping() {` + "\n" + `  err_file="$(mktemp)"` + "\n" + `  trap 'rm -f "${err_file}"' RETURN` + "\n" + `` + "\n" + `  out="$(kubectl -n ${NAMESPACE} exec ${NSC} --container pinger -- /bin/bash -c '\` + "\n" + `    /root/dpdk-pingpong/build/app/pingpong                                       \` + "\n" + `      --no-huge                                                                  \` + "\n" + `      --                                                                         \` + "\n" + `      -n 500                                                                     \` + "\n" + `      -c                                                                         \` + "\n" + `      -C 0a:11:22:33:44:55                                                       \` + "\n" + `      -S 0a:55:44:33:22:11                                                       \` + "\n" + `  ' 2>"${err_file}")"` + "\n" + `` + "\n" + `  if [[ "$?" != 0 ]]; then` + "\n" + `    cat "${err_file}" 1>&2` + "\n" + `    echo "${out}" 1>&2` + "\n" + `    return 1` + "\n" + `  fi` + "\n" + `` + "\n" + `  if ! pong_packets="$(echo "${out}" | grep "rx .* pong packets" | sed -E 's/rx ([0-9]*) pong packets/\1/g')"; then` + "\n" + `    cat "${err_file}" 1>&2` + "\n" + `    echo "${out}" 1>&2` + "\n" + `    return 1` + "\n" + `  fi` + "\n" + `` + "\n" + `  if [[ "${pong_packets}" == 0 ]]; then` + "\n" + `    cat "${err_file}" 1>&2` + "\n" + `    echo "${out}" 1>&2` + "\n" + `    return 1` + "\n" + `  fi` + "\n" + `` + "\n" + `  echo "${out}"` + "\n" + `  return 0` + "\n" + `}`)
	r.Run(`dpdk_ping`)
}
