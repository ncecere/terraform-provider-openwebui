# Model Data Source

Use this data source to get information about a specific OpenWebUI model.

## Example Usage

```hcl
data "openwebui_model" "example" {
  name = "My Custom Model"
}

output "model_capabilities" {
  value = data.openwebui_model.example.meta.capabilities
}
```

## Argument Reference

* `name` - (Required) The name of the model to look up.

## Attribute Reference

* `id` - The ID of the model.
* `user_id` - The ID of the user who created the model.
* `base_model_id` - The ID of the base model used.
* `is_active` - Whether the model is active.
* `created_at` - Timestamp when the model was created.
* `updated_at` - Timestamp when the model was last updated.

### Params Attributes

The `params` block contains:

* `system` - The system prompt for the model.
* `stream_response` - Whether responses are streamed.
* `temperature` - Sampling temperature.
* `top_p` - Nucleus sampling parameter.
* `top_k` - Top-k sampling parameter.
* `min_p` - Minimum probability threshold.
* `max_tokens` - Maximum number of tokens to generate.
* `seed` - Random seed for reproducibility.
* `frequency_penalty` - Penalty for token frequency.
* `repeat_last_n` - Number of tokens considered for repetition penalty.
* `num_ctx` - Context window size.
* `num_batch` - Batch size for processing.
* `num_keep` - Number of tokens kept from the prompt.

### Meta Attributes

The `meta` block contains:

* `description` - Description of the model.
* `profile_image_url` - URL for the model's profile image.
* `capabilities` - Model capabilities configuration.
* `tags` - List of tags associated with the model.

#### Capabilities Attributes

The `capabilities` block contains:

* `vision` - Whether the model supports vision tasks.
* `usage` - Whether usage statistics are tracked.
* `citations` - Whether the model supports citations.

#### Tags Attributes

Each tag in the `tags` list contains:

* `name` - Name of the tag.

### Access Control Attributes

The `access_control` block contains:

* `read` - Read access configuration.
* `write` - Write access configuration.

Both `read` and `write` blocks contain:

* `group_ids` - List of group IDs with access.
* `user_ids` - List of user IDs with access.
