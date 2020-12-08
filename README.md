# integration-tests
integration-tests provides common integration tests for NSM integration repositories. 


## Project struecture

`suites` provides suite based tests that are reusing NSM infrastructure and test some complex scenarious via k8s go api client
`examples` provides suite based tests that are using only shell executor. This package should shouw for the users examples of using.
`tools` provides testing tools for `suites`.

## Versioning

All tests running with deployments from deployment repositories. Currently supported only https://github.com/networkservicemesh/deployments-k8s

Before test running will be checked deployments repositories in $GOPATH/src/github.com/networkservicemesh. If the repositories don't exist that it will be cloned and checkout automatically.