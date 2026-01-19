# UUID API - Service CI Pattern Example

> **Maturity Level**: Emerging - Demonstration of Docker CI for services

A minimal HTTP API that shows how to do Docker-based CI for service-style apps.
The API itself does almost nothing - we're focused on the build and test workflow,
not the application code.

## The Service CI Pattern

When you're building something that runs as a long-lived process, your CI looks
like this:

1. Build the image with tests running inside the Dockerfile
2. Spin up a container from that image
3. Run E2E tests against the running container
4. If everything passes, the image is ready to deploy

This is different from CLI patterns where you're exporting artifacts out of
the container. Here, the image IS the artifact.

## Usage

Build the image and run E2E tests:

```bash
./build/build.sh
```

That script handles everything:

1. Builds the Docker image (unit tests run during build)
2. Starts the container via docker-compose
3. Waits for the health check to pass
4. Runs E2E tests against the running service
5. Tears everything down

## How it works

### Multi-stage Dockerfile

The Dockerfile (`build/Dockerfile`) has two stages:

- **builder**: Downloads dependencies, runs unit tests, compiles the binary
- **runtime**: A scratch image with just the binary (no OS, no shell)

Unit tests run during the build stage via `RUN go test -v ./...`. If tests
fail, the build fails - you don't get an image. Simple as that.

### Build Script

The build script (`build/build.sh`) uses `--rm --no-cache` flags to keep
builds clean and reproducible. Build metadata (version, date, commit) gets
passed as build args and embedded via ldflags.

### Docker Compose E2E Tests

The `tests/docker-compose.yaml` handles E2E testing:

```yaml
services:
  api:
    image: uuid_api:latest
    healthcheck:
      test: ["CMD", "/uuid_api", "--health"]
      interval: 2s
      timeout: 5s
      retries: 5

  tests:
    build:
      context: ./e2e
    depends_on:
      api:
        condition: service_healthy
```

A few patterns worth noting:

- **Health check**: The binary supports a `--health` flag that performs a
  self-check. This works with scratch images since no external tools are needed
- **depends_on with condition**: E2E tests wait for `service_healthy`, not
  just container start (this is key - the container starting isn't the same
  as the service being ready)
- **exit-code-from**: The build script uses `--exit-code-from tests` to
  propagate test failures to CI

## Key Considerations

- Unit tests run inside the Dockerfile during image build
- E2E tests run in a separate container against the running service
- The health check makes sure the service is actually ready before tests start
- Build uses `--no-cache` for reproducible CI builds
- You only get a usable image if both unit and E2E tests pass

## Development Considerations

### Quick Start

```bash
./build/build.sh
```

### Building and running

The build script handles everything:

```bash
# Build with automatic version detection from git
./build/build.sh

# Or override version metadata if you need to
BUILD_VER=v1.0.0 BUILD_DATE=2025-01-01 BUILD_COMMIT=abc123 ./build/build.sh
```

If you want to run the service manually after building:

```bash
docker run --rm -p 8080:8080 uuid_api:latest
```

### Testing

Tests fall into two categories:

- **Unit tests**: Run during docker build via `go test`
- **E2E tests**: Run via docker-compose against the running container

The build script runs both automatically, so you don't have to think about it.

### Versioning

We use git tag-based versioning. The build script pulls version metadata from
git and embeds it in the binary.

### Prerequisites

- Docker 20.10+ -
  [Installation instructions](https://docs.docker.com/get-docker/)
- Docker Compose v2+ (included with Docker Desktop)
- Git (for version metadata)
