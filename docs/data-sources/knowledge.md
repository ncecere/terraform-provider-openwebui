---
layout: data-source
page_title: "openwebui_knowledge Data Source"
sidebar_current: docs-openwebui-data-source-knowledge
description: |-
  Retrieves information about an existing Open WebUI knowledge entry.
---

# openwebui_knowledge (Data Source)

Use this data source to inspect an existing knowledge entry, including access control and associated metadata.

## Example Usage

```hcl
data "openwebui_knowledge" "collection" {
  name = "test_knowledge_1"
}
```

## Argument Reference

* `knowledge_id` (Optional) – Identifier of the knowledge entry to retrieve. When omitted, `name` must be supplied.
* `name` (Optional) – Case-insensitive knowledge name used to resolve the entry when `knowledge_id` is not provided. If multiple entries share the same name the lookup fails and the ID must be specified explicitly.

## Attribute Reference

* `id` – Knowledge entry identifier returned by Open WebUI (matches `knowledge_id`).
* `name` – Knowledge entry name.
* `description` – Description provided for the entry.
* `data_json` – JSON string containing the entry's data payload.
* `meta_json` – JSON string containing metadata associated with the entry.
* `read_groups` – Group names granted read access.
* `write_groups` – Group names granted write access.
* `created_at` – Creation date in `YYYY-MM-DD` format.
* `updated_at` – Last update date in `YYYY-MM-DD` format.
* `user_id` – Identifier of the user who owns the entry.
