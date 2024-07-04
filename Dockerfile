# Using the official Golang image to build the application
FROM golang:1.22.4-alpine as builder

# Set the current working directory
WORKDIR /

# Copy the Go module files and download dependencies if they exist
COPY go.mod ./
RUN if [ -f go.sum ]; then cp go.sum ./; fi
RUN go mod download || true

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -o /liftmetrics ./cmd/liftmetrics

# Use a smaller base image for the final image
FROM alpine:latest

# Set environment variable with a default value
ENV BASE_DIR=/data

# Copy the built binary from the builder
COPY --from=builder /liftmetrics /liftmetrics

# Run the binary
CMD ["/liftmetrics"]