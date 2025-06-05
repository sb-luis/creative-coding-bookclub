# Build stage (bigger debian image for building)
FROM golang:1.24-bookworm AS build 

WORKDIR /app

# Copy go.mod and go.sum for dependency management
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy application code
COPY . .

# Build the Go application
RUN go build -o main ./cmd/server/main.go

# Runtime stage (smaller debian image for runtime)
FROM debian:bookworm-slim

WORKDIR /app

# Install runtime dependencies
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy the binary from build stage
COPY --from=build /app/main .

# Copy web assets 
COPY --from=build /app/web ./web

EXPOSE 4000

CMD ["./main"]
