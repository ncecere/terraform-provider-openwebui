# Knowledge Data Source

Use this data source to retrieve information about an existing knowledge base in OpenWebUI.

## Example Usage

```hcl
# Look up a knowledge base by name
data "openwebui_knowledge" "example" {
  name = "Example Knowledge Base"
}

# Use the data
output "knowledge_id" {
  value = data.openwebui_knowledge.example.id
}

output "knowledge_data" {
  value = data.openwebui_knowledge.example.data
}

# Example with resource reference
resource "openwebui_knowledge" "source" {
  name          = "Source Knowledge Base"
  description   = "This is a source knowledge base"
  access_control = "public"
  
  data = {
    type = "source"
  }
}

data "openwebui_knowledge" "lookup" {
  name = openwebui_knowledge.source.name
}
```

## Schema

### Required

- `name` (String) The name of the knowledge base to look up.

### Read-Only

- `id` (String) The unique identifier of the knowledge base.
- `description` (String) The description of the knowledge base.
- `data` (Map of String) Additional data associated with the knowledge base.
- `access_control` (String) Access control setting of the knowledge base:
  - `"public"` - The knowledge base is publicly accessible
  - `"private"` - The knowledge base is private and requires authentication
- `access_groups` (List of String) List of group IDs that have access to the knowledge base (only applicable for private knowledge bases).
- `access_users` (List of String) List of user IDs that have access to the knowledge base (only applicable for private knowledge bases).
- `last_updated` (String) Timestamp of when the knowledge base was last updated.

## Implementation Details

The knowledge data source is implemented using a modular client architecture:

1. Client Layer (`internal/provider/client/knowledge/`):
   - Handles API communication
   - Implements knowledge-specific operations
   - Manages data serialization/deserialization

2. Data Source Layer (`internal/provider/knowledge_data_source.go`):
   - Implements Terraform data source interface
   - Handles data retrieval and filtering
   - Manages state reading

Benefits of this implementation:
- Consistent interface with the knowledge resource
- Efficient data retrieval
- Type-safe operations
- Clear error handling and reporting
