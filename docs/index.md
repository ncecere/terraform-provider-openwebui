# OpenWebUI Provider

The OpenWebUI provider allows you to manage OpenWebUI resources through Terraform. This provider can be used to automate the management of knowledge bases and other OpenWebUI resources.

## Example Usage

```hcl
terraform {
  required_providers {
    openwebui = {
      source = "ncecere/openwebui"
    }
  }
}

provider "openwebui" {
  endpoint = "http://your-openwebui-instance"  # Optional: can use OPENWEBUI_ENDPOINT env var
  token    = "your-api-token"                  # Optional: can use OPENWEBUI_TOKEN env var
}
```

## Authentication

The provider supports authentication through either provider configuration or environment variables.

### Provider Configuration

```hcl
provider "openwebui" {
  endpoint = "http://your-openwebui-instance"
  token    = "your-api-token"
}
```

### Environment Variables

```bash
export OPENWEBUI_ENDPOINT="http://your-openwebui-instance"
export OPENWEBUI_TOKEN="your-api-token"
```

## Schema

### Required

- `endpoint` (String) The URL of your OpenWebUI instance. Can also be set with the `OPENWEBUI_ENDPOINT` environment variable.
- `token` (String, Sensitive) The API token for authentication. Can also be set with the `OPENWEBUI_TOKEN` environment variable.

## Resources

- [Knowledge Base](./resources/knowledge.md) - Manage knowledge bases in OpenWebUI

## Data Sources

- [Knowledge Base](./data-sources/knowledge.md) - Query existing knowledge bases in OpenWebUI

## Development

The provider is built with a modular architecture:

- Client Layer: Handles API communication and resource-specific operations
  - Base client interface and types
  - Resource-specific client implementations (e.g., knowledge)

- Provider Layer: Implements the Terraform provider interface
  - Resource implementations
  - Data source implementations
  - Provider configuration

This structure allows for easy extension with new resources and features while maintaining clean separation of concerns.
