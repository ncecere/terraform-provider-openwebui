---
layout: data-source
page_title: "openwebui_base_models Data Source"
sidebar_current: docs-openwebui-datasource-base_models
description: |-
  Reads Open WebUI base_models data.
---

# openwebui_base_models (Data Source)

Reads Open WebUI base_models data and exposes the raw JSON payload.

## Example Usage

```hcl
data "openwebui_base_models" "example" {}
```

## Attribute Reference

* `json` – Raw JSON response payload.
