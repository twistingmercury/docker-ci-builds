.PHONY: build test test-local clean help

help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nAvailable targets:\n"} /^[a-zA-Z0-9_-]+:.*##/ { printf "  %-12s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: ## Run library CI (tests in Docker)
	./build/build.sh

test: build ## Alias for build (library CI is just tests)

test-local: ## Run tests locally (not in Docker)
	go test -v -race -cover ./...

clean: ## Remove CI images
	docker rmi timelib-ci:test 2>/dev/null || true
