# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go.mod, go.sum and download dependencies (if they exist)
COPY go.mod go.sum* ./
RUN if [ -f go.sum ]; then go mod download; fi

# Copy source code
COPY . .

# Create dist directory and build the application into it
RUN mkdir -p dist
RUN CGO_ENABLED=0 GOOS=linux go build -o dist/mcp-server .

# Runtime stage
FROM gcr.io/distroless/base-debian11

WORKDIR /

# Copy binary from build stage's dist directory
COPY --from=builder /app/dist/mcp-server /dist/mcp-server

# Command to run with studio subcommand
ENTRYPOINT ["/dist/mcp-server", "studio"]
