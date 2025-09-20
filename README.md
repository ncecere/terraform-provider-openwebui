# Terraform Provider for Open WebUI

This repository contains an experimental Terraform provider that manages Open WebUI resources via its REST API. The initial implementation supports:

- Knowledge bases
- Models
- Prompts
- Groups

> ⚠️ The provider is in an early stage. API compatibility may change as Open WebUI evolves and the provider gains richer coverage and testing.

## What's New in 2.0.0

- Added continuous delivery via GitHub Actions to publish tagged releases directly to the Terraform Registry.
- Normalised prompt commands so Terraform configurations can omit the leading `/` without causing API mismatches.
- Simplified the group resource by removing unsupported `data_json` / `meta_json` arguments and stabilised group member ordering to avoid false-positive plans.
- Updated examples and documentation to use the new structured `params`/`capabilities` attributes and the `~> 2.0` provider constraint.

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
mkdir -p ~/.terraform.d/plugins/local/openwebui/openwebui/2.0.0/
cp terraform-provider-openwebui ~/.terraform.d/plugins/local/openwebui/openwebui/2.0.0/darwin_arm64/
```

Adjust the path and OS/architecture segment to match your environment.

## Publishing a Release

Tagged releases matching `v*.*.*` trigger the GitHub Actions workflow that builds provider artifacts and publishes them to the Terraform Registry. To cut a release:

```bash
git tag -a v2.0.0 -m "Release 2.0.0"
git push origin v2.0.0
```

Ensure the repository is configured with a `TERRAFORM_REGISTRY_TOKEN` secret that has permission to publish to registry.terraform.io.

The workflow expects a signing key so that GoReleaser can sign the checksum file. Add two additional repository secrets before releasing:

- `GPG_PRIVATE_KEY` – ASCII-armoured private key used for signing.
- `PASSPHRASE` – Passphrase for the above key (leave blank if the key is not protected).

## Provider Configuration

```hcl
terraform {
  required_providers {
    openwebui = {
      source  = "local/openwebui/openwebui"
      version = "2.0.0"
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

  read_groups  = ["Support"]
  write_groups = ["Support"]
}
```

### Model

```hcl
resource "openwebui_model" "example" {
  model_id = "custom-rag"
  name     = "Custom Retrieval Model"

  description   = "Retriever tuned for internal knowledge base"
  base_model_id = "gpt-4o"
  is_active     = true

  params = {
    temperature = 0.1
    max_tokens  = 512
  }

  capabilities = {
    web_search = true
  }
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
    "jim@school.edu",
    "john@school.edu",
  ]
  permissions = {
    workspace = {
      models    = true
      knowledge = true
      prompts   = true
      tools     = true
    }

    sharing = {
      public_models = false
    }

    chat = {
      file_upload         = true
      delete              = true
      edit                = true
      continue_response   = true
      regenerate_response = true
      temporary           = true
    }

    features = {
      web_search       = true
      image_generation = true
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
