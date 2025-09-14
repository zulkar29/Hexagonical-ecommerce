# Hexagonal E-commerce SaaS Platform Architecture

📋 **Documentation Navigation**: [📖 Project Home](../README.md) | [🚀 Features](./FEATURES.md) | [📋 Roadmap](./ROADMAP.md) | [🚢 Deployment](./DEPLOYMENT.md)

## Overview
Technical architecture for a single vendor multi-tenant e-commerce SaaS platform using Hexagonal Architecture principles with modular monolith design for optimal performance and maintainability.

**Important Clarification**:
- **"Hexagonal"** refers to the **Hexagonal Architecture pattern** (also known as Ports and Adapters pattern)
- **Multi-tenancy** is achieved through **shared database with tenant_id isolation**
- These are separate architectural concerns that work together

> **Note**: For detailed API specifications and endpoint documentation, see [API_REFERENCE.md](./API_REFERENCE.md)

## System Architecture

### Hexagonal Architecture with Modular Monolith Implementation
**Design**: Hexagonal Architecture implemented as modular monolith for solo developer efficiency, testability, and local market requirements

```
┌─────────────────────────────────────────────────────────────┐
│                    PRESENTATION LAYER                       │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────┬─────────────────────────────┐  │
│  │    Next.js Storefront   │   React.js Dashboard        │  │
│  │   (Customer-facing)     │   (Merchant admin)          │  │
│  └─────────────────────────┴─────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                GOLANG MODULAR MONOLITH                     │
│              (Single service, multiple modules)            │
├─────────────────────────────────────────────────────────────┤
│                      DOMAIN MODULES                        │
│  ┌─────────────┬─────────────┬─────────────┬─────────────┐  │
│  │   Tenant    │   Product   │    Order    │    User     │  │
│  │   Module    │   Module    │   Module    │   Module    │  │
│  └─────────────┴─────────────┴─────────────┴─────────────┘  │
├─────────────────────────────────────────────────────────────┤
│                    APPLICATION LAYER                        │
│  ┌─────────────┬─────────────┬─────────────┬─────────────┐  │
│  │   Tenant    │   Product   │    Order    │    User     │  │
│  │  Service    │   Service   │   Service   │   Service   │  │
│  └─────────────┴─────────────┴─────────────┴─────────────┘  │
├─────────────────────────────────────────────────────────────┤
│                      PORTS & ADAPTERS                      │
│  ┌─────────────┬─────────────┬─────────────┬─────────────┐  │
│  │    HTTP     │  Repository │   Payment   │    Email    │  │
│  │   Handler   │   (GORM)    │(SSLCommerz) │ (SendGrid)  │  │
│  └─────────────┴─────────────┴─────────────┴─────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                  INFRASTRUCTURE                             │
│  ┌─────────────┬─────────────┬─────────────┬─────────────┐  │
│  │  Database   │    Cache    │   Storage   │   Message   │  │
│  │ PostgreSQL  │    Redis    │     S3      │    Queue    │  │
│  └─────────────┴─────────────┴─────────────┴─────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Core Components

### Backend Modules (Golang Monolith)
- **Tenant Module**: Multi-tenant management, subscription plans, settings
- **Product Module**: Catalog management, inventory, variants, categories
- **Order Module**: Cart, checkout, order processing, fulfillment
- **User Module**: Authentication, authorization, customer management
- **Payment Module**: SSLCommerz integration for local payments, subscription billing
- **Notification Module**: Email/SMS notifications, system alerts, webhooks, event-driven notifications
- **Analytics Module**: Sales analytics, customer insights, platform monitoring
- **Security Module**: Multi-tenant data isolation, audit trails, fraud detection
- **Shipping Module**: Shipping zones, rates, tracking, courier integration
- **Support Module**: Contact forms, ticketing system, FAQ management

**Benefits of Modular Monolith for Solo Developer**:
- Single deployment unit (solo-friendly complexity)
- Shared database transactions (easier data management)
- Lower network latency (local infrastructure optimized)
- Easier debugging and testing (single developer can manage)
- Local language features integrated across all modules
- Can be split into microservices later when team expands

### Frontend Applications
- **Customer Storefront (Next.js)**: 
  - Server-side rendered storefronts
  - Dynamic routing based on tenant domains
  - SEO-optimized product pages
  - Mobile-responsive design

- **Merchant Dashboard (React.js)**:
  - Store management interface
  - Analytics and reporting
  - Product and inventory management
  - Order processing
  - Theme customization

## Multi-Tenancy Strategy (Shared Database with tenant_id)

### Initial Implementation: Shared Database Approach
**All Plans (Starter/Professional/Pro/Enterprise)**:
- **Shared PostgreSQL database** with `tenant_id` column for data isolation
- **Row-level security policies** to enforce tenant boundaries
- **Automatic tenant_id injection** in all database queries
- **Cost-effective** approach reducing infrastructure complexity by 80%
- **Scalable** for up to 10,000+ products per tenant initially

### Future Scaling Options:
**When Individual Tenants Exceed Capacity**:
- Option to migrate specific high-volume tenants to dedicated databases
- Hybrid approach available for enterprise customers requiring complete isolation
- Maintains cost efficiency for majority of tenants while allowing custom scaling

### Tenant Context Resolution (Hexagonal Pattern)
```go
// Domain layer - tenant context is injected at the boundary
func (s *ProductService) GetProducts(ctx context.Context, tenantID string) ([]*domain.Product, error) {
    // Business logic with tenant context
    return s.productRepo.FindByTenant(ctx, tenantID)
}

// Repository adapter - automatic tenant_id filtering
func (r *ProductRepository) FindByTenant(ctx context.Context, tenantID string) ([]*domain.Product, error) {
    var products []*Product
    err := r.db.WithContext(ctx).
        Where("tenant_id = ?", tenantID).
        Where("deleted_at IS NULL").
        Find(&products).Error
    return mapToModels(products), err
}
```

### Benefits of Shared Database Approach
- **Cost Efficiency**: Single database infrastructure for all tenants
- **Simplified Operations**: One backup, monitoring, and maintenance process
- **Performance**: No cross-database queries needed
- **Tenant Isolation**: Row-level security ensures complete data separation

## Domain Management
- Custom domain support via DNS CNAME
- Subdomain provisioning (tenant.platform.com)
- SSL certificate automation (Let's Encrypt)
- CDN integration for global performance

## API Design & Versioning
- RESTful API with versioning via URL paths (e.g., /api/v1/)
- OpenAPI/Swagger documentation for all endpoints
- Rate limiting and throttling per tenant

## Security Architecture
- OAuth2/JWT authentication with multi-tenant context
- Multi-factor authentication for admin accounts
- Role-based access control (RBAC) with tenant isolation
- Multi-tenant data isolation with row-level security
- Tenant-aware database queries with automatic tenant_id injection
- Cross-tenant boundary validation and access prevention
- API rate limiting and throttling per tenant
- Input validation and sanitization
- HTTPS everywhere with HSTS
- Comprehensive audit trails for compliance
- Real-time fraud detection and monitoring
- Security vulnerability scanning and prevention

## Scalability Considerations (Optimized)
- **Modular monolith** with clear module boundaries for future microservices split
- **Horizontal scaling** with load balancers (when >100k requests/month)
- **Hybrid database strategy** reduces operational complexity by 70%
- **Redis caching** for product catalogs and session management
- **Database connection pooling** with PgBouncer
- **CDN integration** for static assets and images
- **Read replicas** only for enterprise tenants
- **Microservices migration path** when individual modules exceed 100k requests/month

### Performance Optimizations
- **Tenant-aware caching** with Redis
- **Database indexing** on tenant_id for all shared tables
- **Connection pooling** to prevent pool exhaustion
- **Lazy loading** for large tenant datasets
- **Background jobs** for heavy operations (reports, exports)