#!/usr/bin/env bash
#
# build.sh - Orchestrates Docker build with E2E tests and binary export
#
# This script runs E2E tests first, then exports binaries only if tests pass.
# Configuration is via environment variables (see validate_args function).
#
# Usage:
#   BUILD_VER=v1.0.0 BUILD_DATE=2025-01-01 BUILD_COMMIT=abc12345 ./build.sh
#
set -e

#---------------------------------------------------------------
# Directory Variables
#---------------------------------------------------------------
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJ_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

#---------------------------------------------------------------
# Global Constants
#---------------------------------------------------------------
readonly DOCKERFILE_PATH="${PROJ_ROOT}/build/Dockerfile"
readonly OUTPUT_DIR="${PROJ_ROOT}/.bin"
readonly IMAGE_NAME="cli-build-example"

#---------------------------------------------------------------
# Configuration Defaults
#---------------------------------------------------------------
: "${BUILD_VER:=$(git -C "${PROJ_ROOT}" describe --tags --abbrev=0 2>/dev/null || printf '%s' 'untagged')}"
: "${BUILD_DATE:=$(date +%Y-%m-%d)}"
: "${BUILD_COMMIT:=$(git -C "${PROJ_ROOT}" rev-parse --short=8 HEAD 2>/dev/null || printf '%s' 'unknown')}"

#---------------------------------------------------------------
# Internal Functions
#---------------------------------------------------------------

log_info() {
    printf '[INFO] %s\n' "$1"
}

log_error() {
    printf '[ERROR] %s\n' "$1" >&2
}

log_success() {
    printf '[SUCCESS] %s\n' "$1"
}

validate_args() {
    local has_error=0

    if [ -z "${BUILD_VER}" ]; then
        log_error "BUILD_VER is required but not set"
        has_error=1
    fi

    if [ -z "${BUILD_DATE}" ]; then
        log_error "BUILD_DATE is required but not set"
        has_error=1
    fi

    if [ -z "${BUILD_COMMIT}" ]; then
        log_error "BUILD_COMMIT is required but not set"
        has_error=1
    fi

    if [ ! -f "${DOCKERFILE_PATH}" ]; then
        log_error "Dockerfile not found at: ${DOCKERFILE_PATH}"
        has_error=1
    fi

    if ! command -v docker >/dev/null 2>&1; then
        log_error "docker command not found in PATH"
        has_error=1
    fi

    if [ "${has_error}" -eq 1 ]; then
        return 1
    fi

    return 0
}

print_build_info() {
    log_info "Build Configuration:"
    log_info "  BUILD_VER:    ${BUILD_VER}"
    log_info "  BUILD_DATE:   ${BUILD_DATE}"
    log_info "  BUILD_COMMIT: ${BUILD_COMMIT}"
    log_info "  Dockerfile:   ${DOCKERFILE_PATH}"
    log_info "  Output:       ${OUTPUT_DIR}"
    log_info "  Image Tag:    ${IMAGE_NAME}:${BUILD_VER}"
}

run_e2e_tests() {
    log_info "Running E2E tests (target: e2e_tests)..."

    if ! docker build \
        -f="${DOCKERFILE_PATH}" \
        --build-arg BUILD_VER="${BUILD_VER}" \
        --build-arg BUILD_DATE="${BUILD_DATE}" \
        --build-arg BUILD_COMMIT="${BUILD_COMMIT}" \
        --target e2e_tests \
        "${PROJ_ROOT}"; then
        log_error "E2E tests failed"
        return 1
    fi

    log_success "E2E tests passed"
    return 0
}

export_binaries() {
    log_info "Exporting binaries (target: export)..."

    if ! docker build \
        -f="${DOCKERFILE_PATH}" \
        --build-arg BUILD_VER="${BUILD_VER}" \
        --build-arg BUILD_DATE="${BUILD_DATE}" \
        --build-arg BUILD_COMMIT="${BUILD_COMMIT}" \
        --target export \
        --output "${OUTPUT_DIR}" \
        -t "${IMAGE_NAME}:${BUILD_VER}" \
        "${PROJ_ROOT}"; then
        log_error "Binary export failed"
        return 1
    fi

    log_success "Binaries exported to: ${OUTPUT_DIR}"
    return 0
}

#---------------------------------------------------------------
# Main Entry Point
#---------------------------------------------------------------
main() {
    log_info "Starting build process..."

    if ! validate_args; then
        log_error "Validation failed. Exiting."
        return 1
    fi

    print_build_info

    if ! run_e2e_tests; then
        log_error "Build aborted due to test failure"
        return 1
    fi

    if ! export_binaries; then
        log_error "Build failed during binary export"
        return 1
    fi

    log_success "Build completed successfully"
    log_info "Image tagged as: ${IMAGE_NAME}:${BUILD_VER}"
    return 0
}

main "$@"
