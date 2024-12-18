---
page_title: "openwebui_group Resource - terraform-provider-openwebui"
subcategory: ""
description: |-
  Manages a group in OpenWebUI.
---

# openwebui_group (Resource)

Manages a group in OpenWebUI. Groups allow you to organize users and control their permissions for various features within OpenWebUI.

## Example Usage

```terraform
resource "openwebui_group" "example" {
  name        = "example-group"
  description = "Example group created via Terraform"
  user_ids    = ["user-id-1", "user-id-2"]

  permissions = {
    workspace = {
      models    = true
      knowledge = true
      prompts   = false
      tools     = false
    }
    chat = {
      file_upload = true
      delete      = true
      edit        = true
      temporary   = false
    }
  }
}
```

## Argument Reference

* `name` - (Required) The name of the group.
* `description` - (Optional) A description of the group.
* `user_ids` - (Optional) List of user IDs to include in the group.
* `permissions` - (Optional) Group permissions configuration block.
  * `workspace` - (Required) Workspace-level permissions block.
    * `models` - (Required) Allow access to models.
    * `knowledge` - (Required) Allow access to knowledge bases.
    * `prompts` - (Required) Allow access to prompts.
    * `tools` - (Required) Allow access to tools.
  * `chat` - (Required) Chat-level permissions block.
    * `file_upload` - (Required) Allow file uploads in chat.
    * `delete` - (Required) Allow message deletion.
    * `edit` - (Required) Allow message editing.
    * `temporary` - (Required) Allow temporary messages.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the group.

## Import

Groups can be imported using their ID:

```bash
terraform import openwebui_group.example "group-id"
