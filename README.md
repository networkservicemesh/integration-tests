# integration-tests
integration-tests provides common integration tests for NSM integration repositories


## Implementation details

1. Each test runs in a unique namespace.
2. MSM based suite setups NSM infrastructure once per suite.
3. Spire setups once per suite.
4. Mostly each suite can provide env based config that can help to debug.
5. deployments repo will be automatically cloned if it is not present in `$GOPATH/src/github.com/netwrokservicemesh/deployments-k8s`
6. Steps from each test can be run in the terminal.
7. Each test should use https://kustomize.io/ based configuration.

## How to run tests?

1. Setup kubernetes
2. Run `go test - run TestEntryPoint`
