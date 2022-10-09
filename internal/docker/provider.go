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
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOCKER_REPOSITORY", nil),
				Description: "docker registry",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOCKER_USERNAME", nil),
				Description: "docker username",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOCKER_PASSWORD", nil),
				Description: "docker password",
			},
			"log_caller": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOCKER_LOG_CALLER", nil),
				Description: "Include calling function in log entries",
				Default:     false,
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

	registry := util.ResourceToString(d, "registry")
	username := util.ResourceToString(d, "username")
	password := util.ResourceToString(d, "password")
	conf := providerConfiguration{
		Registry: registry,
		Username: username,
		Password: password,
	}
	return conf, diags
}
