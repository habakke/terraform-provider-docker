package docker

import (
	"context"
	"fmt"
	"github.com/habakke/terraform-provider-docker/internal/registry"
	"github.com/habakke/terraform-provider-docker/internal/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Optional: false,
				Computed: true,
			},
		},
	}
}

func dataSourceDockerRegistryImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conf := meta.(providerConfiguration)

	name := util.ResourceToString(d, "name")
	tag := util.ResourceToString(d, "tag")

	url := fmt.Sprintf("https://%s", conf.Registry)
	r, err := registry.New(ctx, url, conf.Username, conf.Password, util.NewTerraformLogger())
	if err != nil {
		return diag.FromErr(err)
	}

	digest, err := r.ManifestDigest(name, tag)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s:%s", conf.Registry, name, tag))
	_ = d.Set("digest", digest.String())
	return diags
}
