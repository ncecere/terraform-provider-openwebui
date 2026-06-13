---
layout: data-source
page_title: "openwebui_models_export Data Source"
sidebar_current: docs-openwebui-datasource-models_export
description: |-
  Reads Open WebUI models_export data.
---

# openwebui_models_export (Data Source)

Reads Open WebUI models_export data and exposes the raw JSON payload.

## Example Usage

```hcl
data "openwebui_models_export" "example" {}
```

## Attribute Reference

* `json` – Raw JSON response payload.
