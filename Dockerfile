
FROM golang:1.24-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN go build -ldflags="-s -w" -o /out/app .

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /out/app /app/app
COPY validation_messages.yaml /app/validation_messages.yaml

ENV PORT=7032
EXPOSE 7032
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
CMD ["server"]

