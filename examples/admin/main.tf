terraform {
  required_providers {
    openwebui = {
      source  = "nickcecere/openwebui"
      version = "~> 3.0"
    }
  }
}

provider "openwebui" {
  endpoint = var.openwebui_endpoint
  token    = var.openwebui_token
}

data "openwebui_config_export" "current" {}

resource "openwebui_config_import" "restore" {
  count       = var.apply_config_import ? 1 : 0
  config_json = data.openwebui_config_export.current.config_json
}

resource "openwebui_oauth_client" "provider" {
  count = var.oauth_client_url != "" && var.oauth_client_id != "" ? 1 : 0

  url       = var.oauth_client_url
  client_id = var.oauth_client_id

  client_name = var.oauth_client_name != "" ? var.oauth_client_name : null
  type        = var.oauth_client_type != "" ? var.oauth_client_type : null
}

data "openwebui_tool_server_verify" "server" {
  count = var.tool_server_url != "" ? 1 : 0

  url  = var.tool_server_url
  path = var.tool_server_path
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

variable "apply_config_import" {
  type        = bool
  description = "Whether to apply the config import resource."
  default     = false
}

variable "oauth_client_url" {
  type        = string
  description = "OAuth provider URL to register."
  default     = ""
}

variable "oauth_client_id" {
  type        = string
  description = "OAuth client identifier to register."
  default     = ""
}

variable "oauth_client_name" {
  type        = string
  description = "Optional OAuth client display name."
  default     = ""
}

variable "oauth_client_type" {
  type        = string
  description = "Optional OAuth client type query parameter."
  default     = ""
}

variable "tool_server_url" {
  type        = string
  description = "Tool server URL to verify."
  default     = ""
}

variable "tool_server_path" {
  type        = string
  description = "Tool server path to verify."
  default     = "/"
}
