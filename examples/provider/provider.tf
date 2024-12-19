terraform {
  required_providers {
    openwebui = {
      source = "ncecere/openwebui"
    }
  }
}

# Configure the OpenWebUI Provider
# Authentication can be provided in multiple ways:
# 1. Directly in the provider block (not recommended for production)
# 2. Using environment variables:
#    - OPENWEBUI_ENDPOINT
#    - OPENWEBUI_TOKEN
provider "openwebui" {
  endpoint = "http://localhost:8080"  # Optional: can be set via OPENWEBUI_ENDPOINT
  # token = "your-api-token"         # Optional: can be set via OPENWEBUI_TOKEN
}

# Create a group for managing access
resource "openwebui_group" "data_science" {
  name        = "data-science"
  description = "Data Science team group"
  user_ids    = ["user1", "user2"]

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

# Create a model with specific configuration
resource "openwebui_model" "gpt4_custom" {
  name          = "GPT-4 Custom"
  base_model_id = "gpt-4"
  is_active     = true

  params {
    system          = "You are a helpful data science assistant"
    temperature     = 0.7
    top_p           = 0.9
    max_tokens      = 2000
  }

  meta {
    description = "Customized GPT-4 for data science tasks"
    capabilities {
      vision    = false
      usage     = true
      citations = true
    }
  }

  # Grant access to the data science group
  access_control {
    read {
      group_ids = [openwebui_group.data_science.id]
    }
    write {
      group_ids = [openwebui_group.data_science.id]
    }
  }
}

# Create a knowledge base
resource "openwebui_knowledge" "documentation" {
  name        = "Data Science Documentation"
  description = "Knowledge base for data science best practices"
  
  data = {
    department = "Data Science"
    category   = "Documentation"
    version    = "1.0"
  }

  access_control = "private"  # Make it private by default
}

# Use data sources to query existing resources
data "openwebui_group" "existing_admin" {
  name = "administrators"  # Look up existing admin group
}

data "openwebui_model" "base_model" {
  name = "GPT-4"  # Look up existing model
}

# Output some useful information
output "data_science_group_id" {
  value = openwebui_group.data_science.id
}

output "model_status" {
  value = {
    id        = openwebui_model.gpt4_custom.id
    is_active = openwebui_model.gpt4_custom.is_active
  }
}

output "knowledge_base_info" {
  value = {
    id           = openwebui_knowledge.documentation.id
    last_updated = openwebui_knowledge.documentation.last_updated
  }
}
