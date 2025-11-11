# Get all Go module directories
GO_MODULES := $(shell go list -f '{{.Dir}}' -m)

# Run a target in all modules that have a Makefile with that target
define run_in_modules
	@for dir in $(GO_MODULES); do \
		if [ -f $$dir/Makefile ] && grep -q "^$(1):" $$dir/Makefile; then \
			$(MAKE) -C $$dir $(1); \
		fi; \
	done
endef

.DEFAULT_GOAL := help

.PHONY: fmt
fmt: ## Format code in all modules
	@for dir in $(GO_MODULES); do \
		(cd $$dir && go fmt ./...); \
	done
	$(call run_in_modules,fmt)

.PHONY: lint
lint: ## Lint code in all modules
	@for dir in $(GO_MODULES); do \
		(cd $$dir && go vet ./...); \
	done
	$(call run_in_modules,lint)

.PHONY: test
test: ## Run tests in all modules
	@for dir in $(GO_MODULES); do \
		(cd $$dir && go test ./...); \
	done
	$(call run_in_modules,test)

.PHONY: generate
generate: ## Generate code in all modules
	@for dir in $(GO_MODULES); do \
		(cd $$dir && go generate ./...); \
	done
	$(call run_in_modules,generate)

.PHONY: build
build: ## Build all modules
	$(call run_in_modules,build)

.PHONY: clean
clean: ## Clean all modules
	$(call run_in_modules,clean)

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
