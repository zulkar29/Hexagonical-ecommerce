# E-commerce SaaS Backend

Go Fiber backend with hexagonal architecture for multi-tenant e-commerce SaaS platform using PostgreSQL.

## 🏗️ Architecture

```
backend/
├── cmd/                    # Application entry points
│   ├── api/               # Main API server
│   └── migrate/           # Database migration tool
├── internal/              # Private application code
│   ├── domain/           # Business logic (entities, value objects)
│   ├── application/      # Application services
│   ├── ports/            # Interface definitions
│   │   ├── input/        # Use cases (what app can do)
│   │   └── output/       # External dependencies (what app needs)
│   ├── adapters/         # External service implementations
│   │   ├── http/         # HTTP handlers (Fiber)
│   │   ├── database/     # PostgreSQL repositories (GORM)
│   │   ├── cache/        # Redis adapters
│   │   └── email/        # Email service adapters
│   └── infrastructure/   # Framework and external concerns
│       ├── config/       # Configuration management
│       ├── database/     # PostgreSQL connection
│       ├── middleware/   # HTTP middleware
│       └── routes/       # Route definitions
└── pkg/                  # Public utilities
    ├── utils/            # Helper functions
    ├── errors/           # Error handling
    └── validation/       # Input validation
```

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Redis 7+

### Setup
```bash
# Clone and navigate to backend
cd backend

# Copy environment file
cp .env.example .env

# Install dependencies (after implementing go.mod)
# go mod tidy

# Run database migrations (after implementation)
# make migrate

# Start development server (after implementation)
# make run
```

## 📝 Implementation Status

All files contain TODO comments for actual implementation:

### ✅ Completed Structure
- Project structure and directories
- Basic file templates with TODO placeholders
- Hexagonal architecture layout
- Configuration templates

### 🚧 TODO Implementation
- Actual Go Fiber setup
- GORM PostgreSQL models
- JWT authentication
- Business logic implementation
- API handlers with request/response
- Middleware implementation
- Database repositories
- Error handling
- Input validation
- Tests

## 🔧 Development Commands

```bash
# Build application
make build

# Run development server
make run

# Run tests
make test

# Run database migrations
make migrate

# Build Docker image
make docker-build

# Clean build artifacts
make clean
```

## 🐳 Docker

```bash
# Build image
docker build -t ecommerce-saas-backend .

# Run container
docker run -p 8080:8080 ecommerce-saas-backend
```

## 📚 Key Features (To Implement)

- Multi-tenant architecture
- JWT authentication
- Role-based access control
- RESTful API with Fiber
- PostgreSQL database per tenant strategy
- Redis caching
- Stripe payment integration
- Email notifications
- File upload to S3
- Rate limiting
- Request validation
- Structured logging

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests
go test -tags=integration ./...
```

This structure provides a solid foundation for implementing the e-commerce SaaS backend with clean architecture principles.