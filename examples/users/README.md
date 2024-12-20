# OpenWebUI User Data Source Example

This example demonstrates how to use the OpenWebUI provider to look up existing users in your OpenWebUI instance.

## Usage

To run this example:

1. Replace the provider configuration values with your own:
   - `endpoint`: Your OpenWebUI instance URL
   - `token`: Your OpenWebUI API token

2. Initialize Terraform:
   ```bash
   terraform init
   ```

3. Run Terraform:
   ```bash
   terraform plan
   terraform apply
   ```

## Features Demonstrated

The data source supports finding users by:
- Email address
- Name
- ID

You can only use one lookup method at a time. For example:
```hcl
# Find by email
data "openwebui_user" "example" {
  email = "user@example.com"
}

# Or find by name
data "openwebui_user" "example" {
  name = "Example User"
}

# Or find by ID
data "openwebui_user" "example" {
  id = "user-id-123"
}
```

## Available User Information

The data source provides the following user information:
- `id` - The unique identifier of the user
- `name` - The user's display name
- `email` - The user's email address
- `role` - The user's role (pending, admin, or user)
- `profile_image_url` - URL of the user's profile image
- `last_active_at` - Timestamp of the user's last activity
- `created_at` - Timestamp when the user was created
- `updated_at` - Timestamp when the user was last updated
- `api_key` - The user's API key (if any)
- `settings` - User settings (if any)
- `info` - Additional user information (if any)
- `oauth_sub` - OAuth subject identifier (if any)

## Notes

- The OpenWebUI API does not support user creation through the API. Users must be created through the web interface.
- This data source is read-only and cannot modify user information.
- All timestamps are in Unix epoch format.
