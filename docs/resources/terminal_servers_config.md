---
layout: resource
page_title: "openwebui_terminal_servers_config Resource"
sidebar_current: docs-openwebui-resource-terminal-servers-config
description: |-
  Manages Open WebUI terminal server configuration.
---

# openwebui_terminal_servers_config (Resource)

Updates terminal server connection settings in Open WebUI.

## Example Usage

```hcl
resource "openwebui_terminal_servers_config" "default" {
  connections = [
    {
      name      = "Terminal"
      url       = "https://terminal.example.com"
      path      = "/openapi.json"
      enabled   = true
      auth_type = "bearer"
      key       = var.terminal_server_token
    }
  ]
}
```

## Argument Reference

* `connections` (Required) – List of terminal server connections.
  * `url` (Required) – Base URL for the terminal server.
  * `connection_id` (Optional) – Open WebUI terminal server connection ID.
  * `name` (Optional) – Display name.
  * `enabled` (Optional) – Whether the server is enabled.
  * `path` (Optional) – API/config path.
  * `key` (Optional, Sensitive) – Authentication key.
  * `auth_type` (Optional) – Authentication type.
  * `config_json` (Optional) – JSON object with extra configuration.
  * `server_type` (Optional, Computed) – Configured or detected server type.
  * `policy_id` (Optional) – Orchestrator policy ID.
  * `policy_json` (Optional) – JSON object containing policy data.

## Attribute Reference

* `id` – Singleton identifier for the terminal servers config.

## Import

```bash
terraform import openwebui_terminal_servers_config.default terminal_servers
```
