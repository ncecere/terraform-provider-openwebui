# OpenWebUI Model Example

This example demonstrates how to use the OpenWebUI provider to manage models.

## Usage

To run this example, you need to execute:

```bash
$ terraform init
$ terraform plan
$ terraform apply
```

Note that this example may create resources which cost money. Run `terraform destroy` when you don't need these resources.

## Requirements

* Terraform 0.13+
* OpenWebUI instance with API access

## Provider Configuration

The provider needs to be configured with the proper endpoint and credentials before it can be used.

## Resources Created

This example creates:
* A custom model based on GPT-4
* Access control settings for the model
* Custom model parameters and metadata

## Outputs

* `model_id` - The ID of the created model
* `model_capabilities` - The capabilities configured for the model
* `model_tags` - The tags assigned to the model
