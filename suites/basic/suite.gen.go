// Code generated by gotestmd DO NOT EDIT.
package basic

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
	r := s.Runner("../deployments-k8s/examples/basic")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns nsm-system`)
	})
	r.Run(`kubectl create ns nsm-system`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/nsm-system/sa/default \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:nsm-system \` + "\n" + `-selector k8s:sa:default`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/nsm-system/sa/registry-k8s-sa \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:nsm-system \` + "\n" + `-selector k8s:sa:registry-k8s-sa`)
	r.Run(`kubectl apply -k .`)
}
func (s *Suite) TestOpa() {
	r := s.Runner("../deployments-k8s/examples/features/opa")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f ../namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/${NAMESPACE}/sa/default \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:${NAMESPACE} \` + "\n" + `-selector k8s:sa:default`)
	r.Run(`NODE=($(kubectl get nodes -o go-template='{{range .items}}{{ if not .spec.taints  }}{{index .metadata.labels "kubernetes.io/hostname"}} {{end}}{{end}}')[0])`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- ../../../apps/nsc-kernel` + "\n" + `- ../../../apps/nse-kernel` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nsc.yaml` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nsc.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nsc` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nsc` + "\n" + `          env:` + "\n" + `            - name: NSM_MAX_TOKEN_LIFETIME` + "\n" + `              value: -1m` + "\n" + `            - name: NSM_NETWORK_SERVICES` + "\n" + `              value: kernel://icmp-responder/nsm-1` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODE}` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nse` + "\n" + `          env:` + "\n" + `            - name: NSE_CIDR_PREFIX` + "\n" + `              value: 172.16.1.100/31` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODE}` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc -n ${NAMESPACE}`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse -n ${NAMESPACE}`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl logs ${NSC} -n ${NAMESPACE} | grep "PermissionDenied desc = no sufficient privileges"`)
}
func (s *Suite) TestKernel2Kernel() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f ../namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/${NAMESPACE}/sa/default \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:${NAMESPACE} \` + "\n" + `-selector k8s:sa:default`)
	r.Run(`NODE=($(kubectl get nodes -o go-template='{{range .items}}{{ if not .spec.taints  }}{{index .metadata.labels "kubernetes.io/hostname"}} {{end}}{{end}}')[0])`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- ../../../apps/nsc-kernel` + "\n" + `- ../../../apps/nse-kernel` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nsc.yaml` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nsc.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nsc` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nsc` + "\n" + `          env:` + "\n" + `            - name: NSM_NETWORK_SERVICES` + "\n" + `              value: kernel://icmp-responder/nsm-1` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODE}` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nse` + "\n" + `          env:` + "\n" + `            - name: NSE_CIDR_PREFIX` + "\n" + `              value: 172.16.1.100/31` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODE}` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc -n ${NAMESPACE}`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse -n ${NAMESPACE}`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl exec ${NSC} -n ${NAMESPACE} -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl exec ${NSE} -n ${NAMESPACE} -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestKernel2Memif() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2Memif")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f ../namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/${NAMESPACE}/sa/default \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:${NAMESPACE} \` + "\n" + `-selector k8s:sa:default`)
	r.Run(`NODE=($(kubectl get nodes -o go-template='{{range .items}}{{ if not .spec.taints  }}{{index .metadata.labels "kubernetes.io/hostname"}} {{end}}{{end}}')[0])`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- ../../../apps/nsc-kernel` + "\n" + `- ../../../apps/nse-memif` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nsc.yaml` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nsc.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nsc` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nsc` + "\n" + `          env:` + "\n" + `            - name: NSM_NETWORK_SERVICES` + "\n" + `              value: kernel://icmp-responder/nsm-1` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODE}` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nse` + "\n" + `          env:` + "\n" + `            - name: NSE_CIDR_PREFIX` + "\n" + `              value: 172.16.1.100/31` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODE}` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc -n ${NAMESPACE}`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse -n ${NAMESPACE}`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl exec ${NSC} -n ${NAMESPACE} -- ping -c 4 172.16.1.100`)
	r.Run(`result=$(kubectl exec "${NSE}" -n "${NAMESPACE}" -- vppctl ping 172.16.1.101 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}
func (s *Suite) TestKernel2Vxlan2Kernel() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2Vxlan2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f ../namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/${NAMESPACE}/sa/default \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:${NAMESPACE} \` + "\n" + `-selector k8s:sa:default`)
	r.Run(`NODES=($(kubectl get nodes -o go-template='{{range .items}}{{ if not .spec.taints  }}{{index .metadata.labels "kubernetes.io/hostname"}} {{end}}{{end}}'))`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- ../../../apps/nsc-kernel` + "\n" + `- ../../../apps/nse-kernel` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nsc.yaml` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nsc.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nsc` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nsc` + "\n" + `          env:` + "\n" + `            - name: NSM_NETWORK_SERVICES` + "\n" + `              value: kernel://icmp-responder/nsm-1` + "\n" + `` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODES[0]}` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nse` + "\n" + `          env:` + "\n" + `            - name: NSE_CIDR_PREFIX` + "\n" + `              value: 172.16.1.100/31` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODES[1]}` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc -n ${NAMESPACE}`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse -n ${NAMESPACE}`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl exec ${NSC} -n ${NAMESPACE} -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl exec ${NSE} -n ${NAMESPACE} -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestKernel2Vxlan2Memif() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2Vxlan2Memif")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f ../namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/${NAMESPACE}/sa/default \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:${NAMESPACE} \` + "\n" + `-selector k8s:sa:default`)
	r.Run(`NODES=($(kubectl get nodes -o go-template='{{range .items}}{{ if not .spec.taints  }}{{index .metadata.labels "kubernetes.io/hostname"}} {{end}}{{end}}'))`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- ../../../apps/nsc-kernel` + "\n" + `- ../../../apps/nse-memif` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nsc.yaml` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nsc.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nsc` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nsc` + "\n" + `          env:` + "\n" + `            - name: NSM_NETWORK_SERVICES` + "\n" + `              value: kernel://icmp-responder/nsm-1` + "\n" + `` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODES[0]}` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nse` + "\n" + `          env:` + "\n" + `            - name: NSE_CIDR_PREFIX` + "\n" + `              value: 172.16.1.100/31` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODES[1]}` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc -n ${NAMESPACE}`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse -n ${NAMESPACE}`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl exec ${NSC} -n ${NAMESPACE} -- ping -c 4 172.16.1.100`)
	r.Run(`result=$(kubectl exec "${NSE}" -n "${NAMESPACE}" -- vppctl ping 172.16.1.101 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}
func (s *Suite) TestMemif2Kernel() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Memif2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f ../namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/${NAMESPACE}/sa/default \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:${NAMESPACE} \` + "\n" + `-selector k8s:sa:default`)
	r.Run(`NODE=($(kubectl get nodes -o go-template='{{range .items}}{{ if not .spec.taints  }}{{index .metadata.labels "kubernetes.io/hostname"}} {{end}}{{end}}')[0])`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- ../../../apps/nsc-memif` + "\n" + `- ../../../apps/nse-kernel` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nsc.yaml` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nsc.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nsc` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nsc` + "\n" + `          env:` + "\n" + `            - name: NSM_NETWORK_SERVICES` + "\n" + `              value: memif://icmp-responder/nsm-1` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODE}` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nse` + "\n" + `          env:` + "\n" + `            - name: NSE_CIDR_PREFIX` + "\n" + `              value: 172.16.1.100/31` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODE}` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc -n ${NAMESPACE}`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse -n ${NAMESPACE}`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`result=$(kubectl exec "${NSC}" -n "${NAMESPACE}" -- vppctl ping 172.16.1.100 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`kubectl exec ${NSE} -n ${NAMESPACE} -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestMemif2Memif() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Memif2Memif")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f ../namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/${NAMESPACE}/sa/default \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:${NAMESPACE} \` + "\n" + `-selector k8s:sa:default`)
	r.Run(`NODE=($(kubectl get nodes -o go-template='{{range .items}}{{ if not .spec.taints  }}{{index .metadata.labels "kubernetes.io/hostname"}} {{end}}{{end}}')[0])`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- ../../../apps/nsc-memif` + "\n" + `- ../../../apps/nse-memif` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nsc.yaml` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nsc.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nsc` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nsc` + "\n" + `          env:` + "\n" + `            - name: NSM_NETWORK_SERVICES` + "\n" + `              value: memif://icmp-responder/nsm-1` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODE}` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nse` + "\n" + `          env:` + "\n" + `            - name: NSE_CIDR_PREFIX` + "\n" + `              value: 172.16.1.100/31` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODE}` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc -n ${NAMESPACE}`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse -n ${NAMESPACE}`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`result=$(kubectl exec "${NSC}" -n "${NAMESPACE}" -- vppctl ping 172.16.1.100 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`result=$(kubectl exec "${NSE}" -n "${NAMESPACE}" -- vppctl ping 172.16.1.101 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}
func (s *Suite) TestMemif2Vxlan2Kernel() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Memif2Vxlan2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f ../namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/${NAMESPACE}/sa/default \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:${NAMESPACE} \` + "\n" + `-selector k8s:sa:default`)
	r.Run(`NODES=($(kubectl get nodes -o go-template='{{range .items}}{{ if not .spec.taints  }}{{index .metadata.labels "kubernetes.io/hostname"}} {{end}}{{end}}'))`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- ../../../apps/nsc-memif` + "\n" + `- ../../../apps/nse-kernel` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nsc.yaml` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nsc.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nsc` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nsc` + "\n" + `          env:` + "\n" + `            - name: NSM_NETWORK_SERVICES` + "\n" + `              value: memif://icmp-responder/nsm-1` + "\n" + `` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODES[0]}` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nse` + "\n" + `          env:` + "\n" + `            - name: NSE_CIDR_PREFIX` + "\n" + `              value: 172.16.1.100/31` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODES[1]}` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc -n ${NAMESPACE}`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse -n ${NAMESPACE}`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`result=$(kubectl exec "${NSC}" -n "${NAMESPACE}" -- vppctl ping 172.16.1.100 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`kubectl exec ${NSE} -n ${NAMESPACE} -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestMemif2Vxlan2Memif() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Memif2Vxlan2Memif")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f ../namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/${NAMESPACE}/sa/default \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:${NAMESPACE} \` + "\n" + `-selector k8s:sa:default`)
	r.Run(`NODES=($(kubectl get nodes -o go-template='{{range .items}}{{ if not .spec.taints  }}{{index .metadata.labels "kubernetes.io/hostname"}} {{end}}{{end}}'))`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- ../../../apps/nsc-memif` + "\n" + `- ../../../apps/nse-memif` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nsc.yaml` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nsc.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nsc` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nsc` + "\n" + `          env:` + "\n" + `            - name: NSM_NETWORK_SERVICES` + "\n" + `              value: memif://icmp-responder/nsm-1` + "\n" + `` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODES[0]}` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nse` + "\n" + `          env:` + "\n" + `            - name: NSE_CIDR_PREFIX` + "\n" + `              value: 172.16.1.100/31` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODES[1]}` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc -n ${NAMESPACE}`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse -n ${NAMESPACE}`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`result=$(kubectl exec "${NSC}" -n "${NAMESPACE}" -- vppctl ping 172.16.1.100 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`result=$(kubectl exec "${NSE}" -n "${NAMESPACE}" -- vppctl ping 172.16.1.101 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}
