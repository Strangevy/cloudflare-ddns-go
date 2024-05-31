ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN apk add --no-cache build-base && \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags '-s -w' -o /app/cloudflare-ddns-go cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/cloudflare-ddns-go /app/
CMD ["./cloudflare-ddns-go"]