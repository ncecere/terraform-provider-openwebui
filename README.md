# Terraform Provider for Open WebUI

This repository contains an experimental Terraform provider that manages Open WebUI resources via its REST API. The initial implementation supports:

- Knowledge bases
- Models
- Prompts
- Groups

> ⚠️ The provider is in an early stage. API compatibility may change as Open WebUI evolves and the provider gains richer coverage and testing.

## Requirements

- Terraform 1.6 or newer
- Go 1.25 (for building the provider)
- An Open WebUI instance reachable from the machine running Terraform
- An Open WebUI API token (bearer token)

## Building the Provider

Use the provided Makefile targets while developing locally:

```bash
make tidy   # optional; ensures go.mod/go.sum are up to date
make build
make test
```

To install a binary into `./bin` run:

```bash
make install
```

Copy the resulting binary into your Terraform plugin directory, for example on macOS:

```bash
mkdir -p ~/.terraform.d/plugins/local/openwebui/openwebui/0.1.0/
cp terraform-provider-openwebui ~/.terraform.d/plugins/local/openwebui/openwebui/0.1.0/darwin_arm64/
```

Adjust the path and OS/architecture segment to match your environment.

## Provider Configuration

```hcl
terraform {
  required_providers {
    openwebui = {
      source  = "local/openwebui/openwebui"
      version = "0.1.0"
    }
  }
}

provider "openwebui" {
  endpoint = "http://localhost:3000/api/v1"
  token    = var.openwebui_token
}
```

The provider reads the API token from the `token` argument or the `OPENWEBUI_TOKEN` environment variable. The API endpoint defaults to `http://localhost:3000/api/v1` and can be overridden with the `endpoint` argument or `OPENWEBUI_ENDPOINT`.

## Resource Examples

### Knowledge Base

```hcl
resource "openwebui_knowledge" "example" {
  name        = "Support FAQ"
  description = "Knowledge base backing the support chat bot"

  read_groups = ["Support"]
  write_groups = ["Support"]

  data_json = jsonencode({
    category = "support"
  })
}
```

### Model

```hcl
resource "openwebui_model" "example" {
  model_id = "custom-rag"
  name     = "Custom Retrieval Model"

  meta_json = jsonencode({
    description = "Retriever tuned for internal knowledge base"
  })

  params_json = jsonencode({
    temperature = 0.1
  })

  base_model_id = "gpt-4o"
  is_active     = true
}
```

### Prompt

```hcl
resource "openwebui_prompt" "example" {
  command = "triage"
  title   = "Ticket triage"
  content = "You are an assistant that triages support tickets."

  read_groups  = ["Support"]
  write_groups = ["Support"]
}
```

### Group

```hcl
resource "openwebui_group" "example" {
  name        = "Support"
  description = "Support team access"

  users = [
    "cj01@ufl.edu",
    "ebs@ufl.edu",
  ]

  data_json = jsonencode({
    department = "support"
  })

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

## Examples

Reference configurations live under `examples/`. Start with [`examples/basic`](examples/basic) for a provider configuration that exercises all supported resources.

## Known Limitations / Next Steps

- Automated tests are not yet in place; the provider has only been built against API documentation.
- The client currently exchanges opaque JSON fields using raw strings. Typed schemas, validation, and richer Terraform types would improve ergonomics.
- Authentication is limited to bearer tokens. If Open WebUI exposes alternative auth flows they are not yet supported.
- Additional Open WebUI resources (settings, datasets, agents, etc.) can be lifted into Terraform following the patterns used here.

Contributions and feedback are welcome.
