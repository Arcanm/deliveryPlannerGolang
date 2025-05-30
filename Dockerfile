# Build stage
FROM golang:1.24.1-alpine AS builder

# Install git
RUN apk add --no-cache git

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o deliveryPlannerGolang ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /build/deliveryPlannerGolang .

# Expose ports
EXPOSE 8080 50051

# Run the application
CMD ["./deliveryPlannerGolang"]