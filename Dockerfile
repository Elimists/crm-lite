# Use the official Go image for building
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy Go modules files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of your app
COPY . .

# Build the binary
RUN go build -o crm-lite .

# Use a minimal image for runtime
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the built binary
COPY --from=builder /app/crm-lite .

# Expose the port your app listens on
EXPOSE 8080

# Run the app
CMD ["./crm-lite"]
