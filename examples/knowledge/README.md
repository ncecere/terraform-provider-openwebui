# OpenWebUI Knowledge Base Example

This example demonstrates how to use the OpenWebUI provider to manage knowledge bases.

## Prerequisites

- OpenWebUI instance running and accessible
- API token with appropriate permissions
- Terraform installed

## Usage

To run this example:

1. Set up your environment variables:
```bash
export OPENWEBUI_ENDPOINT="http://your-openwebui-instance"
export OPENWEBUI_TOKEN="your-api-token"
```

2. Initialize Terraform:
```bash
terraform init
```

3. Review the execution plan:
```bash
terraform plan
```

4. Apply the configuration:
```bash
terraform apply
```

## Example Resources

This example creates:

1. A public knowledge base:
   - Publicly accessible
   - Contains example metadata
   - Demonstrates basic configuration

2. A private knowledge base:
   - Access controlled
   - Contains custom metadata
   - Demonstrates advanced configuration

3. Data source lookups:
   - Retrieves information about created knowledge bases
   - Demonstrates data source usage
   - Shows how to reference resource attributes

## Implementation Details

The example uses the modular client architecture of the provider:

- Knowledge client handles API communication
- Resource implementation manages state
- Data source implementation retrieves existing resources

## Cleanup

To remove all resources created by this example:

```bash
terraform destroy
```

## Notes

- The example assumes default provider configuration through environment variables
- Access control settings can be modified as needed
- Custom data can be added to suit your use case
