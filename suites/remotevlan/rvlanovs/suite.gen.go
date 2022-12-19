// Code generated by gotestmd DO NOT EDIT.
package rvlanovs

import (
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
	r := s.Runner("../deployments-k8s/examples/remotevlan/rvlanovs")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete -k https://github.com/networkservicemesh/deployments-k8s/examples/remotevlan/rvlanovs?ref=afd86bf3b50d9ab90e5d21b57d08a4eb0fe0d70a`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/remotevlan/rvlanovs?ref=afd86bf3b50d9ab90e5d21b57d08a4eb0fe0d70a`)
	r.Run(`kubectl -n nsm-system wait --for=condition=ready --timeout=2m pod -l app=forwarder-ovs`)
}
func (s *Suite) TestKernel2RVlanBreakout() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2RVlanBreakout")
	s.T().Cleanup(func() {
		r.Run(`docker stop rvm-tester` + "\n" + `docker image rm rvm-tester:latest` + "\n" + `true`)
		r.Run(`kubectl delete ns ns-kernel2rvlan-breakout`)
	})
	r.Run(`kubectl create ns ns-kernel2rvlan-breakout`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Kernel2RVlanBreakout?ref=afd86bf3b50d9ab90e5d21b57d08a4eb0fe0d70a`)
	r.Run(`kubectl -n ns-kernel2rvlan-breakout wait --for=condition=ready --timeout=1m pod -l app=iperf1-s`)
	r.Run(`NSCS=($(kubectl get pods -l app=iperf1-s -n ns-kernel2rvlan-breakout --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'))`)
	r.Run(`cat > Dockerfile <<EOF` + "\n" + `FROM networkstatic/iperf3` + "\n" + `` + "\n" + `RUN apt-get update \` + "\n" + `    && apt-get install -y ethtool iproute2 \` + "\n" + `    && rm -rf /var/lib/apt/lists/*` + "\n" + `` + "\n" + `ENTRYPOINT [ "tail", "-f", "/dev/null" ]` + "\n" + `EOF` + "\n" + `docker build . -t rvm-tester`)
	r.Run(`docker run --cap-add=NET_ADMIN --rm -d --network bridge-2 --name rvm-tester rvm-tester tail -f /dev/null` + "\n" + `docker exec rvm-tester ip link set eth0 down` + "\n" + `docker exec rvm-tester ip link add link eth0 name eth0.100 type vlan id 100` + "\n" + `docker exec rvm-tester ip link set eth0 up` + "\n" + `docker exec rvm-tester ip addr add 172.10.0.254/24 dev eth0.100` + "\n" + `docker exec rvm-tester ethtool -K eth0 tx off`)
	r.Run(`status=0` + "\n" + `    for nsc in "${NSCS[@]}"` + "\n" + `    do` + "\n" + `      IP_ADDRESS=$(kubectl exec ${nsc} -c cmd-nsc -n ns-kernel2rvlan-breakout -- ip -4 addr show nsm-1 | grep -oP '(?<=inet\s)\d+(\.\d+){3}')` + "\n" + `      kubectl exec ${nsc} -c iperf-server -n ns-kernel2rvlan-breakout -- iperf3 -sD -B ${IP_ADDRESS} -1` + "\n" + `      docker exec rvm-tester iperf3 -i0 -t 25 -c ${IP_ADDRESS}` + "\n" + `      if test $? -ne 0` + "\n" + `      then` + "\n" + `        status=1` + "\n" + `      fi` + "\n" + `    done` + "\n" + `    if test ${status} -eq 1` + "\n" + `    then` + "\n" + `      false` + "\n" + `    fi`)
	r.Run(`status=0` + "\n" + `    for nsc in "${NSCS[@]}"` + "\n" + `    do` + "\n" + `      IP_ADDRESS=$(kubectl exec ${nsc} -c cmd-nsc -n ns-kernel2rvlan-breakout -- ip -4 addr show nsm-1 | grep -oP '(?<=inet\s)\d+(\.\d+){3}')` + "\n" + `      kubectl exec ${nsc} -c iperf-server -n ns-kernel2rvlan-breakout -- iperf3 -sD -B ${IP_ADDRESS} -1` + "\n" + `      docker exec rvm-tester iperf3 -i0 -t 5 -u -c ${IP_ADDRESS}` + "\n" + `      if test $? -ne 0` + "\n" + `      then` + "\n" + `        status=1` + "\n" + `      fi` + "\n" + `    done` + "\n" + `    if test ${status} -eq 1` + "\n" + `    then` + "\n" + `      false` + "\n" + `    fi`)
	r.Run(`status=0` + "\n" + `    for nsc in "${NSCS[@]}"` + "\n" + `    do` + "\n" + `      docker exec rvm-tester iperf3 -sD -B 172.10.0.254 -1` + "\n" + `      kubectl exec ${nsc} -c iperf-server -n ns-kernel2rvlan-breakout -- iperf3 -i0 -t 5 -c 172.10.0.254` + "\n" + `      if test $? -ne 0` + "\n" + `      then` + "\n" + `        status=1` + "\n" + `      fi` + "\n" + `    done` + "\n" + `    if test ${status} -eq 1` + "\n" + `    then` + "\n" + `      false` + "\n" + `    fi`)
	r.Run(`status=0` + "\n" + `    for nsc in "${NSCS[@]}"` + "\n" + `    do` + "\n" + `      docker exec rvm-tester iperf3 -sD -B 172.10.0.254 -1` + "\n" + `      kubectl exec ${NSCS[1]} -c iperf-server -n ns-kernel2rvlan-breakout -- iperf3 -i0 -t 5 -u -c 172.10.0.254` + "\n" + `      if test $? -ne 0` + "\n" + `      then` + "\n" + `        status=1` + "\n" + `      fi` + "\n" + `    done` + "\n" + `    if test ${status} -eq 1` + "\n" + `    then` + "\n" + `      false` + "\n" + `    fi`)
}
func (s *Suite) TestKernel2RVlanInternal() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2RVlanInternal")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-kernel2rvlan-internal`)
	})
	r.Run(`kubectl create ns ns-kernel2rvlan-internal`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Kernel2RVlanInternal?ref=afd86bf3b50d9ab90e5d21b57d08a4eb0fe0d70a`)
	r.Run(`kubectl -n ns-kernel2rvlan-internal wait --for=condition=ready --timeout=1m pod -l app=iperf1-s`)
	r.Run(`NSCS=($(kubectl get pods -l app=iperf1-s -n ns-kernel2rvlan-internal --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'))`)
	r.Run(`IP_ADDR=$(kubectl exec ${NSCS[0]} -c cmd-nsc -n ns-kernel2rvlan-internal -- ip -4 addr show nsm-1 | grep -oP '(?<=inet\s)\d+(\.\d+){3}')` + "\n" + `    kubectl exec ${NSCS[0]} -c iperf-server -n ns-kernel2rvlan-internal -- iperf3 -sD -B ${IP_ADDR} -1` + "\n" + `    kubectl exec ${NSCS[1]} -c iperf-server -n ns-kernel2rvlan-internal -- iperf3 -i0 -t 5 -c ${IP_ADDR}`)
	r.Run(`IP_ADDR=$(kubectl exec ${NSCS[1]} -c cmd-nsc -n ns-kernel2rvlan-internal -- ip -4 addr show nsm-1 | grep -oP '(?<=inet\s)\d+(\.\d+){3}')` + "\n" + `    kubectl exec ${NSCS[1]} -c iperf-server -n ns-kernel2rvlan-internal -- iperf3 -sD -B ${IP_ADDR} -1` + "\n" + `    kubectl exec ${NSCS[0]} -c iperf-server -n ns-kernel2rvlan-internal -- iperf3 -i0 -t 5 -u -c ${IP_ADDR}`)
	r.Run(`IP_ADDR=$(kubectl exec ${NSCS[0]} -c cmd-nsc -n ns-kernel2rvlan-internal -- ip -6 a s nsm-1 scope global | grep -oP '(?<=inet6\s)([0-9a-f:]+:+)+[0-9a-f]+')` + "\n" + `    kubectl exec ${NSCS[0]} -c iperf-server -n ns-kernel2rvlan-internal -- iperf3 -sD -B ${IP_ADDR} -1` + "\n" + `    kubectl exec ${NSCS[1]} -c iperf-server -n ns-kernel2rvlan-internal -- iperf3 -i0 -t 5 -6 -c ${IP_ADDR}`)
	r.Run(`IP_ADDR=$(kubectl exec ${NSCS[1]} -c cmd-nsc -n ns-kernel2rvlan-internal -- ip -6 a s nsm-1 scope global | grep -oP '(?<=inet6\s)([0-9a-f:]+:+)+[0-9a-f]+')` + "\n" + `    kubectl exec ${NSCS[1]} -c iperf-server -n ns-kernel2rvlan-internal -- iperf3 -sD -B ${IP_ADDR} -1` + "\n" + `    kubectl exec ${NSCS[0]} -c iperf-server -n ns-kernel2rvlan-internal -- iperf3 -i0 -t 5 -6 -u -c ${IP_ADDR}`)
}
func (s *Suite) TestKernel2RVlanMultiNS() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2RVlanMultiNS")
	s.T().Cleanup(func() {
		r.Run(`docker stop rvm-tester && \` + "\n" + `docker image rm rvm-tester:latest` + "\n" + `true`)
		r.Run(`kubectl delete --namespace=nsm-system -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/afd86bf3b50d9ab90e5d21b57d08a4eb0fe0d70a/examples/use-cases/Kernel2RVlanMultiNS/client.yaml`)
		r.Run(`kubectl delete ns ns-kernel2vlan-multins-1`)
		r.Run(`kubectl delete ns ns-kernel2vlan-multins-2`)
	})
	r.Run(`kubectl create ns ns-kernel2vlan-multins-1` + "\n" + `kubectl create ns ns-kernel2vlan-multins-2`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Kernel2RVlanMultiNS/ns-1?ref=afd86bf3b50d9ab90e5d21b57d08a4eb0fe0d70a`)
	r.Run(`kubectl apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/afd86bf3b50d9ab90e5d21b57d08a4eb0fe0d70a/examples/use-cases/Kernel2RVlanMultiNS/ns-2/netsvc.yaml` + "\n" + `kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Kernel2RVlanMultiNS/ns-2?ref=afd86bf3b50d9ab90e5d21b57d08a4eb0fe0d70a`)
	r.Run(`kubectl apply -n nsm-system -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/afd86bf3b50d9ab90e5d21b57d08a4eb0fe0d70a/examples/use-cases/Kernel2RVlanMultiNS/client.yaml`)
	r.Run(`kubectl -n ns-kernel2vlan-multins-1 wait --for=condition=ready --timeout=1m pod -l app=nse-remote-vlan`)
	r.Run(`kubectl -n ns-kernel2vlan-multins-1 wait --for=condition=ready --timeout=1m pod -l app=alpine-1`)
	r.Run(`kubectl -n ns-kernel2vlan-multins-2 wait --for=condition=ready --timeout=1m pod -l app=nse-remote-vlan`)
	r.Run(`kubectl -n ns-kernel2vlan-multins-2 wait --for=condition=ready --timeout=1m pod -l app=alpine-2`)
	r.Run(`kubectl -n ns-kernel2vlan-multins-2 wait --for=condition=ready --timeout=1m pod -l app=alpine-3`)
	r.Run(`kubectl -n nsm-system wait --for=condition=ready --timeout=1m pod -l app=alpine-4`)
	r.Run(`cat > Dockerfile <<EOF` + "\n" + `FROM alpine:3.15.0` + "\n" + `` + "\n" + `RUN apk add ethtool tcpdump iproute2` + "\n" + `` + "\n" + `ENTRYPOINT [ "tail", "-f", "/dev/null" ]` + "\n" + `EOF` + "\n" + `docker build . -t rvm-tester`)
	r.Run(`docker run --cap-add=NET_ADMIN --rm -d --network bridge-2 --name rvm-tester rvm-tester tail -f /dev/null` + "\n" + `docker exec rvm-tester ip link set eth0 down` + "\n" + `docker exec rvm-tester ip link add link eth0 name eth0.100 type vlan id 100` + "\n" + `docker exec rvm-tester ip link add link eth0 name eth0.300 type vlan id 300` + "\n" + `docker exec rvm-tester ip link set eth0 up` + "\n" + `docker exec rvm-tester ip addr add 172.10.0.254/24 dev eth0.100` + "\n" + `docker exec rvm-tester ip addr add 172.10.1.254/24 dev eth0` + "\n" + `docker exec rvm-tester ip addr add 172.10.2.254/24 dev eth0.300` + "\n" + `docker exec rvm-tester ethtool -K eth0 tx off`)
	r.Run(`NSCS=($(kubectl get pods -l app=alpine-1 -n ns-kernel2vlan-multins-1 --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'))`)
	r.Run(`status=0` + "\n" + `LINK_MTU=$(docker exec kind-worker cat /sys/class/net/ext_net1/mtu)` + "\n" + `for nsc in "${NSCS[@]}"` + "\n" + `do` + "\n" + `  MTU=$(kubectl exec ${nsc} -c cmd-nsc -n ns-kernel2vlan-multins-1 -- cat /sys/class/net/nsm-1/mtu)` + "\n" + `` + "\n" + `  echo "$LINK_MTU vs $MTU"` + "\n" + `` + "\n" + `  if test "${MTU}" = ""` + "\n" + `    then` + "\n" + `      status=1` + "\n" + `  fi` + "\n" + `  if test $MTU -ne $LINK_MTU` + "\n" + `    then` + "\n" + `      status=2` + "\n" + `  fi` + "\n" + `done` + "\n" + `if test ${status} -ne 0` + "\n" + `  then` + "\n" + `    false` + "\n" + `fi`)
	r.Run(`declare -A IP_ADDR` + "\n" + `for nsc in "${NSCS[@]}"` + "\n" + `do` + "\n" + `  IP_ADDR[$nsc]=$(kubectl exec ${nsc} -n ns-kernel2vlan-multins-1 -c alpine -- ip -4 addr show nsm-1 | grep -oP '(?<=inet\s)\d+(\.\d+){3}')` + "\n" + `done`)
	r.Run(`status=0` + "\n" + `for nsc in "${NSCS[@]}"` + "\n" + `do` + "\n" + `  for vlan_if_name in eth0.100 eth0.300` + "\n" + `  do` + "\n" + `    docker exec rvm-tester ping -w 1 -c 1 ${IP_ADDR[$nsc]} -I ${vlan_if_name}` + "\n" + `    if test $? -eq 0` + "\n" + `      then` + "\n" + `        status=2` + "\n" + `    fi` + "\n" + `  done` + "\n" + `  docker exec rvm-tester ping -c 1 ${IP_ADDR[$nsc]} -I eth0` + "\n" + `  if test $? -ne 0` + "\n" + `    then` + "\n" + `      status=1` + "\n" + `  fi` + "\n" + `done` + "\n" + `if test ${status} -eq 1` + "\n" + `  then` + "\n" + `    false` + "\n" + `fi`)
	r.Run(`NSCS_BLUE=($(kubectl get pods -l app=alpine-2 -n ns-kernel2vlan-multins-2 --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'))` + "\n" + `NSCS_GREEN=($(kubectl get pods -l app=alpine-3 -n ns-kernel2vlan-multins-2 --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'))`)
	r.Run(`status=0` + "\n" + `for nsc in "${NSCS_BLUE[@]} ${NSCS_GREEN[@]}"` + "\n" + `do` + "\n" + `  MTU=$(kubectl exec ${nsc} -c cmd-nsc -n ns-kernel2vlan-multins-2 -- cat /sys/class/net/nsm-1/mtu)` + "\n" + `` + "\n" + `  echo "$LINK_MTU vs $MTU"` + "\n" + `` + "\n" + `  if test "${MTU}" = ""` + "\n" + `    then` + "\n" + `      status=1` + "\n" + `  fi` + "\n" + `  if test $MTU -ne $LINK_MTU` + "\n" + `    then` + "\n" + `      status=2` + "\n" + `  fi` + "\n" + `done` + "\n" + `if test ${status} -ne 0` + "\n" + `  then` + "\n" + `    false` + "\n" + `fi`)
	r.Run(`declare -A IP_ADDR_BLUE` + "\n" + `for nsc in "${NSCS_BLUE[@]}"` + "\n" + `do` + "\n" + `  IP_ADDR_BLUE[$nsc]=$(kubectl exec ${nsc} -n ns-kernel2vlan-multins-2 -c alpine -- ip -4 addr show nsm-1 | grep -oP '(?<=inet\s)\d+(\.\d+){3}')` + "\n" + `done` + "\n" + `declare -A IP_ADDR_GREEN` + "\n" + `for nsc in "${NSCS_GREEN[@]}"` + "\n" + `do` + "\n" + `  IP_ADDR_GREEN[$nsc]=$(kubectl exec ${nsc} -n ns-kernel2vlan-multins-2 -c alpine -- ip -4 addr show nsm-1 | grep -oP '(?<=inet\s)\d+(\.\d+){3}')` + "\n" + `done`)
	r.Run(`status=0` + "\n" + `for nsc in "${NSCS_BLUE[@]}"` + "\n" + `do` + "\n" + `  for vlan_if_name in eth0.100 eth0` + "\n" + `  do` + "\n" + `    docker exec rvm-tester ping -w 1 -c 1 ${IP_ADDR_BLUE[$nsc]} -I ${vlan_if_name}` + "\n" + `    if test $? -eq 0` + "\n" + `      then` + "\n" + `        status=2` + "\n" + `    fi` + "\n" + `  done` + "\n" + `  docker exec rvm-tester ping -c 1 ${IP_ADDR_BLUE[$nsc]} -I eth0.300` + "\n" + `  if test $? -ne 0` + "\n" + `    then` + "\n" + `      status=1` + "\n" + `  fi` + "\n" + `done` + "\n" + `for nsc in "${NSCS_GREEN[@]}"` + "\n" + `do` + "\n" + `  for vlan_if_name in eth0.100 eth0` + "\n" + `  do` + "\n" + `    docker exec rvm-tester ping -w 1 -c 1 ${IP_ADDR_GREEN[$nsc]} -I ${vlan_if_name}` + "\n" + `    if test $? -eq 0` + "\n" + `      then` + "\n" + `        status=2` + "\n" + `    fi` + "\n" + `  done` + "\n" + `  docker exec rvm-tester ping -c 1 ${IP_ADDR_GREEN[$nsc]} -I eth0.300` + "\n" + `  if test $? -ne 0` + "\n" + `    then` + "\n" + `      status=1` + "\n" + `  fi` + "\n" + `done` + "\n" + `if test ${status} -eq 1` + "\n" + `  then` + "\n" + `    false` + "\n" + `fi`)
	r.Run(`kubectl delete deployment alpine-2-bg -n ns-kernel2vlan-multins-2`)
	r.Run(`status=0` + "\n" + `for nsc in "${NSCS_GREEN[@]}"` + "\n" + `do` + "\n" + `  docker exec rvm-tester ping -c 1 ${IP_ADDR_GREEN[$nsc]} -I eth0.300` + "\n" + `  if test $? -ne 0` + "\n" + `    then` + "\n" + `      status=1` + "\n" + `  fi` + "\n" + `done` + "\n" + `if test ${status} -eq 1` + "\n" + `  then` + "\n" + `    false` + "\n" + `fi`)
	r.Run(`NSCS=($(kubectl get pods -l app=alpine-4 -n nsm-system --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'))`)
	r.Run(`status=0` + "\n" + `for nsc in "${NSCS[@]}"` + "\n" + `do` + "\n" + `  MTU=$(kubectl exec ${nsc} -c cmd-nsc -n nsm-system -- cat /sys/class/net/nsm-1/mtu)` + "\n" + `` + "\n" + `  echo "$LINK_MTU vs $MTU"` + "\n" + `` + "\n" + `  if test "${MTU}" = ""` + "\n" + `    then` + "\n" + `      status=1` + "\n" + `  fi` + "\n" + `  if test $MTU -ne $LINK_MTU` + "\n" + `    then` + "\n" + `      status=2` + "\n" + `  fi` + "\n" + `done` + "\n" + `if test ${status} -ne 0` + "\n" + `  then` + "\n" + `    false` + "\n" + `fi`)
	r.Run(`declare -A IP_ADDR` + "\n" + `for nsc in "${NSCS[@]}"` + "\n" + `do` + "\n" + `  IP_ADDR[$nsc]=$(kubectl exec ${nsc} -n nsm-system -c alpine -- ip -4 addr show nsm-1 | grep -oP '(?<=inet\s)\d+(\.\d+){3}')` + "\n" + `done`)
	r.Run(`status=0` + "\n" + `for nsc in "${NSCS[@]}"` + "\n" + `do` + "\n" + `  for vlan_if_name in eth0 eth0.300` + "\n" + `  do` + "\n" + `    docker exec rvm-tester ping -w 1 -c 1 ${IP_ADDR[$nsc]} -I ${vlan_if_name}` + "\n" + `    if test $? -eq 0` + "\n" + `      then` + "\n" + `        status=2` + "\n" + `    fi` + "\n" + `  done` + "\n" + `  docker exec rvm-tester ping -c 1 ${IP_ADDR[$nsc]} -I eth0.100` + "\n" + `  if test $? -ne 0` + "\n" + `    then` + "\n" + `      status=1` + "\n" + `  fi` + "\n" + `done` + "\n" + `if test ${status} -eq 1` + "\n" + `  then` + "\n" + `    false` + "\n" + `fi`)
}
