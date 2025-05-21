# --- Stage 1: Builder ---
    FROM golang:1.23-alpine AS builder
    WORKDIR /app
    
    # install git for go mod and air
    RUN apk add --no-cache git
    
    # download deps
    COPY go.mod go.sum ./
    RUN go mod download
    
    # copy everything (for builder only)
    COPY . .
    
    # install Air CLI
    RUN go install github.com/air-verse/air@latest
    
    # build your server into /app/bin/app
    RUN CGO_ENABLED=0 GOOS=linux go build -o bin/app ./cmd/server
    
    # --- Stage 2: Development (hot reload with Air) ---
    FROM golang:1.23-alpine AS dev
    WORKDIR /app
    
    # only need the Air binary here
    COPY --from=builder /go/bin/air /usr/local/bin/air
    
    EXPOSE 8080
    CMD ["air", "-c", "air.toml"]
    
    # --- Stage 3: Production ---
    FROM alpine:3.19 AS prod
    RUN apk --no-cache add ca-certificates
    
    WORKDIR /root
    
    # copy in the compiled binary
    COPY --from=builder /app/bin/app .
    
    EXPOSE 8080
    CMD ["./app"]
    