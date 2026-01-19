.PHONY: build test vet staticcheck lint coverage all help

GIT_COMMIT=$(shell git rev-parse --short=8 HEAD)
GIT_TAG=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "untagged")
BUILD_DATE=$(shell date +%Y-%m-%d)

help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\n\033[1mAvailable targets:\033[0m\n"} /^[a-zA-Z0-9_-]+:.*##/ { printf "  %-12s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""

local: ## Build locally using workspace
	 go build -o .bin/amd64/darwin/cli-tool cmd/main/main.go

build: ## Build binaries (runs E2E tests first)
	BUILD_VER="${GIT_TAG}" BUILD_DATE="${BUILD_DATE}" BUILD_COMMIT="${GIT_COMMIT}" ./build/build.sh
