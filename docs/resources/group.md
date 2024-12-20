# Group Resource

Manages a group in OpenWebUI. Groups allow you to organize users and control their access to various features within OpenWebUI through a comprehensive permissions system.

## Example Usage

### Basic Group

```hcl
resource "openwebui_group" "basic" {
  name        = "basic-users"
  description = "Basic user group with limited permissions"
  user_ids    = ["user1", "user2"]

  permissions = {
    workspace = {
      models    = false
      knowledge = true
      prompts   = true
      tools     = false
    }
    chat = {
      file_upload = true
      delete      = false
      edit        = true
      temporary   = true
    }
  }
}
```

### Admin Group

```hcl
resource "openwebui_group" "admin" {
  name        = "administrators"
  description = "Administrative group with full permissions"

  permissions = {
    workspace = {
      models    = true
      knowledge = true
      prompts   = true
      tools     = true
    }
    chat = {
      file_upload = true
      delete      = true
      edit        = true
      temporary   = true
    }
  }
}
```

### Custom Permissions

```hcl
resource "openwebui_group" "knowledge_managers" {
  name        = "knowledge-managers"
  description = "Group for managing knowledge bases"

  permissions = {
    workspace = {
      models    = false
      knowledge = true  # Only knowledge base access
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

## Schema

### Required

- `name` (String) The name of the group.
- `permissions` (Block) Group permissions configuration. Must include both `workspace` and `chat` blocks.

### Optional

- `description` (String) A description of the group.
- `user_ids` (List of String) List of user IDs to include in the group.

### Permissions Configuration

The `permissions` block consists of two required nested blocks:

#### Workspace Permissions (`workspace` Block)

Controls access to different workspace features:

- `models` (Bool, Required) - When true, allows access to:
  - View and select models
  - Modify model parameters
  - Create custom models

- `knowledge` (Bool, Required) - When true, allows access to:
  - View knowledge bases
  - Create and modify knowledge bases
  - Upload documents

- `prompts` (Bool, Required) - When true, allows access to:
  - View saved prompts
  - Create and modify prompts
  - Share prompts with others

- `tools` (Bool, Required) - When true, allows access to:
  - View available tools
  - Configure tool settings
  - Create custom tools

#### Chat Permissions (`chat` Block)

Controls chat-related features:

- `file_upload` (Bool, Required) - When true, allows:
  - Uploading files in chat
  - Sharing files with others

- `delete` (Bool, Required) - When true, allows:
  - Deleting chat messages
  - Clearing chat history

- `edit` (Bool, Required) - When true, allows:
  - Editing sent messages
  - Modifying chat settings

- `temporary` (Bool, Required) - When true, allows:
  - Creating temporary chats
  - Using ephemeral messages

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` (String) The unique identifier of the group.

## Import

Groups can be imported using their ID:

```bash
terraform import openwebui_group.example <group-id>
```

## Implementation Notes

- Groups are created with basic information first, then updated with full permissions and user assignments
- Changes to permissions take effect immediately for all group members
- Removing a user from a group immediately revokes their group-based permissions
- Group names must be unique within an OpenWebUI instance
