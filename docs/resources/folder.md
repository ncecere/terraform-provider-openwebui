---
layout: resource
page_title: "openwebui_folder Resource"
sidebar_current: docs-openwebui-resource-folder
description: |-
  Manages Open WebUI folders.
---

# openwebui_folder (Resource)

Manages an Open WebUI folder.

## Example Usage

```hcl
resource "openwebui_folder" "example" {
  name        = "Terraform Folder"
  is_expanded = true
  meta_json = jsonencode({
    icon = "folder"
  })
}
```

## Argument Reference

* `name` (Required) – Folder name.
* `parent_id` (Optional) – Parent folder ID.
* `data_json` (Optional) – Folder data JSON.
* `meta_json` (Optional) – Folder metadata JSON.
* `is_expanded` (Optional) – Whether the folder is expanded.
* `delete_contents` (Optional) – Whether to delete folder contents on destroy. Defaults to `true`.

## Attribute Reference

* `id` – Folder ID.
* `user_id` – Owning user ID.
* `created_at` – Creation timestamp.
* `updated_at` – Last update timestamp.
