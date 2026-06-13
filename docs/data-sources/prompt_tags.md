---
layout: data-source
page_title: "openwebui_prompt_tags Data Source"
sidebar_current: docs-openwebui-datasource-prompt_tags
description: |-
  Reads Open WebUI prompt_tags data.
---

# openwebui_prompt_tags (Data Source)

Reads Open WebUI prompt_tags data and exposes the raw JSON payload.

## Example Usage

```hcl
data "openwebui_prompt_tags" "example" {}
```

## Attribute Reference

* `json` – Raw JSON response payload.
