---
layout: resource
page_title: "openwebui_model Resource"
sidebar_current: docs-openwebui-resource-model
description: |-
  Manages custom models registered with Open WebUI.
---

# openwebui_model (Resource)

Registers and manages models that Open WebUI can serve. The resource mirrors the behaviour of other Open WebUI resources by exposing structured arguments, automatic access control translation, and computed metadata preservation.

## Example Usage

```hcl
resource "openwebui_model" "volt_answer" {
  model_id = "volt-answer"
  name     = "Volt Answer"

  base_model_id = "gpt-oss-20b"
  is_active     = true

  params {
    system          = "You will only say I don't know and you should look it up your self.  be super sassy."
    stream_response = true
    temperature     = 0.8
    max_tokens      = 128
    top_p           = 0.9
    top_k           = 40
    num_thread      = 2
    num_batch       = 512
    num_ctx         = 2048
    custom_params = {
      custom_param_name = "custom_param_value"
    }
  }

  profile_image_url = "/static/favicon.png"
  tags              = ["tag1"]
  suggestion_prompts = [
    "I'm a promt suggestion",
  ]

  capabilities {
    vision           = true
    file_upload      = true
    web_search       = true
    image_generation = true
    code_interpreter = true
    citations        = true
  }

  tool_ids            = ["deep_research", "weather"]
  default_feature_ids = ["web_search", "image_generation", "code_interpreter"]

  read_groups  = ["Knowledge Managers"]
  write_groups = ["Model Curators"]
}
```

Additional API fields that are not modelled explicitly are preserved and exposed via the computed `*_additional_json` attributes.

## Argument Reference

Top-level arguments:

* `model_id` (Required) – The Open WebUI identifier for the model. This is sent to the API when creating or updating the model.
* `name` (Required) – Friendly name displayed in Open WebUI.
* `params` (Required) – Single nested block specifying model parameters. See [Params Block](#params-block) for details.
* `base_model_id` (Optional) – Identifier of the base model to extend.
* `is_active` (Optional) – Whether the model should be marked active. Defaults to the value returned by the API when omitted.
* `profile_image_url`, `description`, `suggestion_prompts`, `tags`, `tool_ids`, `default_feature_ids`, `capabilities` (Optional) – Presentation metadata. See [Metadata Arguments](#metadata-arguments).
* `read_groups` (Optional) – Group names or IDs granted read access. When populated, the provider resolves names to IDs using the Open WebUI API.
* `write_groups` (Optional) – Group names or IDs granted write access. Groups listed here automatically receive read access.
* `params_additional_json` (Optional) – Extra JSON merged into the params payload. This field is also populated automatically when the API returns unsupported keys.
* `meta_additional_json` (Optional) – Extra JSON merged into the metadata payload. This field is also populated automatically to preserve API-only fields.

### Params Block

The `params` block mirrors the payload expected by the Open WebUI model API. Every argument is optional unless specifically noted.

* `system` – System prompt injected ahead of user messages.
* `stream_response` – Enables incremental streaming of tokens when set to `true`.
* `stream_delta_chunk_size` – Number of streamed tokens per chunk (defaults to `1`).
* `function_calling` – Function calling strategy (`native`, `auto`, `none`, etc.).
* `reasoning_tags` – List of tags that opt the model into reasoning workflows.
* `seed` – Integer seed that enforces deterministic sampling when supported.
* `temperature` – Softmax temperature used for sampling (higher = more random).
* `keep_alive` – Idle duration (for example `5m`) before the model is evicted from memory.
* `num_gpu` – Number of GPU devices allocated to the model.
* `num_thread` – CPU thread count used while generating tokens.
* `num_batch` – Batch size applied during generation.
* `num_ctx` – Context window size in tokens.
* `num_keep` – Prefix tokens preserved from repetition penalty trimming.
* `format` – Response format hint such as `json`.
* `think` – When `true`, toggles the extended “think” mode used by some models.
* `use_mlock` – Locks model pages in memory to reduce swapping.
* `use_mmap` – Reads model weights via memory-mapped files when supported.
* `repeat_penalty` – Penalty factor applied to repeated tokens.
* `tfs_z` – Tail-free sampling Z parameter.
* `repeat_last_n` – Number of most recent tokens considered by the repeat penalty.
* `mirostat_tau` – Target entropy for the Mirostat sampler.
* `mirostat_eta` – Learning rate for the Mirostat sampler.
* `mirostat` – Mirostat mode selector (`0`, `1`, or `2`).
* `presence_penalty` – Penalty applied for introducing new tokens.
* `frequency_penalty` – Penalty applied to frequently seen tokens.
* `min_p` – Minimum probability mass preserved during sampling.
* `top_p` – Nucleus sampling probability mass cap.
* `top_k` – Limits sampling pool to the top-k candidates.
* `max_tokens` – Hard cap on the number of tokens generated per response.
* `reasoning_effort` – Reasoning workload hint (`low`, `medium`, or `high`).
* `custom_params` – Arbitrary key/value pairs forwarded verbatim to Open WebUI.

### Metadata Arguments

* `profile_image_url` – Profile image displayed for the model.
* `description` – Human-readable description.
* `suggestion_prompts` – Prompt suggestions surfaced to end users when selecting the model.
* `tags` – Simple list of tag names.
* `tool_ids` – Tool identifiers made available to this model.
* `default_feature_ids` – Feature identifiers enabled by default.
* `capabilities` – Nested block of boolean capability flags with the attributes `vision`, `file_upload`, `web_search`, `image_generation`, `code_interpreter`, `citations`, `status_updates`, and `usage`.

## Attribute Reference

* `id` – Open WebUI generated model identifier.
* `user_id` – Identifier of the user that created the model.
* `created_at` – Unix timestamp when the model was created.
* `updated_at` – Unix timestamp for the latest update.
* `meta_additional_json` / `params_additional_json` – Computed JSON preserving fields returned by the API that are not modelled directly.

## Import

Models can be imported using the model `id` returned by Open WebUI:

```bash
terraform import openwebui_model.custom_rag custom-rag
```
