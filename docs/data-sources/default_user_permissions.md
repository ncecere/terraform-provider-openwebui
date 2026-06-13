---
layout: data-source
page_title: "openwebui_default_user_permissions Data Source"
sidebar_current: docs-openwebui-datasource-default_user_permissions
description: |-
  Reads Open WebUI default_user_permissions.
---

# openwebui_default_user_permissions (Data Source)

Reads Open WebUI default_user_permissions as raw JSON.

## Example Usage

```hcl
data "openwebui_default_user_permissions" "current" {}
```

## Attribute Reference

* `config_json` – Raw JSON configuration payload.
