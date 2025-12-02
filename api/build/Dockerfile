#------------------------------------------------------------
# Stage 1: Builder - download deps, build binary
#------------------------------------------------------------
FROM golang:1.23-alpine AS builder

ENV GOTOOLCHAIN=auto

RUN apk add --no-cache git ca-certificates tzdata

ARG BUILD_VER=dev
ARG BUILD_DATE=unknown
ARG BUILD_COMMIT=unknown

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG LD_FLAGS="-s -w \
    -X 'main.Version=${BUILD_VER}' \
    -X 'main.BuildDate=${BUILD_DATE}' \
    -X 'main.GitCommit=${BUILD_COMMIT}'"

RUN go test -v ./...
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux \
    go build -ldflags "${LD_FLAGS}" -o /out/time-api ./main.go


#------------------------------------------------------------
# Stage 2: Runtime - minimal production image
#------------------------------------------------------------
FROM alpine AS runtime

RUN apk add --no-cache ca-certificates tzdata wget

COPY --from=builder /out/time-api /time-api

EXPOSE 8080

ENTRYPOINT ["/time-api"]
