provider "docker" {
}

data "docker_registry_image" "ubuntu" {
  name     = "{{.name}}"
  tag      = "{{.tag}}"
}
