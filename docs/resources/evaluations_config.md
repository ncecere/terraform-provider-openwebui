---
layout: resource
page_title: "openwebui_evaluations_config Resource"
sidebar_current: docs-openwebui-resource-evaluations_config
description: |-
  Applies Open WebUI evaluations_config.
---

# openwebui_evaluations_config (Resource)

Applies Open WebUI evaluations_config from raw JSON.

> This resource stores the desired JSON in Terraform state to avoid drift from server-added defaults. Use the matching data source to round-trip current server config when needed.

## Example Usage

```hcl
data "openwebui_evaluations_config" "current" {}

resource "openwebui_evaluations_config" "default" {
  config_json = data.openwebui_evaluations_config.current.config_json
}
```

## Argument Reference

* `config_json` (Required) – Raw JSON configuration payload.

## Attribute Reference

* `id` – Singleton config identifier.
