package main

import (
	"context"
	"fmt"
	"github.com/habakke/terraform-provider-docker/internal/docker"
	"github.com/habakke/terraform-provider-docker/internal/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"os"
)

var (
	version   string // build version number
	commit    string // sha1 revision used to build the program
	buildTime string // when the executable was built
	buildBy   string
)

func getVersionString(name string) string {
	return fmt.Sprintf("%s %s (%s at %s by %s)", name, version, commit, buildTime, buildBy)
}

func main() {
	ctx := context.Background()
	logger := util.NewTerraformLogger()
	path, err := os.Getwd()
	if err != nil {
		logger.Errorf(ctx, "failed to initialize provider: %s", err.Error())
	}

	logger.Infof(ctx, "%s", getVersionString("terraform-docker-provider"))
	logger.Infof(ctx, "%s", path)
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return docker.Provider()
		},
	})
}
