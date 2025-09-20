terraform {
  required_providers {
    openwebui = {
      source  = "nickcecere/openwebui"
      version = "~> 2.0"
    }
  }
}

provider "openwebui" {
  endpoint = var.openwebui_endpoint
  token    = var.openwebui_token
}

resource "openwebui_knowledge" "support_faq" {
  name        = "Support FAQ"
  description = "Knowledge base backing the support chatbot"

  read_groups  = ["Support"]
  write_groups = ["Support"]
}

resource "openwebui_model" "custom_rag" {
  model_id = "custom-rag"
  name     = "Custom Retrieval Model"

  description         = "Retriever tuned for internal knowledge base"
  base_model_id       = "gpt-4o"
  is_active           = true
  read_groups         = ["Support"]
  write_groups        = ["Support"]
  default_feature_ids = ["web_search"]

  params = {
    temperature     = 0.1
    max_tokens      = 512
    stream_response = true
  }

  capabilities = {
    web_search = true
  }
}

resource "openwebui_prompt" "triage" {
  command = "triage"
  title   = "Ticket triage"
  content = <<-EOT
    You are an assistant that triages inbound support tickets.
  EOT

  read_groups  = ["Support"]
  write_groups = ["Support"]
}

resource "openwebui_group" "support" {
  name        = "Support"
  description = "Support team access"

  users = [
    "jim@school.edu",
    "bob@school.edu",
  ]

  permissions = {
    workspace = {
      models    = true
      knowledge = true
      prompts   = true
      tools     = true
    }

    sharing = {
      public_models = false
    }

    chat = {
      file_upload         = true
      delete              = true
      edit                = true
      continue_response   = true
      regenerate_response = true
      temporary           = true
    }

    features = {
      web_search       = true
      image_generation = true
    }
  }
}

variable "openwebui_endpoint" {
  type        = string
  description = "Base URL for the Open WebUI API"
  default     = "http://localhost:3000/api/v1"
}

variable "openwebui_token" {
  type        = string
  description = "API token for Open WebUI"
  sensitive   = true
}
