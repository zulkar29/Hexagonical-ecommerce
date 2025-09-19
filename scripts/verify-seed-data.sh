#!/bin/bash
# Seed Data Verification Script

set -e

echo "🔍 Verifying seed data in the database..."

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to run SQL query and display results
run_query() {
    local query="$1"
    local description="$2"
    
    echo -e "\n${YELLOW}📊 $description${NC}"
    echo "----------------------------------------"
    
    docker-compose -f docker-compose.dev.yml exec -T postgres psql -U postgres -d ecommerce_saas_dev -c "$query" || {
        echo -e "${RED}❌ Failed to execute query${NC}"
        return 1
    }
}

# Check if database is accessible
echo -e "${GREEN}🚀 Starting seed data verification...${NC}"

# Verify tenants
run_query "SELECT id, name, subdomain, plan, status FROM tenants ORDER BY name;" "Tenants Overview"

# Verify users by tenant
run_query "
SELECT 
    t.name as tenant,
    u.email,
    u.role,
    u.first_name || ' ' || u.last_name as full_name,
    u.status
FROM users u 
JOIN tenants t ON u.tenant_id = t.id 
ORDER BY t.name, u.role;
" "Users by Tenant"

# Verify categories and products
run_query "
SELECT 
    t.name as tenant,
    c.name as category,
    COUNT(p.id) as product_count
FROM tenants t
LEFT JOIN categories c ON t.id = c.tenant_id
LEFT JOIN products p ON c.id = p.category_id
GROUP BY t.name, c.name
ORDER BY t.name, c.name;
" "Categories and Product Count"

# Verify products with pricing
run_query "
SELECT 
    t.name as tenant,
    p.name as product,
    p.price,
    p.inventory_quantity,
    p.status
FROM products p
JOIN tenants t ON p.tenant_id = t.id
ORDER BY t.name, p.price DESC;
" "Products with Pricing"

# Verify product variants
run_query "
SELECT 
    p.name as product,
    pv.name as variant,
    pv.sku,
    pv.price,
    pv.inventory_quantity
FROM product_variants pv
JOIN products p ON pv.product_id = p.id
ORDER BY p.name, pv.name;
" "Product Variants"

# Verify orders
run_query "
SELECT 
    t.name as tenant,
    o.order_number,
    u.first_name || ' ' || u.last_name as customer,
    o.status,
    o.total_amount,
    o.currency,
    o.payment_status
FROM orders o
JOIN tenants t ON o.tenant_id = t.id
JOIN users u ON o.user_id = u.id
ORDER BY o.created_at DESC;
" "Orders Overview"

# Verify order items
run_query "
SELECT 
    o.order_number,
    oi.product_name,
    oi.quantity,
    oi.unit_price,
    oi.total_price
FROM order_items oi
JOIN orders o ON oi.order_id = o.id
ORDER BY o.order_number, oi.product_name;
" "Order Items Details"

# Verify payments
run_query "
SELECT 
    o.order_number,
    p.amount,
    p.currency,
    p.status,
    p.gateway,
    p.processed_at
FROM payments p
JOIN orders o ON p.order_id = o.id
ORDER BY p.processed_at DESC;
" "Payment Records"

# Verify permissions
run_query "
SELECT 
    role,
    COUNT(*) as permission_count
FROM role_permissions
GROUP BY role
ORDER BY role;
" "Permissions by Role"

# Verify analytics data
run_query "
SELECT 
    event_type,
    event_name,
    COUNT(*) as event_count
FROM analytics_events
GROUP BY event_type, event_name
ORDER BY event_type, event_name;
" "Analytics Events Summary"

# Final summary
run_query "
SELECT 
    'Tenants' as entity,
    COUNT(*) as count
FROM tenants
UNION ALL
SELECT 
    'Users' as entity,
    COUNT(*) as count
FROM users
UNION ALL
SELECT 
    'Categories' as entity,
    COUNT(*) as count
FROM categories
UNION ALL
SELECT 
    'Products' as entity,
    COUNT(*) as count
FROM products
UNION ALL
SELECT 
    'Product Variants' as entity,
    COUNT(*) as count
FROM product_variants
UNION ALL
SELECT 
    'Orders' as entity,
    COUNT(*) as count
FROM orders
UNION ALL
SELECT 
    'Order Items' as entity,
    COUNT(*) as count
FROM order_items
UNION ALL
SELECT 
    'Payments' as entity,
    COUNT(*) as count
FROM payments
UNION ALL
SELECT 
    'Permissions' as entity,
    COUNT(*) as count
FROM permissions
UNION ALL
SELECT 
    'Role Permissions' as entity,
    COUNT(*) as count
FROM role_permissions
UNION ALL
SELECT 
    'Analytics Events' as entity,
    COUNT(*) as count
FROM analytics_events
ORDER BY entity;
" "Database Summary"

echo -e "\n${GREEN}✅ Seed data verification complete!${NC}"
echo -e "${GREEN}🎉 Your e-commerce SaaS platform is ready with comprehensive demo data!${NC}"

echo -e "\n${YELLOW}📝 Demo Login Credentials:${NC}"
echo "Email: admin@techhub.com"
echo "Password: password123"
echo ""
echo "API Base URL: http://localhost:8080/api/v1"
echo "Storefront: http://localhost:3002"
echo "Admin Dashboard: http://localhost:3001"