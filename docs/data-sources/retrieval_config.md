---
layout: data-source
page_title: "openwebui_retrieval_config Data Source"
sidebar_current: docs-openwebui-datasource-retrieval_config
description: |-
  Reads Open WebUI retrieval_config.
---

# openwebui_retrieval_config (Data Source)

Reads Open WebUI retrieval_config as raw JSON.

## Example Usage

```hcl
data "openwebui_retrieval_config" "current" {}
```

## Attribute Reference

* `config_json` – Raw JSON configuration payload.
