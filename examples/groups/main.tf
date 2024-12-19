# Configure the OpenWebUI Provider
terraform {
  required_providers {
    openwebui = {
      source = "ncecere/openwebui"
    }
  }
}

provider "openwebui" {
  endpoint = "https://your-openwebui-instance"
  # Token can be provided via OPENWEBUI_TOKEN environment variable
  # token = "your-token"
}

# Example 1: Admin Group
resource "openwebui_group" "administrators" {
  name        = "administrators"
  description = "Administrative group with full system access"
  user_ids    = ["admin1", "admin2"]

  permissions = {
    workspace = {
      models    = true    # Can manage all models
      knowledge = true    # Can manage all knowledge bases
      prompts   = true    # Can manage prompts
      tools     = true    # Can manage tools
    }
    chat = {
      file_upload = true  # Can upload files
      delete      = true  # Can delete messages
      edit        = true  # Can edit messages
      temporary   = true  # Can use temporary chats
    }
  }
}

# Example 2: Read-Only Users
resource "openwebui_group" "viewers" {
  name        = "viewers"
  description = "Users with read-only access"
  user_ids    = ["viewer1", "viewer2", "viewer3"]

  permissions = {
    workspace = {
      models    = false   # Cannot manage models
      knowledge = false   # Cannot manage knowledge bases
      prompts   = false   # Cannot manage prompts
      tools     = false   # Cannot manage tools
    }
    chat = {
      file_upload = false # Cannot upload files
      delete      = false # Cannot delete messages
      edit        = false # Cannot edit messages
      temporary   = true  # Can use temporary chats
    }
  }
}

# Example 3: Model Developers
resource "openwebui_group" "model_developers" {
  name        = "model-developers"
  description = "Team responsible for model development and training"
  user_ids    = ["dev1", "dev2"]

  permissions = {
    workspace = {
      models    = true    # Can manage models
      knowledge = true    # Can manage knowledge bases
      prompts   = true    # Can manage prompts
      tools     = false   # Cannot manage tools
    }
    chat = {
      file_upload = true  # Can upload files
      delete      = true  # Can delete messages
      edit        = true  # Can edit messages
      temporary   = true  # Can use temporary chats
    }
  }
}

# Example 4: Content Managers
resource "openwebui_group" "content_managers" {
  name        = "content-managers"
  description = "Team managing knowledge bases and documentation"
  user_ids    = ["cm1", "cm2"]

  permissions = {
    workspace = {
      models    = false   # Cannot manage models
      knowledge = true    # Can manage knowledge bases
      prompts   = true    # Can manage prompts
      tools     = false   # Cannot manage tools
    }
    chat = {
      file_upload = true  # Can upload files
      delete      = false # Cannot delete messages
      edit        = true  # Can edit messages
      temporary   = true  # Can use temporary chats
    }
  }
}

# Create a model with group-based access control
resource "openwebui_model" "managed_model" {
  name          = "Managed Model"
  base_model_id = "gpt-4"
  is_active     = true

  params {
    system = "You are a helpful assistant."
  }

  access_control {
    read {
      # All groups can read
      group_ids = [
        openwebui_group.administrators.id,
        openwebui_group.model_developers.id,
        openwebui_group.content_managers.id,
        openwebui_group.viewers.id
      ]
    }
    write {
      # Only admins and model developers can modify
      group_ids = [
        openwebui_group.administrators.id,
        openwebui_group.model_developers.id
      ]
    }
  }
}

# Create a knowledge base with group-based access
resource "openwebui_knowledge" "managed_kb" {
  name           = "Managed Knowledge Base"
  description    = "Knowledge base with group-based access control"
  access_control = "private"

  data = {
    managed_by = "content-managers"
    type       = "documentation"
  }

  depends_on = [
    openwebui_group.administrators,
    openwebui_group.content_managers
  ]
}

# Use data sources to look up existing groups
data "openwebui_group" "existing_admin" {
  name = openwebui_group.administrators.name
}

data "openwebui_group" "existing_devs" {
  name = openwebui_group.model_developers.name
}

# Outputs for verification and reference
output "group_ids" {
  value = {
    administrators   = openwebui_group.administrators.id
    viewers         = openwebui_group.viewers.id
    model_developers = openwebui_group.model_developers.id
    content_managers = openwebui_group.content_managers.id
  }
}

output "group_permissions" {
  value = {
    admin_permissions = data.openwebui_group.existing_admin.permissions
    dev_permissions   = data.openwebui_group.existing_devs.permissions
  }
}

output "resource_access" {
  value = {
    model_id = openwebui_model.managed_model.id
    kb_id    = openwebui_knowledge.managed_kb.id
  }
}

# Example of checking group membership
output "admin_users" {
  value = {
    group_name = openwebui_group.administrators.name
    user_count = length(openwebui_group.administrators.user_ids)
    users      = openwebui_group.administrators.user_ids
  }
}
