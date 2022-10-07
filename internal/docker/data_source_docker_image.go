package docker

import (
	"context"
	"github.com/habakke/terraform-provider-docker/internal/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heroku/docker-registry-client/registry"
)

func dataSourceDockerRegistryImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDockerRegistryImageRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Required: true,
			},
			"digest": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDockerRegistryImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	//conf := meta.(providerConfiguration)

	name := util.ResourceToString(d, "name")
	tag := util.ResourceToString(d, "tag")

	url := "https://registry-1.docker.io/"
	username := "" // anonymous
	password := "" // anonymous

	r, err := registry.New(url, username, password)
	if err != nil {
		return diag.FromErr(err)
	}

	digest, err := r.ManifestDigest(name, tag)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("digest", digest.Hex())
	return diags
}
