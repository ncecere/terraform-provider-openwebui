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

# Create groups for access control
resource "openwebui_group" "developers" {
  name        = "developers"
  description = "Development team with access to technical documentation"
  
  permissions = {
    workspace = {
      models    = false
      knowledge = true  # Can access knowledge bases
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

resource "openwebui_group" "researchers" {
  name        = "researchers"
  description = "Research team with access to research documentation"
  
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

# Example 1: Technical Documentation Knowledge Base
resource "openwebui_knowledge" "tech_docs" {
  name           = "Technical Documentation"
  description    = "Comprehensive technical documentation for our systems"
  access_control = "private"  # Restricted access

  data = {
    category     = "technical"
    department   = "engineering"
    version      = "1.0"
    last_review  = "2024-01-01"
    maintainer   = "devops-team"
    format       = "markdown"
    source_repo  = "github.com/company/tech-docs"
  }

  # Associated with a model for better context
  depends_on = [openwebui_model.documentation_assistant]
}

# Example 2: Research Papers Knowledge Base
resource "openwebui_knowledge" "research_papers" {
  name           = "Research Papers"
  description    = "Collection of research papers and findings"
  access_control = "private"

  data = {
    category     = "research"
    department   = "r&d"
    type         = "academic"
    field        = "machine-learning"
    curator      = "research-lead"
    updated      = "2024-02-01"
    citation_format = "IEEE"
  }
}

# Example 3: Public Documentation
resource "openwebui_knowledge" "public_docs" {
  name           = "Public API Documentation"
  description    = "Public-facing API documentation and guides"
  access_control = "public"

  data = {
    category     = "api-documentation"
    version      = "2.0"
    status       = "published"
    target       = "external-developers"
    format       = "openapi"
    base_url     = "api.example.com/v2"
  }
}

# Example 4: Training Materials
resource "openwebui_knowledge" "training" {
  name           = "Employee Training Materials"
  description    = "Internal training documentation and resources"
  access_control = "private"

  data = {
    category     = "training"
    department   = "hr"
    type         = "onboarding"
    format       = "mixed"
    required     = "true"
    validity     = "2024"
    compliance   = "required"
  }
}

# Create a specialized model for documentation
resource "openwebui_model" "documentation_assistant" {
  name          = "Documentation Assistant"
  base_model_id = "gpt-4"
  is_active     = true

  params {
    system = <<-EOT
      You are a documentation specialist. Help users understand and navigate
      technical documentation, research papers, and training materials.
      Provide clear, concise explanations and relevant examples.
    EOT
    temperature = 0.3
  }

  meta {
    description = "Specialized model for documentation assistance"
    capabilities = {
      vision    = false
      usage     = true
      citations = true
    }
  }
}

# Data source examples
data "openwebui_knowledge" "tech_docs" {
  name = openwebui_knowledge.tech_docs.name
}

data "openwebui_knowledge" "research" {
  name = openwebui_knowledge.research_papers.name
}

# Outputs for verification and reference
output "knowledge_base_ids" {
  value = {
    tech_docs      = openwebui_knowledge.tech_docs.id
    research       = openwebui_knowledge.research_papers.id
    public_docs    = openwebui_knowledge.public_docs.id
    training       = openwebui_knowledge.training.id
  }
}

output "knowledge_metadata" {
  value = {
    tech_docs_data = data.openwebui_knowledge.tech_docs.data
    research_data  = data.openwebui_knowledge.research.data
  }
}

output "access_control" {
  value = {
    developers_group = openwebui_group.developers.id
    researchers_group = openwebui_group.researchers.id
  }
}

# Example of using knowledge base with a model
output "documentation_setup" {
  value = {
    model_id = openwebui_model.documentation_assistant.id
    knowledge_bases = [
      openwebui_knowledge.tech_docs.id,
      openwebui_knowledge.public_docs.id
    ]
  }
}
