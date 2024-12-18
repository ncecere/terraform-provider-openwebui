# Terraform Provider OpenWebUI

This Terraform provider allows you to manage OpenWebUI resources through Terraform.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
```shell
git clone git@github.com:ncecere/terraform-provider-openwebui.git
```

2. Enter the repository directory
```shell
cd terraform-provider-openwebui
```

3. Build the provider
```shell
make build
```

## Installing The Provider

To install the provider for local development:

```shell
make install
```

This will build and install the provider into your `~/.terraform.d/plugins` directory.

## Using the Provider

To use the provider, you'll need:
1. An OpenWebUI instance
2. An API token for authentication

Configure the provider in your Terraform configuration:

```hcl
terraform {
  required_providers {
    openwebui = {
      source = "ncecere/openwebui"
    }
  }
}

provider "openwebui" {
  endpoint = "http://your-openwebui-instance"  # Or use OPENWEBUI_ENDPOINT env var
  token    = "your-api-token"                  # Or use OPENWEBUI_TOKEN env var
}
```

### Example Usage

Creating a knowledge base:

```hcl
resource "openwebui_knowledge" "example" {
  name        = "Example Knowledge Base"
  description = "This is an example knowledge base"
  
  data = {
    source = "terraform"
    type   = "example"
  }
  
  access_control = "public"  # or "private"
}
```

Looking up an existing knowledge base:

```hcl
data "openwebui_knowledge" "lookup" {
  name = "Example Knowledge Base"
}
```

## Project Structure

The provider is organized into several key components:

```
terraform-provider-openwebui/
├── docs/                    # Provider documentation
├── examples/               # Example configurations
├── internal/
│   └── provider/
│       ├── client/         # API client implementations
│       │   ├── knowledge/  # Knowledge-specific client
│       │   └── ...        # Other resource clients
│       └── ...            # Provider and resource implementations
└── local_testing/         # Local development test configurations
```

## Development

### Adding a New Resource

1. Create a new client package in `internal/provider/client/`
2. Implement the resource client interface
3. Add resource implementation in `internal/provider/`
4. Update provider to include the new resource
5. Add documentation and examples

### Running Tests

```shell
make test
```

### Local Testing

1. Build and install the provider:
```shell
make build install
```

2. Use the example configurations in `local_testing/`:
```shell
cd local_testing
terraform init
terraform apply
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This provider is licensed under the MIT License. See the LICENSE file for details.
