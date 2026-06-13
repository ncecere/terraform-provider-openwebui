---
layout: data-source
page_title: "openwebui_prompt_history_entry Data Source"
sidebar_current: docs-openwebui-datasource-prompt-history-entry
---

# openwebui_prompt_history_entry (Data Source)

Reads a specific Open WebUI prompt history entry as raw JSON.

```hcl
data "openwebui_prompt_history" "example" {
  command = openwebui_prompt.example.command
}

locals {
  history_id = jsondecode(data.openwebui_prompt_history.example.history_json)[0].id
}

data "openwebui_prompt_history_entry" "example" {
  command    = openwebui_prompt.example.command
  history_id = local.history_id
}
```

* `command` – Prompt command.
* `history_id` – History entry ID.
* `history_json` – Raw history entry JSON.
