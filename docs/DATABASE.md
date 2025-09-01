# Database Design Strategy

Comprehensive database schema and design patterns for the multi-tenant e-commerce SaaS platform with optimized hybrid approach combining cost efficiency and scalability.

## Database Strategy (Optimized)

### Hybrid Multi-tenant Approach
**Optimization**: Combines shared database efficiency with dedicated database security

**Shared Database (Default)**:
- Row-level tenant isolation using `tenant_id` column
- Cost-effective for Free/Starter/Professional plans (<5k products)  
- Single backup/maintenance strategy
- 60-80% cost reduction vs dedicated databases

**Dedicated Database (Enterprise)**:
- Complete data separation for high-volume tenants
- Triggered when tenant exceeds 5,000 products or chooses Enterprise plan
- Independent scaling and backup strategies
- Full compliance isolation for enterprise customers

### Database Technology
- **Primary**: PostgreSQL 15+
- **Features Used**: JSONB, UUID, Full-text search, Partitioning, Row-level Security
- **Connection Pooling**: PgBouncer with tenant-aware routing
- **Monitoring**: pganalyze or similar

## Platform Database Schema

Contains tenant metadata and platform-level data:

```sql
-- Platform tenants registry (Hybrid approach design)
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    subdomain VARCHAR(50) UNIQUE NOT NULL,
    custom_domain VARCHAR(255),
    plan VARCHAR(20) NOT NULL DEFAULT 'free', -- free, starter, professional, pro, enterprise
    status VARCHAR(20) DEFAULT 'active',
    tenant_type VARCHAR(20) DEFAULT 'shared', -- shared, dedicated
    database_name VARCHAR(100), -- Only for dedicated tenants
    product_count INTEGER DEFAULT 0,
    monthly_requests INTEGER DEFAULT 0,
    settings JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Subscription plans with detailed configurations
CREATE TABLE subscription_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL, -- free, starter, professional, pro, enterprise
    price DECIMAL(10,2) NOT NULL,
    features JSONB NOT NULL,
    limits JSONB NOT NULL, -- {"max_products": 1000, "max_requests_per_month": 10000}
    database_type VARCHAR(20) DEFAULT 'shared', -- shared, dedicated
    created_at TIMESTAMP DEFAULT NOW()
);

-- Insert default plans with standardized pricing structure
INSERT INTO subscription_plans (name, price, features, limits, database_type) VALUES
('free', 0.00, 
 '[
    "basic_storefront", 
    "community_support", 
    "standard_templates",
    "ssl_certificate"
  ]', 
 '{
    "max_products": 10,
    "max_requests_per_month": 1000,
    "max_storage_gb": 0.5,
    "max_bandwidth_gb": 5,
    "support_response_hours": 72,
    "branding_removal": false
  }', 
 'shared'),

('starter', 1990.00, 
 '[
    "basic_analytics", 
    "email_support", 
    "standard_templates",
    "ssl_certificate",
    "local_payment_gateways"
  ]', 
 '{
    "max_products": 500,
    "max_requests_per_month": 10000,
    "max_storage_gb": 5,
    "max_bandwidth_gb": 50,
    "support_response_hours": 48,
    "branding_removal": true
  }', 
 'shared'),
 
('professional', 4990.00, 
 '[
    "basic_analytics", 
    "advanced_analytics", 
    "priority_support", 
    "custom_domain", 
    "advanced_templates",
    "email_marketing",
    "inventory_management",
    "basic_api_access",
    "social_integrations"
  ]', 
 '{
    "max_products": 5000,
    "max_requests_per_month": 100000,
    "max_storage_gb": 100,
    "max_bandwidth_gb": 500,
    "support_response_hours": 24,
    "branding_removal": true
  }', 
 'shared'),

('pro', 7990.00, 
 '[
    "all_professional_features", 
    "advanced_analytics", 
    "priority_support", 
    "custom_domain", 
    "premium_templates",
    "advanced_email_marketing",
    "inventory_management",
    "api_access",
    "social_integrations",
    "abandoned_cart_recovery"
  ]', 
 '{
    "max_products": 10000,
    "max_requests_per_month": 250000,
    "max_storage_gb": 250,
    "max_bandwidth_gb": 1000,
    "support_response_hours": 12,
    "branding_removal": true
  }', 
 'shared'),
 
('enterprise', 12990.00, 
 '[
    "all_features", 
    "dedicated_database", 
    "full_api_access", 
    "sla_support",
    "white_labeling",
    "custom_integrations",
    "custom_development",
    "priority_onboarding",
    "dedicated_support"
  ]', 
 '{
    "max_products": -1,
    "max_requests_per_month": -1,
    "max_storage_gb": -1,
    "max_bandwidth_gb": -1,
    "support_response_hours": 4,
    "sla_uptime_percent": 99.9,
    "branding_removal": true
  }', 
 'dedicated');

-- Add monitoring and usage tracking tables
CREATE TABLE tenant_usage_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    metric_date DATE NOT NULL,
    requests_count INTEGER DEFAULT 0,
    bandwidth_used_gb DECIMAL(10,3) DEFAULT 0,
    storage_used_gb DECIMAL(10,3) DEFAULT 0,
    products_count INTEGER DEFAULT 0,
    orders_count INTEGER DEFAULT 0,
    active_users INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(tenant_id, metric_date)
);

CREATE INDEX idx_usage_metrics_tenant_date ON tenant_usage_metrics(tenant_id, metric_date DESC);

-- Migration tracking table
CREATE TABLE tenant_migrations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    migration_type VARCHAR(50) NOT NULL, -- 'shared_to_dedicated', 'plan_upgrade', etc.
    from_state JSONB NOT NULL,
    to_state JSONB NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'in_progress', 'completed', 'failed'
    started_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,
    error_message TEXT,
    rollback_data JSONB
);

CREATE INDEX idx_migrations_tenant_id ON tenant_migrations(tenant_id, started_at DESC);
```

## Shared Database Schema

For tenants using shared database (Free/Starter/Professional plans) - all tables include tenant_id:

```sql
-- Products table with tenant isolation
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    sku VARCHAR(100),
    price DECIMAL(10,2) NOT NULL,
    compare_price DECIMAL(10,2),
    cost_price DECIMAL(10,2),
    status VARCHAR(20) DEFAULT 'draft',
    inventory JSONB,
    seo JSONB,
    attributes JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    -- Tenant-scoped unique constraints
    UNIQUE(tenant_id, sku)
);

-- Optimized indexes for multi-tenant queries
CREATE INDEX idx_products_tenant_id ON products(tenant_id);
CREATE INDEX idx_products_tenant_status ON products(tenant_id, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_tenant_created ON products(tenant_id, created_at DESC);

-- Row Level Security for automatic tenant isolation
ALTER TABLE products ENABLE ROW LEVEL SECURITY;
CREATE POLICY products_tenant_isolation ON products
    FOR ALL
    TO application_role
    USING (tenant_id = current_setting('app.current_tenant_id')::UUID)
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id')::UUID);

-- Product categories with hierarchy support
CREATE TABLE product_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    description TEXT,
    parent_id UUID REFERENCES product_categories(id),
    image_url VARCHAR(500),
    position INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(tenant_id, slug)
);

-- Orders with comprehensive tracking
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    order_number VARCHAR(50) NOT NULL,
    customer_email VARCHAR(255) NOT NULL,
    customer_id UUID,
    status VARCHAR(20) DEFAULT 'pending',
    payment_status VARCHAR(20) DEFAULT 'pending',
    fulfillment_status VARCHAR(20) DEFAULT 'unfulfilled',
    total_amount DECIMAL(12,2) NOT NULL,
    subtotal_amount DECIMAL(12,2) NOT NULL,
    tax_amount DECIMAL(12,2) DEFAULT 0,
    shipping_amount DECIMAL(12,2) DEFAULT 0,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'BDT',
    items JSONB NOT NULL,
    shipping_address JSONB,
    billing_address JSONB,
    payment_details JSONB,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(tenant_id, order_number)
);

-- Optimized indexes for order management
CREATE INDEX idx_orders_tenant_id ON orders(tenant_id);
CREATE INDEX idx_orders_tenant_status ON orders(tenant_id, status, created_at DESC);
CREATE INDEX idx_orders_tenant_customer ON orders(tenant_id, customer_email);

-- Customer management
CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    email VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    date_of_birth DATE,
    default_address JSONB,
    addresses JSONB,
    marketing_consent BOOLEAN DEFAULT false,
    tags JSONB,
    notes TEXT,
    total_orders INTEGER DEFAULT 0,
    total_spent DECIMAL(12,2) DEFAULT 0,
    last_order_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(tenant_id, email)
);

-- Performance optimization: Partitioning for large tables
CREATE TABLE orders_partitioned (
    LIKE orders INCLUDING ALL
) PARTITION BY HASH (tenant_id);

-- Create partitions (expand as needed)
CREATE TABLE orders_part_0 PARTITION OF orders_partitioned
    FOR VALUES WITH (modulus 4, remainder 0);
CREATE TABLE orders_part_1 PARTITION OF orders_partitioned  
    FOR VALUES WITH (modulus 4, remainder 1);
CREATE TABLE orders_part_2 PARTITION OF orders_partitioned
    FOR VALUES WITH (modulus 4, remainder 2);
CREATE TABLE orders_part_3 PARTITION OF orders_partitioned
    FOR VALUES WITH (modulus 4, remainder 3);
```

## Dedicated Database Schema

For enterprise tenants with dedicated databases - same structure but WITHOUT tenant_id columns:

```sql
-- Products table (no tenant_id needed)
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    sku VARCHAR(100) UNIQUE,
    price DECIMAL(10,2) NOT NULL,
    compare_price DECIMAL(10,2),
    cost_price DECIMAL(10,2),
    status VARCHAR(20) DEFAULT 'draft',
    inventory JSONB,
    seo JSONB,
    attributes JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Simpler indexes without tenant_id
CREATE INDEX idx_products_status ON products(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_created ON products(created_at DESC);

-- No Row Level Security needed for dedicated databases
```

## Cost-Benefit Analysis

### Infrastructure Cost Comparison

**Traditional Database-Per-Tenant Approach**:
- 100 tenants = 100 database instances
- Average cost per database: ৳5,500/month
- Total monthly cost: ৳5,50,000
- Backup storage: 100 separate backup strategies
- Monitoring overhead: 100 databases to monitor

**Optimized Hybrid Approach**:
- 95 tenants on shared database: 1 instance (৳3,000/month)
- 5 enterprise tenants on dedicated: 5 instances (৳3,500/month each)  
- Total monthly cost: ৳20,500
- **Cost Reduction: 96%** for same tenant count
- Backup storage: 6 backup strategies total
- Monitoring overhead: 6 databases to monitor

### Operational Overhead Reduction

**Backup & Recovery**:
- **Before**: 100 separate backup jobs, 100 recovery procedures
- **After**: 6 backup jobs, streamlined recovery procedures
- **Time Savings**: 85% reduction in backup management

**Database Migrations**:
- **Before**: Run migration on each database individually
- **After**: Single migration for shared tenants, selective for dedicated
- **Migration Time**: 90% reduction for schema changes

**Monitoring & Alerting**:
- **Before**: 100 database monitoring dashboards
- **After**: 6 monitoring dashboards with tenant-aware metrics
- **Alert Noise**: 80% reduction in database alerts

### Connection Pool Efficiency

**Shared Database**: Efficient connection reuse across tenants
- **Connection Utilization**: 70-90% vs 10-30% in dedicated approach
- **Memory Usage**: 60% reduction in connection overhead

## Migration Decision Matrix

| Trigger Condition | Action | Timeline | Risk Level |
|-------------------|--------|-----------| ---------- |
| Product count > 10k | Migrate to dedicated | 48-72 hours | Medium |
| Monthly requests > 1M | Migrate to dedicated | 24-48 hours | Medium |
| Plan upgrade to Enterprise | Migrate to dedicated | 24-48 hours | Low |
| Compliance requirement | Immediate dedicated | 4-8 hours | High |
| Performance issues | Evaluate migration | Variable | Medium |
| Data size > 10GB | Migrate to dedicated | 72+ hours | High |

## Implementation Best Practices

### Database Connection Management
```go
// Smart connection routing
func (db *DatabaseManager) GetConnection(tenantID uuid.UUID) (*gorm.DB, error) {
    tenant, err := db.getTenant(tenantID)
    if err != nil {
        return nil, err
    }
    
    if tenant.Type == TenantTypeDedicated {
        return db.getDedicatedConnection(tenant)
    }
    
    return db.getSharedConnection(tenantID)
}

func (db *DatabaseManager) getSharedConnection(tenantID uuid.UUID) (*gorm.DB, error) {
    conn := db.sharedPool.Get()
    // Set tenant context for row-level security
    return conn.Exec("SET app.current_tenant_id = ?", tenantID)
}
```

### Query Pattern Optimization
```sql
-- Always include tenant_id in WHERE clauses for shared database
-- Good:
SELECT * FROM products WHERE tenant_id = ? AND status = 'active';

-- Bad (will fail with RLS):
SELECT * FROM products WHERE status = 'active';

-- Efficient batch operations
UPDATE products 
SET status = 'inactive' 
WHERE tenant_id = ? AND category_id = ?;
```

## Compliance & Security

### Data Isolation Guarantees

**Shared Database Security**:
- Row-level security policies enforce tenant isolation
- Application-level tenant context validation
- Database-level access controls and audit logging
- Encrypted connections and at-rest encryption

**Dedicated Database Security**:
- Complete physical and logical separation
- Independent access controls and audit trails
- Customizable encryption and compliance settings
- Isolated backup and disaster recovery

### Regulatory Compliance

**GDPR Compliance**:
- **Shared**: Data deletion via tenant_id filtering
- **Dedicated**: Complete database removal option
- **Data Portability**: Tenant-scoped export capabilities
- **Right to be Forgotten**: Automated data purging

## Implementation Strategy

### Migration Decision Logic
```go
// Auto-migration when tenant grows beyond shared database limits
func (s *TenantService) CheckAndMigrateTenant(tenantID uuid.UUID) error {
    tenant, err := s.GetTenant(tenantID)
    if err != nil {
        return err
    }
    
    // Check if migration needed based on thresholds
    if tenant.ShouldUseDedicatedDatabase() && tenant.Type == TenantTypeShared {
        return s.MigrateToDedicatedDatabase(tenant)
    }
    
    return nil
}

const (
    ProductCountThreshold = 10000
    RequestVolumeThreshold = 100000 // per month
)

func (t *Tenant) ShouldMigrate() bool {
    return t.ProductCount >= ProductCountThreshold ||
           t.MonthlyRequests >= RequestVolumeThreshold ||
           t.Plan == PlanEnterprise
}
```

### Connection Management Implementation
```go
func (db *DatabaseManager) GetConnection(tenantID uuid.UUID) (*gorm.DB, error) {
    tenant, err := db.getTenant(tenantID)
    if err != nil {
        return nil, err
    }
    
    if tenant.Type == TenantTypeDedicated {
        return db.getDedicatedConnection(tenant)
    }
    
    return db.getSharedConnection(tenantID)
}

func (db *DatabaseManager) getSharedConnection(tenantID uuid.UUID) (*gorm.DB, error) {
    conn := db.sharedPool.Get()
    // Set tenant context for row-level security
    return conn.Exec("SET app.current_tenant_id = ?", tenantID)
}

func (db *DatabaseManager) getDedicatedConnection(tenant *Tenant) (*gorm.DB, error) {
    if tenant.DatabaseName == nil {
        return nil, errors.New("no dedicated database configured")
    }
    
    dsn := fmt.Sprintf("postgres://user:pass@host/%s", *tenant.DatabaseName)
    return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
```

### Migration Scenarios

**Startup → Growth Path**:
1. **Starter Plan**: Shared database (0-500 products)
2. **Professional Plan**: Shared database (500-5k products)
3. **Pro Plan**: Shared database (5k-10k products)  
4. **Enterprise Plan**: Dedicated database (10k+ products)

**Migration Triggers**:
- Product count exceeds 10,000
- Monthly API requests exceed 100,000
- Plan upgrade to Enterprise
- Compliance requirements
- Performance issues

### Implementation Checklist

**Phase 1: Foundation**
- [x] Update Tenant entity with Type and DatabaseName fields
- [x] Create hybrid connection management
- [x] Implement Row Level Security policies  
- [x] Add tenant-scoped unique constraints

**Phase 2: Migration System**
- [x] Create migration service for dedicated databases
- [ ] Implement background job for usage monitoring
- [ ] Add automated migration triggers
- [ ] Create monitoring dashboard for tenant metrics

**Phase 3: Operations**
- [ ] Implement backup strategies for both approaches
- [ ] Add performance monitoring and alerting
- [ ] Create disaster recovery procedures
- [ ] Implement compliance reporting

## Future Considerations

### Microservices Evolution
When individual modules exceed 100k requests/month:
1. Extract module to separate service
2. Maintain hybrid database strategy per service
3. Use event-driven communication between services

### Multi-Region Support
- **Shared Database**: Regional read replicas
- **Dedicated Database**: Region-specific instances  
- **Data Residency**: Compliance with local regulations

This hybrid database strategy provides the optimal balance of cost efficiency, operational simplicity, and enterprise-grade scalability for the multi-tenant e-commerce SaaS platform.