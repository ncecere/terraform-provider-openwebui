.PHONY: build test tidy install

BIN ?= terraform-provider-openwebui
GO ?= go

build:
	$(GO) build ./...

test:
	$(GO) test ./...

tidy:
	$(GO) mod tidy

install: build
	$(GO) build -o bin/$(BIN)
