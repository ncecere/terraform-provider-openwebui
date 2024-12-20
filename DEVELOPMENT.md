# Development Guide for OpenWebUI Terraform Provider

This guide provides detailed instructions for developing and testing the OpenWebUI Terraform Provider.

## Development Environment Setup

### Prerequisites

1. **Go Installation**
   - Install [Go](https://golang.org/doc/install) >= 1.19
   - Set up your GOPATH and GOBIN environment variables
   ```bash
   export GOPATH=$HOME/go
   export GOBIN=$GOPATH/bin
   export PATH=$PATH:$GOBIN
   ```

2. **Terraform Installation**
   - Install [Terraform](https://www.terraform.io/downloads.html) >= 1.0
   - Verify installation: `terraform -v`

3. **Git Setup**
   - Install [Git](https://git-scm.com/downloads)
   - Configure your Git user:
   ```bash
   git config --global user.name "Your Name"
   git config --global user.email "your.email@example.com"
   ```

### Repository Setup

1. **Clone the Repository**
   ```bash
   git clone git@github.com:ncecere/terraform-provider-openwebui.git
   cd terraform-provider-openwebui
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   ```

## Building the Provider

### Local Build

1. **Basic Build**
   ```bash
   make build
   ```
   This creates a provider binary in the `dist` directory.

2. **Install for Local Testing**
   ```bash
   make install
   ```
   This builds and installs the provider to your local Terraform plugin directory.

### Cross-Platform Building

To build for multiple platforms:

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 make build

# Build for Windows
GOOS=windows GOARCH=amd64 make build

# Build for macOS
GOOS=darwin GOARCH=amd64 make build
```

## Testing

### Running Tests

1. **Unit Tests**
   ```bash
   make test
   ```

2. **Acceptance Tests**
   ```bash
   # Set required environment variables
   export OPENWEBUI_ENDPOINT="http://your-test-instance"
   export OPENWEBUI_TOKEN="your-test-token"

   # Run acceptance tests
   make testacc
   ```

### Local Testing

1. **Setup Test Environment**
   ```bash
   cd local_testing
   ```

2. **Initialize Terraform**
   ```bash
   terraform init
   ```

3. **Run Test Configuration**
   ```bash
   terraform plan
   terraform apply
   ```

4. **Cleanup**
   ```bash
   terraform destroy
   ```

## Development Workflow

### Adding a New Resource

1. **Create Client Implementation**
   - Add new client package in `internal/provider/client/`
   - Implement API operations
   - Add types and models

   Example structure:
   ```go
   // internal/provider/client/newresource/client.go
   package newresource

   type Client struct {
       endpoint string
       token    string
   }

   func NewClient(endpoint, token string) *Client {
       return &Client{
           endpoint: endpoint,
           token:    token,
       }
   }

   // Add CRUD operations
   ```

2. **Implement Resource**
   - Create resource implementation in `internal/provider/`
   - Implement CRUD functions
   - Add schema definition

   Example:
   ```go
   // internal/provider/newresource_resource.go
   func NewNewResourceResource() resource.Resource {
       return &NewResourceResource{}
   }

   func (r *NewResourceResource) Schema(...) {...}
   func (r *NewResourceResource) Create(...) {...}
   func (r *NewResourceResource) Read(...) {...}
   func (r *NewResourceResource) Update(...) {...}
   func (r *NewResourceResource) Delete(...) {...}
   ```

3. **Add to Provider**
   - Register resource in provider.go
   ```go
   func (p *OpenWebUIProvider) Resources(...) {
       return []func() resource.Resource{
           // Add your new resource
           NewNewResourceResource,
       }
   }
   ```

4. **Add Documentation**
   - Create resource documentation in `docs/resources/`
   - Add examples in `examples/`
   - Update provider documentation

5. **Add Tests**
   - Add unit tests
   - Add acceptance tests
   - Add example configurations

### Code Style and Standards

1. **Go Formatting**
   ```bash
   # Format code
   go fmt ./...

   # Run linter
   golangci-lint run
   ```

2. **Documentation Standards**
   - Use complete sentences
   - Include examples for all resources
   - Document all schema attributes
   - Follow [Terraform documentation standards](https://www.terraform.io/docs/registry/providers/docs.html)

3. **Commit Messages**
   - Use conventional commits format
   - Include relevant issue numbers
   - Be descriptive but concise

## Debugging

### Common Issues

1. **Plugin Installation**
   - Check plugin directory: `~/.terraform.d/plugins/`
   - Verify provider version in terraform configuration
   - Clear terraform plugin cache: `rm -rf ~/.terraform.d/plugin-cache`

2. **API Communication**
   - Enable debug logging:
   ```bash
   export TF_LOG=DEBUG
   export TF_LOG_PATH=terraform.log
   ```
   - Check API responses in logs
   - Verify endpoint and token configuration

3. **Build Issues**
   - Clean build artifacts: `make clean`
   - Update dependencies: `go mod tidy`
   - Check Go version compatibility

### Debugging Tools

1. **Delve Debugger**
   ```bash
   # Install Delve
   go install github.com/go-delve/delve/cmd/dlv@latest

   # Debug tests
   dlv test ./internal/provider/...
   ```

2. **Terraform Logs**
   ```bash
   # Set log level
   export TF_LOG=TRACE

   # Save logs to file
   export TF_LOG_PATH=./terraform.log
   ```

## Release Process

1. **Version Update**
   - Update version in `main.go`
   - Update CHANGELOG.md
   - Update documentation if needed

2. **Testing**
   ```bash
   # Run all tests
   make test
   make testacc
   ```

3. **Build Release**
   ```bash
   # Build for all platforms
   make release
   ```

4. **Create Release**
   - Tag the release
   - Create GitHub release
   - Upload built artifacts

## Contributing

1. **Fork and Clone**
   ```bash
   git clone git@github.com:your-username/terraform-provider-openwebui.git
   ```

2. **Create Branch**
   ```bash
   git checkout -b feature/your-feature
   ```

3. **Make Changes**
   - Write code
   - Add tests
   - Update documentation

4. **Submit PR**
   - Push changes
   - Create pull request
   - Respond to review comments

## Additional Resources

- [Terraform Plugin Framework Documentation](https://developer.hashicorp.com/terraform/plugin/framework)
- [Go Documentation](https://golang.org/doc/)
- [OpenWebUI API Documentation](https://your-openwebui-docs-url)
