.PHONY: build test clean help

GIT_COMMIT := $(shell git rev-parse --short=8 HEAD 2>/dev/null || echo "unknown")
GIT_TAG := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nAvailable targets:\n"} /^[a-zA-Z0-9_-]+:.*##/ { printf "  %-12s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: ## Build Docker image
	./build/build.sh
	docker system prune -f

# test: build ## Run E2E tests via docker-compose
# 	docker compose up --build --abort-on-container-exit --exit-code-from e2e-tests
# 	docker compose down

clean: ## Remove built images and containers
	docker compose down --rmi local 2>/dev/null || true
	docker rmi time-api:latest 2>/dev/null || true
