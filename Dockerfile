# 1. Start from the official Go image
FROM golang:1.24-bookworm As builder

# 2. Set the working directory inside the container
WORKDIR /app

# 3. Copy the Go code into the container
COPY go.mod ./
RUN go mod download
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:3.20

# Set the working directory
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Change ownership to non-root user
RUN chown appuser /app/main

# Switch to non-root user some changes to test CI job
USER appuser

EXPOSE 8080

# Start app
CMD ["./main", "8080"]
