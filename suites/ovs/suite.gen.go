// Code generated by gotestmd DO NOT EDIT.
package ovs

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
	r := s.Runner("../deployments-k8s/examples/ovs")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete mutatingwebhookconfiguration --all` + "\n" + `kubectl delete ns nsm-system`)
	})
	r.Run(`kubectl create ns nsm-system`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/ovs?ref=ae5fa6673066122db59787885e09954dbc5147ff`)
	r.Run(`WH=$(kubectl get pods -l app=admission-webhook-k8s -n nsm-system --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `kubectl wait --for=condition=ready --timeout=1m pod ${WH} -n nsm-system`)
}
func (s *Suite) TestWebhook_smartvf() {
	r := s.Runner("../deployments-k8s/examples/features/webhook-smartvf")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`WH=$(kubectl get pods -l app=admission-webhook-k8s -n nsm-system --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `kubectl wait --for=condition=ready --timeout=1m pod ${WH} -n nsm-system`)
	r.Run(`NAMESPACE=($(kubectl create -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/ae5fa6673066122db59787885e09954dbc5147ff/examples/features/namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`NODES=($(kubectl get nodes -o go-template='{{range .items}}{{ if not .spec.taints  }}{{index .metadata.labels "kubernetes.io/hostname"}} {{end}}{{end}}'))`)
	r.Run(`cat > postgres-cl.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: v1` + "\n" + `kind: Pod` + "\n" + `metadata:` + "\n" + `  name: postgres-cl` + "\n" + `  annotations:` + "\n" + `    # Add in the sriovToken label your own SmartVF interface pool` + "\n" + `    networkservicemesh.io: kernel://my-postgres-service/nsm-1?sriovToken=worker.domain/100G` + "\n" + `  labels:` + "\n" + `    app: postgres-cl` + "\n" + `    "spiffe.io/spiffe-id": "true"` + "\n" + `spec:` + "\n" + `  containers:` + "\n" + `  - name: postgres-cl` + "\n" + `    image: postgres` + "\n" + `    imagePullPolicy: IfNotPresent` + "\n" + `    env:` + "\n" + `      - name: POSTGRES_HOST_AUTH_METHOD` + "\n" + `        value: trust` + "\n" + `  nodeSelector:` + "\n" + `    kubernetes.io/hostname: ${NODES[0]}` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse-kernel` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: postgres` + "\n" + `          image: postgres` + "\n" + `          imagePullPolicy: IfNotPresent` + "\n" + `          ports:` + "\n" + `            - containerPort: 5432` + "\n" + `          env:` + "\n" + `            - name: POSTGRES_DB` + "\n" + `              value: test` + "\n" + `            - name: POSTGRES_USER` + "\n" + `              value: admin` + "\n" + `            - name: POSTGRES_PASSWORD` + "\n" + `              value: admin` + "\n" + `        - name: nse` + "\n" + `          env:` + "\n" + `            - name: NSM_LABELS` + "\n" + `              # Add your own serviceDomain` + "\n" + `              value: serviceDomain:worker.domain` + "\n" + `            - name: NSM_SERVICE_NAMES` + "\n" + `              value: my-postgres-service` + "\n" + `            - name: NSM_CIDR_PREFIX` + "\n" + `              value: 172.16.1.100/31` + "\n" + `          resources:` + "\n" + `            limits:` + "\n" + `              # Add your own SmartVF interface pool` + "\n" + `              worker.domain/100G: 1` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODES[1]}` + "\n" + `EOF`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- https://github.com/networkservicemesh/deployments-k8s/apps/nse-kernel?ref=ae5fa6673066122db59787885e09954dbc5147ff` + "\n" + `` + "\n" + `resources:` + "\n" + `- postgres-cl.yaml` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=5m pod -l app=nse-kernel -n ${NAMESPACE}`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod postgres-cl -n ${NAMESPACE}`)
	r.Run(`NSC=$(kubectl get pods -l app=postgres-cl -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-kernel -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl exec ${NSC} -n ${NAMESPACE} -c postgres-cl -- sh -c 'PGPASSWORD=admin psql -h 172.16.1.100 -p 5432 -U admin test'`)
}
func (s *Suite) TestKernel2Kernel() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/ae5fa6673066122db59787885e09954dbc5147ff/examples/use-cases/namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`NODE=($(kubectl get nodes -o go-template='{{range .items}}{{ if not .spec.taints  }}{{index .metadata.labels "kubernetes.io/hostname"}} {{end}}{{end}}')[0])`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `resources: ` + "\n" + `- client.yaml` + "\n" + `bases:` + "\n" + `- https://github.com/networkservicemesh/deployments-k8s/apps/nse-kernel?ref=ae5fa6673066122db59787885e09954dbc5147ff` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > client.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: v1` + "\n" + `kind: Pod` + "\n" + `metadata:` + "\n" + `  name: alpine` + "\n" + `  labels:` + "\n" + `    app: alpine    ` + "\n" + `  annotations:` + "\n" + `    networkservicemesh.io: kernel://icmp-responder/nsm-1` + "\n" + `spec:` + "\n" + `  containers:` + "\n" + `  - name: alpine` + "\n" + `    image: alpine:3.15.0` + "\n" + `    imagePullPolicy: IfNotPresent` + "\n" + `    stdin: true` + "\n" + `    tty: true` + "\n" + `  nodeSelector:` + "\n" + `    kubernetes.io/hostname: ${NODE}` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse-kernel` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nse` + "\n" + `          env:` + "\n" + `            - name: NSM_CIDR_PREFIX` + "\n" + `              value: 172.16.1.100/31` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODE}` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ${NAMESPACE}`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ${NAMESPACE}`)
	r.Run(`NSC=$(kubectl get pods -l app=alpine -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-kernel -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl exec ${NSC} -n ${NAMESPACE} -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl exec ${NSE} -n ${NAMESPACE} -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestKernel2KernelVLAN() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2KernelVLAN")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/ae5fa6673066122db59787885e09954dbc5147ff/examples/use-cases/namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`NODE=($(kubectl get nodes -o go-template='{{range .items}}{{ if not .spec.taints  }}{{index .metadata.labels "kubernetes.io/hostname"}} {{end}}{{end}}')[0])`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- https://github.com/networkservicemesh/deployments-k8s/apps/nsc-kernel?ref=ae5fa6673066122db59787885e09954dbc5147ff` + "\n" + `- https://github.com/networkservicemesh/deployments-k8s/apps/nse-vlan-vpp?ref=ae5fa6673066122db59787885e09954dbc5147ff` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nsc.yaml` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nsc.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nsc-kernel` + "\n" + `spec:` + "\n" + `  replicas: 2` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nsc` + "\n" + `          env:` + "\n" + `            - name: NSM_NETWORK_SERVICES` + "\n" + `              value: kernel://vlan-vpp-responder/nsm-1` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODE}` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse-kernel` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nse` + "\n" + `          env:` + "\n" + `            - name: NSM_CIDR_PREFIX` + "\n" + `              value: 172.16.1.100/30` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODE}` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-kernel -n ${NAMESPACE}`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ${NAMESPACE}`)
	r.Run(`NSC=$((kubectl get pods -l app=nsc-kernel -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{" "}}{{end}}') | cut -d' ' -f1)` + "\n" + `TARGET_IP=$(kubectl exec -ti ${NSC} -n ${NAMESPACE} -- ip route show | grep 172.16 | cut -d' ' -f1)`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-kernel -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl exec ${NSC} -n ${NAMESPACE} -- ping -c 4 ${TARGET_IP}`)
}
func (s *Suite) TestSmartVF2SmartVF() {
	r := s.Runner("../deployments-k8s/examples/use-cases/SmartVF2SmartVF")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/ae5fa6673066122db59787885e09954dbc5147ff/examples/use-cases/namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- https://github.com/networkservicemesh/deployments-k8s/apps/nsc-kernel?ref=ae5fa6673066122db59787885e09954dbc5147ff` + "\n" + `- https://github.com/networkservicemesh/deployments-k8s/apps/nse-kernel?ref=ae5fa6673066122db59787885e09954dbc5147ff` + "\n" + `` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nsc.yaml` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nsc.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nsc-kernel` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nsc` + "\n" + `          env:` + "\n" + `            - name: NSM_NETWORK_SERVICES` + "\n" + `              value: kernel://icmp-responder/nsm-1?sriovToken=worker.domain/100G` + "\n" + `          resources:` + "\n" + `            limits:` + "\n" + `              worker.domain/100G: 1` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse-kernel` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nse` + "\n" + `          env:` + "\n" + `            - name: NSM_LABELS` + "\n" + `              value: serviceDomain:worker.domain` + "\n" + `            - name: NSM_CIDR_PREFIX` + "\n" + `              value: 172.16.1.100/31` + "\n" + `          resources:` + "\n" + `            limits:` + "\n" + `              worker.domain/100G: 1` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl -n ${NAMESPACE} wait --for=condition=ready --timeout=1m pod -l app=nsc-kernel`)
	r.Run(`kubectl -n ${NAMESPACE} wait --for=condition=ready --timeout=1m pod -l app=nse-kernel`)
	r.Run(`NSC=$(kubectl -n ${NAMESPACE} get pods -l app=nsc-kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl -n ${NAMESPACE} exec ${NSC} -- ping -c 4 172.16.1.100`)
}
