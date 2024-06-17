# Use the official Golang image as a build stage
FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o ./bin/cmd ./cmd

# Use a smaller base image for the final container
FROM ubuntu:22.04

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/bin/cmd /app/bin/cmd

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./bin/cmd", "start", "--layout", "plain"]
