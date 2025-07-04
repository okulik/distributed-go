# Go related variables.
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
TOOLSBIN := $(GOBASE)/.tools
DIR=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))

export PATH := $(TOOLSBIN):$(PATH)

MAKEFLAGS += --silent

.PHONY: default
default: help

.PHONY: tools
tools: $(TOOLSBIN)/golangci_lint ## Install tools used for development
	@echo "Done."

.PHONY: clean
clean: ## Clean build files and artifacts
	GOBIN=$(GOBIN) go clean ./...
	rm -rfv $(GOBIN)

.PHONY: lint
lint: tools ## Runs linter on all source code
	$(DIR).tools/golangci-lint run -v
	go vet ./...

.PHONY: test
test: ## Runs unit tests
	go test --timeout 30s ./...

.PHONY: test-coverage
test-coverage: ## Runs a test coverage reporting
	go test -cover ./...

.PHONY: test-coverage-visualize
test-coverage-visualize: ## Creates a test coverage visualization
	go test -coverprofile=c.out ./... && go tool cover -html=c.out

$(TOOLSBIN)/golangci_lint:  ## Installs golang linter runner
	$(DIR)scripts/install_golangci_lint $(TOOLSBIN) latest

.PHONY: help
help: ## Creates this help message
	grep -hE '^[/a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-26s\033[0m %s\n", $$1, $$2}'
