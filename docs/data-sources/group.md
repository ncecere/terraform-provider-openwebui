---
page_title: "openwebui_group Data Source - terraform-provider-openwebui"
subcategory: ""
description: |-
  Data source for retrieving OpenWebUI group information.
---

# openwebui_group (Data Source)

Use this data source to retrieve information about an OpenWebUI group by its name.

## Example Usage

```terraform
data "openwebui_group" "example" {
  name = "Example Group"
}

output "group_info" {
  value = {
    id = data.openwebui_group.example.id
    description = data.openwebui_group.example.description
    user_ids = data.openwebui_group.example.user_ids
    permissions = data.openwebui_group.example.permissions
    created_at = data.openwebui_group.example.created_at
    updated_at = data.openwebui_group.example.updated_at
  }
}
```

## Argument Reference

* `name` - (Required) The name of the group to look up.

## Attribute Reference

* `id` - The ID of the group.
* `description` - The description of the group.
* `user_ids` - List of user IDs in the group.
* `permissions` - Group permissions configuration block.
  * `workspace` - Workspace-level permissions block.
    * `models` - Whether access to models is allowed.
    * `knowledge` - Whether access to knowledge bases is allowed.
    * `prompts` - Whether access to prompts is allowed.
    * `tools` - Whether access to tools is allowed.
  * `chat` - Chat-level permissions block.
    * `file_upload` - Whether file uploads in chat are allowed.
    * `delete` - Whether message deletion is allowed.
    * `edit` - Whether message editing is allowed.
    * `temporary` - Whether temporary messages are allowed.
* `created_at` - Timestamp when the group was created.
* `updated_at` - Timestamp when the group was last updated.
