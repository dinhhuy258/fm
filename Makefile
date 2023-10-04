.ONESHELL:

PACKAGES := $(shell go list ./... | grep -v /vendor/)
GO_LINT_RUN = golangci-lint run

.PHONY: fmt
fmt:
	@go fmt $(PACKAGES)

.PHONY: lint
lint: ## Run golang linter
	$(GO_LINT_RUN)