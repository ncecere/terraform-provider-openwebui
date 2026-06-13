---
page_title: "Minimal Open WebUI Provider Example"
description: |-
  Minimal Terraform configuration for Open WebUI.
---

# Minimal Example

Use this when you want to verify provider authentication and manage one simple object.

```hcl
terraform {
  required_providers {
    openwebui = {
      source  = "nickcecere/openwebui"
      version = "~> 3.0"
    }
  }
}

provider "openwebui" {
  # Uses OPENWEBUI_ENDPOINT and OPENWEBUI_TOKEN when omitted.
}

resource "openwebui_prompt" "hello" {
  command = "hello-terraform"
  title   = "Hello Terraform"
  content = "You are a concise assistant managed by Terraform."
}

data "openwebui_prompt" "hello" {
  command = openwebui_prompt.hello.command
}

output "prompt_command" {
  value = data.openwebui_prompt.hello.command
}
```

Run it with:

```bash
export OPENWEBUI_ENDPOINT="https://your-openwebui.example.com/api/v1"
export OPENWEBUI_TOKEN="..."
terraform init
terraform apply
```
