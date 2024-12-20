---
page_title: "openwebui_user Data Source - terraform-provider-openwebui"
subcategory: ""
description: |-
  Data source for retrieving OpenWebUI user information.
---

# openwebui_user (Data Source)

Use this data source to retrieve information about an existing OpenWebUI user. Users can be looked up by their ID, email address, or name. This data source is particularly useful when you need to reference existing users for group assignments, permissions, or when you need to access user metadata for integration purposes.

## Example Usage

### Basic Usage

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

### Complete Example with Outputs

```hcl
# Query user information
data "openwebui_user" "full_example" {
  email = "user@example.com"
}

# Output all available user information
output "user_details" {
  value = {
    # Basic Information
    id                = data.openwebui_user.full_example.id
    name              = data.openwebui_user.full_example.name
    email             = data.openwebui_user.full_example.email
    role              = data.openwebui_user.full_example.role
    profile_image_url = data.openwebui_user.full_example.profile_image_url
    
    # Timestamps
    last_active_at    = data.openwebui_user.full_example.last_active_at
    created_at        = data.openwebui_user.full_example.created_at
    updated_at        = data.openwebui_user.full_example.updated_at
    
    # Advanced Properties
    api_key           = data.openwebui_user.full_example.api_key
    settings          = data.openwebui_user.full_example.settings
    info              = data.openwebui_user.full_example.info
    oauth_sub         = data.openwebui_user.full_example.oauth_sub
  }
}

# Example of using user data in other resources
resource "openwebui_group" "example" {
  name = "Example Group"
  members = [data.openwebui_user.full_example.id]
}
```

## Argument Reference

-> **Note** Only one of `id`, `email`, or `name` can be specified. Using multiple lookup methods will result in an error.

* `id` - (Optional) The unique identifier of the user. Use this when you have a specific user ID from another resource or data source.
* `email` - (Optional) The email address of the user. This is often the most reliable way to look up users as email addresses are unique in the system.
* `name` - (Optional) The display name of the user. Note that names may not be unique across all users.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `role` - The role assigned to the user. Possible values:
  * `pending` - User account is created but not yet activated
  * `admin` - User has administrative privileges
  * `user` - Standard user account
* `profile_image_url` - URL of the user's profile image. Can be empty if no custom image is set.
* `last_active_at` - Unix timestamp (in seconds) of the user's last activity. Useful for monitoring user engagement.
* `created_at` - Unix timestamp (in seconds) when the user account was created.
* `updated_at` - Unix timestamp (in seconds) when the user account was last modified.
* `api_key` - The user's API key for programmatic access. Only available if an API key has been generated for the user.
* `settings` - A nested map containing user-specific settings:
  * `ui` - A map of UI preferences and configurations. Common settings include:
    * `theme` - User's preferred UI theme
    * `language` - User's preferred language
    * `notifications` - Notification preferences
* `info` - A map containing additional metadata about the user. This can include custom fields or integration-specific data.
* `oauth_sub` - The OAuth subject identifier if the user was authenticated through an OAuth provider.

## Import

This is a data source and does not support import operations.

## Notes

* The OpenWebUI API does not support user creation through the API. Users must be created through the web interface or OAuth authentication.
* This data source is read-only and cannot modify user information. Any changes to user data must be made through the OpenWebUI interface.
* All timestamps are in Unix epoch format (seconds since January 1, 1970 UTC).
* The `settings` and `info` maps may contain different keys depending on your OpenWebUI configuration and any custom fields that have been defined.
* API keys are sensitive information and should be handled securely. Consider using Terraform's sensitive output feature when exposing API keys in outputs.

## Related Resources

* [openwebui_group Resource](../resources/group.md) - Manage groups that users can be members of
* [openwebui_group Data Source](group.md) - Query existing groups and their members
