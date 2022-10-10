package docker

import (
	"github.com/habakke/terraform-provider-docker/internal/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestDockerRegistryImage_DataSource_Basic(t *testing.T) {
	resourceName := "docker_registry_image.ubuntu"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testDockerPreCheck(t, resourceName) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: util.LoadTestTemplateConfig(t, "docker_registry_image_basic.tf", map[string]string{
					"name": os.Getenv("DOCKER_IMAGE"),
					"tag":  os.Getenv("DOCKER_IMAGE_TAG"),
				}),
				Check: resource.ComposeTestCheckFunc(
				//resource.TestCheckResourceAttrSet(resourceName, "digest"),
				),
			},
		},
	})
}
