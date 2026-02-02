# STAGE 1: Build the binary
FROM golang:1.25-alpine AS builder

# Install git (needed for fetching some go modules)
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the app into a binary called 'kozm-server'
RUN go build -o /kozm-server ./cmd/server/main.go

# STAGE 2: Create the final lean image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy only the binary from the builder stage
COPY --from=builder /kozm-server .

# Expose the port your server runs on
EXPOSE 8080

# Run it!
CMD ["./kozm-server"]