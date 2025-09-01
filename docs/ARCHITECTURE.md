# E-commerce SaaS Platform Architecture

## Overview
Technical architecture for a multi-tenant e-commerce SaaS platform using hexagonal (clean) architecture principles.

## System Architecture

### Optimized Modular Monolith Architecture
**Updated**: Simplified from microservices to modular monolith for better initial scalability

```
┌─────────────────────────────────────────────────────────────┐
│                    PRESENTATION LAYER                       │
├─────────────────────────────────────────────────────────────┤
│              Next.js (Customer + Dashboard)                 │
│         Unified frontend with route-based separation        │
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
│  │   Handler   │   (GORM)    │  (Stripe)   │ (SendGrid)  │  │
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
- **Payment Module**: Stripe integration, subscription billing
- **Notification Module**: Email notifications, webhooks
- **Analytics Module**: Basic reporting and metrics

**Benefits of Modular Monolith**:
- Single deployment unit (reduced complexity)
- Shared database transactions
- Lower network latency
- Easier debugging and testing
- Can be split into microservices later when needed

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

## Multi-Tenancy Strategy (Optimized)

### Hybrid Database Approach
**Small-Medium Tenants (Basic/Professional Plans)**:
- Shared PostgreSQL database with `tenant_id` column
- Row-level security for data isolation
- Cost-effective for up to 10,000 products per tenant
- 60-80% reduction in infrastructure costs

**Enterprise Tenants (Enterprise Plan)**:
- Dedicated database per tenant
- Complete isolation for large-scale operations
- Custom scaling and backup strategies
- Triggered when tenant exceeds 10,000 products

### Tenant Type Selection Logic
```go
func (t *Tenant) ShouldUseDedicatedDatabase() bool {
    return t.Plan == PlanEnterprise || t.ProductCount >= 10000
}
```

### Benefits of Hybrid Approach
- **Cost Efficiency**: Shared DB for 95% of tenants
- **Scalability**: Dedicated DB for high-volume tenants
- **Simplified Operations**: Reduced backup/monitoring overhead
- **Gradual Migration**: Seamless upgrade path

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
- OAuth2/JWT authentication
- Multi-factor authentication
- Role-based access control (RBAC)
- API rate limiting and throttling
- Input validation and sanitization
- HTTPS everywhere with HSTS

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