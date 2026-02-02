# Docker CI Build Examples

> **Maturity Level**: Emerging - Demonstration project for Docker CI patterns

This project shows you two ways to set up Docker-based CI builds - for services
and CLI tools. The example apps are dead simple on purpose. We're here to talk
about CI patterns, not application code.

## Usage

Build and test the API:

```bash
make api
```

Build and test the CLI (requires API to be built first):

```bash
make cli
```

Both builds accept version metadata overrides:

```bash
BUILD_VER=v1.0.0 BUILD_DATE=2025-01-01 BUILD_COMMIT=abc123 make api
BUILD_VER=v1.0.0 BUILD_DATE=2025-01-01 BUILD_COMMIT=abc123 make cli
```

If you want to run the API manually after building:

```bash
docker run --rm -p 8080:8080 uuid_api:latest
```

## How it works

### Why Docker for CI?

Here's the thing: when your build logic lives in Dockerfiles and shell scripts,
your CI config becomes almost boring. The CI platform just orchestrates - it
doesn't own your build process. That buys you:

- **Portability**: Switching CI platforms? Your build logic stays the same.
- **Consistency**: What builds locally builds the same way in CI. No more
  "works on my machine."
- **Reproducibility**: Your build environment is code, not some runner config
  you have to remember to update.

### Two CI Patterns

We've got two distinct patterns here, depending on what you're building.

#### Service Pattern (api/)

This one's for things that run as long-lived processes - web servers, APIs,
daemons, that sort of thing:

1. Build the Docker image (tests run during build)
2. Start a container from that image
3. Run E2E tests against the running container
4. The image itself is your deployable artifact

The API example is a minimal HTTP service that generates UUIDs. The build
process runs quality checks (goimports, golangci-lint, govulncheck, gosec)
and unit tests during the Docker build. If anything fails, you don't get
an image.

**Multi-stage Dockerfile:**

- **builder**: Downloads dependencies, runs linters/security tools, runs unit
  tests, compiles the binary
- **runtime**: A scratch image with just the binary (no OS, no shell)

**Health check pattern:**

The binary supports a `--health` flag for self-checks, which works with scratch
images since no external tools are needed. The docker-compose setup uses
`depends_on` with `condition: service_healthy` to ensure E2E tests wait for
the service to actually be ready, not just started.

**Build and test:**

```bash
make api
```

The build script runs with `--no-cache` for reproducibility, starts the
container via docker-compose, waits for health checks, runs E2E tests, and
tears down. Version metadata (from git tags) gets embedded via ldflags.

#### CLI Pattern (cli/)

This one's for standalone binaries - CLI tools, utilities, anything you
distribute as a file:

1. Build binaries inside Docker
2. Export them to your host using the `--output` flag
3. Run E2E tests against the exported binary
4. The binaries are what you ship

The CLI example is a minimal tool that calls the UUID API. The build process
runs the same quality checks as the API, then cross-compiles binaries for
darwin, linux, and windows (all amd64).

**Multi-stage Dockerfile:**

- **builder**: Downloads dependencies, runs linters/security tools, runs unit
  tests, cross-compiles binaries
- **export**: A scratch stage that copies binaries for extraction

**Binary export:**

The `--output` flag extracts files from the container to your host:

```bash
docker build --target export --output .bin/ .
```

Cross-compilation uses Go environment variables:

```dockerfile
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build ...
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build ...
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build ...
```

**E2E testing:**

The docker-compose setup tests the actual exported linux binary (not a fresh
build) against the running API service. The test container copies the binary
from `.bin/amd64/linux/`.

**Build and test:**

```bash
make cli
```

This builds all platform binaries and runs E2E tests. You'll find binaries in
`.bin/amd64/{darwin,linux,windows}/`. The API must be built first since E2E
tests need `uuid_api:latest` to exist locally.

### Key Differences

| Aspect       | Service           | CLI                      |
| ------------ | ----------------- | ------------------------ |
| Artifact     | Docker image      | Binary files             |
| Export       | None              | `--output`               |
| E2E tests    | Running container | Exported binary          |
| Distribution | Registry          | File downloads           |
| Example      | UUID API service  | CLI tool calling the API |

### CI/CD Workflows

The project uses separate GitHub Actions workflows for the API and CLI, with
artifact passing to ensure the CLI tests against the exact API image that
passed its own tests.

#### Workflow Separation

**API CI** (`api-ci.yaml`):

- Builds and tests the API service
- Saves the tested image as artifact `uuid-api-image` (1-day retention)
- Only runs when API code changes

**CLI CI** (`cli-ci.yaml`):

- Runs automatically after API CI succeeds (via `workflow_run` trigger)
- Also runs independently when CLI code changes (via push/PR to CLI paths)
- When triggered by API CI completion: downloads and loads the API artifact
- When triggered by CLI changes: requires `uuid_api:latest` image to exist locally
- Builds and tests the CLI against the API image

#### Artifact Passing

The workflows coordinate through GitHub Actions artifacts:

1. API CI saves `uuid_api:latest` as artifact `uuid-api-image`
2. CLI CI downloads and loads this artifact when triggered by API CI completion
3. When CLI CI is triggered by CLI code changes, it expects `uuid_api:latest` to
   exist (either from a previous run or local build)

This pattern ensures the CLI tests run against a known-good API build when
triggered by API changes, while allowing independent CLI development when
the API hasn't changed.

#### Production Note

In production environments, you'd skip artifact passing. Instead, the CLI
workflow would pull the latest tested API image from a container registry
(Docker Hub, Azure Container Registry, or GitHub Container Registry). The
artifact pattern here demonstrates CI coordination without requiring registry
infrastructure.

### Build Process Details

Both projects use multi-stage Dockerfiles that run quality checks before
compilation:

- Code formatting with goimports
- Linting with golangci-lint
- Vulnerability scanning with govulncheck
- Security scanning with gosec
- Unit tests with go test

All checks must pass before binaries are compiled. Builds use `--no-cache`
for reproducible CI builds.

### Testing Approach

- **Unit tests**: Run inside the Dockerfile during image build
- **E2E tests**: Run in a separate container via docker-compose

For the API, E2E tests run against the running service container. For the CLI,
E2E tests run against the exported linux binary. Both test setups use health
checks to ensure services are ready before tests start.

### Version Metadata

Both projects use git tag-based versioning. Build scripts extract version from
`git describe --tags`, build date from the current date, and commit hash from
`git rev-parse`. These values get embedded into binaries via ldflags.

### Local Development

The repository uses a `go.work` file for local development across modules.
This lets you work on both modules simultaneously without publishing. When
running `go build` or `go test` locally, Go automatically uses the local
versions. Docker builds don't use go.work - they handle dependencies explicitly
in each Dockerfile.

## Key Considerations

- **API must be built first**: The CLI's E2E tests require `uuid_api:latest`
  to exist locally. Always build the API before building the CLI.
- **Artifact retention**: GitHub Actions artifacts are retained for 1 day in
  this demo. Production workflows should adjust based on deployment cadence.
- **Production registry usage**: The artifact passing pattern is for demo
  purposes. Production environments should pull tested images from a container
  registry instead.
- **No-cache builds**: Both projects use `--no-cache` for reproducibility in
  CI. This ensures clean builds but increases build time.

## Development Considerations

### Prerequisites

- Docker 20.10+ - [Installation instructions](https://docs.docker.com/get-docker/)
- Docker Compose v2+ (included with Docker Desktop)
- GNU Make 3.81+ (for CLI project)
- Git (for version metadata)

### Quick Start

1. Clone the repository
2. Build the API: `make api`
3. Build the CLI: `make cli`

### Building & Running

**API service:**

```bash
make api
```

This builds the Docker image, runs all quality checks and tests, then starts
the service and runs E2E tests. The resulting `uuid_api:latest` image is ready
to use.

To run the API manually:

```bash
docker run --rm -p 8080:8080 uuid_api:latest
```

**CLI tool:**

```bash
make cli
```

This builds binaries for all platforms (darwin, linux, windows) and runs E2E
tests against the linux binary. Find the binaries in `.bin/amd64/`.

### Testing

**Unit tests** run automatically during Docker builds. To run them separately:

```bash
# API tests
cd api && go test ./...

# CLI tests
cd cli && go test ./...
```

**E2E tests** run via docker-compose after builds complete. To run them
separately:

```bash
# API E2E tests (requires built image)
cd api && docker compose -f build/docker-compose.yaml up --abort-on-container-exit

# CLI E2E tests (requires built API image and CLI binary)
cd cli && docker compose -f build/docker-compose.yaml up --abort-on-container-exit
```

### Versioning

Both projects use git tag-based versioning with semantic versioning (e.g.,
`v1.2.3`). Version metadata is automatically embedded during builds via
ldflags.

To override version metadata:

```bash
BUILD_VER=v1.0.0 BUILD_DATE=2025-01-01 BUILD_COMMIT=abc123 make api
BUILD_VER=v1.0.0 BUILD_DATE=2025-01-01 BUILD_COMMIT=abc123 make cli
```
