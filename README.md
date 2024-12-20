# Terraform Provider for OpenWebUI

This Terraform provider enables infrastructure-as-code management of OpenWebUI resources, allowing you to automate the configuration of users, groups, knowledge bases, and models in your OpenWebUI instance.

## Features

- **User Management**: Query user information and integrate with other resources
- **Group Management**: Create and manage user groups with granular permissions
- **Knowledge Base Management**: Create and configure knowledge bases with access controls
- **Model Management**: Deploy and configure AI models with custom parameters

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19
- An OpenWebUI instance with API access

## Quick Start

1. **Install the Provider**

   ```hcl
   terraform {
     required_providers {
       openwebui = {
         source = "ncecere/openwebui"
       }
     }
   }
   ```

2. **Configure Provider Authentication**

   ```hcl
   provider "openwebui" {
     endpoint = "http://your-openwebui-instance"  # Or use OPENWEBUI_ENDPOINT env var
     token    = "your-api-token"                  # Or use OPENWEBUI_TOKEN env var
   }
   ```

3. **Start Managing Resources**

   See the examples below for common use cases.

## Installation

### From Terraform Registry

The provider is available on the [Terraform Registry](https://registry.terraform.io/providers/ncecere/openwebui/latest). Terraform will automatically download the provider when you run `terraform init`.

### Local Development Build

1. Clone the repository
   ```shell
   git clone git@github.com:ncecere/terraform-provider-openwebui.git
   cd terraform-provider-openwebui
   ```

2. Build and install locally
   ```shell
   make build install
   ```

   This will build and install the provider into your `~/.terraform.d/plugins` directory.

For detailed development instructions, see [DEVELOPMENT.md](DEVELOPMENT.md).

## Example Usage

### User and Group Management

```hcl
# Look up user information
data "openwebui_user" "admin" {
  email = "admin@example.com"
}

# Create an administrative group
resource "openwebui_group" "admins" {
  name        = "administrators"
  description = "Administrative group with full permissions"
  user_ids    = [data.openwebui_user.admin.id]

  permissions = {
    workspace = {
      models    = true
      knowledge = true
      prompts   = true
      tools     = true
    }
    chat = {
      file_upload = true
      delete      = true
      edit        = true
      temporary   = true
    }
  }
}
```

### Knowledge Base Configuration

```hcl
# Create a knowledge base with access control
resource "openwebui_knowledge" "team_docs" {
  name          = "Team Documentation"
  description   = "Internal team documentation and resources"
  access_control = "private"
  
  data = {
    source = "terraform"
    type   = "documentation"
    team   = "engineering"
  }
}

# Query existing knowledge base
data "openwebui_knowledge" "existing" {
  name = "Existing Knowledge Base"
}
```

### Model Deployment

```hcl
# Deploy a custom model
resource "openwebui_model" "custom_assistant" {
  base_model_id = "gpt-4"
  name          = "Engineering Assistant"
  is_active     = true

  params {
    system          = "You are a helpful engineering assistant"
    temperature     = 0.7
    max_tokens      = 2000
    num_ctx         = 4096
  }

  meta {
    description = "Specialized assistant for engineering tasks"
    capabilities {
      vision    = false
      usage     = true
      citations = true
    }
    tags {
      name = "engineering"
    }
  }

  access_control {
    read {
      group_ids = [openwebui_group.admins.id]
    }
  }
}
```

## Documentation

- [Provider Configuration](docs/index.md)
- [User Data Source](docs/data-sources/user.md)
- [Group Resource](docs/resources/group.md)
- [Group Data Source](docs/data-sources/group.md)
- [Knowledge Base Resource](docs/resources/knowledge.md)
- [Knowledge Base Data Source](docs/data-sources/knowledge.md)
- [Model Resource](docs/resources/model.md)
- [Model Data Source](docs/data-sources/model.md)
- [Development Guide](DEVELOPMENT.md)

## Project Structure

```
terraform-provider-openwebui/
├── docs/                    # Provider and resource documentation
├── examples/               # Example configurations for each resource
├── internal/              # Provider implementation
│   └── provider/
│       ├── client/        # API client implementations
│       │   ├── groups/    # Group-specific client
│       │   ├── knowledge/ # Knowledge-specific client
│       │   ├── models/    # Model-specific client
│       │   └── users/     # User-specific client
│       └── ...           # Provider and resource implementations
└── local_testing/        # Local development test configurations
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This provider is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
