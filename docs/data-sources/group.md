# Group Data Source

Use this data source to retrieve information about an existing OpenWebUI group. This is useful for referencing group details in other resources or for querying group configurations.

## Example Usage

### Basic Group Lookup

```hcl
data "openwebui_group" "example" {
  name = "example-group"
}

output "group_details" {
  value = {
    id          = data.openwebui_group.example.id
    description = data.openwebui_group.example.description
    user_count  = length(data.openwebui_group.example.user_ids)
  }
}
```

### Using Group Permissions in Conditions

```hcl
data "openwebui_group" "admin" {
  name = "administrators"
}

resource "openwebui_model" "restricted" {
  name = "Restricted Model"
  # ... other model configuration ...

  # Only allow access if the group has model permissions
  count = data.openwebui_group.admin.permissions.workspace.models ? 1 : 0
}
```

### Referencing Group Users

```hcl
data "openwebui_group" "team" {
  name = "development-team"
}

resource "openwebui_knowledge" "team_docs" {
  name        = "Team Documentation"
  description = "Documentation for ${data.openwebui_group.team.name}"
  
  access_control = {
    read = {
      user_ids = data.openwebui_group.team.user_ids
    }
  }
}
```

## Schema

### Required

- `name` (String) The name of the group to look up.

### Read-Only

- `id` (String) The unique identifier of the group.
- `description` (String) The description of the group.
- `user_ids` (List of String) List of user IDs that are members of the group.
- `created_at` (Number) Unix timestamp when the group was created.
- `updated_at` (Number) Unix timestamp when the group was last updated.

### Permissions Structure

The `permissions` block (Read-only) contains detailed group permissions:

#### Workspace Permissions (`permissions.workspace`)

Controls access to workspace features:

- `models` (Bool) - Access to model management
- `knowledge` (Bool) - Access to knowledge base management
- `prompts` (Bool) - Access to prompt management
- `tools` (Bool) - Access to tools management

#### Chat Permissions (`permissions.chat`)

Controls chat feature access:

- `file_upload` (Bool) - Ability to upload files in chat
- `delete` (Bool) - Ability to delete messages
- `edit` (Bool) - Ability to edit messages
- `temporary` (Bool) - Access to temporary chat features

## Using Group Data

### In Outputs

```hcl
output "group_permissions" {
  value = {
    can_manage_models     = data.openwebui_group.example.permissions.workspace.models
    can_manage_knowledge = data.openwebui_group.example.permissions.workspace.knowledge
    chat_capabilities    = data.openwebui_group.example.permissions.chat
  }
}
```

### In Resource Configurations

```hcl
resource "openwebui_model" "shared" {
  name = "Shared Model"
  
  # Use group data in access control
  access_control = {
    read = {
      group_ids = [data.openwebui_group.example.id]
    }
  }
}
```

### For Dynamic Configurations

```hcl
locals {
  is_admin_group = data.openwebui_group.example.permissions.workspace.models && 
                  data.openwebui_group.example.permissions.workspace.knowledge &&
                  data.openwebui_group.example.permissions.workspace.tools
}

# Use the information for conditional resource creation
resource "openwebui_knowledge" "admin_docs" {
  count = local.is_admin_group ? 1 : 0
  
  name        = "Administrative Documentation"
  description = "Documentation for administrators"
  # ... other configuration ...
}
