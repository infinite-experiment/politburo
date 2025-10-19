# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Politburo** is a Go-based backend for the Infinite Experiment Discord bot and web client. It provides REST APIs for managing a virtual airline system, integrating with Infinite Flight Live API, Airtable for data storage, and PostgreSQL for persistent storage.

The application serves flight tracking data, user management, role-based access control, and synchronization with external services (Airtable, Live API).

## Development Commands

### Local Development (with hot reload)

```bash
# Start with Docker Compose (recommended)
docker-compose -f docker-compose.local.yml up --build

# Stop services
docker-compose -f docker-compose.local.yml down

# Rebuild single service
docker-compose -f docker-compose.prod.yml up -d --build api
```

### Native Development (without Docker)

```bash
# Install Air for hot reloading (if not installed)
go install github.com/air-verse/air@latest

# Run with hot reload using Air
air

# Manual build and run
go build -o .air_tmp/main ./cmd/server
.air_tmp/main

# Using the run script (builds + generates swagger + runs)
./run-script.sh
```

### Building

```bash
# Development build
go build -buildvcs=false -o .air_tmp/main ./cmd/server

# Production Docker build
docker build --target prod -t politburo:latest .
```

### Swagger Documentation

```bash
# Generate Swagger docs (required after API changes)
swag init -g cmd/server/main.go --output docs

# Access Swagger UI at http://localhost:8080/swagger/index.html
```

### Utilities

```bash
# Generate new API key
go run ./cmd/api_key_gen/main.go
```

### Go Module Management

```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy
```

## Architecture

### Entry Points

- **cmd/server/main.go**: Main HTTP server, connects to PostgreSQL and initializes routing
- **cmd/api_key_gen/main.go**: Utility to generate API keys for authentication

### Routing & Middleware Stack (internal/routes/router.go)

The application uses **Chi router** with nested route groups implementing a hierarchical role-based access control system:

1. **Public routes**: No authentication required (`/public/*`, `/healthCheck`, `/swagger/*`)
2. **API v1 routes** (`/api/v1`): All require authentication via `AuthMiddleware`
   - **Registered users**: Base access level (registration, server initialization)
   - **Member**: Access to live flights, sessions
   - **Staff**: User flight queries, user sync operations
   - **Admin**: Role management, VA configuration, debug endpoints
   - **God**: Special role with delete permissions

### Authentication System

**Two authentication methods** (internal/middleware/auth.go):

1. **JWT (Bearer tokens)**: Currently stubbed out, returns 401
2. **API Keys**: Production method using headers:
   - `X-API-Key`: The API key
   - `X-Server-Id`: Discord server ID
   - `X-Discord-Id`: Discord user ID

**Claims Interface** (internal/auth/claims.go):

The `UserClaims` interface abstracts authentication sources with implementations:
- `JWTClaims`: For future JWT support
- `APIKeyClaims`: For API key authentication

Claims are stored in request context and accessed via helper functions.

### Role System (internal/constants/roles.go)

Three-tier role hierarchy stored as Postgres ENUM:
- `pilot`: Basic member access
- `staff`: Airline manager permissions
- `admin`: Full administrative access

The `VARole` type implements `sql.Scanner` and `driver.Valuer` for seamless database integration.

### Service Layer Architecture

Services are initialized in `RegisterRoutes` and follow dependency injection:

**Core Services** (internal/common/):
- `CacheService`: In-memory cache using `go-cache` (60000ms default, 600s cleanup)
- `LiveAPIService`: Integration with Infinite Flight Live API
- `VAConfigService`: Virtual airline configuration management
- `AirtableApiService`: Airtable API integration for external data sync

**Business Services** (internal/services/):
- `RegistrationService`: User and server registration workflows
- `VAManagementService`: VA role and user management
- `AtSyncService`: Airtable synchronization logic
- `FlightsService`: Flight data aggregation and queries

**Repositories** (internal/db/repositories/):
- `UserRepository`: User CRUD operations
- `ApiKeysRepo`: API key validation
- `SyncRepository`: Sync history tracking
- `VARepository`: Virtual airline data access

### Background Workers

Two workers start as goroutines in `RegisterRoutes`:

1. **LogbookWorker** (internal/workers/logbook_worker.go):
   - Consumes from `LogbookQueue` channel (buffer: 100)
   - Fetches flight routes from Live API
   - Processes waypoints, calculates stats (max speed, altitude)
   - Caches complete flight info for ~7 days

2. **StartCacheFiller** (internal/workers/meta_cache_worker.go):
   - Pre-fills cache with frequently accessed data
   - Reduces API calls to external services

### Database Layer

**Connection** (internal/db/postgres.go):
- Uses `sqlx` for enhanced SQL operations
- Implements retry logic (10 attempts, 500ms intervals)
- Connection configured via environment variables

**Migrations** (internal/db/migrations/):
Applied manually in sequential order:
1. `001_init_database.sql`: Base schema
2. `002_setup_users.sql`: User tables
3. `003_va_setup_base.sql`: Virtual airline setup
4. `004_config_setup.sql`: Configuration tables
5. `005_pireps.sql`: PIREP (flight reports) tables

### Synchronization Jobs (internal/jobs/at_sync_job.go)

Two pagination-based sync jobs for Airtable integration:
- `SyncPilotsJob`: Syncs pilot records with optional last-modified filtering
- `SyncRoutesJob`: Syncs route records with offset-based pagination

Both use the claims context to determine which VA server to sync for.

## Environment Configuration

Required environment variables (see `.env.local` example):

```bash
APP_ENV=local          # Environment: local/production
DEBUG=true             # Debug mode
PORT=8080              # HTTP server port
PG_HOST=db             # PostgreSQL host
PG_PORT=5432           # PostgreSQL port
PG_USER=ieuser         # Database user
PG_DB=infinite         # Database name
PG_PASSWORD=iepass     # Database password
```

## Key Patterns & Conventions

### Dependency Injection
Services and repositories are constructed in `RegisterRoutes` and injected into handlers. Avoid global service instances.

### Context-Based Claims
User claims are attached to `http.Request` context via `auth.SetUserClaims` and retrieved with `auth.GetUserClaims`.

### Channel-Based Workers
Use buffered channels (e.g., `LogbookQueue`) for async processing. Workers run for the application lifetime.

### Cache Keys
Flight cache keys follow pattern: `{cacheKey}` constructed per flight to ensure uniqueness.

### Error Handling
HTTP errors use standard `http.Error()`. Database retries are handled in connection initialization only.

### Swagger Annotations
API handlers use Swaggo annotations. Regenerate docs after changes with `swag init -g cmd/server/main.go --output docs`.

## Important Files

- **internal/routes/router.go**: Complete routing definition and service wiring
- **internal/middleware/auth.go**: Authentication logic
- **internal/auth/claims.go**: Claims interface and implementations
- **internal/constants/roles.go**: Role definitions with DB adapters
- **internal/workers/logbook_worker.go**: Flight data processing worker
- **cmd/server/main.go**: Application entry point
