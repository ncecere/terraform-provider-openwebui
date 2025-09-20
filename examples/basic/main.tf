terraform {
  required_providers {
    openwebui = {
      source  = "nickcecere/openwebui"
      version = "~> 0.1"
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

  read_groups = ["Support"]
  write_groups = ["Support"]

  data_json = jsonencode({
    category = "support"
  })
}

resource "openwebui_model" "custom_rag" {
  model_id = "custom-rag"
  name     = "Custom Retrieval Model"

  meta_json = jsonencode({
    description = "Retriever tuned for internal knowledge base"
  })

  params_json = jsonencode({
    temperature = 0.1
  })

  base_model_id = "gpt-4o"
  is_active     = true
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
    "cj01@ufl.edu",
    "ebs@ufl.edu",
  ]

  data_json = jsonencode({
    department = "support"
  })

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
