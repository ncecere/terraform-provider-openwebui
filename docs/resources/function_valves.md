---
layout: resource
page_title: "openwebui_function_valves Resource"
sidebar_current: docs-openwebui-resource-function-valves
description: |-
  Manages Open WebUI function valves.
---

# openwebui_function_valves (Resource)

Manages valve settings for a function.

```hcl
resource "openwebui_function_valves" "example" {
  function_id = openwebui_function.example.id
  valves_json = jsonencode({
    enabled = true
  })
}
```
