# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies (none in this case, but good practice)
RUN go mod download

# Copy source code
COPY main.go ./

# Build the binary
RUN go build -o test-runner .

# Runtime stage
FROM alpine:3.19

WORKDIR /

# Copy the binary from builder
COPY --from=builder /build/test-runner /test-runner

# Run the test
ENTRYPOINT ["/test-runner"]
