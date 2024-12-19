# Model Resource

Manages a model in OpenWebUI.

## Example Usage

```hcl
resource "openwebui_model" "example" {
  base_model_id = "gpt-4"
  name          = "My Custom Model"
  is_active     = true

  params {
    system           = "You are a helpful assistant"
    stream_response  = true
    temperature     = 0.7
    top_p           = 0.9
    max_tokens      = 2000
    seed            = 42
    frequency_penalty = 1
    repeat_last_n    = 64
    num_ctx          = 4096
    num_batch        = 512
  }

  meta {
    description      = "A customized GPT-4 model"
    profile_image_url = "/static/favicon.png"
    
    capabilities {
      vision    = false
      usage     = true
      citations = true
    }

    tags {
      name = "production"
    }
  }

  access_control {
    read {
      group_ids = ["group1", "group2"]
      user_ids  = ["user1"]
    }
    write {
      group_ids = ["admin-group"]
      user_ids  = ["admin-user"]
    }
  }
}
```

## Argument Reference

* `base_model_id` - (Required) The ID of the base model to use.
* `name` - (Required) The name of the model.
* `is_active` - (Optional) Whether the model is active. Defaults to true.

### Params Configuration

The `params` block supports:

* `system` - (Optional) The system prompt for the model.
* `stream_response` - (Optional) Whether to stream responses. Defaults to true.
* `temperature` - (Optional) Sampling temperature. Range: 0.0-1.0.
* `top_p` - (Optional) Nucleus sampling parameter. Range: 0.0-1.0.
* `top_k` - (Optional) Top-k sampling parameter.
* `min_p` - (Optional) Minimum probability threshold.
* `max_tokens` - (Optional) Maximum number of tokens to generate.
* `seed` - (Optional) Random seed for reproducibility.
* `frequency_penalty` - (Optional) Penalty for token frequency.
* `repeat_last_n` - (Optional) Number of tokens to consider for repetition penalty.
* `num_ctx` - (Optional) Context window size.
* `num_batch` - (Optional) Batch size for processing.
* `num_keep` - (Optional) Number of tokens to keep from the prompt.

### Meta Configuration

The `meta` block supports:

* `description` - (Optional) Description of the model.
* `profile_image_url` - (Optional) URL for the model's profile image.
* `capabilities` - (Optional) Model capabilities configuration block.
* `tags` - (Optional) List of tags associated with the model.

#### Capabilities Configuration

The `capabilities` block supports:

* `vision` - (Optional) Whether the model supports vision tasks.
* `usage` - (Optional) Whether to track usage statistics.
* `citations` - (Optional) Whether the model supports citations.

#### Tags Configuration

The `tags` block supports:

* `name` - (Required) Name of the tag.

### Access Control Configuration

The `access_control` block supports:

* `read` - (Optional) Read access configuration block.
* `write` - (Optional) Write access configuration block.

Both `read` and `write` blocks support:

* `group_ids` - (Optional) List of group IDs with access.
* `user_ids` - (Optional) List of user IDs with access.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the model.
* `user_id` - The ID of the user who created the model.
* `created_at` - Timestamp when the model was created.
* `updated_at` - Timestamp when the model was last updated.

## Import

Models can be imported using their ID:

```shell
terraform import openwebui_model.example <model_id>
