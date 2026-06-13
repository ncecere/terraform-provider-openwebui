---
layout: data-source
page_title: "openwebui_models Data Source"
sidebar_current: docs-openwebui-datasource-models
description: |-
  Reads Open WebUI models data.
---

# openwebui_models (Data Source)

Reads Open WebUI models data and exposes the raw JSON payload.

## Example Usage

```hcl
data "openwebui_models" "example" {}
```

## Attribute Reference

* `json` – Raw JSON response payload.
