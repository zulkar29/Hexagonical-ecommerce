# Hexagonal E-commerce SaaS Docker Development Environment

📋 **Documentation Navigation**: [📖 Project Home](../README.md) | [🏗️ Architecture](./ARCHITECTURE.md) | [🚢 Deployment](./DEPLOYMENT.md) | [🚀 Features](./FEATURES.md)

Complete Docker setup for the hexagonal e-commerce SaaS platform with all services containerized for easy development and deployment.

## 🚀 Quick Start

### Prerequisites
- Docker Desktop 4.0+ or Docker Engine 20.0+
- Docker Compose 2.0+
- At least 4GB RAM available for containers
- Ports 3000, 3001, 5432, 6379, 8080 available

### Start Full Application
```bash
# Clone the repository
git clone <repository-url>
cd hexagonal-ecommerce

# Start all services with one command
make dev-up

# Check status of all services  
make status

# View logs from all services
make dev-logs

# Stop all services
make dev-down
```

## 📊 Service Overview

| Service | Port | Description | URL |
|---------|------|-------------|-----|
| **Backend API** | 8080 | Go/Gin REST API | http://localhost:8080 |
| **Storefront** | 3000 | Next.js customer site | http://localhost:3000 |
| **Admin Panel** | 3001 | React SaaS admin | http://localhost:3001 |
| **PostgreSQL** | 5432 | Primary database | postgresql://localhost:5432 |
| **Redis** | 6379 | Cache & sessions | redis://localhost:6379 |

## 🛠 Available Commands

### Main Development Commands
```bash
make dev-up          # Start all services
make dev-down        # Stop all services
make dev-logs        # View logs from all services
make dev-clean       # Stop and remove all data
make dev-rebuild     # Rebuild and restart everything
make status          # Check service status
```

### Service-Specific Commands
```bash
# Backend
make backend-logs    # Backend logs only
make backend-shell   # Access backend container

# Storefront  
make storefront-logs # Storefront logs only

# Database
make db-migrate      # Run database migrations
make db-shell        # PostgreSQL shell access
make db-reset        # Reset database (WARNING: deletes data)
```

### Core Services Only
```bash
make services-up     # Start only DB, Redis, cache
make services-down   # Stop only core services
```

## 📂 Project Structure

```
hexagonal-ecommerce/
├── docker-compose.dev.yml    # Development services
├── docker-compose.prod.yml   # Production configuration
├── Makefile                  # Development commands
├── backend/
│   ├── Dockerfile.dev       # Backend dev container
│   ├── Dockerfile.prod      # Backend production
│   └── .air.toml           # Hot reload config
├── storefront/
│   ├── Dockerfile.dev       # Next.js dev container
│   └── Dockerfile.prod      # Next.js production
├── shop-keeper-dashboard/
│   ├── Dockerfile.dev       # React admin dev
│   └── Dockerfile.prod      # React admin prod
└── database/
    └── init/                # DB initialization scripts
```

## 🔧 Configuration

### Environment Variables

**Backend (.env)**
```bash
ENVIRONMENT=development
PORT=8080
DATABASE_URL=postgres://postgres:postgres123@postgres:5432/ecommerce_saas_dev?sslmode=disable
REDIS_URL=redis://redis:6379
JWT_SECRET=dev-jwt-secret-key
```

**Storefront (.env.local)**
```bash
NODE_ENV=development
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_APP_URL=http://localhost:3000
```

**Admin Panel (.env)**
```bash
NODE_ENV=development
VITE_API_URL=http://localhost:8080/api/v1
PORT=3001
```

### Database Configuration
- **Database**: `ecommerce_saas_dev`
- **User**: `postgres`
- **Password**: `postgres123`
- **Host**: `localhost:5432` (from host machine)
- **Host**: `postgres:5432` (from containers)

## 🎯 Development Workflow

### 1. First Time Setup
```bash
# Start the application
make dev-up

# Wait for services to start (30-60 seconds)
make status

# Run database migrations
make db-migrate

# Access the application
open http://localhost:3000  # Storefront
open http://localhost:3001  # Admin Panel
```

### 2. Daily Development
```bash
# Start development
make dev-up

# View logs while coding
make dev-logs

# Access specific service logs
make backend-logs
make storefront-logs

# Stop when done
make dev-down
```

### 3. Debugging
```bash
# Check service status
make status

# Access backend shell for debugging
make backend-shell

# Access database directly
make db-shell

# Reset database if needed
make db-reset
```

## 🔍 Troubleshooting

### Common Issues

**Port Already in Use**
```bash
# Check what's using the port
lsof -i :3000
lsof -i :8080

# Stop existing processes
kill <PID>

# Or use different ports in docker-compose.dev.yml
```

**Database Connection Issues**
```bash
# Check if PostgreSQL is running
make status

# Check database logs
docker-compose -f docker-compose.dev.yml logs postgres

# Reset database
make db-reset
```

**Hot Reload Not Working**
```bash
# Rebuild containers
make dev-rebuild

# Check volumes are properly mounted
docker-compose -f docker-compose.dev.yml ps
```

**Memory Issues**
```bash
# Check container resource usage
docker stats

# Increase Docker Desktop memory to 6GB+
# Or stop unused containers
make dev-clean
```

### Service Health Checks

**Check Backend API**
```bash
curl http://localhost:8080/health
# Should return: {"status": "ok"}
```

**Check Database Connection**
```bash
make db-shell
\l  # List databases
\q  # Quit
```

**Check Redis**
```bash
docker-compose -f docker-compose.dev.yml exec redis redis-cli ping
# Should return: PONG
```

## 🚀 Production Deployment

### Production Docker Compose
```bash
# Use production configuration
docker-compose -f docker-compose.prod.yml up -d

# Check production status
docker-compose -f docker-compose.prod.yml ps
```

### Production Environment Variables
Set these in production:
```bash
# Backend
DATABASE_URL=postgresql://user:pass@prod-db:5432/ecommerce_saas
REDIS_URL=redis://prod-redis:6379
JWT_SECRET=<secure-random-secret>

# Frontend
NEXT_PUBLIC_API_URL=https://api.yourdomain.com/api/v1
NEXT_PUBLIC_APP_URL=https://yourdomain.com
```

## 📈 Performance Optimization

### Development Performance
```bash
# Use specific services only
make services-up
# Then run backend/frontend manually

# Limit log output
make backend-logs --tail=100

# Clean unused containers
make dev-clean
```

### Resource Limits
Add to docker-compose.dev.yml:
```yaml
services:
  backend:
    deploy:
      resources:
        limits:
          memory: 512M
        reservations:
          memory: 256M
```

## 🔐 Security Notes

### Development Security
- Default passwords are for development only
- Change all credentials in production
- Use environment files for sensitive data
- Never commit .env files to version control

### Network Security
```yaml
networks:
  ecommerce-network:
    driver: bridge
    internal: true  # For production isolation
```

## 🎓 Learning Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Reference](https://docs.docker.com/compose/)
- [Go Hot Reload with Air](https://github.com/cosmtrek/air)
- [Next.js Docker Guide](https://nextjs.org/docs/deployment#docker-image)

---

**Need Help?**
- Check service logs: `make dev-logs`
- Reset everything: `make dev-clean && make dev-up`
- Contact: development team