.PHONY: build test clean help local api cli

GIT_COMMIT := $(shell git rev-parse --short=8 HEAD 2>/dev/null || echo "unknown")
GIT_TAG := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\n\033[1mAvailable targets:\033[0m\n"} /^[a-zA-Z0-9_-]+:.*##/ { printf "  %-12s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""

api: ## Builds the api
	api/build/build.sh
	docker system prune -f

cli: ## Builds the cli
	cli/build/build.sh
	docker system prune -f