# Configure the OpenWebUI Provider
terraform {
  required_providers {
    openwebui = {
      source = "ncecere/openwebui"
    }
  }
}

provider "openwebui" {
  # Configuration options - can be provided by environment variables:
  # endpoint = "http://your-openwebui-instance"  # OPENWEBUI_ENDPOINT
  # token    = "your-api-token"                  # OPENWEBUI_TOKEN
}

# Create a public knowledge base
resource "openwebui_knowledge" "public_example" {
  name           = "Public Knowledge Base"
  description    = "This is a public knowledge base managed by Terraform"
  access_control = "public"

  data = {
    source      = "terraform"
    type        = "documentation"
    environment = "example"
    managed_by  = "terraform"
  }
}

# Create a private knowledge base
resource "openwebui_knowledge" "private_example" {
  name           = "Private Knowledge Base"
  description    = "This is a private knowledge base managed by Terraform"
  access_control = "private"

  data = {
    source      = "terraform"
    type        = "internal"
    environment = "example"
    managed_by  = "terraform"
    visibility  = "private"
  }
}

# Look up the public knowledge base
data "openwebui_knowledge" "public_lookup" {
  name = openwebui_knowledge.public_example.name

  depends_on = [openwebui_knowledge.public_example]
}

# Look up the private knowledge base
data "openwebui_knowledge" "private_lookup" {
  name = openwebui_knowledge.private_example.name

  depends_on = [openwebui_knowledge.private_example]
}

# Output the resource IDs
output "public_knowledge_id" {
  value = openwebui_knowledge.public_example.id
}

output "private_knowledge_id" {
  value = openwebui_knowledge.private_example.id
}

# Output the lookup results
output "public_lookup_data" {
  value = {
    id             = data.openwebui_knowledge.public_lookup.id
    name           = data.openwebui_knowledge.public_lookup.name
    description    = data.openwebui_knowledge.public_lookup.description
    access_control = data.openwebui_knowledge.public_lookup.access_control
    data          = data.openwebui_knowledge.public_lookup.data
  }
}

output "private_lookup_data" {
  value = {
    id             = data.openwebui_knowledge.private_lookup.id
    name           = data.openwebui_knowledge.private_lookup.name
    description    = data.openwebui_knowledge.private_lookup.description
    access_control = data.openwebui_knowledge.private_lookup.access_control
    data          = data.openwebui_knowledge.private_lookup.data
  }
}
