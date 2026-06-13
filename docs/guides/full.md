---
page_title: "Full Open WebUI Provider Example"
description: |-
  Representative Terraform configuration for Open WebUI v0.9.6.
---

# Full Example

The full example shows the main v3.0.0 workflows without depending on external fixtures or MCP/tool servers:

- folders
- prompts and prompt history
- knowledge bases, knowledge directories, files, and knowledge-file attachments
- functions
- raw catalog data sources

```hcl
provider "openwebui" {}

resource "openwebui_folder" "workspace" {
  name = "Terraform Workspace"
}

resource "openwebui_prompt" "triage" {
  command        = "tf-triage"
  title          = "Terraform Triage"
  content        = "You are an assistant that triages inbound support requests."
  commit_message = "Initial Terraform-managed prompt"
  tags           = ["terraform", "example"]
}

data "openwebui_prompt_history" "triage" {
  command    = openwebui_prompt.triage.command
  depends_on = [openwebui_prompt.triage]
}

resource "openwebui_knowledge" "support" {
  name        = "Terraform Support"
  description = "Example knowledge base managed by Terraform"
}

resource "openwebui_knowledge_dir" "runbooks" {
  knowledge_id = openwebui_knowledge.support.id
  name         = "Runbooks"
}

resource "openwebui_file" "readme" {
  target_filename = "terraform-example.txt"
  content         = "This file was uploaded by Terraform."
}

resource "openwebui_knowledge_file" "readme" {
  knowledge_id = openwebui_knowledge.support.id
  file_id      = openwebui_file.readme.id
  directory_id = openwebui_knowledge_dir.runbooks.id
  process_file = false
}

resource "openwebui_function" "filter" {
  function_id = "tf_example_filter"
  name        = "Terraform Example Filter"
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

data "openwebui_models" "all" {}
data "openwebui_prompts" "all" {}
data "openwebui_tools" "all" {}
```

A runnable version lives in [`examples/full`](../../examples/full).
