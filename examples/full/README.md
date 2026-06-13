# Full Example

This example demonstrates a representative Open WebUI v0.9.6 workflow:

- folder creation
- prompt creation and prompt history readback
- knowledge base, knowledge directory, file upload, and knowledge-file attachment
- function creation
- raw catalog data sources

Run it against a non-production or test Open WebUI instance first:

```bash
export OPENWEBUI_ENDPOINT="https://your-openwebui.example.com/api/v1"
export OPENWEBUI_TOKEN="..."
terraform init
terraform apply -var='name_suffix=demo'
```

Destroy the example resources when done:

```bash
terraform destroy -var='name_suffix=demo'
```
