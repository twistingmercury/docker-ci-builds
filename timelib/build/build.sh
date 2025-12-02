#!/usr/bin/env bash
#
# build.sh - Library CI: runs tests in Docker
#
# This is the simplest CI pattern: just verify tests pass.
# No binary export, no E2E tests, no docker-compose.
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJ_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

readonly IMAGE_NAME="timelib-ci"

log_info() {
    printf '[INFO] %s\n' "$1"
}

log_success() {
    printf '[SUCCESS] %s\n' "$1"
}

log_error() {
    printf '[ERROR] %s\n' "$1" >&2
}

run_tests() {
    log_info "Building and testing timelib..."

    if ! docker build \
        --rm \
        --file "${SCRIPT_DIR}/Dockerfile" \
        --tag "${IMAGE_NAME}:test" \
        "${PROJ_ROOT}"; then
        log_error "Tests failed"
        return 1
    fi

    log_success "All tests passed"
    return 0
}

main() {
    log_info "Starting library CI..."

    if ! command -v docker >/dev/null 2>&1; then
        log_error "docker command not found in PATH"
        return 1
    fi

    if ! run_tests; then
        return 1
    fi

    log_success "Library CI completed successfully"
    return 0
}

main "$@"
