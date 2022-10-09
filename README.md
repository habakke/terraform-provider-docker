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

## How to build
To build an test the plugin locally first create a `~/.terraformrc` file

```shell
provider_installation {

  dev_overrides {
    "habakke/docker" = "/Users/habakke/.terraform.d/plugins"
  }
  direct {}
}
```

Then build and install the plugin locally using

```shell
make install
```

## Running tests
To run the internal unit tests run test `test` make target

```shell
make test
```

To run terraform acceptance tests, the `TF_ACC` env variable must be set to true before making the
`test` make target, or the `testacc` make target can be used

```shell
make testacc
```

## TODO
