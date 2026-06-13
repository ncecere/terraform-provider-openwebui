---
layout: data-source
page_title: "openwebui_folder Data Source"
sidebar_current: docs-openwebui-datasource-folder
description: |-
  Reads an Open WebUI folder.
---

# openwebui_folder (Data Source)

Reads an Open WebUI folder by ID.

```hcl
data "openwebui_folder" "example" {
  id = openwebui_folder.example.id
}
```
