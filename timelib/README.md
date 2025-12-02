# timelib - Go Library CI Pattern

> **Maturity Level**: Emerging - Demonstration of Docker CI for Go libraries

This package demonstrates the simplest Docker CI pattern for Go libraries.

## What Makes Library CI Different?

Unlike applications (like `api/` or `cli/`), libraries don't produce runnable
artifacts. There's no binary to ship, no container image to deploy. The CI goal
is simpler:

1. **Tests pass** - Your code works
2. **Code compiles** - Consumers can import it

That's it. One Dockerfile stage. No docker-compose. No E2E tests.

## The Pattern

```text
Application CI:          Library CI:

  build stage              test stage (only stage)
      |                        |
  test stage               [done]
      |
  runtime stage
      |
  e2e tests
      |
  [artifact]
```

## Usage

```bash
# Run CI (tests in Docker)
make build

# Run tests locally (faster iteration)
make test-local
```

## Importing This Library

```go
import "github.com/twistingmercury/docker-ci-build/timelib"

func handler(c *gin.Context) {
    c.JSON(200, gin.H{"time": timelib.Now()})
}
```

For local development across modules, use the `go.work` file at the repo root.

## Key Considerations

- This is a demonstration library, not production code
- The library exports a single function for simplicity
- The API project imports this library, demonstrating cross-module dependencies

## Development Considerations

### Prerequisites

- Docker 20.10+ - [Installation instructions](https://docs.docker.com/get-docker/)
- Go 1.23+ (for local development only) - [Installation instructions](https://go.dev/doc/install)

### Testing

```bash
# Run CI (tests in Docker)
make build

# Run tests locally (faster iteration)
make test-local
```

### Versioning

This library uses the repository's git tags for versioning. For local development,
go.work handles module resolution. Docker builds use `go mod edit -replace` for
local module resolution.
