---
layout: data-source
page_title: "openwebui_file_process_status Data Source"
sidebar_current: docs-openwebui-datasource-file-process-status
description: |-
  Reads Open WebUI file processing status.
---

# openwebui_file_process_status (Data Source)

Reads processing status for an Open WebUI file.

```hcl
data "openwebui_file_process_status" "example" {
  file_id = openwebui_file.example.id
}
```

## Argument Reference

* `file_id` (Required) – File ID.

## Attribute Reference

* `status_json` – JSON response from Open WebUI's file process status endpoint.
