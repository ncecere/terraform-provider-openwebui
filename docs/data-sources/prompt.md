---
layout: data-source
page_title: "openwebui_prompt Data Source"
sidebar_current: docs-openwebui-data-source-prompt
description: |-
  Retrieves details for an existing Open WebUI prompt.
---

# openwebui_prompt (Data Source)

Use this data source to read the definition of an existing prompt, including its access control lists.

## Example Usage

```hcl
data "openwebui_prompt" "triage" {
  command = "triage"
}
```

## Argument Reference

* `command` (Required) – Command string that identifies the prompt.

## Attribute Reference

* `id` – Identifier for the prompt (mirrors `command`).
* `title` – Prompt title displayed in Open WebUI.
* `content` – Prompt body text.
* `read_groups` – Group names granted read access.
* `write_groups` – Group names granted write access.
* `timestamp` – Prompt timestamp formatted as `YYYY-MM-DD`.
* `user_id` – Identifier of the user who owns the prompt.
