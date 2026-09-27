PROVIDER_NAME := terraform-provider-hitechcloud
VERSION ?= 1.0.0
GOFMT_FILES := $(shell find . -name '*.go' -not -path './vendor/*')

default: build

.PHONY: build
build: ## Build the provider binary for the current platform
	go build -o bin/$(PROVIDER_NAME) .

.PHONY: install
install: build ## Install the provider into the local Terraform plugin cache (dev override alternative)
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/hitechcloud-vietnam/hitechcloud/$(VERSION)/$$(go env GOOS)_$$(go env GOARCH)
	cp bin/$(PROVIDER_NAME) ~/.terraform.d/plugins/registry.terraform.io/hitechcloud-vietnam/hitechcloud/$(VERSION)/$$(go env GOOS)_$$(go env GOARCH)/$(PROVIDER_NAME)_v$(VERSION)

.PHONY: fmt
fmt: ## Format Go source files
	gofmt -s -w $(GOFMT_FILES)

.PHONY: fmtcheck
fmtcheck: ## Check Go source formatting
	@unformatted="$$(gofmt -l $(GOFMT_FILES))"; if [ -n "$$unformatted" ]; then echo "The following files are not gofmt-formatted:"; echo "$$unformatted"; exit 1; fi

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint if available
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || echo "golangci-lint not installed, skipping"

.PHONY: test
test: ## Run unit tests (no Terraform CLI or live API required)
	go test ./... -count=1 -timeout=10m

.PHONY: test-race
test-race: ## Run unit tests with the race detector
	go test ./... -race -count=1 -timeout=10m

.PHONY: testacc
testacc: ## Run acceptance tests against a mock API (requires terraform CLI on PATH)
	TF_ACC=1 go test ./... -v -count=1 -timeout=30m

.PHONY: cover
cover: ## Run unit tests with coverage
	go test ./... -count=1 -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out

.PHONY: validate
validate: fmtcheck vet test ## Run formatting checks, vet and unit tests

.PHONY: tools
tools: ## Install development tooling
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

.PHONY: generate
generate: ## Regenerate provider documentation (requires tfplugindocs)
	go generate ./...

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf bin/ dist/ coverage.out
