# Start with a base Go image
FROM golang:1.22 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -o service-catalog ./cmd/app

# Use a minimal base image to reduce the image size
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/service-catalog .

# Expose port 80 to the outside world
EXPOSE 80

# Command to run the executable
CMD ["./service-catalog"]