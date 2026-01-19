# CLI Tool - CLI CI Pattern Example

> **Maturity Level**: Emerging - Demonstration of Docker CI for CLI tools

A minimal CLI tool that shows how to do Docker-based CI for command-line apps.
The CLI itself is trivial - we're focused on the build and test workflow, not
the application code.

## The CLI CI Pattern

When you're building CLI tools, your CI workflow looks different from services.
You need to get binaries OUT of the container:

1. Build binaries for multiple platforms inside Docker
2. Run unit tests during the build
3. Export binaries to the host filesystem using `--output`
4. Run E2E tests against the exported binary in a separate container
5. If everything passes, your binaries are ready to ship

This is different from service patterns where the image itself is what you deploy.
Here, you're extracting files.

## Usage

Build all platform binaries:

```bash
make build
```

You'll find the binaries in `.bin/amd64/{darwin,linux,windows}/`.

## How it works

### Multi-stage Dockerfile with Export

The Dockerfile (`build/Dockerfile`) has two stages:

- **builder**: Downloads dependencies, runs unit tests, cross-compiles binaries
  for darwin, linux, and windows (all amd64)
- **export**: A `scratch` stage that just copies binaries for extraction

The `--output` flag is what gets files out of the container and onto your host:

```bash
docker build --target export --output .bin/ .
```

### Cross-compilation

Go makes cross-compilation pretty straightforward with environment variables:

```dockerfile
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build ...
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build ...
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build ...
```

### Build Script

The build script (`build/build.sh`) orchestrates two things:

1. Export binaries via `docker build --target export --output`
2. Run E2E tests via docker-compose

Build metadata (version, date, commit) gets passed as build args.

### Docker Compose E2E Tests

The `tests/docker-compose.yaml` orchestrates E2E testing:

- Starts the `uuid_api` service (from the api/ project)
- Builds a test container that includes the exported linux binary
- Tests the CLI against the running API

The E2E test container copies the binary from `.bin/amd64/linux/`:

```dockerfile
COPY .bin/amd64/linux/* /go/bin/
```

This is important: you're testing the actual compiled binary, not rebuilding
something fresh.

## Key Considerations

- Unit tests run inside the Dockerfile during build
- E2E tests run against the exported binary, not a fresh build
- The export stage uses `scratch` as a base (no OS, just files)
- E2E tests need the uuid_api service running - you must build the API project
  first (`cd ../api && ./build/build.sh`) before running CLI E2E tests
- The docker-compose.yaml references the `uuid_api:latest` image, so it must
  exist locally
- The API healthcheck uses the binary's built-in `--health` flag:
  `["CMD", "/uuid_api", "--health"]`

## Development Considerations

### Quick Start

```bash
make build
```

This builds all binaries and runs both unit and E2E tests.

### Building and running

The Makefile wraps the build script:

```bash
# Build with automatic version detection from git tags
make build

# Or run the build script directly if you want custom values
BUILD_VER=v1.0.0 BUILD_DATE=2025-01-01 BUILD_COMMIT=abc123 ./build/build.sh
```

### Testing

Tests are baked into the Docker build:

- **Unit tests**: Run in the builder stage via `go test`
- **E2E tests**: Run in a separate container against the compiled linux binary

If any tests fail, the build fails. No half-broken artifacts.

### Versioning

We use git tag-based versioning. The build script extracts:

- Version from `git describe --tags`
- Build date from the current date
- Commit hash from `git rev-parse`

These values get embedded into binaries via ldflags.

### Prerequisites

- Docker 20.10+ -
  [Installation instructions](https://docs.docker.com/get-docker/)
- Docker Compose v2+ (included with Docker Desktop)
- GNU Make 3.81+
- Git (for version metadata)
