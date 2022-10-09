package docker

import (
	"github.com/habakke/terraform-provider-docker/internal/util"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var providerFactories = map[string]func() (*schema.Provider, error){
	"docker": func() (*schema.Provider, error) {
		return New(), nil
	},
}

func TestMain(m *testing.M) {
	util.ConfigureLogging(util.GetEnv("LOGLEVEL", "info"), true)
	if os.Getenv("TF_ACC") == "" {
		os.Exit(m.Run())
	}
	resource.TestMain(m)
}

func init() {
}

//nolint:deadcode,unused
func testDockerPreCheck(t *testing.T, resourceName string) {
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("provider internal validation failed: %v", err)
	}
}
