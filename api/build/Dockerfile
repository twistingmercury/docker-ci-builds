#------------------------------------------------------------
# Stage 1: Builder - download deps, build binary
#------------------------------------------------------------
FROM ghcr.io/twistingmercury/golang-tooling:v1.0-go1.25-alpine AS builder

ARG BUILD_VER=dev
ARG BUILD_DATE=unknown
ARG BUILD_COMMIT=unknown
ARG BIN_NAME=uuid_api

WORKDIR /workspace

# Copy api module files and source
COPY api/go.mod api/go.sum ./
COPY api/ .

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
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "${LD_FLAGS}" -o /out/${BIN_NAME} ./cmd/main/main.go


#------------------------------------------------------------
# Stage 2: Runtime - minimal production image
#------------------------------------------------------------
FROM scratch AS runtime

COPY --from=builder /out/${BIN_NAME}  /${BIN_NAME}

EXPOSE 8080

ENTRYPOINT ["/uuid_api"]
