---
layout: data-source
page_title: "openwebui_terminal_server_verify Data Source"
sidebar_current: docs-openwebui-datasource-terminal-server-verify
description: |-
  Verifies Open WebUI terminal server connectivity.
---

# openwebui_terminal_server_verify (Data Source)

Verifies that Open WebUI can reach a terminal server connection.

## Example Usage

```hcl
data "openwebui_terminal_server_verify" "example" {
  url       = "https://terminal.example.com"
  path      = "/openapi.json"
  auth_type = "bearer"
  key       = var.terminal_server_token
}
```

## Argument Reference

* `url` (Required) – Terminal server base URL.
* `path` (Optional) – Terminal server path.
* `key` (Optional, Sensitive) – Authentication key.
* `auth_type` (Optional) – Authentication type.
* `config_json` (Optional) – JSON object with extra configuration.

## Attribute Reference

* `verified` – Whether verification succeeded.
