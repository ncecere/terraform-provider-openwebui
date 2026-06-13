---
layout: resource
page_title: "openwebui_knowledge_file Resource"
sidebar_current: docs-openwebui-resource-knowledge-file
description: |-
  Manages file attachments for Open WebUI knowledge bases.
---

# openwebui_knowledge_file (Resource)

Attaches a file to a knowledge base in Open WebUI.

## Example Usage

### Minimal

```hcl
resource "openwebui_knowledge_file" "faq" {
  knowledge_id = openwebui_knowledge.support_faq.id
  file_id      = openwebui_file.support_doc.id
}
```

### Full

```hcl
resource "openwebui_knowledge_file" "faq" {
  knowledge_id = openwebui_knowledge.support_faq.id
  file_id      = openwebui_file.support_doc.id
  directory_id = openwebui_knowledge_dir.docs.directory_id
  delete_file  = true
}
```

## Argument Reference

* `knowledge_id` (Required) – Knowledge base identifier.
* `file_id` (Required) – File identifier to attach.
* `directory_id` (Optional) – Knowledge directory ID to place the file in. Updating this moves the attachment.
* `delete_file` (Optional) – Whether the file should be deleted when detached (defaults to `true`).

## Attribute Reference

* `id` – Composite identifier in the form `knowledge_id:file_id`.
* `file_json` – Raw JSON describing the attached file.

## Import

Knowledge file attachments can be imported using the composite ID:

```bash
terraform import openwebui_knowledge_file.faq 65e5e86e-0e23-4cd8-8eee-447c6923f632:fe12b9e3-9f49-40c4-8c6d-1cb8c7f4d9a0
```
