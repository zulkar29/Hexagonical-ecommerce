# Seed Data Setup Guide

This guide explains how to set up comprehensive demo data for the e-commerce SaaS platform using Docker.

## Quick Start

The easiest way to get started with demo data is to use the provided Make commands:

```bash
# Start the entire development environment
make dev-up

# Wait for services to be ready (about 30 seconds), then seed the database
make db-seed
```

## Available Make Commands

### Database Setup
- `make db-migrate` - Run database migrations only
- `make db-seed` - Populate database with comprehensive demo data
- `make db-setup` - Run migrations and seed data (complete setup)
- `make db-reset` - Reset database and recreate with seed data
- `make db-status` - Check database status and tables

### Development Environment
- `make dev-up` - Start all services (PostgreSQL, Redis, Backend, Frontend)
- `make dev-down` - Stop all services
- `make services-up` - Start only database and cache services

## Seed Data Overview

The seed data creates a comprehensive multi-tenant e-commerce environment with:

### 🏢 **3 Demo Tenants**
1. **TechHub Electronics** (`techhub`)
   - Plan: Professional
   - Focus: Electronics and gadgets
   - Custom domain: techhub.example.com

2. **Fashion Forward** (`fashionforward`)
   - Plan: Premium  
   - Focus: Fashion and lifestyle
   - Custom domain: fashionforward.example.com

3. **Home & Garden** (`homegarden`)
   - Plan: Starter
   - Focus: Home decor and garden tools
   - Subdomain only: homegarden.platform.com

### 👥 **User Accounts** (password: `password123`)
- **Admins**: `admin@techhub.com`, `admin@fashionforward.com`, `admin@homegarden.com`
- **Merchants**: `merchant@techhub.com`
- **Customers**: `customer1@example.com`, `customer2@example.com`, `customer3@example.com`

### 🛍️ **Product Catalog**
- **10+ Categories** with hierarchical structure
- **8+ Products** with realistic pricing and inventory
- **Product Variants** (colors, sizes, storage options)
- **Complete product data** (SKUs, descriptions, images, SEO)

### 📦 **Order Data**
- **3 Sample orders** with different statuses
- **Order items** and **shipping addresses**
- **Order fulfillment** tracking

### 💳 **Payment Data**
- **Payment records** for orders
- **Payment methods** (cards, mobile wallets)
- **SSLCommerz** gateway integration

### 📊 **Analytics Data**
- **Page views** and **product views**
- **User behavior events** (add to cart, purchases)
- **UTM tracking** data

### 🔐 **Permission System**
- **Comprehensive permissions** for all resources
- **Role-based access control** (RBAC)
- **Role assignments** for different user types

## Manual Seed Commands

If you prefer to run commands manually:

### Using Docker Compose
```bash
# Run migrations
docker-compose -f docker-compose.dev.yml exec backend go run cmd/migrate/main.go

# Seed database
docker-compose -f docker-compose.dev.yml exec backend go run cmd/migrate/main.go -action=seed
```

### Direct Go Commands (if running locally)
```bash
cd backend

# Run migrations
go run cmd/migrate/main.go -action=up

# Seed database  
go run cmd/migrate/main.go -action=seed
```

## Database Schema

The seed data populates these key tables:

### Core Tables
- `tenants` - multi-tenant stores
- `users` - User accounts with roles
- `categories` - Product categories (hierarchical)
- `products` - Product catalog
- `product_variants` - Product variations
- `orders` - Customer orders
- `order_items` - Order line items
- `shipping_addresses` - Delivery addresses

### Advanced Tables
- `permissions` - System permissions
- `role_permissions` - Role-based access control
- `payments` - Payment transactions
- `payment_methods` - Saved payment methods
- `analytics_events` - User behavior tracking
- `page_views` - Website analytics
- `product_views` - Product analytics

## Verification

After seeding, you can verify the data:

```bash
# Check database tables
make db-status

# Connect to database shell
make db-shell

# View tenant data
SELECT name, subdomain, plan FROM tenants;

# View user accounts
SELECT email, role, first_name, last_name FROM users;

# View products by tenant
SELECT t.name as tenant, p.name as product, p.price 
FROM products p 
JOIN tenants t ON p.tenant_id = t.id;
```

## Environment Variables

Make sure these environment variables are set in your `.env` file:

```env
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ecommerce_saas_dev

# Application
JWT_SECRET=your-super-secret-jwt-key-change-in-production
ENVIRONMENT=development
```

## Troubleshooting

### Services not starting
```bash
# Check service status
docker-compose -f docker-compose.dev.yml ps

# View logs
make dev-logs
```

### Database connection issues
```bash
# Restart database service
docker-compose -f docker-compose.dev.yml restart postgres

# Check PostgreSQL logs
docker-compose -f docker-compose.dev.yml logs postgres
```

### Migration errors
```bash
# Reset and recreate database
make db-reset
```

### Seed data conflicts
The seed functions use `ON CONFLICT DO NOTHING` to prevent duplicate data, so it's safe to run multiple times.

## API Testing

Once seeded, you can test the API:

```bash
# Health check
curl http://localhost:8080/health

# Login as admin
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@techhub.com", "password": "password123"}'

# Get products (requires authentication)
curl http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <jwt_token>" \
  -H "X-Tenant-ID: 11111111-1111-1111-1111-111111111111"
```

## Next Steps

After setting up seed data:

1. **Access the dashboards** at the URLs shown after `make dev-up`
2. **Login with the demo accounts** to explore features
3. **Test the multi-tenant isolation** by switching between tenants
4. **Explore the API** using the provided demo data
5. **Build new features** on top of the comprehensive data structure

## Data Customization

To customize the seed data:

1. Edit `/backend/internal/shared/database/seed.go`
2. Modify the data structures in each `seed*` function  
3. Re-run `make db-seed` to apply changes

The seed data is designed to be comprehensive yet realistic, providing a solid foundation for development and testing of the e-commerce SaaS platform.