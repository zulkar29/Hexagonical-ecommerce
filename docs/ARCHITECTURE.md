# Hexagonal E-commerce SaaS Platform Architecture

## Overview
Technical architecture for a multi-tenant e-commerce SaaS platform implementing **Hexagonal Architecture** (Ports and Adapters pattern) with modular monolith design for optimal testability, maintainability, and business logic isolation.

> **Note**: For detailed API specifications and endpoint documentation, see [API_REFERENCE.md](./API_REFERENCE.md)

## Hexagonal Architecture Pattern

### Core Principles
1. **Domain-Centric Design**: Business logic is at the center, isolated from external concerns
2. **Dependency Inversion**: Dependencies point inward toward the domain
3. **Ports and Adapters**: Clean interfaces separate business logic from implementation details
4. **Testability**: Easy to test business logic in isolation
5. **Technology Independence**: Can swap databases, frameworks, or external services

### Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│                    EXTERNAL ACTORS                          │
│  ┌─────────────────────────┬─────────────────────────────┐  │
│  │    Next.js Storefront   │   React.js Dashboard        │  │
│  │   (Customer-facing)     │   (Merchant admin)          │  │
│  │                         │   SaaS Admin Panel          │  │
│  └─────────────────────────┴─────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │ HTTP Adapters
┌─────────────────────────────────────────────────────────────┐
│                 ADAPTERS (INFRASTRUCTURE)                   │
│  ┌─────────────┬─────────────┬─────────────┬─────────────┐  │
│  │   HTTP      │  Database   │   Payment   │    Email    │  │
│  │  Handlers   │ Repository  │(SSLCommerz) │ (SendGrid)  │  │
│  │(gin/fiber)  │   (GORM)    │   Adapter   │   Adapter   │  │
│  └─────────────┴─────────────┴─────────────┴─────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │ Ports (Interfaces)
┌─────────────────────────────────────────────────────────────┐
│                    APPLICATION LAYER                        │
│  ┌─────────────┬─────────────┬─────────────┬─────────────┐  │
│  │   Product   │    User     │   Payment   │    Cart     │  │
│  │   Service   │   Service   │   Service   │   Service   │  │
│  │(Use Cases)  │(Use Cases)  │(Use Cases)  │(Use Cases)  │  │
│  └─────────────┴─────────────┴─────────────┴─────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                      DOMAIN LAYER                           │
│  ┌─────────────┬─────────────┬─────────────┬─────────────┐  │
│  │   Product   │    User     │   Payment   │    Cart     │  │
│  │  Entities   │  Entities   │  Entities   │  Entities   │  │
│  │+ Business   │+ Business   │+ Business   │+ Business   │  │
│  │   Rules     │   Rules     │   Rules     │   Rules     │  │
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

## Module Structure (Hexagonal Pattern)

Each business module follows the same hexagonal pattern:

```
internal/
├── product/                    # Product Domain Module
│   ├── product.go             # Domain Entities + Business Rules
│   ├── service.go             # Application Services (Use Cases)
│   ├── repository.go          # Port Interface + Adapter Implementation
│   ├── handler.go             # HTTP Adapter (Primary Port)
│   ├── module.go              # Dependency Injection Container
│   └── product_test.go        # Comprehensive Testing
├── user/                      # User Domain Module
│   ├── user.go               # Domain Entities + Business Rules
│   ├── service.go            # Application Services (Use Cases)
│   ├── repository.go         # Port Interface + Adapter Implementation
│   ├── handler.go            # HTTP Adapter (Primary Port)
│   ├── security.go           # Security Domain Logic
│   ├── security_service.go   # Security Use Cases
│   └── service_test.go       # Comprehensive Testing
├── payment/                   # Payment Domain Module
├── cart/                     # Cart Domain Module
├── security/                 # Security Domain Module
└── shared/                   # Shared Utilities
    └── testhelpers/          # Testing Infrastructure
```

### File Responsibilities in Hexagonal Architecture

#### 1. **Domain Layer** (`*.go` entity files)
```go
// product.go - Pure business entities and rules
type Product struct {
    ID          uuid.UUID
    TenantID    uuid.UUID
    Name        string
    Price       float64
    // ... other fields
}

// Business rules embedded in domain
func (p *Product) CalculateDiscount() float64 {
    return p.ComparePrice - p.Price
}

func (p *Product) IsInStock() bool {
    return p.InventoryQuantity > 0 || p.AllowBackorder
}
```

#### 2. **Application Layer** (`service.go`)
```go
// service.go - Use cases and orchestration
type Service struct {
    repo      Repository  // Port (interface)
    validator *validator.Validate
}

func (s *Service) CreateProduct(tenantID uuid.UUID, product *Product) (*Product, error) {
    // Use case orchestration
    if err := s.validator.Struct(product); err != nil {
        return nil, err
    }
    
    product.TenantID = tenantID
    return s.repo.CreateProduct(product)
}
```

#### 3. **Ports** (Interfaces in `repository.go`)
```go
// repository.go - Port definitions
type Repository interface {
    CreateProduct(product *Product) (*Product, error)
    GetProductByID(tenantID, productID uuid.UUID) (*Product, error)
    UpdateProduct(product *Product) (*Product, error)
    DeleteProduct(tenantID, productID uuid.UUID) error
    // ... other methods
### Testing Strategy for Hexagonal Architecture

#### Integration Testing Pattern
```go
// product_test.go - Comprehensive testing approach
func TestProductService_Integration(t *testing.T) {
    // Real database container for integration testing
    container := startTestDB(t)
    defer container.Terminate(context.Background())
    
    db := connectTestDB(t, container)
    module := product.NewModule(db)
    
    // Test business logic with real dependencies
    t.Run("CreateProduct", func(t *testing.T) {
        // Test with actual database, validation, business rules
        product, err := module.Service.CreateProduct(tenantID, &Product{
            Name:  "Test Product",
            Price: 99.99,
        })
        
        assert.NoError(t, err)
        assert.NotEmpty(t, product.ID)
        assert.Equal(t, tenantID, product.TenantID)
    })
}
```

#### Why Integration Testing for Hexagonal Architecture
1. **Port Validation**: Tests ensure ports (interfaces) work correctly with adapters
2. **Business Logic**: Domain rules are tested with real dependencies
3. **Database Constraints**: Foreign key relationships and constraints validated
4. **Multi-tenancy**: Tenant isolation tested with real data
5. **Performance**: Real database performance metrics captured

## Multi-Tenancy with Shared Database

### Tenant Isolation Strategy
```sql
-- Every table includes tenant_id for data isolation
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_products_tenant 
        FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

-- Row Level Security (RLS) enforces tenant isolation
ALTER TABLE products ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_products ON products
    USING (tenant_id = current_setting('app.current_tenant')::UUID);
```

### Database Schema Organization

```sql
-- Core tenant management
tenants                 -- SaaS tenant registration
users                   -- User authentication + tenant association

-- E-commerce domain
products               -- Product catalog with tenant_id
categories             -- Product categories with tenant_id
orders                 -- Order management with tenant_id
order_items           -- Order line items
payments              -- Payment processing with tenant_id

-- Supporting modules
notifications         -- Email/SMS notifications
analytics            -- Tenant-specific analytics
billing              -- SaaS subscription billing
```

## Module Communication Patterns

### 1. **Direct Service Injection** (Preferred)
```go
// Order service needs product information
type OrderService struct {
    orderRepo    OrderRepository
    productService *product.Service  // Direct injection
}

func (s *OrderService) CreateOrder(tenantID uuid.UUID, req *CreateOrderRequest) error {
    // Validate products exist in same tenant
    for _, item := range req.Items {
        product, err := s.productService.GetProduct(tenantID, item.ProductID)
        if err != nil {
            return err
        }
        // Business logic continues...
    }
}
```

### 2. **Event-Driven Communication** (For decoupling)
```go
// Events for async operations
type OrderCreatedEvent struct {
    TenantID uuid.UUID
    OrderID  uuid.UUID
    Items    []OrderItem
}

// Publisher/Subscriber pattern
func (s *OrderService) CreateOrder(order *Order) error {
    // Create order
    if err := s.repo.CreateOrder(order); err != nil {
        return err
    }
    
    // Publish event for inventory updates
    s.eventBus.Publish(OrderCreatedEvent{
        TenantID: order.TenantID,
        OrderID:  order.ID,
        Items:    order.Items,
    })
    
    return nil
}
```

## Technology Stack Alignment

### Backend (Go with Hexagonal Architecture)
- **Framework**: Gin/Fiber for HTTP adapters
- **ORM**: GORM for database adapters  
- **Validation**: go-playground/validator for domain validation
- **Testing**: testcontainers for integration testing
- **Migration**: golang-migrate for database versioning

### Frontend Applications
- **Customer Storefront**: Next.js with App Router
- **Merchant Dashboard**: React.js SPA  
- **SaaS Admin Panel**: React.js SPA
- **Styling**: Tailwind CSS across all frontends

### Infrastructure
- **Database**: PostgreSQL with RLS for tenant isolation
- **Cache**: Redis for session and application caching
- **Storage**: AWS S3 compatible for file uploads
- **Deployment**: Docker containers with docker-compose

## Advantages of This Architecture

### 1. **Testability**
- Business logic isolated from infrastructure
- Easy to mock external dependencies
- Integration tests with real databases via testcontainers
- Clear separation of concerns

### 2. **Maintainability**  
- Each module is self-contained
- Changes in one module don't affect others
- Clear interfaces between layers
- Easy to refactor individual components

### 3. **Technology Independence**
- Can replace GORM with raw SQL without changing business logic
- Can switch from Gin to Fiber without changing services
- Can change payment providers without touching domain code
- Database migrations don't affect business rules

### 4. **Solo Developer Efficiency**
- Single codebase for all backend logic
- Shared utilities and testing infrastructure
- No microservice complexity overhead
- Local development simplicity

### 5. **Multi-tenant Scalability**
- Row-level security enforces tenant isolation
- Single database instance supports thousands of tenants
- Efficient resource utilization
- Cost-effective for SaaS business model

## Development Workflow

### 1. **Adding New Features**
```bash
# 1. Define domain entities
touch internal/newmodule/newmodule.go

# 2. Define application services (use cases)  
touch internal/newmodule/service.go

# 3. Define ports (interfaces)
touch internal/newmodule/repository.go

# 4. Implement adapters
# - Database adapter in repository.go
# - HTTP adapter in handler.go

# 5. Wire dependencies
touch internal/newmodule/module.go

# 6. Write comprehensive tests
touch internal/newmodule/newmodule_test.go

# 7. Add migration for database changes
touch migrations/001_create_newmodule.sql
```

### 2. **Testing New Modules**
```bash
# Run integration tests with real database
cd backend
go test ./internal/newmodule -v

# Run all tests
make test

# Run with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Error Handling Strategy

### Domain Errors
```go
// Define domain-specific errors
type ProductError struct {
    Code    string
    Message string
    Field   string
}

func (e ProductError) Error() string {
    return e.Message
}

// Use in domain layer
func (p *Product) Validate() error {
    if p.Price < 0 {
        return ProductError{
            Code:    "INVALID_PRICE",
            Message: "Product price cannot be negative",
            Field:   "price",
        }
    }
    return nil
}
```

### HTTP Error Translation
```go
// HTTP adapter translates domain errors to HTTP responses
func (h *Handler) CreateProduct(c *gin.Context) {
    product, err := h.service.CreateProduct(tenantID, req)
    if err != nil {
        switch e := err.(type) {
        case ProductError:
            c.JSON(400, gin.H{
                "error":   e.Code,
                "message": e.Message,
                "field":   e.Field,
            })
        default:
            c.JSON(500, gin.H{"error": "INTERNAL_ERROR"})
        }
        return
    }
    
    c.JSON(201, product)
}
```

## Security Implementation

### Tenant Context Injection
```go
// Middleware extracts tenant from JWT and injects into context
func TenantMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractJWT(c)
        claims := validateJWT(token)
        
        // Inject tenant ID into context
        ctx := context.WithValue(c.Request.Context(), "tenant_id", claims.TenantID)
        c.Request = c.Request.WithContext(ctx)
        
        c.Next()
    }
}

// Repository uses tenant context for queries
func (r *repository) GetProducts(ctx context.Context) ([]Product, error) {
    tenantID := ctx.Value("tenant_id").(uuid.UUID)
    
    var products []Product
    err := r.db.Where("tenant_id = ?", tenantID).Find(&products).Error
    return products, err
}
```

## Monitoring and Observability

### Metrics Collection
```go
// Service layer instruments business metrics
func (s *Service) CreateProduct(product *Product) (*Product, error) {
    start := time.Now()
    defer func() {
        metrics.ProductCreationDuration.Observe(time.Since(start).Seconds())
        metrics.ProductCreationTotal.Inc()
    }()
    
    result, err := s.repo.CreateProduct(product)
    if err != nil {
        metrics.ProductCreationErrors.Inc()
    }
    
    return result, err
}
```

### Health Checks
```go
// Health check adapter tests all external dependencies
func (h *HealthHandler) CheckHealth(c *gin.Context) {
    status := gin.H{"status": "healthy"}
    
    // Test database connectivity
    if err := h.db.Ping(); err != nil {
        status["database"] = "unhealthy"
        status["status"] = "unhealthy"
    }
    
    // Test external services
    if err := h.paymentService.Ping(); err != nil {
        status["payment"] = "unhealthy"
        status["status"] = "unhealthy"
    }
    
    c.JSON(200, status)
}
```

## Core Components & Business Modules

### Backend Modules (Golang Hexagonal Architecture)
- **Tenant Module**: Multi-tenant management, subscription plans, settings
- **Product Module**: Catalog management, inventory, variants, categories
- **Order Module**: Cart, checkout, order processing, fulfillment
- **User Module**: Authentication, authorization, customer management
- **Payment Module**: SSLCommerz integration for local payments, subscription billing
- **Notification Module**: Email/SMS notifications, system alerts, webhooks
- **Analytics Module**: Sales analytics, customer insights, platform monitoring
- **Security Module**: Multi-tenant data isolation, audit trails, fraud detection
- **Shipping Module**: Shipping zones, rates, tracking, courier integration
- **Support Module**: Contact forms, ticketing system, FAQ management

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

This Hexagonal Architecture implementation provides a robust foundation for the e-commerce SaaS platform, ensuring clean separation of concerns, excellent testability, and maintainable code that can evolve with business requirements.

## Reference Documentation

For implementation-specific details, see:
- **[INFRASTRUCTURE.md](./INFRASTRUCTURE.md)** - Infrastructure setup, database configuration, caching strategies
- **[DEPLOYMENT.md](./DEPLOYMENT.md)** - Docker deployment, VPS setup, domain management
- **[MONITORING.md](./MONITORING.md)** - Performance monitoring, health checks, observability
- **[API_REFERENCE.md](./API_REFERENCE.md)** - Complete API documentation and endpoints

## Reference Documentation

For implementation-specific details, see:
- **[INFRASTRUCTURE.md](./INFRASTRUCTURE.md)** - Infrastructure setup, database configuration, caching strategies
- **[DEPLOYMENT.md](./DEPLOYMENT.md)** - Docker deployment, VPS setup, domain management
- **[MONITORING.md](./MONITORING.md)** - Performance monitoring, health checks, observability
- **[API_REFERENCE.md](./API_REFERENCE.md)** - Complete API documentation and endpoints
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

## Technical Implementation Details

### Database Connection Pooling Configuration
```yaml
PostgreSQL Connection Pool (PgBouncer):
  max_client_conn: 200
  default_pool_size: 25
  min_pool_size: 5
  reserve_pool_size: 5
  reserve_pool_timeout: 5
  max_db_connections: 100
  pool_mode: transaction

Connection Pool Per Plan:
  Free: 2 max connections
  Starter: 10 max connections
  Professional: 25 max connections
  Pro: 50 max connections
  Enterprise: 100 max connections
```

### Redis Cache Configuration
```yaml
Redis Cache Policies:
  Eviction Policy: allkeys-lru
  Max Memory: 2GB
  Max Memory Policy: allkeys-lru

Cache TTL Settings:
  Product catalog: 1 hour
  User sessions: 24 hours
  API rate limits: 1 hour
  Search results: 30 minutes

Tenant Isolation:
  Key pattern: "tenant:{tenant_id}:{cache_type}:{key}"
  Example: "tenant:shop123:products:category_electronics"
```

### File Upload Limits & Processing
```yaml
File Upload Configuration:
  Max file size by plan:
    Free: 2MB
    Starter: 10MB
    Professional: 25MB
    Pro: 50MB
    Enterprise: 100MB

  Supported formats:
    Images: JPEG, PNG, WebP, SVG
    Documents: PDF, DOC, DOCX
    Archives: ZIP, RAR

  Processing pipeline:
    1. Virus scan (ClamAV)
    2. File type validation
    3. Size compression
    4. CDN upload
    5. Tenant storage quota check
```

### WebSocket Connection Management
```yaml
WebSocket Limits:
  Free: 10 concurrent connections
  Starter: 100 concurrent connections
  Professional: 500 concurrent connections
  Pro: 2,000 concurrent connections
  Enterprise: 10,000 concurrent connections

Connection Features:
  - Real-time order updates
  - Inventory synchronization
  - Customer chat support
  - Admin notifications
  - Live analytics updates

Timeout Settings:
  Idle timeout: 30 minutes
  Heartbeat interval: 30 seconds
  Reconnection attempts: 3
  Backoff strategy: Exponential
```