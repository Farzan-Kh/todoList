# Stage 1: Build
FROM golang:1.20-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -o main .

# Stage 2: Runtime
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the port the service runs on
EXPOSE 3002

# Command to run the application
CMD ["./main"]