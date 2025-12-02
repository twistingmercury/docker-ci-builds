#------------------------------------------------------------
# Library CI: Test Only (no binary artifact)
#------------------------------------------------------------
FROM golang:1.23-alpine

RUN apk add --no-cache git gcc musl-dev

WORKDIR /workspace

# Copy dependency files first for layer caching
COPY go.mod go.sum* ./
RUN go mod download

# Copy source
COPY . .

# Run tests with coverage (CGO_ENABLED=1 for race detector)
RUN CGO_ENABLED=1 go test -v -race -coverprofile=coverage.out ./...
RUN go tool cover -func=coverage.out

# Verify the package builds (but don't output anything)
RUN go build ./...
