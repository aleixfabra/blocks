# Use the official Go image for building
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy files
COPY go.mod .
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -o main

# Use a smaller base image for running
FROM alpine:3.17

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the port the application runs on
EXPOSE 8080

# Define the command to run your server
CMD ["./main"]