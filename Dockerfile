# syntax=docker/dockerfile:1
FROM golang:1.22.1

# Set the current working directory inside the container
WORKDIR /api

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
# Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Build the Go app
RUN go build -o main cmd/api/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]

# Health check to verify API is up and running
HEALTHCHECK --interval=30s --timeout=30s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1
  