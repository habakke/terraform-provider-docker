package docker

import (
	"context"
	"github.com/habakke/terraform-provider-docker/internal/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type providerConfiguration struct {
}

// Provider represents a terraform provider definition
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"log_caller": {
				Type:        schema.TypeBool,
				Optional:    true,
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

	conf := providerConfiguration{}
	return conf, diags
}
