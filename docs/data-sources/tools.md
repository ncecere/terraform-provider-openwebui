---
layout: data-source
page_title: "openwebui_tools Data Source"
sidebar_current: docs-openwebui-datasource-tools
description: |-
  Reads Open WebUI tools data.
---

# openwebui_tools (Data Source)

Reads Open WebUI tools data and exposes the raw JSON payload.

## Example Usage

```hcl
data "openwebui_tools" "example" {}
```

## Attribute Reference

* `json` – Raw JSON response payload.
