---
layout: data-source
page_title: "openwebui_prompt_history_diff Data Source"
sidebar_current: docs-openwebui-datasource-prompt-history-diff
---

# openwebui_prompt_history_diff (Data Source)

Reads the raw diff between two Open WebUI prompt history entries.

```hcl
data "openwebui_prompt_history_diff" "example" {
  command = "/example"
  from_id = "history-entry-a"
  to_id   = "history-entry-b"
}
```

* `command` – Prompt command.
* `from_id` – Source history entry ID.
* `to_id` – Target history entry ID.
* `diff_json` – Raw diff JSON.
