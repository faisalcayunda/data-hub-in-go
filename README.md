# Portal Data Backend

A **Modular Monolith by Feature** backend application built with **Go**, following **Clean Architecture** principles (based on [bxcodec/go-clean-arch](https://github.com/bxcodec/go-clean-arch)). Each feature encapsulates its own domain, use case, repository, and delivery layer.

## Architecture

This project uses **Go Clean Architecture** where:
- **Domain Layer** - Pure Go structs with no external dependencies
- **Usecase Layer** - Business logic with interface definitions
- **Repository Layer** - Interface in domain, implementation in infrastructure
- **Delivery Layer** - HTTP handlers for REST API
- **Infrastructure Layer** - Database, config, security, HTTP server

### Dependency Flow

```
┌─────────────┐
│   cmd/      │ Application entry point
└──────┬──────┘
       │
       ├─────────────────────────────────────┐
       ▼                                     ▼
┌──────────────┐                      ┌──────────────┐
│  internal/   │                      │  pkg/        │
│              │                      │  errors/     │
│  auth/       │                      │  validator/  │
│  user/       │                      │  utils/      │
│  dataset/    │                      └──────────────┘
│  tag/        │
│  ...         │
│              │
│ domain/      │ ←── Core business logic
│ usecase/     │
│ repository/  │
│ delivery/    │
└──────────────┘
       │
       ▼
┌──────────────────┐
│ infrastructure/  │ External dependencies
├──────────────────┤
│  db/            │ PostgreSQL
│  http/          │ Chi router
│  security/      │ JWT, Password
│  config/        │ Configuration
│  logger/        │ Logging
└──────────────────┘
```

### Directory Structure

```
portal-data-backend/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
│
├── internal/
│   ├── auth/                    # Authentication feature
│   │   ├── domain/              # Entities, repository interfaces
│   │   ├── usecase/             # Business logic
│   │   ├── repository/          # Repository implementations
│   │   └── delivery/http/       # HTTP handlers
│   │
│   ├── user/                    # User management
│   ├── organization/            # Organization management
│   ├── dataset/                 # Dataset management
│   ├── tag/                     # Tag management
│   └── ...
│
├── infrastructure/
│   ├── config/                  # Configuration loading
│   ├── db/                      # Database connection
│   ├── http/                    # HTTP server & middleware
│   ├── security/                # JWT & Password hashing
│   └── logger/                  # Structured logging
│
├── pkg/                         # Public reusable packages
│   ├── errors/                  # Sentinel errors
│   ├── validator/               # Request validation
│   └── utils/                   # Utility functions
│
├── migrations/                  # Database migrations
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL (for database)
- Redis (for caching - optional)

### Installation

```bash
# Clone repository
git clone <repository-url>
cd portal-data-backend

# Download dependencies
go mod download

# Setup environment
cp env.example .env
# Edit .env with your configuration

# Run database migrations (using golang-migrate)
migrate -path migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up

# Build the application
go build -o bin/portal-data-backend cmd/server/main.go

# Run the application
./bin/portal-data-backend
# Or: go run cmd/server/main.go
```

Visit `http://localhost:8080/health` for health check.

### Using Makefile

```bash
# Download dependencies
make deps

# Run the application
make run

# Build binary
make build

# Run tests
make test
make test-unit
make test-integration
make test-coverage

# Lint and format
make fmt
make vet
```

## API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/auth/register` | Register new user | No |
| POST | `/auth/login` | Login user | No |
| POST | `/auth/logout` | Logout user | No |
| POST | `/auth/refresh` | Refresh access token | No |
| POST | `/auth/revoke-all` | Revoke all user tokens | Yes |
| GET | `/me` | Get current user | Yes |

### Users

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/users` | List users (paginated) | Yes |
| GET | `/users/{id}` | Get user by ID | Yes |
| PUT | `/users/{id}` | Update user | Yes |
| DELETE | `/users/{id}` | Delete user | Yes |
| PATCH | `/users/{id}/status` | Update user status | Yes |

### Organizations

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/organizations` | List organizations | No |
| GET | `/organizations/{id}` | Get organization by ID | No |
| GET | `/organizations/code/{code}` | Get organization by code | No |
| POST | `/organizations` | Create organization | Yes |
| PUT | `/organizations/{id}` | Update organization | Yes |
| DELETE | `/organizations/{id}` | Delete organization | Yes |
| PATCH | `/organizations/{id}/status` | Update status | Yes |

### Datasets

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/datasets` | List datasets | No |
| GET | `/datasets/{id}` | Get dataset by ID | No |
| GET | `/datasets/slug/{slug}` | Get dataset by slug | No |
| POST | `/datasets` | Create dataset | Yes |
| PUT | `/datasets/{id}` | Update dataset | Yes |
| DELETE | `/datasets/{id}` | Delete dataset | Yes |
| PATCH | `/datasets/{id}/status` | Update status | Yes |

### Tags

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/tags` | List tags | No |
| GET | `/tags/{id}` | Get tag by ID | No |
| POST | `/tags` | Create tag | Yes |
| PUT | `/tags/{id}` | Update tag | Yes |
| DELETE | `/tags/{id}` | Delete tag | Yes |

## Example Request/Response

### Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "code": "OPERATION_SUCCESSFUL",
  "message": "Login successful",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "organization_id": "660e8400-e29b-41d4-a716-446655440000",
      "role_id": "admin",
      "name": "John Doe",
      "username": "johndoe",
      "email": "user@example.com"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900,
    "token_type": "Bearer"
  }
}
```

### Authenticated Request

```bash
curl -X GET http://localhost:8080/users \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

## Testing

### Run Tests

```bash
# Run all tests
go test ./... -v

# Run unit tests only
go test ./internal/... -short -v

# Run with coverage
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Test Structure

```
internal/
├── auth/
│   ├── usecase/
│   │   └── auth_usecase_test.go    # Unit tests with mocks
│   └── repository/
│       └── postgres_test.go         # Integration tests
```

## Adding a New Feature

1. Create domain entities in `internal/<feature>/domain/entity.go`
2. Define repository interface in `internal/<feature>/domain/repository.go`
3. Implement usecases in `internal/<feature>/usecase/`
4. Implement repository in `internal/<feature>/repository/postgres.go`
5. Create HTTP handlers in `internal/<feature>/delivery/http/handler.go`
6. Wire dependencies in `cmd/server/main.go`

## Configuration

Configuration is managed via environment variables (see `env.example`):

```bash
# App
APP_NAME=portal-data-backend
APP_ENV=development
APP_DEBUG=true

# Server
SERVER_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=portal_data

# JWT
JWT_SECRET=your-secret-key
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h
```

## Deployment

### Build for Production

```bash
# Build optimized binary
go build -ldflags="-s -w" -o bin/portal-data-backend cmd/server/main.go

# Run
./bin/portal-data-backend
```

### Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o bin/portal-data-backend cmd/server/main.go

FROM alpine:latest
COPY --from=builder /app/bin/portal-data-backend /usr/local/bin/
EXPOSE 8080
CMD ["portal-data-backend"]
```

## License

MIT License - see LICENSE file for details
