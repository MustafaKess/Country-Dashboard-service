# Use golang base image to build the application
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the entire application into the container
COPY . .

# Run `go mod tidy` to ensure all dependencies are available
RUN go mod tidy

# Build the Go application with static linking (CGO_ENABLED=0 to avoid dynamic linking to libc)
WORKDIR /go/src/app/main
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o server

# Use a minimal Debian image as the base for the final stage
FROM debian:bullseye-slim

# Install necessary dependencies (if any are needed in the final image)
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy the statically built binary from the builder image
COPY --from=builder /go/src/app/main/server /app/

# Copy Firebase credentials file to the correct path inside the container
COPY .env/firebaseKey.json /app/.env/firebaseKey.json

# Set the working directory inside the container
WORKDIR /app

# Expose the port the application will listen on
EXPOSE 8080

# Command to run the application when the container starts
CMD ["./server"]

