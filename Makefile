BINARY=terraform-provider-labplatform
VERSION=0.1.0
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)
PLUGIN_DIR=~/.terraform.d/plugins/registry.terraform.io/desotech-it/labplatform/$(VERSION)/$(OS_ARCH)

.PHONY: build install clean test

build:
	go build -o $(BINARY) .

install: build
	mkdir -p $(PLUGIN_DIR)
	cp $(BINARY) $(PLUGIN_DIR)/

clean:
	rm -f $(BINARY)

test:
	go test ./... -v

fmt:
	go fmt ./...

tidy:
	go mod tidy
