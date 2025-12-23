# Build stage
FROM golang:1.25.5-alpine3.23 AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata


# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies (cached layer)
RUN go mod download

# Copy source code
COPY . .

# Build arguments
ARG VERSION=0.0.1
ARG BUILD_DATE
ARG COMMIT_SHA

# Build the binary
# -ldflags: strip debug info, set version info
# -tags: build tags if needed
# CGO_ENABLED=0: static binary, no C dependencies
# TARGETOS and TARGETARCH are automatically set by Buildx for multi-arch builds
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X github.com/jubblin/omni-api/internal/api/handlers.Version=${VERSION} -X main.buildDate=${BUILD_DATE} -X main.commitSha=${COMMIT_SHA}" \
    -a -installsuffix cgo \
    -o omni-api \
    main.go

# Runtime stage - using distroless for minimal attack surface
FROM gcr.io/distroless/static-debian12:nonroot

# Copy CA certificates from builder for TLS connections
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /build/omni-api /usr/local/bin/omni-api

# Set working directory
WORKDIR /app

# Use the nonroot user (uid:gid = 65532:65532)
USER nonroot:nonroot

# Expose default port
EXPOSE 8080

# Health check
# Note: Distroless doesn't include wget/curl, so health checks should be done
# externally via HTTP requests to /health endpoint
# HEALTHCHECK can be configured at runtime or via docker-compose

# Labels for metadata
LABEL org.opencontainers.image.title="Omni API" \
      org.opencontainers.image.description="A REST API to interface with Sidero Omni" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.revision="${COMMIT_SHA}" \
      org.opencontainers.image.source="https://github.com/jubblin/omni-api" \
      org.opencontainers.image.licenses="Apache-2.0"

# Run the binary
ENTRYPOINT ["/usr/local/bin/omni-api"]
