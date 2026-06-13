---
layout: resource
page_title: "openwebui_default_user_permissions Resource"
sidebar_current: docs-openwebui-resource-default_user_permissions
description: |-
  Applies Open WebUI default_user_permissions.
---

# openwebui_default_user_permissions (Resource)

Applies Open WebUI default_user_permissions from raw JSON.

> This resource stores the desired JSON in Terraform state to avoid drift from server-added defaults. Use the matching data source to round-trip current server config when needed.

## Example Usage

```hcl
data "openwebui_default_user_permissions" "current" {}

resource "openwebui_default_user_permissions" "default" {
  config_json = data.openwebui_default_user_permissions.current.config_json
}
```

## Argument Reference

* `config_json` (Required) – Raw JSON configuration payload.

## Attribute Reference

* `id` – Singleton config identifier.
