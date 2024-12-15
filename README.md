# Terraform Provider OpenWebUI

This Terraform provider allows you to manage resources in an OpenWebUI instance.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
```shell
git clone https://github.com/ncecere/terraform-provider-openwebui
```

2. Enter the repository directory
```shell
cd terraform-provider-openwebui
```

3. Build the provider
```shell
make build
```

## Installing the Provider

To use the provider in your Terraform configuration, add the following terraform block:

```hcl
terraform {
  required_providers {
    openwebui = {
      source = "ncecere/openwebui"
    }
  }
}
```

## Using the Provider

To configure the provider:

```hcl
provider "openwebui" {
  endpoint = "http://localhost:8080"  # OpenWebUI API endpoint
  # token = "your-api-token"         # Token can be provided here or via OPENWEBUI_TOKEN env var
}
```

### Managing Knowledge Bases

Create a knowledge base:

```hcl
resource "openwebui_knowledge" "example" {
  name        = "Example Knowledge Base"
  description = "This is an example knowledge base"
  
  # Optional: Additional data
  data = {
    key = "value"
  }
  
  # Optional: Access control settings
  access_control = {
    visibility = "private"
  }
}
```

## Environment Variables

The following environment variables can be used to configure the provider:

- `OPENWEBUI_ENDPOINT`: The endpoint URL of the OpenWebUI instance
- `OPENWEBUI_TOKEN`: The API token to authenticate with OpenWebUI

## Development

### Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

### Building

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the `make build` command

```shell
make build
```

### Testing

To run the tests:

```shell
make test
```

### Documentation

Documentation is generated from the provider schema. The schema is defined in the provider code and is used to generate both the documentation and the provider's configuration interface.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This provider is licensed under the MIT License. See the LICENSE file for details.
