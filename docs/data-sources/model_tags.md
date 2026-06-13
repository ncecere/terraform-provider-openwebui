---
layout: data-source
page_title: "openwebui_model_tags Data Source"
sidebar_current: docs-openwebui-datasource-model_tags
description: |-
  Reads Open WebUI model_tags data.
---

# openwebui_model_tags (Data Source)

Reads Open WebUI model_tags data and exposes the raw JSON payload.

## Example Usage

```hcl
data "openwebui_model_tags" "example" {}
```

## Attribute Reference

* `json` – Raw JSON response payload.
