# Use an official Golang image as a builder
FROM golang:1.24 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the entire project
COPY . .

# Build the Go application (update path if main.go is inside a subfolder)
RUN go build -o model-consumer ./cmd/app

# Use a smaller image for the final container
FROM golang:1.24

# Set working directory
WORKDIR /app

# Copy the compiled binary
COPY --from=builder /app/model-consumer .

# Expose port 8081
EXPOSE 8081

# Run the application
CMD ["./model-consumer"]