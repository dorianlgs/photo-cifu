# Build stage
FROM node:20-alpine AS frontend-builder

WORKDIR /app/ui
COPY ui/package*.json ui/yarn.lock* ./
RUN yarn install --frozen-lockfile

COPY ui/ ./
RUN yarn run build

# Go build stage
FROM golang:1.24-alpine AS go-builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Copy built frontend from previous stage
COPY --from=frontend-builder /app/ui/build ./ui/build

# Generate embedded assets and build the application
RUN go generate ./...
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o photo-cifu .

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary from builder stage
COPY --from=go-builder /app/photo-cifu .

# Create directory for PocketBase data
RUN mkdir -p pb_data

# Expose PocketBase default port
EXPOSE 8090

# Set up health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8090/api/health || exit 1

# Run the application
CMD ["./photo-cifu", "serve", "--http=0.0.0.0:8090"]