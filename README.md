# Docker CI Build Examples

> **Maturity Level**: Emerging - Demonstration project for Docker CI patterns

This project shows you two ways to set up Docker-based CI builds - one for services,
one for CLI tools. The example apps are dead simple on purpose. We're here to talk
about CI patterns, not application code.

## Why Docker for CI?

Here's the thing: when your build logic lives in Dockerfiles and shell scripts,
your CI config becomes almost boring. The CI platform just orchestrates - it
doesn't own your build process. That buys you:

- **Portability**: Switching CI platforms? Your build logic stays the same.
- **Consistency**: What builds locally builds the same way in CI. No more
  "works on my machine."
- **Reproducibility**: Your build environment is code, not some runner config
  you have to remember to update.

## Two CI Patterns

We've got two distinct patterns here, depending on what you're building.

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

## Key Differences

| Aspect       | Service Pattern          | CLI Pattern        |
| ------------ | ------------------------ | ------------------ |
| Artifact     | Docker image             | Binary files       |
| Export       | None (image is artifact) | `--output` to host |
| E2E target   | Running container        | Exported binary    |
| Distribution | Container registry       | File downloads     |

## Quick Start

Build and test the API:

```bash
cd api && ./build/build.sh
```

Build and test the CLI:

```bash
cd cli && make build
```

## Prerequisites

- Docker 20.10+ -
  [Installation instructions](https://docs.docker.com/get-docker/)
- Docker Compose v2+ (included with Docker Desktop)
- GNU Make 3.81+ (for CLI project)
- Git (for version metadata)
