# Terraform Provider for Open WebUI

This repository contains an experimental Terraform provider that manages Open WebUI resources via its REST API. The current implementation supports:

- Knowledge bases
- Knowledge base file attachments
- Knowledge base directories
- Models
- Prompts
- Groups
- Tools and tool valves
- Functions and function valves
- Pipelines and pipeline valves
- Files
- Folders
- Admin configs (connections, tool servers, terminal servers, code execution, models, suggestions, banners, retrieval, evaluations, default user permissions)
- Config import/export
- Raw catalog data sources for models, prompts, tools, and tags
- Prompt history data sources
- OAuth client registration

> ⚠️ The provider is in an early stage. API compatibility may change as Open WebUI evolves and the provider gains richer coverage and testing.

## What's New in 3.0.0

- Updated compatibility for Open WebUI v0.9.6 and refreshed the bundled OpenAPI schema.
- Added safe live E2E coverage under `dev_testing/e2e` for core resources, catalog data sources, and singleton config round-trips.
- Added folders, knowledge directories, file content/process status, terminal server config/verify, functions/function valves, prompt history entry/diff, and broad raw catalog data sources.
- Improved v0.9.6 prompt, model, file-list, and group compatibility fixes.
- See [`CHANGELOG.md`](CHANGELOG.md) for the full release notes.

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
mkdir -p ~/.terraform.d/plugins/local/openwebui/openwebui/3.0.0/
cp terraform-provider-openwebui ~/.terraform.d/plugins/local/openwebui/openwebui/3.0.0/darwin_arm64/
```

Adjust the path and OS/architecture segment to match your environment.

## Acceptance Tests

Acceptance tests require a live Open WebUI instance and an admin API token.

Set the following environment variables before running tests:

- `TF_ACC=1`
- `OPENWEBUI_TOKEN` (required)
- `OPENWEBUI_ENDPOINT` (optional; defaults to `http://localhost:3000/api/v1`)

Optional acceptance tests require additional variables:

- `OPENWEBUI_TEST_CONFIG_IMPORT=1` to enable config import tests
- `OPENWEBUI_OAUTH_CLIENT_URL` to enable OAuth client registration tests
- `OPENWEBUI_TOOL_SERVER_URL` to enable tool server verification tests

Run the acceptance tests:

```bash
go test ./internal/provider -run TestAcc -v
```

## Publishing a Release

Tagged releases matching `v*.*.*` trigger the GitHub Actions workflow that builds provider artifacts and publishes them to the Terraform Registry. To cut a release:

```bash
git tag -a v3.0.0 -m "Release 3.0.0"
git push origin v3.0.0
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
      version = "3.0.0"
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

### Tool

```hcl
resource "openwebui_tool" "calculator" {
  tool_id = "calculator"
  name    = "Calculator"
  content = file("./tools/calculator.py")

  description   = "Internal calculator tool"
  manifest_json = jsonencode({
    version = "1.0.0"
  })
}

resource "openwebui_tool_valves" "calculator" {
  tool_id = openwebui_tool.calculator.id
  valves_json = jsonencode({
    enabled = true
  })
}
```

### File

```hcl
resource "openwebui_file" "support_doc" {
  source_path = "./docs/support_faq.txt"
  metadata_json = jsonencode({
    category = "support"
  })
}
```

## Examples

Reference configurations live under `examples/`:

- [`examples/minimal`](examples/minimal) – smallest practical configuration; creates and reads one prompt.
- [`examples/full`](examples/full) – representative v3.0.0 workflow covering folders, prompts/history, knowledge directories, files, functions, and catalog data sources.
- [`examples/basic`](examples/basic) – legacy core-resource example for knowledge, models, prompts, groups, and tools.
- [`examples/admin`](examples/admin) – admin-focused flows that are intentionally guarded behind variables because they can affect global server settings.

Provider documentation also includes [minimal](docs/guides/minimal.md) and [full](docs/guides/full.md) walkthroughs.

## Known Limitations / Next Steps

- Live E2E coverage is broad but intentionally avoids fixture-backed/default-destructive cases; pipelines, pipeline valves, OAuth client registration, config import, and external tool/terminal server verification remain opt-in or deferred.
- MCP/tool server item add/remove testing is deferred; current server config resources manage singleton configuration payloads.
- Several Open WebUI v0.9.6 endpoints expose loose schemas, so the provider intentionally uses raw JSON fields/data sources for unstable catalog, history, and config payloads.
- Authentication is limited to bearer tokens. If Open WebUI exposes alternative auth flows they are not yet supported.
- The suggestions config endpoint does not expose a read API, so state is maintained from the last apply.
- Config export/import and admin config payloads may include sensitive values; treat Terraform state accordingly.
- Additional Open WebUI resources (settings, datasets, agents, etc.) can be lifted into Terraform following the patterns used here.

Contributions and feedback are welcome.
