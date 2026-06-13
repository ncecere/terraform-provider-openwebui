---
layout: data-source
page_title: "openwebui_prompt_history Data Source"
sidebar_current: docs-openwebui-datasource-prompt-history
description: |-
  Reads Open WebUI prompt version history.
---

# openwebui_prompt_history (Data Source)

Reads version history for a prompt by command.

```hcl
data "openwebui_prompt_history" "example" {
  command = openwebui_prompt.example.command
}
```

## Argument Reference

* `command` (Required) – Prompt command.
* `page` (Optional) – History page.

## Attribute Reference

* `history_json` – Raw JSON history payload.
