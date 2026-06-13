---
layout: resource
page_title: "openwebui_function Resource"
sidebar_current: docs-openwebui-resource-function
description: |-
  Manages Open WebUI functions.
---

# openwebui_function (Resource)

Creates and manages an Open WebUI function.

```hcl
resource "openwebui_function" "example" {
  function_id = "terraform_filter"
  name        = "Terraform Filter"
  description = "Managed by Terraform"
  content     = file("./filter.py")
  is_active   = true
}
```

## Argument Reference

* `function_id` (Required) – Function identifier. Forces replacement.
* `name` (Required) – Display name.
* `content` (Required) – Python function source.
* `description` (Optional) – Description stored in function metadata.
* `manifest_json` (Optional) – Function manifest JSON.
* `is_active` (Optional) – Whether the function is active.
* `is_global` (Optional) – Whether the function is global.

## Attribute Reference

* `id` – Function ID.
* `type` – Function type detected by Open WebUI.
* `user_id`, `created_at`, `updated_at` – Server metadata.
