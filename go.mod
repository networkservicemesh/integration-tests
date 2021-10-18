module github.com/networkservicemesh/integration-tests

go 1.16

require (
	github.com/google/uuid v1.2.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/networkservicemesh/gotestmd v0.0.0-20210616071812-739f61445c2f
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/sys v0.0.0-20210616045830-e2b7044e8c71 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
	k8s.io/api v0.20.5
	k8s.io/apimachinery v0.20.5
	k8s.io/client-go v0.20.5
)

replace github.com/networkservicemesh/gotestmd => github.com/Mixaster995/gotestmd v0.0.0-20211022102757-ca2cb4b9f76d
