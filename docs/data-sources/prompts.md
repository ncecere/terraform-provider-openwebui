---
layout: data-source
page_title: "openwebui_prompts Data Source"
sidebar_current: docs-openwebui-datasource-prompts
description: |-
  Reads Open WebUI prompts data.
---

# openwebui_prompts (Data Source)

Reads Open WebUI prompts data and exposes the raw JSON payload.

## Example Usage

```hcl
data "openwebui_prompts" "example" {}
```

## Attribute Reference

* `json` – Raw JSON response payload.
