package images_test

import (
	"fmt"
	"strings"

	"github.com/networkservicemesh/integration-tests/extensions/prefetch/images"
)

var anyFileMatch = func(string) bool { return true }
var yamlFileMatch = func(s string) bool { return strings.HasSuffix(s, ".yaml") }

func Example_GithubRawContentsURL() {
	const sampleURL = "https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/v1.0.0/apps/nsc-kernel/nsc.yaml"

	var list = images.ReteriveList([]string{sampleURL}, anyFileMatch)

	fmt.Println(list.Images[0])
	// Output:
	// ghcr.io/networkservicemesh/cmd-nsc:v1.0.0
}

func Example_GithubAPIFileURL() {
	const sampleURL = "https://api.github.com/repos/networkservicemesh/deployments-k8s/contents/apps/nsc-kernel/nsc.yaml?ref=v1.0.0"

	var list = images.ReteriveList([]string{sampleURL}, anyFileMatch)

	fmt.Println(list.Images[0])
	// Output:
	// ghcr.io/networkservicemesh/cmd-nsc:v1.0.0
}

func Example_GithubAPIDirectoryURL() {
	const sampleURL = "https://api.github.com/repos/networkservicemesh/deployments-k8s/contents/apps/nsc-kernel?ref=v1.0.0"

	var list = images.ReteriveList([]string{sampleURL}, anyFileMatch)

	fmt.Println(list.Images[0])
	// Output:
	// ghcr.io/networkservicemesh/cmd-nsc:v1.0.0
}

func Example_LocalDirectory1() {
	var list = images.ReteriveList([]string{"file://samples"}, yamlFileMatch)
	fmt.Println(list.Images[0])
	fmt.Println(list.Images[1])
	fmt.Println(list.Images[2])
	// Output:
	// alpine
	// image1
	// image2
}

func Example_LocalDirectory2() {
	var list = images.ReteriveList([]string{"file://./"}, yamlFileMatch)
	fmt.Println(list.Images[0])
	fmt.Println(list.Images[1])
	fmt.Println(list.Images[2])
	// Output:
	// alpine
	// image1
	// image2
}

func Example_LocalFile1() {
	var list = images.ReteriveList([]string{"file://samples/alpine.yaml"}, yamlFileMatch)
	fmt.Println(list.Images[0])
	// Output:
	// alpine
}

func Example_Prefetch() {
	var list = images.ReteriveList([]string{"file://samples/prefetch.yaml"}, yamlFileMatch)
	fmt.Println(list.Images[0])
	fmt.Println(list.Images[1])
	// Output:
	// image1
	// image2
}
