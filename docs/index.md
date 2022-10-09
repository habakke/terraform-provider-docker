# Docker provider for terraform
This document describes how to build and use the docker terraform provider.

This provider lets you fetch information about docker images from a docker registry without
requiring a local docker daemon

## Example usage

```terraform
terraform {
  required_providers {
    docker = "~> 1.0"
  }
}

provider "docker" {
  registry   = "registry.docker.com"
  username   = ""
  password   = ""
  log_caller = false
}

data "docker_registry_image" "debian" {
  name = "debian"
  tag  = "latest"
}
```
