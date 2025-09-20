---
layout: data-source
page_title: "openwebui_model Data Source"
sidebar_current: docs-openwebui-data-source-model
description: |-
  Retrieves details for an existing Open WebUI model.
---

# openwebui_model (Data Source)

Use this data source to read metadata and parameter settings for an existing model registered with Open WebUI.

## Example Usage

```hcl
data "openwebui_model" "volt" {
  model_id = "volt-answer"
}
```

## Argument Reference

* `model_id` (Required) – Identifier of the model to retrieve. This matches the `model_id` supplied when the model resource is created.

## Attribute Reference

* `id` – Open WebUI model identifier.
* `name` – Human-readable model name.
* `base_model_id` – Optional base model identifier.
* `is_active` – Whether the model is active.
* `user_id` – Identifier of the user who owns the model.
* `created_at` – Unix timestamp when the model was created.
* `updated_at` – Unix timestamp of the last update.
* `read_groups` – Group names with read access to the model.
* `write_groups` – Group names with write access to the model.
* `meta_additional_json` – JSON string preserving metadata returned by the API that is not otherwise exposed.
* `params_additional_json` – JSON string preserving parameter keys not otherwise exposed.

### `params` Block

The nested `params` block mirrors the Open WebUI model parameters and contains the same attributes described for the `openwebui_model` resource, including sampling controls (`temperature`, `top_p`, `repeat_penalty`, etc.), runtime flags (`use_mlock`, `use_mmap`, `think`), and tuning values (`num_ctx`, `max_tokens`, `reasoning_effort`).

### `capabilities` Block

The `capabilities` block exposes boolean flags that describe the features enabled for the model: `vision`, `file_upload`, `web_search`, `image_generation`, `code_interpreter`, `citations`, `status_updates`, and `usage`.
