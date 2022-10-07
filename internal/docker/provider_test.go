package docker

import (
	"github.com/habakke/terraform-provider-docker/internal/util"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testDockerProviders map[string]*schema.Provider
var testDockerProvider *schema.Provider

func TestMain(m *testing.M) {
	util.ConfigureLogging(util.GetEnv("LOGLEVEL", "info"), true)
	code := m.Run()
	os.Exit(code)
}

func init() {

	testDockerProvider = Provider()
	testDockerProviders = map[string]*schema.Provider{
		"docker": testDockerProvider,
	}
}

//nolint:deadcode,unused
func testDockerPreCheck(t *testing.T, resourceName string) {
	r := testDockerProvider.ResourcesMap[strings.Split(resourceName, ".")[0]]
	d := testDockerProvider.DataSourcesMap[strings.Split(resourceName, ".")[0]]

	if r != nil && d != nil {
		t.Fatalf("missing resource '%s'", resourceName)
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("provider internal validation failed: %v", err)
	}
}
