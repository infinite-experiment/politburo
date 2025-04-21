# Stage 1: Builder
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Install Git (needed for go install of air)
RUN apk add --no-cache git

# Copy Go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code (includes cmd/, internal/, pkg/, etc.)
COPY . .

# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest

# Build binary for production
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/app ./cmd/server

# Stage 2: Development with Air
FROM golang:1.23-alpine AS dev
WORKDIR /app

# Copy Air binary from builder
COPY --from=builder /go/bin/air /usr/local/bin/air

# Copy source code (this will be overridden by the volume in docker-compose)
COPY . .

EXPOSE 8080

# Entrypoint for dev
CMD ["air"]

# Stage 3: Production Image
FROM alpine:latest AS prod
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/bin/app .

# Optionally copy prod .env
COPY .env.production .env

EXPOSE 8080
CMD ["./app"]
