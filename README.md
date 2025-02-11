# integration-tests
integration-tests provides common integration tests for NSM integration repositories. 
The tests are a result of generating from examples from deployments repositories. Currently, we are using only https://github.com/networkservicemesh/deployments-k8s.

## How re-generate tests manually?

### Prerequisite

Install gotestmd

``` bash
go install github.com/networkservicemesh/gotestmd@main
```

Install goimports

```bash
go install golang.org/x/tools/cmd/goimports@latest
```

```
go generate ./...
```
