# Use a base image with build dependencies
FROM golang:1.23 AS builder

# Set environment variables to enable CGO for x86_64
ENV CGO_ENABLED=1 \
    GOOS=linux

# Set the working directory inside the container
WORKDIR /app

# Copy the Go project files into the container
COPY . .

# Build the Go binary with CGO enabled
RUN go build -o ./services/users/build/users ./services/users/cmd/users/main.go

ENTRYPOINT ["/app/services/users/build/users"]