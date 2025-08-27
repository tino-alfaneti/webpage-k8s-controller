# -------------------------
# Stage 1: Build Go binary
# -------------------------
FROM golang:1.24 AS builder

WORKDIR /workspace

# Download dependencies first (caching benefit)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . ./

# Build a statically linked Linux binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o main


# -------------------------
# Stage 2: Minimal runtime
# -------------------------
FROM gcr.io/distroless/base:nonroot


# Copy binary from builder stage
COPY --from=builder /workspace/main /workspace/main

# Run as non-root user (distroless base uses 65532 by default)
USER 65532:65532

# Set entrypoint
CMD ["/workspace/main"]
