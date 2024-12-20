---
page_title: "openwebui_user Data Source - terraform-provider-openwebui"
subcategory: ""
description: |-
  Data source for retrieving OpenWebUI user information.
---

# openwebui_user (Data Source)

Use this data source to retrieve information about an existing OpenWebUI user. Users can be looked up by their ID, email address, or name.

## Example Usage

```hcl
# Find user by email
data "openwebui_user" "example" {
  email = "user@example.com"
}

# Find user by name
data "openwebui_user" "by_name" {
  name = "Example User"
}

# Find user by ID
data "openwebui_user" "by_id" {
  id = "user-id-123"
}
```

## Argument Reference

-> **Note** Only one of `id`, `email`, or `name` can be specified.

* `id` - (Optional) The ID of the user to look up.
* `email` - (Optional) The email address of the user to look up.
* `name` - (Optional) The name of the user to look up.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `role` - The role of the user (pending, admin, or user).
* `profile_image_url` - URL of the user's profile image.
* `last_active_at` - Timestamp of the user's last activity (Unix epoch format).
* `created_at` - Timestamp when the user was created (Unix epoch format).
* `updated_at` - Timestamp when the user was last updated (Unix epoch format).
* `api_key` - The user's API key (if any).
* `settings` - A map containing user settings.
  * `ui` - A map containing UI-specific settings.
* `info` - A map containing additional user information.
* `oauth_sub` - OAuth subject identifier (if any).

## Import

This is a data source and does not support import operations.

## Notes

* The OpenWebUI API does not support user creation through the API. Users must be created through the web interface.
* This data source is read-only and cannot modify user information.
* All timestamps are in Unix epoch format.
