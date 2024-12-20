# Configure the OpenWebUI Provider
provider "openwebui" {
  endpoint = "https://chat.example.com"  # Replace with your OpenWebUI endpoint
  # token = "your-api-token"            # Token can be provided via OPENWEBUI_TOKEN env var
}

# Create a group for model access control
resource "openwebui_group" "ml_team" {
  name        = "ml-team"
  description = "Machine Learning team with model management permissions"
  
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

# Example 1: DevOps-focused GPT-4 Model
resource "openwebui_model" "devops_gpt4" {
  base_model_id = "gpt-4"
  name          = "DevOps GPT-4"
  is_active     = true

  params {
    # Specialized system prompt for DevOps tasks
    system          = <<-EOT
      You are a DevOps expert specialized in:
      - Infrastructure as Code
      - CI/CD pipelines
      - Cloud architecture
      - Container orchestration
      - Monitoring and logging
      Provide practical, security-conscious advice.
    EOT
    stream_response = true
    temperature    = 0.7    # Balanced creativity and consistency
    top_p          = 0.9    # Nucleus sampling for diverse responses
    max_tokens     = 2000   # Longer responses for detailed explanations
    seed           = 42     # Fixed seed for reproducibility
  }

  meta {
    description       = "A customized GPT-4 model optimized for DevOps and infrastructure tasks"
    profile_image_url = "/static/devops-icon.png"
    
    capabilities {
      vision    = false     # No vision capabilities needed
      usage     = true      # Track usage statistics
      citations = true      # Enable source citations
    }

    tags {
      name = "infrastructure"
    }
    tags {
      name = "devops"
    }
  }

  access_control {
    read {
      group_ids = [openwebui_group.ml_team.id]
    }
    write {
      group_ids = [openwebui_group.ml_team.id]
    }
  }
}

# Example 2: Research Assistant Model
resource "openwebui_model" "research_assistant" {
  base_model_id = "gpt-4"
  name          = "Research Assistant"
  is_active     = true

  params {
    # Academic research-focused system prompt
    system          = <<-EOT
      You are a research assistant specialized in:
      - Academic paper analysis
      - Literature review
      - Research methodology
      - Data analysis
      - Scientific writing
      Provide well-referenced, academically rigorous responses.
    EOT
    stream_response = false   # Complete responses preferred
    temperature    = 0.3     # Lower temperature for more focused responses
    top_p          = 0.95    # High precision in academic context
    max_tokens     = 4000    # Extended responses for detailed analysis
    num_ctx        = 8192    # Larger context window for research papers
  }

  meta {
    description = "Academic research assistant model with emphasis on scientific rigor"
    
    capabilities {
      vision    = true      # Enable vision for analyzing graphs and figures
      usage     = true      # Track usage for research groups
      citations = true      # Critical for academic work
    }

    tags {
      name = "academic"
    }
    tags {
      name = "research"
    }
  }

  access_control {
    read {
      group_ids = [openwebui_group.ml_team.id]
    }
    write {
      group_ids = [openwebui_group.ml_team.id]
    }
  }
}

# Example 3: Code Review Model
resource "openwebui_model" "code_reviewer" {
  base_model_id = "gpt-4"
  name          = "Code Review Assistant"
  is_active     = true

  params {
    system          = <<-EOT
      You are a code review assistant specialized in:
      - Code quality assessment
      - Security vulnerability detection
      - Performance optimization suggestions
      - Best practices enforcement
      - Documentation improvements
      Focus on providing actionable, specific feedback.
    EOT
    stream_response = true
    temperature    = 0.2    # Low temperature for consistent code analysis
    top_p          = 0.8
    max_tokens     = 1500
    frequency_penalty = 1    # Reduce repetitive suggestions
  }

  meta {
    description = "Specialized model for code review and analysis"
    
    capabilities {
      vision    = false
      usage     = true
      citations = true
    }

    tags {
      name = "code-review"
    }
    tags {
      name = "development"
    }
  }

  access_control {
    read {
      group_ids = [openwebui_group.ml_team.id]
    }
    write {
      group_ids = [openwebui_group.ml_team.id]
    }
  }
}

# Create a knowledge base for model documentation
resource "openwebui_knowledge" "model_docs" {
  name        = "Model Documentation"
  description = "Documentation and best practices for custom models"
  
  data = {
    category = "technical-documentation"
    version  = "1.0"
  }
}

# Use data sources to query models
data "openwebui_model" "devops" {
  name = openwebui_model.devops_gpt4.name
}

# Outputs for verification and reference
output "model_ids" {
  value = {
    devops    = openwebui_model.devops_gpt4.id
    research  = openwebui_model.research_assistant.id
    code      = openwebui_model.code_reviewer.id
  }
}

output "model_capabilities" {
  value = {
    devops_capabilities   = data.openwebui_model.devops.meta.capabilities
    research_capabilities = openwebui_model.research_assistant.meta.capabilities
    code_capabilities    = openwebui_model.code_reviewer.meta.capabilities
  }
}

output "model_access" {
  value = {
    ml_team_id = openwebui_group.ml_team.id
    docs_id    = openwebui_knowledge.model_docs.id
  }
}
