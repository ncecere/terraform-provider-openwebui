# Configure the OpenWebUI Provider
provider "openwebui" {
  endpoint = "https://chat.example.com"  # Replace with your OpenWebUI endpoint
  token    = "your-api-token"           # Replace with your API token
}

# Create a model with custom configuration
resource "openwebui_model" "custom_gpt4" {
  base_model_id = "gpt-4"
  name          = "Custom GPT-4"
  is_active     = true

  params {
    system          = "You are a helpful AI assistant specialized in infrastructure and DevOps."
    stream_response = true
    temperature    = 0.7
    top_p          = 0.9
    max_tokens     = 2000
    seed           = 42
  }

  meta {
    description       = "A customized GPT-4 model for infrastructure tasks"
    profile_image_url = "/static/favicon.png"
    
    capabilities {
      vision    = false
      usage     = true
      citations = true
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
      group_ids = ["devops-team"]
      user_ids  = []
    }
    write {
      group_ids = ["admin-team"]
      user_ids  = []
    }
  }
}

# Use the model data source to read model information
data "openwebui_model" "custom_gpt4" {
  name = openwebui_model.custom_gpt4.name
}

# Output model information
output "model_id" {
  value = data.openwebui_model.custom_gpt4.id
}

output "model_capabilities" {
  value = data.openwebui_model.custom_gpt4.meta.capabilities
}

output "model_tags" {
  value = data.openwebui_model.custom_gpt4.meta.tags
}
