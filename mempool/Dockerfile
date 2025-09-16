# Use the official Go image for building
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy files
COPY go.mod go.sum .
RUN go mod download

# Define the command to run your server
CMD ["go", "run", "main.go"]