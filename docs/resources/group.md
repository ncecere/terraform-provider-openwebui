---
layout: resource
page_title: "openwebui_group Resource"
sidebar_current: docs-openwebui-resource-group
description: |-
  Manages groups and their permissions in Open WebUI.
---

# openwebui_group (Resource)

Provisions a group and controls its membership and feature permissions.

## Example Usage

```hcl
resource "openwebui_group" "support" {
  name        = "Support"
  description = "Support team access"

  users = [
    "cj01@ufl.edu",
    "ebs@ufl.edu",
  ]

  permissions = {
    workspace = {
      models    = true
      knowledge = true
      prompts   = true
      tools     = true
    }

    sharing = {
      public_models = false
    }

    chat = {
      file_upload         = true
      delete              = true
      edit                = true
      continue_response   = true
      regenerate_response = true
      temporary           = true
    }

    features = {
      web_search       = true
      image_generation = true
    }
  }
}
```

## Argument Reference

* `name` (Required) – Group name.
* `description` (Required) – Description visible within Open WebUI.
* `users` (Optional) – List of usernames or email addresses. The provider resolves them to the required user IDs automatically when creating or updating the group.
* `permissions` (Optional) – Nested block defining category-specific permissions. Each map only accepts recognised keys:
  * `workspace` – `models`, `knowledge`, `prompts`, `tools`
  * `sharing` – `public_models`, `public_knowledge`, `public_prompts`, `public_tools`
  * `chat` – `controls`, `valves`, `system_prompt`, `params`, `file_upload`, `delete`, `delete_message`, `continue_response`, `regenerate_response`, `rate_response`, `edit`, `share`, `export`, `stt`, `tts`, `call`, `multiple_models`, `temporary`, `temporary_enforced`
  * `features` – `direct_tool_servers`, `web_search`, `image_generation`, `code_interpreter`, `notes`

## Attribute Reference

* `id` – Unique group identifier assigned by Open WebUI.
* `created_at` – Creation date in `YYYY-MM-DD` format.
* `updated_at` – Last update date in `YYYY-MM-DD` format.
* `user_id` – Identifier of the user that created the group.
* `users` – Resolved usernames/email addresses currently associated with the group (sorted by email/username).

## Import

Groups can be imported using the group ID:

```bash
terraform import openwebui_group.support 65e5e86e-0e23-4cd8-8eee-447c6923f632
```
