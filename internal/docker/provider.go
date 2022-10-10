package docker

import (
	"context"
	"github.com/habakke/terraform-provider-docker/internal/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type providerConfiguration struct {
	Registry string
	Username string
	Password string
}

// Provider represents a terraform provider definition
func Provider() *schema.Provider {
	return New()
}

func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"registry": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOCKER_REGISTRY", "index.docker.io"),
				Description: "docker registry host",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOCKER_USERNAME", nil),
				Description: "docker repository username",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOCKER_PASSWORD", nil),
				Description: "docker repository password",
			},
			"log_caller": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOCKER_LOG_CALLER", false),
				Description: "include calling function in log entries",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"docker_registry_image": dataSourceDockerRegistryImage(),
		},
		ResourcesMap:         map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// configure logging
	logCaller := util.ResourceToBool(d, "log_caller")
	util.ConfigureTerraformProviderLogging(logCaller)

	conf := providerConfiguration{
		Registry: util.ResourceToString(d, "registry"),
		Username: util.ResourceToString(d, "username"),
		Password: util.ResourceToString(d, "password"),
	}
	return conf, diags
}
