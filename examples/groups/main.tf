# Configure the OpenWebUI Provider
provider "openwebui" {
  endpoint = "https://your-openwebui-instance"
  # Token can also be provided via OPENWEBUI_TOKEN environment variable
  token = "your-token"
}

# Create a group
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
