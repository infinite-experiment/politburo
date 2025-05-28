# ─── Stage 1: Build ───────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder
WORKDIR /app

# git is needed for go mod
RUN apk add --no-cache git

# download deps
COPY go.mod go.sum ./
RUN go mod download

# copy source & build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/app ./cmd/server

# ─── Stage 2: Production ──────────────────────────────────────────────────────
FROM alpine:3.19
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/bin/app .

# expose your port
EXPOSE 8080

# simple healthcheck
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget --spider --quiet http://localhost:8080/healthCheck || exit 1

# run the compiled binary
ENTRYPOINT ["./app"]
