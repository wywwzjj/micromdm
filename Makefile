.PHONY: build

WORKSPACE := $(dir $(shell go env GOMOD))
BUILD_DIR := $(WORKSPACE)/build

export GOBIN=$(BUILD_DIR)

download:
	@echo "Download dependencies"
	@go mod download

install-tools: download
	@echo "Installing tools"
	@cat cmd/tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %
