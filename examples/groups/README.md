# OpenWebUI Groups Example

This example demonstrates how to use the OpenWebUI Terraform provider to manage groups.

## Usage

To run this example, you need to execute:

```bash
$ export OPENWEBUI_ENDPOINT="https://your-openwebui-instance"
$ export OPENWEBUI_TOKEN="your-token"
$ terraform init
$ terraform plan
$ terraform apply
```

Note that this example may create resources which cost money. Run `terraform destroy` when you don't need these resources.

## Requirements

* You must have an OpenWebUI instance running and accessible
* You must have a valid API token with appropriate permissions
* You must configure the endpoint and token either via environment variables or in the provider configuration

## Provider Configuration

The provider can be configured with the following attributes:

* `endpoint` - (Required) The URL of your OpenWebUI instance. Can also be set with the `OPENWEBUI_ENDPOINT` environment variable.
* `token` - (Required) The API token for authentication. Can also be set with the `OPENWEBUI_TOKEN` environment variable.

## Resource Configuration

The `openwebui_group` resource supports the following attributes:

* `name` - (Required) The name of the group
* `description` - (Optional) A description of the group
* `user_ids` - (Optional) List of user IDs to include in the group
* `permissions` - (Optional) Group permissions configuration
  * `workspace` - (Required) Workspace-level permissions
    * `models` - (Required) Allow access to models
    * `knowledge` - (Required) Allow access to knowledge bases
    * `prompts` - (Required) Allow access to prompts
    * `tools` - (Required) Allow access to tools
  * `chat` - (Required) Chat-level permissions
    * `file_upload` - (Required) Allow file uploads in chat
    * `delete` - (Required) Allow message deletion
    * `edit` - (Required) Allow message editing
    * `temporary` - (Required) Allow temporary messages
