# E-commerce SaaS Platform

A multi-tenant e-commerce SaaS platform built with hexagonal architecture, enabling businesses to create and customize their own online stores with custom domains/subdomains.

## 📋 Documentation

📚 **Complete Documentation**: See individual docs below

### **Quick Links**
- [🏗️ Architecture](./docs/ARCHITECTURE.md) - System architecture and hexagonal design patterns
- [📡 API Documentation](./docs/API_REFERENCE.md) - Complete REST API documentation (450+ endpoints)
- [🚀 API Architecture](./docs/API_ARCHITECTURE.md) - GraphQL, SSE, advanced pagination
- [🗄️ Database Strategy](./docs/DATABASE.md) - Hybrid multi-tenant database design
- [✨ Features](./docs/FEATURES.md) - Complete feature list and pricing tiers  
- [🐳 Docker Setup](./docs/README-DOCKER.md) - Development environment with Docker

## 🚀 Quick Start

### Development with Docker (Recommended)
```bash
# Start all services
make dev-up

# Check status
make status

# View logs  
make dev-logs

# Stop services
make dev-down
```

### Manual Setup
```bash
# Install dependencies
make install

# Start core services
make services-up

# Run each service manually
cd backend && go run cmd/api/main.go
cd storefront && npm run dev  
cd dashboard && npm start
```

## 🏗️ Project Structure

```
├── docs/                    # All documentation
├── backend/                 # Go Fiber API
├── storefront/              # Next.js customer store  
├── dashboard/               # React.js merchant admin
├── docker-compose.dev.yml   # Development environment
└── Makefile                 # Development commands
```

## 🎯 Tech Stack

- **Backend**: Go 1.21+ with Fiber framework
- **Frontend**: Next.js 14+ (storefront) + React 18+ (dashboard)
- **Database**: PostgreSQL 15+ with tenant isolation
- **Cache**: Redis 7+
- **Storage**: MinIO (S3-compatible)
- **State**: Jotai for client state management
- **Styling**: Tailwind CSS
- **Deployment**: Docker containers

## 📊 Services Overview

| Service | Port | Description |
|---------|------|-------------|
| Backend API | 8080 | Go Fiber REST API |
| Storefront | 3000 | Customer-facing store |
| Dashboard | 3001 | Merchant admin panel |
| PostgreSQL | 5432 | Primary database |
| Redis | 6379 | Cache & sessions |
| MinIO | 9000 | File storage |
| MailHog | 8025 | Email testing |

## 🔧 Development

See [Docker Setup Guide](./docs/README-DOCKER.md) for complete development environment setup.

### Available Commands
```bash
make help           # Show all commands
make dev-up         # Start development environment
make dev-down       # Stop all services
make dev-logs       # View logs
make install        # Install dependencies
make test           # Run tests
make format         # Format code
```

## 📚 Additional Resources

- **Architecture**: Hexagonal (Clean) Architecture with Go
- **Multi-tenancy**: Database-per-tenant isolation
- **Authentication**: JWT with role-based access
- **Payments**: Stripe integration
- **Deployment**: AWS ECS with Docker

---

**Status**: Development Setup Complete  
**Next**: Implement core business logic