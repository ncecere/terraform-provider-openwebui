terraform {
  required_providers {
    openwebui = {
      source = "ncecere/openwebui"
    }
  }
}

# Configure the OpenWebUI Provider
provider "openwebui" {
  endpoint = "https://chat.example.com"  # Your OpenWebUI endpoint
  token    = "your-api-token"           # Your OpenWebUI API token
}

# Example: Find user by email
data "openwebui_user" "example" {
  email = "user@example.com"  # You can use email to find a user
  # Or use name = "Example User" to find by name
  # Or use id = "user-id-123" to find by ID
  # Note: Only use one of: email, name, or id
}

# Output user information
output "user_info" {
  value = {
    id                = data.openwebui_user.example.id
    name              = data.openwebui_user.example.name
    email             = data.openwebui_user.example.email
    role              = data.openwebui_user.example.role
    profile_image_url = data.openwebui_user.example.profile_image_url
    last_active_at    = data.openwebui_user.example.last_active_at
    created_at        = data.openwebui_user.example.created_at
    updated_at        = data.openwebui_user.example.updated_at
  }
}
