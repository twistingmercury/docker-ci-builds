FROM golang:1.23-alpine

ENV GOTOOLCHAIN=auto

RUN apk add --no-cache curl bash

WORKDIR /e2e

COPY go.mod go.sum ./
RUN go mod download

COPY api_test.go ./
COPY test-runner.sh /test-runner.sh

RUN chmod +x /test-runner.sh

ENTRYPOINT ["/test-runner.sh"]
