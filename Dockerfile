# --- Stage 1: Builder ---
    FROM golang:1.23-alpine AS builder
    WORKDIR /app
    
    # Install git for go mod and air install
    RUN apk add --no-cache git
    
    # Copy go mod files and download deps
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copy the entire source (cmd/, internal/, etc.)
    COPY . .
    
    # Install Air for hot reload (dev only)
    RUN go install github.com/air-verse/air@latest
    
    # Build the production binary
    RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/app ./cmd/server
    
    # --- Stage 2: Development (hot-reload with Air) ---
    FROM golang:1.23-alpine AS dev
    WORKDIR /app
    
    # Copy air binary from builder
    COPY --from=builder /go/bin/air /usr/local/bin/air
    
    # Copy source (will be volume-mounted in dev)
    COPY . .
    
    EXPOSE 8080
    
    # Entrypoint for development
    CMD ["air"]
    
    # --- Stage 3: Production ---
    FROM alpine:3.19 AS prod
    RUN apk --no-cache add ca-certificates
    
    WORKDIR /root/
    
    COPY --from=builder /app/bin/app .
    
    # If you need to copy a prod .env file, do it here.
    # COPY .env.production .env
    
    EXPOSE 8080
    
    CMD ["./app"]
    