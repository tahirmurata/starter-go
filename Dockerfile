# Use the official Bun image as the base for the builder stage
FROM oven/bun:1 AS bun
WORKDIR /usr/src/app

# Install dependencies with better caching strategy
FROM bun AS bun-deps
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

# Build Tailwind CSS
FROM bun-deps AS bun-build
COPY scripts/build.ts ./scripts/
COPY tailwind.config.js postcss.config.js ./ 
COPY src/ ./src/
RUN bun run build:tailwindcss

# Use the official Go image for building
FROM golang:1.24-alpine3.21 AS go-builder
WORKDIR /usr/src/app

# Install required build tools and security scanning
RUN apk add --no-cache git ca-certificates tzdata && \
    update-ca-certificates

# Install dependencies first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Generate code from templates and SQL
COPY . .
COPY --from=bun-build /usr/src/app/public ./public
RUN go run github.com/a-h/templ/cmd/templ generate && \
    cd sqlc && go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

# Build the application with security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static'" \
    -o main cmd/api/main.go

# Create a minimal production image
FROM alpine:3.21 AS prod
WORKDIR /usr/src/app

# Add non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chown -R appuser:appgroup /usr/src/app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary and config from builder stage
COPY --from=go-builder /usr/src/app/main ./
COPY --from=go-builder /usr/src/app/config.yml ./
COPY --from=bun-build /usr/src/app/public ./public

# Set proper permissions
RUN chmod +x ./main

# Set metadata
LABEL maintainer="Development Team"
LABEL version="1.0"
LABEL description="Go application with Bun frontend tooling"

# Add health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Configure the container
EXPOSE 8080
USER appuser
CMD ["./main"]
