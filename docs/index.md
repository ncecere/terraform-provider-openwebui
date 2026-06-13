---
layout: provider
page_title: "OpenWebUI Provider"
sidebar_current: docs-openwebui-index
description: |-
  Interact with Open WebUI knowledge bases, models, prompts, and groups using Terraform.
---

# OpenWebUI Provider

The OpenWebUI provider lets you manage Open WebUI workspaces through Terraform, including knowledge bases, models, prompts, groups, tools, functions, folders, files, and selected admin configuration. It communicates with an Open WebUI deployment via the REST API and requires a bearer token for authentication.

> **3.0.0** – Updated for Open WebUI v0.9.6 with expanded file/folder/knowledge/function/prompt-history coverage, raw catalog data sources, and broad safe E2E validation. See `CHANGELOG.md` for release notes.

## Example Usage

Start with the [minimal example](guides/minimal) to validate provider authentication, then review the [full example](guides/full) for a representative Open WebUI v0.9.6 workflow. Runnable versions live in [`examples/minimal`](../examples/minimal) and [`examples/full`](../examples/full).

```hcl
terraform {
  required_providers {
    openwebui = {
      source  = "nickcecere/openwebui"
      version = "~> 3.0"
    }
  }
}

provider "openwebui" {
  endpoint = "https://openwebui.example.com/api/v1"
  token    = var.openwebui_token
}
```

## Authentication

Authentication uses an HTTP bearer token. Supply it either directly with the `token` argument or through the `OPENWEBUI_TOKEN` environment variable.

## Configuration Reference

The provider supports the following configuration arguments:

* `endpoint` (Optional) – Base URL for your Open WebUI instance (defaults to `http://localhost:3000/api/v1`).
* `token` (Optional, Sensitive) – API token for authenticating requests. Can also be set via `OPENWEBUI_TOKEN`.

## Environment Variables

* `OPENWEBUI_ENDPOINT` – Overrides the API endpoint.
* `OPENWEBUI_TOKEN` – Supplies the API token when the provider block omits `token`.

## Available Resources

* [`openwebui_knowledge`](resources/knowledge)
* [`openwebui_model`](resources/model)
* [`openwebui_prompt`](resources/prompt)
* [`openwebui_group`](resources/group)
* [`openwebui_tool`](resources/tool)
* [`openwebui_tool_valves`](resources/tool_valves)
* [`openwebui_function`](resources/function)
* [`openwebui_function_valves`](resources/function_valves)
* [`openwebui_pipeline`](resources/pipeline)
* [`openwebui_pipeline_valves`](resources/pipeline_valves)
* [`openwebui_file`](resources/file)
* [`openwebui_knowledge_file`](resources/knowledge_file)
* [`openwebui_knowledge_dir`](resources/knowledge_dir)
* [`openwebui_folder`](resources/folder)
* [`openwebui_config_import`](resources/config_import)
* [`openwebui_connections_config`](resources/connections_config)
* [`openwebui_tool_servers_config`](resources/tool_servers_config)
* [`openwebui_terminal_servers_config`](resources/terminal_servers_config)
* [`openwebui_code_execution_config`](resources/code_execution_config)
* [`openwebui_models_config`](resources/models_config)
* [`openwebui_suggestions_config`](resources/suggestions_config)
* [`openwebui_banners_config`](resources/banners_config)
* [`openwebui_oauth_client`](resources/oauth_client)
* [`openwebui_retrieval_config`](resources/retrieval_config)
* [`openwebui_evaluations_config`](resources/evaluations_config)
* [`openwebui_default_user_permissions`](resources/default_user_permissions)

## Available Data Sources

* [`openwebui_knowledge`](data-sources/knowledge)
* [`openwebui_knowledge_files`](data-sources/knowledge_files)
* [`openwebui_model`](data-sources/model)
* [`openwebui_models`](data-sources/models)
* [`openwebui_base_models`](data-sources/base_models)
* [`openwebui_model_tags`](data-sources/model_tags)
* [`openwebui_models_export`](data-sources/models_export)
* [`openwebui_prompt`](data-sources/prompt)
* [`openwebui_prompts`](data-sources/prompts)
* [`openwebui_prompt_tags`](data-sources/prompt_tags)
* [`openwebui_prompt_history`](data-sources/prompt_history)
* [`openwebui_prompt_history_entry`](data-sources/prompt_history_entry)
* [`openwebui_prompt_history_diff`](data-sources/prompt_history_diff)
* [`openwebui_group`](data-sources/group)
* [`openwebui_tool`](data-sources/tool)
* [`openwebui_function`](data-sources/function)
* [`openwebui_functions`](data-sources/functions)
* [`openwebui_tools`](data-sources/tools)
* [`openwebui_tools_export`](data-sources/tools_export)
* [`openwebui_pipeline`](data-sources/pipeline)
* [`openwebui_file`](data-sources/file)
* [`openwebui_files`](data-sources/files)
* [`openwebui_file_content`](data-sources/file_content)
* [`openwebui_file_process_status`](data-sources/file_process_status)
* [`openwebui_config_export`](data-sources/config_export)
* [`openwebui_user`](data-sources/user)
* [`openwebui_retrieval_config`](data-sources/retrieval_config)
* [`openwebui_evaluations_config`](data-sources/evaluations_config)
* [`openwebui_default_user_permissions`](data-sources/default_user_permissions)
* [`openwebui_tool_server_verify`](data-sources/tool_server_verify)
* [`openwebui_folder`](data-sources/folder)
* [`openwebui_terminal_server_verify`](data-sources/terminal_server_verify)

## Import

Most resources support import; refer to each resource page for the supported ID format. Singleton config resources accept any ID, and composite resources (like knowledge attachments) document their composite ID format.

Use the `terraform import` command with the relevant resource type and identifier, for example:

```bash
terraform import openwebui_group.example 65e5e86e-0e23-4cd8-8eee-447c6923f632
```

## Coverage

See [coverage](coverage) for a summary of supported API areas.

## Limitations

This provider is experimental. It reflects the REST API behaviour captured in the supplied `openapi.json` and may require adjustments for other Open WebUI versions. Acceptance tests are available but require a live Open WebUI instance and admin credentials.
