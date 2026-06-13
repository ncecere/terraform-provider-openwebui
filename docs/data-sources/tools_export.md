---
layout: data-source
page_title: "openwebui_tools_export Data Source"
sidebar_current: docs-openwebui-datasource-tools_export
description: |-
  Reads Open WebUI tools_export data.
---

# openwebui_tools_export (Data Source)

Reads Open WebUI tools_export data and exposes the raw JSON payload.

## Example Usage

```hcl
data "openwebui_tools_export" "example" {}
```

## Attribute Reference

* `json` – Raw JSON response payload.
