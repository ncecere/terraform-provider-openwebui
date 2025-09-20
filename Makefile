
.PHONY: build test tidy install clean plugin

BIN ?= terraform-provider-openwebui
GO ?= go
BIN_DIR ?= $(CURDIR)/bin
VERSION ?= 2.0.0
OS ?= $(shell $(GO) env GOOS)
ARCH ?= $(shell $(GO) env GOARCH)
LDFLAGS ?= -X github.com/nickcecere/terraform-provider-openwebui/internal/provider.Version=$(VERSION)

LOCAL_REGISTRY_ROOT ?= $(BIN_DIR)/registry.terraform.io/nickcecere/openwebui/$(VERSION)/$(OS)_$(ARCH)
LOCAL_PLUGIN_NAME ?= terraform-provider-openwebui_v$(VERSION)

build:
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BIN)

test:
	$(GO) test ./...

tidy:
	$(GO) mod tidy

install: plugin

plugin: build
	mkdir -p $(LOCAL_REGISTRY_ROOT)
	cp $(BIN_DIR)/$(BIN) $(LOCAL_REGISTRY_ROOT)/$(LOCAL_PLUGIN_NAME)

clean:
	rm -rf $(BIN_DIR)
