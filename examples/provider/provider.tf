terraform {
  required_providers {
    openwebui = {
      source = "ncecere/openwebui"
    }
  }
}

provider "openwebui" {
  endpoint = "http://localhost:8080"  # OpenWebUI API endpoint
  # token = "your-api-token"         # Token can be provided here or via OPENWEBUI_TOKEN env var
}

# Example knowledge resource
resource "openwebui_knowledge" "example" {
  name        = "Example Knowledge Base"
  description = "This is an example knowledge base"
}
