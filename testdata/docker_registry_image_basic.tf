provider "docker" {
  registry = "{{.registry}}"
  username = "{{.username}}"
  password = "{{.password}}"
}

data "docker_registry_image" "ubuntu" {
  name     = "{{.name}}"
  tag      = "{{.tag}}"
}
