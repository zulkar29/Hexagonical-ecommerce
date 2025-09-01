# Docker Development Setup

Complete Docker development environment for the E-commerce SaaS platform.

## Quick Start

```bash
# Start all development services
make dev-up

# Stop all services
make dev-down

# View logs
make dev-logs

# Clean everything
make dev-clean
```

## Services Included

### **Core Application Services**
- **Backend API** (Go Fiber) - `localhost:8080`
- **Storefront** (Next.js) - `localhost:3000` 
- **Dashboard** (React.js) - `localhost:3001`

### **Infrastructure Services**
- **PostgreSQL** - `localhost:5432`
- **Redis** - `localhost:6379`
- **MinIO** (S3-compatible) - `localhost:9000`
- **MailHog** (Email testing) - `localhost:8025`

## Development Features

### **Hot Reloading**
- ✅ **Backend**: Air for Go hot reloading
- ✅ **Storefront**: Next.js fast refresh
- ✅ **Dashboard**: React hot reloading

### **Volume Mounting**
- Source code mounted for live editing
- Dependencies cached for faster rebuilds
- Database data persisted

### **Environment Variables**
- Development-optimized settings
- Service discovery via container names
- Local file system access

## Available Commands

```bash
# Development Environment
make dev-up          # Start all services
make dev-down        # Stop all services  
make dev-logs        # Show all logs
make dev-clean       # Clean up + remove volumes
make dev-rebuild     # Rebuild and restart

# Core Services Only
make services-up     # Start DB, Redis, MinIO only
make services-down   # Stop core services

# Individual Service Logs
make backend-logs    # Backend API logs
make storefront-logs # Storefront logs
make dashboard-logs  # Dashboard logs

# Database Operations
make db-migrate      # Run migrations
make db-shell        # PostgreSQL shell
make db-reset        # Reset database (WARNING: deletes data)

# Development Utilities
make install         # Install all dependencies
make format          # Format all code
make test            # Run all tests
make status          # Check service status
```

## Database Setup

```bash
# Start services
make dev-up

# Run migrations (after implementing)
make db-migrate

# Access database shell
make db-shell
```

## 📧 Email Testing

MailHog captures all outgoing emails in development:
- **Web UI**: http://localhost:8025
- **SMTP**: localhost:1025

## 📁 File Storage

MinIO provides S3-compatible storage:
- **Console**: http://localhost:9001
- **Credentials**: minioadmin/minioadmin123
- **Endpoint**: http://localhost:9000

## 🔧 Configuration

### **Environment Variables**
Copy and modify `.env.dev`:
```bash
cp .env.dev .env
```

### **Port Mapping**
- Backend: 8080
- Storefront: 3000
- Dashboard: 3001
- PostgreSQL: 5432
- Redis: 6379
- MinIO API: 9000
- MinIO Console: 9001
- MailHog: 8025

## 🐛 Troubleshooting

### **Port Conflicts**
```bash
# Check what's using ports
lsof -i :3000
lsof -i :8080

# Stop conflicting services
make dev-down
```

### **Volume Issues**
```bash
# Clean volumes and restart
make dev-clean
make dev-up
```

### **Permission Issues**
```bash
# Fix file permissions
sudo chown -R $USER:$USER .
```

### **Build Issues**
```bash
# Rebuild without cache
make dev-rebuild
```

## 📝 Development Workflow

1. **Start Environment**:
   ```bash
   make dev-up
   ```

2. **Check Status**:
   ```bash
   make status
   ```

3. **View Logs**:
   ```bash
   make dev-logs
   ```

4. **Develop**:
   - Edit code in your IDE
   - Changes auto-reload in containers
   - Access services via localhost

5. **Test**:
   ```bash
   make test
   ```

6. **Stop When Done**:
   ```bash
   make dev-down
   ```

## 🎯 Next Steps

1. Implement actual Go Fiber backend
2. Complete Next.js storefront
3. Finish React dashboard
4. Add database migrations
5. Integrate with external APIs

This Docker setup provides a complete development environment with hot reloading, database persistence, and all necessary services for building the e-commerce SaaS platform!