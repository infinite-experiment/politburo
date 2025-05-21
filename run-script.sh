#!/usr/bin/env bash
set -e

# 1) load .env if present
if [ -f .env ]; then
  set -a
  source ./.env
  set +a
fi

# 2) build (regular invocation, not exec)
go build -o .air_tmp/main ./cmd/server

# 3) Build swagger
~/go/bin/swag init -g cmd/server/main.go --output docs

# 4) now replace this shell with your binary so logs stick
exec .air_tmp/main
