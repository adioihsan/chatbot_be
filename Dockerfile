# Dockerfile.dev
FROM golang:1.25-alpine

RUN apk add --no-cache bash git ca-certificates build-base curl && \
    adduser -D -u 1000 vscode

# Make sure Go is on PATH for all users
ENV PATH="/usr/local/go/bin:/go/bin:${PATH}"
ENV CGO_ENABLED=0 GOFLAGS=-mod=mod

# Air + Delve (Air moved to air-verse)
RUN GOBIN=/usr/local/bin go install github.com/air-verse/air@latest && \
    GOBIN=/usr/local/bin go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /work
USER vscode
