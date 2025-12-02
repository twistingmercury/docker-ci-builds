FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git bash

ARG BUILD_VER
ARG BUILD_DATE
ARG BUILD_COMMIT

ARG LD_FLAGS="-s -w -X 'github.com/twistingmercury/cli-tool-build/cmd/version.version=${BUILD_VER}' -X 'github.com/twistingmercury/cli-tool-build/cmd/version.buildDate=${BUILD_DATE}' -X 'github.com/twistingmercury/cli-tool-build/cmd/version.gitCommit=${BUILD_COMMIT}'"

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go test -v ./cmd/version/version_test.go

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build --ldflags "${LD_FLAGS}" -o /out/amd64/darwin/cli-tool ./cmd/main.go
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build --ldflags "${LD_FLAGS}" -o /out/amd64/linux/cli-tool ./cmd/main.go
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build --ldflags "${LD_FLAGS}" -o /out/amd64/windows/cli-tool.exe ./cmd/main.go


#---------------------------------------------------------------------
FROM scratch AS export

COPY --from=builder /out/ .