---
layout: resource
page_title: "openwebui_knowledge_dir Resource"
sidebar_current: docs-openwebui-resource-knowledge-dir
description: |-
  Manages directories within Open WebUI knowledge bases.
---

# openwebui_knowledge_dir (Resource)

Creates and manages a directory inside an Open WebUI knowledge base.

```hcl
resource "openwebui_knowledge_dir" "docs" {
  knowledge_id = openwebui_knowledge.example.id
  name         = "Policies"
}
```

## Argument Reference

* `knowledge_id` (Required) – Knowledge base ID. Forces replacement.
* `name` (Required) – Directory name.
* `parent_id` (Optional) – Parent directory ID.
* `move_files_on_delete` (Optional) – Whether files are moved to the parent directory when the directory is deleted. Defaults to `true`.

## Attribute Reference

* `id` – Composite ID in the form `knowledge_id:directory_id`.
* `directory_id` – Open WebUI directory ID.
* `user_id` – Owning user ID.
* `created_at` – Creation timestamp.
* `updated_at` – Last update timestamp.

## Import

```bash
terraform import openwebui_knowledge_dir.docs knowledge_id:directory_id
```
