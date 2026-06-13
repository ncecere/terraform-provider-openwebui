terraform {
  required_providers {
    openwebui = {
      source  = "nickcecere/openwebui"
      version = "~> 3.0"
    }
  }
}

provider "openwebui" {
  # Prefer OPENWEBUI_ENDPOINT and OPENWEBUI_TOKEN for local use so secrets do not
  # end up in shell history or committed tfvars files.
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
