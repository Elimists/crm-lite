# Builder stage
FROM golang:1.25-alpine AS builder

# Install build tools for cgo
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Enable CGO for go-sqlite3
ENV CGO_ENABLED=1
RUN go build -o crm-lite .

# Final stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/crm-lite .
COPY allowed_origins.json .

# Copy database folder
COPY database ./database

EXPOSE 8080
CMD ["./crm-lite"]
