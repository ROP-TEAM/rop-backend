# rop-backend

Go REST API for vehicle routing optimization — Fiber v3 + GORM + PostgreSQL.

## Requirements

- Go 1.24+
- PostgreSQL (GORM auto-migrates on startup)

## Quick Start

```bash
git clone https://github.com/ROP-TEAM/rop-backend.git
cd rop-backend

# Create .env (see below)
cp .env.example .env

# Run
go run main.go          # :8080
# or: air                # hot-reload
```

## Environment Setup

`.env` (not committed):

```env
APP_PORT=8080
JWT_SECRET=change-me
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME=ROP_DB
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=

# Distance Matrix Provider
DISTANCE_MATRIX_PROVIDER=google    # "google" (default) or "osrm"
GOOGLE_MAPS_API_KEY=               # required for google
OSRM_BASE_URL=http://localhost:5000 # required for osrm

# Solver
SOLVER_BINARY_PATH=                # empty → StubSolver; path → gRPC solver

# OTP (SMS)
OTP_APP_KEY=
OTP_APP_SECRET=
```

## Architecture

```
main.go → config/ → database/ → router/
              ├── handlers/     parse request, call service; Swagger annotations
              ├── services/     business logic (planning, solver calls, matrix factory)
              ├── repository/   GORM queries
              ├── middleware/   JWT auth, rate limiter, OTP limiter
              ├── models/       GORM structs
              ├── dto/          request/response DTOs
              └── utils/        shared helpers
```

## Distance Matrix Provider

Switch via `.env` — no code change:

```go
// internal/services/matrix_factory.go
matrix, err := services.NewDistanceMatrix(cfg)
// returns gmap.GoogleMapsMatrix or osrm.OSRMMatrix based on DISTANCE_MATRIX_PROVIDER
```

| Provider | Env value | API cost |
|---|---|---|
| Google Maps | `google` | Paid per request |
| OSRM | `osrm` | Free (self-hosted) |

## Development Workflow

### Dev Flow — daily, PR, local

Uses Go `replace` directive — clone both repos side by side:

```
projects/
├── rop-algorithm/    ← github.com/ROP-TEAM/rop-algorithm
└── rop-backend/      ← github.com/ROP-TEAM/rop-backend
                        go.mod: replace => ../rop-algorithm
```

```bash
git clone https://github.com/ROP-TEAM/rop-algorithm.git
git clone https://github.com/ROP-TEAM/rop-backend.git
cd rop-backend && go build ./...
```

No tag, no GOPRIVATE, no PAT needed. Changes to rop-algorithm are immediately visible.

CI (`.github/workflows/ci.yml`): checks out both repos → `go build` on every PR.

### Release Flow — production deploy

1. **Tag rop-algorithm:**
   ```bash
   cd rop-algorithm
   git tag v0.1.0 && git push origin v0.1.0
   ```
   → CI creates GitHub Release

2. **Push to rop-backend dev (or tag):**
   ```bash
   git push origin dev
   ```
   → CI fetches latest release tag → Docker build → push GHCR → kubectl deploy

Workflow: `.github/workflows/build-and-deploy.yml`

Note: rop-algorithm must be **public** for release flow (or configure GOPRIVATE + PAT).

## CI

| Workflow | Trigger | Action |
|---|---|---|
| `ci.yml` | PR, push dev | Checkout both repos → `go build` |
| `build-and-deploy.yml` | push dev, tag `v*` | Pin release tag → Docker build → deploy |

## API Docs

```bash
swag init --parseDependency --parseInternal   # regenerate docs
```

Swagger UI: `http://localhost:8080/api/swagger/index.html`

## Commands

```bash
go run main.go                              # dev server :8080
air                                         # hot-reload (uses .air.toml)
go build -o ROP_Backend.exe .
go mod tidy
swag init --parseDependency --parseInternal # regenerate Swagger
go test ./...                               # (no tests yet)
```
