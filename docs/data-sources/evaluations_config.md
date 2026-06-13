---
layout: data-source
page_title: "openwebui_evaluations_config Data Source"
sidebar_current: docs-openwebui-datasource-evaluations_config
description: |-
  Reads Open WebUI evaluations_config.
---

# openwebui_evaluations_config (Data Source)

Reads Open WebUI evaluations_config as raw JSON.

## Example Usage

```hcl
data "openwebui_evaluations_config" "current" {}
```

## Attribute Reference

* `config_json` – Raw JSON configuration payload.
