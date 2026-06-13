---
layout: resource
page_title: "openwebui_retrieval_config Resource"
sidebar_current: docs-openwebui-resource-retrieval_config
description: |-
  Applies Open WebUI retrieval_config.
---

# openwebui_retrieval_config (Resource)

Applies Open WebUI retrieval_config from raw JSON.

> This resource stores the desired JSON in Terraform state to avoid drift from server-added defaults. Use the matching data source to round-trip current server config when needed.

## Example Usage

```hcl
data "openwebui_retrieval_config" "current" {}

resource "openwebui_retrieval_config" "default" {
  config_json = data.openwebui_retrieval_config.current.config_json
}
```

## Argument Reference

* `config_json` (Required) – Raw JSON configuration payload.

## Attribute Reference

* `id` – Singleton config identifier.
