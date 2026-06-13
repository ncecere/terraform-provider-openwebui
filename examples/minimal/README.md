# Minimal Example

This is the smallest practical configuration for the provider. It creates one prompt and reads it back.

```bash
export OPENWEBUI_ENDPOINT="https://your-openwebui.example.com/api/v1"
export OPENWEBUI_TOKEN="..."
terraform init
terraform apply
```

The token should be an Open WebUI API token with permission to manage prompts.
