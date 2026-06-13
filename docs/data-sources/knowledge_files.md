---
layout: data-source
page_title: "openwebui_knowledge_files Data Source"
sidebar_current: docs-openwebui-datasource-knowledge-files
description: |-
  Lists files and directories in an Open WebUI knowledge base.
---

# openwebui_knowledge_files (Data Source)

Lists files for a knowledge base, optionally filtered by directory and query.

```hcl
data "openwebui_knowledge_files" "example" {
  knowledge_id = openwebui_knowledge.example.id
}
```

## Argument Reference

* `knowledge_id` (Required) – Knowledge base ID.
* `query` (Optional) – Search query.
* `directory_id` (Optional) – Directory filter. Use an empty string for root.
* `include_content` (Optional) – Include content in the response.
* `view_option` (Optional)
* `order_by` (Optional)
* `direction` (Optional)
* `page` (Optional) – Page number.

## Attribute Reference

* `files_json` – JSON array of matching files.
* `directories_json` – JSON array of directories.
* `breadcrumbs_json` – JSON array of breadcrumb directories.
* `total` – Total matching files.
