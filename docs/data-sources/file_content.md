---
layout: data-source
page_title: "openwebui_file_content Data Source"
sidebar_current: docs-openwebui-datasource-file-content
description: |-
  Reads extracted Open WebUI file content.
---

# openwebui_file_content (Data Source)

Reads extracted data content for an Open WebUI file.

```hcl
data "openwebui_file_content" "example" {
  file_id = openwebui_file.example.id
}
```

## Argument Reference

* `file_id` (Required) – File ID.

## Attribute Reference

* `content_json` – JSON response from Open WebUI's file data content endpoint.
