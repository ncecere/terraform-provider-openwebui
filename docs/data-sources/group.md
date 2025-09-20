---
layout: data-source
page_title: "openwebui_group Data Source"
sidebar_current: docs-openwebui-data-source-group
description: |-
  Retrieves information about an existing Open WebUI group.
---

# openwebui_group (Data Source)

Use this data source to read an existing group's membership, permissions, and metadata.

## Example Usage

```hcl
data "openwebui_group" "admins" {
  name = "Administrators"
}
```

## Argument Reference

* `group_id` (Optional) – Identifier of the group to retrieve. When omitted, `name` must be provided.
* `name` (Optional) – Case-insensitive group name used to resolve the group when `group_id` is not supplied. If the name is not unique the lookup fails and the ID must be specified.

## Attribute Reference

* `id` – Open WebUI identifier for the group (matches `group_id`).
* `name` – Group name.
* `description` – Group description text.
* `users` – List of user labels (emails, usernames, or names) that belong to the group.
* `permissions` – Nested block exposing maps of boolean flags for `workspace`, `sharing`, `chat`, and `features` permissions.
* `meta_json` – JSON string containing metadata attached to the group.
* `data_json` – JSON string containing additional group data.
* `user_id` – Identifier of the user who created the group.
* `created_at` – Creation date in `YYYY-MM-DD` format.
* `updated_at` – Last update date in `YYYY-MM-DD` format.
