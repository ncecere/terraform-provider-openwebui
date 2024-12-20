default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Build provider
.PHONY: build
build:
	go build -o terraform-provider-openwebui

# Install provider locally
.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/ncecere/openwebui/0.1.0/$(shell go env GOOS)_$(shell go env GOARCH)
	cp terraform-provider-openwebui ~/.terraform.d/plugins/registry.terraform.io/ncecere/openwebui/0.1.0/$(shell go env GOOS)_$(shell go env GOARCH)

# Clean build artifacts
.PHONY: clean
clean:
	rm -f terraform-provider-openwebui

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Run tests
.PHONY: test
test:
	go test ./... -v

# Generate documentation
.PHONY: docs
docs:
	go generate ./...

# Run all pre-commit checks
.PHONY: pre-commit
pre-commit: fmt test

# Run linter
.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: all
all: clean fmt lint test build
