# Use the official Golang image as a build stage
FROM golang:1.23.2-alpine AS builder

# update and install dependencies
RUN apk add --no-cache gcc musl-dev openssl

# Install gcc and build-essential (which includes gcc)
# RUN apt-get update && apt-get install -y \
#     gcc \
#     build-essential \
#     && rm -rf /var/lib/apt/lists/*


# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Generate a self-signed certificate and private key
# You can customize the subject as needed
RUN openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout server.key -out server.crt \
    -subj "/C=US/ST=State/L=City/O=Organization/OU=Unit/CN=localhost:8080"

# Build the Go app
# RUN go build -o forum cmd/app/main.go
# RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o forum cmd/app/main.go
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o forum cmd/app/main.go

# Update
# RUN apk update
# RUN apt-get update && apt-get install -y \
#     gcc \
#     sqlite3 libsqlite3-dev \
#     && rm -rf /var/lib/apt/lists/*

# Start a new stage from scratch
FROM alpine:latest
# FROM debian:latest

# update and install dependencies
RUN apk update && apk add --no-cache sqlite openssl

# RUN apt-get update && apt-get install -y \
#     openssl \
#     && rm -rf /var/lib/apt/lists/*

# Set the Current Working Directory inside the container
WORKDIR /app

RUN echo ">>>>debug"
# debug
RUN ls -lah

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/forum .

# Set executable permissions
RUN chmod +x /app/forum

# COPY --from=builder /app/forum .
# COPY --from=builder /app/.env .
# COPY --from=builder /app/internal/db internal/db
# COPY --from=builder /app/static static/
# COPY --from=builder /app/web web/
# COPY --from=builder /app/generate_cert.sh .
# COPY --from=builder /app/openssl.conf .

COPY --from=builder /app/ .

# Copy the SSL certificate from the root of the project
# COPY server.crt /app/server.crt
# COPY . /app/

# Copy the .env file
# COPY .env /app/.env


# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./forum"]
