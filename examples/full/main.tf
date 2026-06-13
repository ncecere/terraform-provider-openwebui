terraform {
  required_providers {
    openwebui = {
      source  = "nickcecere/openwebui"
      version = "~> 3.0"
    }
  }
}

provider "openwebui" {}

locals {
  suffix = var.name_suffix
}

resource "openwebui_folder" "workspace" {
  name = "Terraform Workspace ${local.suffix}"
}

resource "openwebui_prompt" "triage" {
  command        = "tf-triage-${local.suffix}"
  title          = "Terraform Triage ${local.suffix}"
  content        = "You are an assistant that triages inbound support requests."
  commit_message = "Initial Terraform-managed prompt"
  tags           = ["terraform", "example"]

  meta_json = jsonencode({
    description = "Managed by Terraform"
  })
}

data "openwebui_prompt_history" "triage" {
  command    = openwebui_prompt.triage.command
  depends_on = [openwebui_prompt.triage]
}

resource "openwebui_knowledge" "support" {
  name        = "Terraform Support ${local.suffix}"
  description = "Example knowledge base managed by Terraform"
}

resource "openwebui_knowledge_dir" "runbooks" {
  knowledge_id = openwebui_knowledge.support.id
  name         = "Runbooks"
}

resource "openwebui_file" "readme" {
  target_filename = "terraform-example-${local.suffix}.txt"
  content         = "This file was uploaded by Terraform."

  metadata_json = jsonencode({
    source = "terraform-provider-openwebui examples/full"
  })
}

resource "openwebui_knowledge_file" "readme" {
  knowledge_id = openwebui_knowledge.support.id
  file_id      = openwebui_file.readme.id
  directory_id = openwebui_knowledge_dir.runbooks.id
  process_file = false
}

resource "openwebui_function" "filter" {
  function_id = "tf_example_filter_${replace(local.suffix, "-", "_")}"
  name        = "Terraform Example Filter ${local.suffix}"
  description = "A minimal filter function managed by Terraform"
  is_active   = true

  content = <<-PY
    class Filter:
        def __init__(self):
            pass

        async def inlet(self, body: dict, __user__: dict = None) -> dict:
            return body

        async def outlet(self, body: dict, __user__: dict = None) -> dict:
            return body
  PY
}

data "openwebui_functions" "all" {
  depends_on = [openwebui_function.filter]
}

data "openwebui_models" "all" {}
data "openwebui_prompts" "all" {}
data "openwebui_tools" "all" {}

variable "name_suffix" {
  type        = string
  description = "Unique suffix to avoid name collisions in shared Open WebUI instances."
  default     = "example"
}

output "folder_id" {
  value = openwebui_folder.workspace.id
}

output "prompt_history_json" {
  value = data.openwebui_prompt_history.triage.history_json
}

output "knowledge_id" {
  value = openwebui_knowledge.support.id
}

output "function_id" {
  value = openwebui_function.filter.id
}
