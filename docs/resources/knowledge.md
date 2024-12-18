# Knowledge Resource

Manages a knowledge base in OpenWebUI. Knowledge bases can be configured with different access controls and custom data.

## Example Usage

```hcl
# Create a public knowledge base
resource "openwebui_knowledge" "public_example" {
  name          = "Public Knowledge Base"
  description   = "This is a public knowledge base"
  access_control = "public"
  
  data = {
    source = "terraform"
    type   = "documentation"
  }
}

# Create a private knowledge base
resource "openwebui_knowledge" "private_example" {
  name          = "Private Knowledge Base"
  description   = "This is a private knowledge base"
  access_control = "private"
  
  data = {
    source = "terraform"
    type   = "internal"
  }
}
```

## Schema

### Required

- `name` (String) The name of the knowledge base.
- `description` (String) A description of the knowledge base.

### Optional

- `data` (Map of String) Additional data associated with the knowledge base. This can be used to store custom metadata.
- `access_control` (String) Access control setting for the knowledge base. Valid values are:
  - `"public"` (default) - The knowledge base is publicly accessible
  - `"private"` - The knowledge base is private and requires authentication

### Read-Only

- `id` (String) The unique identifier of the knowledge base.
- `last_updated` (String) Timestamp of when the knowledge base was last updated.

## Import

Knowledge bases can be imported using their ID:

```shell
terraform import openwebui_knowledge.example <knowledge-base-id>
```

## Implementation Details

The knowledge resource is implemented using a modular client architecture:

1. Client Layer (`internal/provider/client/knowledge/`):
   - Handles API communication
   - Implements knowledge-specific operations
   - Manages data serialization/deserialization

2. Resource Layer (`internal/provider/knowledge_resource.go`):
   - Implements Terraform resource interface
   - Manages resource lifecycle (CRUD operations)
   - Handles state management

This modular approach ensures:
- Clean separation of concerns
- Easy maintenance and updates
- Consistent error handling
- Type safety through strong typing
