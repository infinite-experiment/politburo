# Stage 1: Builder Stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Install Git (required for 'go install')
RUN apk add --no-cache git

# Copy dependency files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Install Air for hot reloading (Air will be installed in /go/bin/air by default)
RUN go install github.com/cosmtrek/air@latest

# Build the production binary (compiled without Air)
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/app ./app

# Stage 2: Development Image (with Hot Reload)
FROM golang:1.23-alpine AS dev
WORKDIR /app

# Copy Air binary from the builder stage
COPY --from=builder /go/bin/air /usr/local/bin/air

# Copy source code (this will be overridden by docker-compose volume mount)
COPY . .

EXPOSE 8080

# Use Air as the entrypoint for hot reloading
CMD ["air"]

# Stage 3: Production Image
FROM alpine:latest AS prod
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/bin/app .

# Optionally copy production environment file into the image
COPY .env.production .env

EXPOSE 8080
CMD ["./app"]
