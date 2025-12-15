#------------------------------------------------------------
# Stage 1: Builder - download deps, build binary
#------------------------------------------------------------
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

## Install any tooling needed for the build
RUN go install github.com/securego/gosec/v2/cmd/gosec@latest
RUN go install golang.org/x/vuln/cmd/govulncheck@latest
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

ARG BUILD_VER=dev
ARG BUILD_DATE=unknown
ARG BUILD_COMMIT=unknown

WORKDIR /workspace

COPY timelib/ /workspace/timelib/

# Copy api module files and source
COPY api/go.mod api/go.sum ./
COPY api/ .

# Tell Go where to find timelib locally and update checksums
RUN go mod edit -replace github.com/twistingmercury/docker-ci-build/timelib=/workspace/timelib
RUN go mod tidy

## Run linters, formatters, security scanners, etc
RUN goimports -w .
RUN golangci-lint run
RUN govulncheck ./...
RUN gosec ./...


ARG LD_FLAGS="-s -w \
    -X 'main.Version=${BUILD_VER}' \
    -X 'main.BuildDate=${BUILD_DATE}' \
    -X 'main.GitCommit=${BUILD_COMMIT}'"

RUN go test -v ./...
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "${LD_FLAGS}" -o /out/time-api ./main.go


#------------------------------------------------------------
# Stage 2: Runtime - minimal production image
#------------------------------------------------------------
FROM alpine AS runtime

RUN apk add --no-cache ca-certificates tzdata wget

COPY --from=builder /out/time-api /time-api

EXPOSE 8080

ENTRYPOINT ["/time-api"]
