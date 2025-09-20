---
layout: resource
page_title: "openwebui_knowledge Resource"
sidebar_current: docs-openwebui-resource-knowledge
description: |-
  Manages knowledge base entries in Open WebUI.
---

# openwebui_knowledge (Resource)

Creates and manages knowledge base entries inside Open WebUI.

## Example Usage

```hcl
resource "openwebui_knowledge" "support_faq" {
  name        = "Support FAQ"
  description = "Knowledge base backing the support chatbot"

  read_groups = [
    "Support",
  ]

  write_groups = [
    "Support",
  ]

  data_json = jsonencode({
    category = "support"
  })
}
```

If both `read_groups` and `write_groups` are omitted (or empty), the knowledge entry remains public.

## Argument Reference

* `name` (Required) – Human readable name of the knowledge entry.
* `description` (Required) – Description shown in Open WebUI.
* `read_groups` (Optional) – List of group names or IDs granted read access. Leave unset (or empty) for public knowledge.
* `write_groups` (Optional) – List of group names or IDs granted write access. Groups here automatically receive read access.
* `data_json` (Optional) – JSON object string for additional metadata sent during create/update.
* `meta_json` (Optional) – JSON object string persisted in the knowledge entry metadata. The API may enrich this field and it is surfaced in state.

## Attribute Reference

* `id` – Unique knowledge identifier assigned by Open WebUI.
* `created_at` – Creation date in `YYYY-MM-DD` format.
* `updated_at` – Last update date in `YYYY-MM-DD` format.
* `meta_json` – JSON metadata returned by the API.
* `user_id` – Identifier of the user that owns the knowledge entry.
* `read_groups` / `write_groups` – Resolved group names currently applied to the entry.

## Import

Knowledge entries can be imported using the knowledge ID, for example:

```bash
terraform import openwebui_knowledge.support_faq 65e5e86e-0e23-4cd8-8eee-447c6923f632
```
