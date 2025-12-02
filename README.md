# Docker CI Build Examples

> **Maturity Level**: Emerging - Demonstration project for Docker CI patterns

This project shows you three ways to set up Docker-based CI builds - for services,
CLI tools, and libraries. The example apps are dead simple on purpose. We're here
to talk about CI patterns, not application code.

## Why Docker for CI?

Here's the thing: when your build logic lives in Dockerfiles and shell scripts,
your CI config becomes almost boring. The CI platform just orchestrates - it
doesn't own your build process. That buys you:

- **Portability**: Switching CI platforms? Your build logic stays the same.
- **Consistency**: What builds locally builds the same way in CI. No more
  "works on my machine."
- **Reproducibility**: Your build environment is code, not some runner config
  you have to remember to update.

## Three CI Patterns

We've got three distinct patterns here, depending on what you're building.

### Service Pattern (api/)

This one's for things that run as long-lived processes - web servers, APIs,
daemons, that sort of thing:

1. Build the Docker image (tests run during build)
2. Start a container from that image
3. Run E2E tests against the running container
4. The image itself is your deployable artifact

Check out [api/README.md](api/README.md) for the details.

### CLI Pattern (cli/)

This one's for standalone binaries - CLI tools, utilities, anything you
distribute as a file:

1. Build binaries inside Docker
2. Export them to your host using the `--output` flag
3. Run E2E tests against the exported binary
4. The binaries are what you ship

See [cli/README.md](cli/README.md) for how it works.

### Library Pattern (timelib/)

This one's for Go packages that other code imports - no binary, no container,
just code that needs to compile and pass tests:

1. Run tests inside Docker
2. That's it. No artifacts, no exports, no E2E.

The simplest pattern of the three. Check out
[timelib/README.md](timelib/README.md) for the details.

## Key Differences

| Aspect       | Service           | CLI              | Library         |
| ------------ | ----------------- | ---------------- | --------------- |
| Artifact     | Docker image      | Binary files     | None            |
| Export       | None              | `--output`       | None            |
| E2E tests    | Running container | Exported binary  | None (unit only)|
| Distribution | Registry          | File downloads   | Go modules      |

## Quick Start

Build and test the API:

```bash
cd api && ./build/build.sh
```

Build and test the CLI:

```bash
cd cli && make build
```

Test the library:

```bash
cd timelib && make build
```

## Local Development

The repository uses a `go.work` file for local development across modules.
This lets you modify timelib and immediately use those changes in the API
without publishing the module.

The go.work file includes all three modules:

- api/
- cli/
- timelib/

When running `go build` or `go test` locally, Go automatically uses the
local versions. Docker builds don't use go.work - they handle dependencies
explicitly in each Dockerfile.

## Prerequisites

- Docker 20.10+ -
  [Installation instructions](https://docs.docker.com/get-docker/)
- Docker Compose v2+ (included with Docker Desktop)
- GNU Make 3.81+ (for CLI project)
- Git (for version metadata)
