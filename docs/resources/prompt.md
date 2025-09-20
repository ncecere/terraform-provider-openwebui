---
layout: resource
page_title: "openwebui_prompt Resource"
sidebar_current: docs-openwebui-resource-prompt
description: |-
  Manages prompt definitions within Open WebUI.
---

# openwebui_prompt (Resource)

Creates and manages prompts exposed to Open WebUI users.

## Example Usage

```hcl
resource "openwebui_prompt" "triage" {
  command = "triage"
  title   = "Ticket triage"
  content = <<-EOT
    You are an assistant that triages inbound support tickets.
  EOT

  read_groups = ["Support"]
  write_groups = ["Support"]
}
```

If both `read_groups` and `write_groups` are omitted (or empty), the prompt remains public.

## Argument Reference

* `command` (Required) – Unique identifier for the prompt.
* `title` (Required) – Display name inside Open WebUI.
* `content` (Required) – Prompt body text.
* `read_groups` (Optional) – List of group names or IDs granted read access.
* `write_groups` (Optional) – List of group names or IDs granted write access. Groups listed here automatically receive read access.

## Attribute Reference

* `id` – Mirrors the `command` string.
* `timestamp` – Date (YYYY-MM-DD) recorded by Open WebUI when the prompt was last updated.
* `user_id` – Identifier of the user that owns the prompt.

## Import

Prompts can be imported using the command string:

```bash
terraform import openwebui_prompt.triage triage
```
