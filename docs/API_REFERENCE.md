# Hexagonal E-commerce SaaS API Documentation

📋 **Documentation Navigation**: [📖 Project Home](../README.md) | [🏗️ Architecture](./ARCHITECTURE.md) | [🚀 Features](./FEATURES.md) | [🔄 User Flows](./USER_FLOWS.md)

Comprehensive REST API specification for the hexagonal e-commerce SaaS platform with **271 documented endpoints** covering all business operations across 27 active modules with multi-tenant architecture, authentication, and WebSocket real-time capabilities.

## Base URL
```
Development: http://localhost:8080/api/v1
Production: https://api.yourplatform.com/api/v1
```

## Naming Conventions

### URL Structure
- **Lowercase**: All URLs use lowercase letters
- **Hyphens**: Use hyphens for multi-word paths (`/forgot-password`, not `/forgotPassword`)
- **Plurals**: Resource collections use plural nouns (`/products`, `/orders`)
- **Nested Resources**: Use hierarchical structure (`/products/:id/variants`)

### HTTP Methods
- **GET**: Retrieve resources (read-only)
- **POST**: Create new resources
- **PUT**: Update entire resources or specific actions
- **PATCH**: Partial updates (status changes)
- **DELETE**: Remove resources

### Path Parameters
- **Consistent naming**: Use descriptive parameter names
  - `:id` for primary resource IDs
  - `:product-id`, `:customer-id` for related resource IDs
  - `:variant-id`, `:address-id` for nested resources

### Query Parameters
- **Filtering**: `?status=active&category=electronics`
- **Pagination**: `?page=1&limit=20`
- **Sorting**: `?sort=created_at&order=desc`
- **Search**: `?search=keyword`
- **Views**: `?view=stats` (statistics), `?view=trends` (trends analysis)
- **Stock Filtering**: `?stock=low` (low stock items)
- **User Context**: `?customer=current` (current user's data)

### Flexible Endpoint Design
Most list endpoints support multiple query parameters for flexible data retrieval:
```
GET /products?search=laptop&stock=low&view=stats
GET /orders?customer=current&status=shipped  
GET /reviews?status=pending&sort=created_at&order=desc
```

## Authentication
All API requests require authentication via JWT token in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

## Multi-tenant Context
Tenant context is resolved from:
1. Custom domain (e.g., store.example.com)
2. Subdomain (e.g., store.platform.com)
3. X-Tenant-ID header

## Error Handling
All API endpoints return standard HTTP status codes and JSON error responses:

### Standard Status Codes
- **200 OK**: Successful request
- **201 Created**: Resource successfully created
- **400 Bad Request**: Invalid request parameters
- **401 Unauthorized**: Invalid or missing authentication
- **403 Forbidden**: Insufficient permissions
- **404 Not Found**: Resource not found
- **422 Unprocessable Entity**: Validation errors
- **429 Too Many Requests**: Rate limit exceeded
- **500 Internal Server Error**: Server error

### Error Response Format
```json
{
  "error": "error_code",
  "message": "Human readable error message",
  "details": {
    "field": "validation error details"
  }
}
```

## Core Endpoints

### System Health
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/health` | Basic health check | ❌ | ❌ |

## User Module (23 endpoints)

### Authentication (Token-based)
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/auth/register` | User registration (admin/customer, returns JWT) | ❌ | ✅ |
| POST | `/auth/login` | User login (admin/customer, returns JWT) | ❌ | ✅ |
| POST | `/auth/refresh` | Refresh JWT token | ❌ | ✅ |
| POST | `/auth/logout` | Invalidate JWT token | ✅ | ✅ |
| POST | `/auth/forgot-password` | Request password reset token | ❌ | ✅ |
| POST | `/auth/reset-password` | Reset password with token | ❌ | ✅ |
| POST | `/auth/verify-email` | Verify email with token | ❌ | ✅ |
| POST | `/auth/resend-verification` | Resend verification token | ❌ | ✅ |
| POST | `/auth/verify-phone` | Verify phone number with OTP | ❌ | ✅ |
| POST | `/auth/resend-phone-otp` | Resend phone verification OTP | ❌ | ✅ |

### User Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/users/profile` | Get user profile | ✅ | ❌ |
| PUT | `/users/profile` | Update user profile | ✅ | ❌ |
| POST | `/users/change-password` | Change user password | ✅ | ❌ |
| PATCH | `/users/account` | Manage account (?action=deactivate\|anonymize\|delete\|export_data) | ✅ | ❌ |
| GET | `/users/preferences` | Get user preferences | ✅ | ❌ |
| PUT | `/users/preferences` | Update user preferences | ✅ | ❌ |
| GET | `/users` | List users (admin) | ✅ | ❌ |
| GET | `/users/:id` | Get user by ID (admin) | ✅ | ❌ |
| GET | `/users/:id/activity` | Get user activity logs | ✅ | ❌ |
| PATCH | `/users/:id` | Update user (status, role, profile, etc.) | ✅ | ❌ |
| POST | `/users/bulk-import` | Bulk import users | ✅ | ❌ |
| POST | `/users/export` | Export user data | ✅ | ❌ |
| GET | `/users/:id/orders` | Get user's orders | ✅ | ❌ |
| GET | `/users/:id/addresses` | Get user's addresses | ✅ | ❌ |

## Admin Dashboard Module (8 endpoints)

### Tenant Admin Management (Store Management)
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/admin/dashboard` | Get tenant admin dashboard (?period=day\|week\|month, ?metrics=sales\|orders\|customers) | ✅ | ✅ |
| GET | `/admin/quick-stats` | Get quick statistics (orders, revenue, customers, products) | ✅ | ✅ |
| GET | `/admin/staff` | List tenant staff (?role=admin\|manager\|staff, ?status=active\|inactive) | ✅ | ✅ |
| PATCH | `/admin/staff/:id` | Manage tenant staff (?action=create\|update\|delete\|assign_roles\|change_status) | ✅ | ✅ |
| GET | `/admin/roles` | List tenant roles & permissions (?include_permissions=true) | ✅ | ✅ |
| PATCH | `/admin/roles/:id` | Manage tenant roles (?action=create\|update\|delete\|assign_permissions) | ✅ | ✅ |
| GET | `/admin/activity-logs` | Get tenant admin activity logs (?user_id=id, ?action=login\|update\|delete, ?date_from=date, ?date_to=date) | ✅ | ✅ |
| GET | `/admin/system-health` | Get tenant system health & performance metrics | ✅ | ✅ |

## Platform Admin Module (15 endpoints)

### Super Admin Operations (Platform Management)
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/platform/dashboard` | Get platform admin dashboard (?period=day\|week\|month, ?metrics=tenants\|revenue\|users\|system) | ✅ | ❌ |
| GET | `/platform/stats` | Get platform statistics (total tenants, revenue, growth metrics) | ✅ | ❌ |
| GET | `/platform/system-status` | Get overall system status and health | ✅ | ❌ |
| GET | `/platform/audit-logs` | Get platform audit logs (?user_id=id, ?tenant_id=id, ?action=create\|update\|delete, ?date_from=date, ?date_to=date) | ✅ | ❌ |

### Platform User & Access Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/platform/admins` | List platform administrators (?role=super_admin\|platform_admin\|support) | ✅ | ❌ |
| PATCH | `/platform/admins/:id` | Manage platform admins (?action=create\|update\|delete\|assign_roles\|change_status) | ✅ | ❌ |
| GET | `/platform/roles` | List platform roles & permissions | ✅ | ❌ |
| PATCH | `/platform/roles/:id` | Manage platform roles (?action=create\|update\|delete\|assign_permissions) | ✅ | ❌ |

### Platform Configuration
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/platform/settings` | Get platform settings (?category=billing\|notifications\|security\|features) | ✅ | ❌ |
| PATCH | `/platform/settings` | Update platform settings | ✅ | ❌ |

### Platform Tenant Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/platform/tenants` | Create new tenant (platform onboarding) | ✅ | ❌ |
| GET | `/platform/tenants` | List all tenants (super admin) (?view=stats, ?status=active\|pending\|suspended, ?include=users\|revenue\|usage) | ✅ | ❌ |
| GET | `/platform/tenants/:id` | Get tenant details (super admin view) (?include=subscription\|usage\|settings\|analytics) | ✅ | ❌ |
| PATCH | `/platform/tenants/:id` | Update tenant (super admin) (?action=activate\|suspend\|deactivate\|update_plan\|force_billing) | ✅ | ❌ |
| DELETE | `/platform/tenants/:id` | Delete tenant (with data retention policy) | ✅ | ❌ |

## Security Module (16 endpoints)

### Security & Audit Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/security/audit-logs` | Get audit trail (?user_id=id, ?action=create\|update\|delete, ?resource=type, ?date_from=date, ?date_to=date) | ✅ | ✅ |
| GET | `/security/failed-logins` | Get failed login attempts (?ip=address, ?user_id=id, ?threshold=count) | ✅ | ✅ |
| PATCH | `/security/user/:id/lock` | Lock/unlock user account (?action=lock\|unlock\|temporary_lock) | ✅ | ✅ |
| GET | `/security/fraud-detection` | Get fraud detection alerts (?status=pending\|resolved, ?risk_level=high\|medium\|low) | ✅ | ✅ |
| PATCH | `/security/fraud-alerts/:id` | Manage fraud alert (?action=investigate\|resolve\|escalate) | ✅ | ✅ |
| GET | `/security/tenant-boundaries` | Validate tenant data isolation (?tenant_id=id, ?check_type=data\|access\|query) | ✅ | ❌ |
| POST | `/security/vulnerability-scan` | Run security vulnerability scan | ✅ | ✅ |
| GET | `/security/permissions/:user_id` | Get user permissions matrix | ✅ | ✅ |
| POST | `/security/2fa/setup` | Setup two-factor authentication | ✅ | ✅ |
| POST | `/security/2fa/verify` | Verify 2FA code | ✅ | ✅ |

### API Key Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/security/api-keys` | List tenant API keys | ✅ | ✅ |
| POST | `/security/api-keys` | Create new API key | ✅ | ✅ |
| PATCH | `/security/api-keys/:id` | Manage API key (?action=regenerate\|revoke\|activate\|deactivate) | ✅ | ✅ |
| GET | `/security/api-keys/:id/usage` | Get API key usage statistics | ✅ | ✅ |

### Rate Limiting
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/security/rate-limits` | Get current rate limits for tenant | ✅ | ✅ |
| PATCH | `/security/rate-limits` | Update rate limits (?endpoint=specific\|global, ?action=adjust\|reset) | ✅ | ✅ |

## Product Module (31 endpoints)

*Note: Stock management is handled within product endpoints for this single-vendor platform*

### Products
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/products` | Create product | ✅ | ✅ |
| GET | `/products` | List products (?search=keyword, ?stock=low\|out, ?status=active\|draft, ?category=id, ?sort=name\|price\|created_at, ?view=stats\|analytics, pagination) | ✅ | ✅ |
| GET | `/products/:id` | Get product details | ✅ | ✅ |
| PATCH | `/products/:id` | Update product (?action=duplicate\|archive\|publish\|update_inventory\|adjust_stock) | ✅ | ✅ |
| DELETE | `/products/:id` | Delete product | ✅ | ✅ |
| POST | `/products/bulk` | Bulk operations (?operation=import\|export\|update\|delete) | ✅ | ✅ |
| GET | `/products/slug/:slug` | Get product by slug (storefront) | ❌ | ✅ |
| POST | `/products/:id/images` | Upload product images | ✅ | ✅ |
| DELETE | `/products/:id/images/:image_id` | Delete product image | ✅ | ✅ |
| GET | `/products/:id/analytics` | Get product analytics (?type=related\|history\|performance) | ✅ | ✅ |

### Product Tags
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/products/tags` | List product tags (?search=keyword, ?popular=true) | ✅ | ✅ |
| POST | `/products/tags` | Create product tag | ✅ | ✅ |
| PATCH | `/products/tags/:id` | Update tag (?action=update\|delete) | ✅ | ✅ |
| POST | `/products/:id/tags` | Assign tags to product | ✅ | ✅ |
| DELETE | `/products/:id/tags/:tag_id` | Remove tag from product | ✅ | ✅ |

### Stock Reservation
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/products/reserve` | Reserve stock for cart items | ✅ | ✅ |
| DELETE | `/products/reserve/:reservation_id` | Release stock reservation | ✅ | ✅ |
| GET | `/products/reservations` | List active reservations | ✅ | ✅ |

### Stock Management
| Method | Endpoint | Description | Auth | Public |
|--------|-----|-------------|------|--------|
| GET | `/products/low-stock` | Get low stock alerts (?threshold=10, ?category=id, ?urgent_only=true) | ✅ | ✅ |
| POST | `/products/stock/audit` | Generate stock audit report (?format=csv\|json) | ✅ | ✅ |

### Product Variants
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/products/:id/variants` | Create product variant | ✅ | ✅ |
| GET | `/products/:id/variants` | Get product variants | ✅ | ✅ |
| PATCH | `/products/:id/variants/:variant_id` | Update variant (?action=update\|delete\|upload_image) | ✅ | ✅ |

### Categories
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/categories` | Create category | ✅ | ✅ |
| GET | `/categories` | List categories (?parent_id=id, ?view=tree\|flat\|stats, ?include_products=true) | ✅ | ✅ |
| GET | `/categories/:id` | Get category details (?include=children\|products\|stats) | ✅ | ✅ |
| PATCH | `/categories/:id` | Update category (?action=move\|reorder\|upload_image\|delete_image) | ✅ | ✅ |
| DELETE | `/categories/:id` | Delete category | ✅ | ✅ |

### Public Product Access (Storefront)
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/public/products` | Browse products (?search=keyword, ?category=id, ?sort=price\|name\|rating, pagination) | ❌ | ✅ |
| GET | `/public/products/:id` | Get product details (?include=variants\|reviews\|related) | ❌ | ✅ |
| GET | `/public/categories` | Browse categories (?view=tree\|flat, ?include_products=true) | ❌ | ✅ |

## Order Module (11 endpoints)

### Orders
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/orders` | Create order | ✅ | ✅ |
| GET | `/orders` | List orders (?customer=current\|id, ?status=pending\|processing\|shipped, ?view=stats\|analytics, ?search=number, pagination) | ✅ | ✅ |
| GET | `/orders/:id` | Get order details (?include=items\|notes\|history\|invoice\|documents) | ✅ | ✅ |
| PATCH | `/orders/:id` | Update order (?action=cancel\|fulfill\|refund\|add_item\|remove_item\|add_note\|update_status) | ✅ | ✅ |
| POST | `/orders/bulk` | Bulk operations (?operation=export\|update_status\|cancel\|assign_fulfillment) | ✅ | ✅ |
| GET | `/orders/lookup/:number` | Get order by number (?public=true for customer access) | ❌/✅ | ✅ |
| GET | `/orders/:id/tracking` | Track order (?public=true for customer access) | ❌/✅ | ✅ |

### Order Disputes
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/orders/:id/disputes` | Create order dispute | ✅ | ✅ |
| GET | `/orders/disputes` | List order disputes (?status=pending\|resolved, ?customer=current) | ✅ | ✅ |
| GET | `/orders/disputes/:id` | Get dispute details | ✅ | ✅ |
| PATCH | `/orders/disputes/:id` | Update dispute (?action=resolve\|escalate\|close\|add_evidence) | ✅ | ✅ |

## Payment Module (6 endpoints)

### Payment Processing
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/payments` | Create payment | ✅ | ✅ |
| GET | `/payments` | List payments (?status=pending\|completed\|failed, ?method=sslcommerz\|cod, ?view=stats) | ✅ | ✅ |
| GET | `/payments/:id` | Get payment details | ✅ | ✅ |
| PATCH | `/payments/:id` | Update payment (?action=process\|refund\|capture\|void) | ✅ | ✅ |
| GET | `/payments/methods` | List payment methods (?active=true) | ✅ | ✅ |
| PATCH | `/payments/methods/:id` | Update payment method (?action=enable\|disable\|update) | ✅ | ✅ |

### Payment Webhooks
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/webhooks/payment/:provider` | Payment provider webhook (sslcommerz primary, others as fallback) | ❌ | ❌ |

## Shipping Module (11 endpoints)

### Shipping Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/shipping/zones` | List shipping zones (?view=stats) | ✅ | ✅ |
| PATCH | `/shipping/zones/:id` | Manage shipping zone (?action=create\|update\|delete) | ✅ | ✅ |
| GET | `/shipping/rates` | List shipping rates (?zone_id=id) | ✅ | ✅ |
| PATCH | `/shipping/rates/:id` | Manage shipping rate (?action=create\|update\|delete) | ✅ | ✅ |
| POST | `/shipping/calculate` | Calculate shipping cost | ✅ | ✅ |
| GET | `/shipping/labels` | List shipping labels (?order_id=id) | ✅ | ✅ |
| PATCH | `/shipping/labels/:id` | Manage shipping label (?action=create\|update\|delete) | ✅ | ✅ |
| GET | `/shipping/track/:tracking_number` | Track shipment | ✅ | ✅ |
| GET | `/shipping/providers` | List shipping providers (?active=true) | ✅ | ✅ |
| PATCH | `/shipping/providers/:id` | Manage provider (?action=create\|update\|delete\|enable\|disable) | ✅ | ✅ |

### Shipping Webhooks
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/webhooks/shipping/:provider` | Shipping provider webhook (pathao, redx, paperfly, dhl, fedex) | ❌ | ❌ |

## Notification Module (8 endpoints)

### Notification Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/notifications` | Send notification (?type=email\|sms\|push) | ✅ | ✅ |
| GET | `/notifications` | List notifications (?view=stats, ?status=read\|unread) | ✅ | ✅ |
| GET | `/notifications/:id` | Get notification details | ✅ | ✅ |
| PATCH | `/notifications/:id` | Update notification (?action=mark_read\|mark_unread\|delete) | ✅ | ✅ |
| GET | `/notifications/templates` | List notification templates (?type=email\|sms) | ✅ | ✅ |
| PATCH | `/notifications/templates/:id` | Manage template (?action=create\|update\|delete) | ✅ | ✅ |
| GET | `/notifications/preferences` | Get notification preferences | ✅ | ✅ |
| PATCH | `/notifications/preferences` | Update notification preferences | ✅ | ✅ |


## Analytics Module (10 endpoints)

### Analytics
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/analytics/track` | Track events (?type=event\|page_view\|product_view\|purchase) | ❌ | ✅ |
| GET | `/analytics/dashboard` | Get dashboard analytics (?type=traffic\|sales\|realtime\|overview, ?period=day\|week\|month) | ✅ | ✅ |
| GET | `/analytics/insights` | Get insights (?type=top_products\|top_pages\|referrers\|cohorts\|funnel\|retention\|clv) | ✅ | ✅ |
| GET | `/analytics/reports` | List reports (?status=scheduled\|completed) | ✅ | ✅ |
| PATCH | `/analytics/reports/:id` | Manage reports (?action=generate\|schedule\|export\|delete) | ✅ | ✅ |

### Customer Segmentation
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/analytics/segments` | Create customer segment (?criteria=purchase_amount\|frequency\|category\|geography) | ✅ | ✅ |
| GET | `/analytics/segments` | List customer segments | ✅ | ✅ |
| GET | `/analytics/segments/:id` | Get segment details (?include=customers\|stats) | ✅ | ✅ |
| GET | `/analytics/segments/:id/customers` | Get customers in segment (?export=csv, pagination) | ✅ | ✅ |
| PATCH | `/analytics/segments/:id` | Update segment (?action=update\|refresh\|delete) | ✅ | ✅ |

## Marketing Module (10 endpoints)

### Marketing
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/marketing/campaigns` | List campaigns (?view=stats\|performance) | ✅ | ✅ |
| GET | `/marketing/campaigns/:id` | Get campaign details (?include=stats\|emails) | ✅ | ✅ |
| PATCH | `/marketing/campaigns/:id` | Manage campaign (?action=create\|update\|delete\|schedule\|send) | ✅ | ✅ |
| GET | `/marketing/templates` | List email templates | ✅ | ✅ |
| PATCH | `/marketing/templates/:id` | Manage template (?action=create\|update\|delete) | ✅ | ✅ |
| GET | `/marketing/segments` | List customer segments | ✅ | ✅ |
| PATCH | `/marketing/segments/:id` | Manage segment (?action=create\|update\|delete\|refresh) | ✅ | ✅ |
| POST | `/marketing/newsletter` | Manage newsletter (?action=subscribe\|unsubscribe) | ❌ | ✅ |
| GET | `/marketing/subscribers` | List newsletter subscribers | ✅ | ✅ |
| GET | `/marketing/abandoned-carts` | Get abandoned carts | ✅ | ✅ |

## Discount Module (9 endpoints)

### Discount Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/discounts` | List discounts (?type=coupon\|auto, ?view=stats\|performance, ?status=active\|expired) | ✅ | ✅ |
| GET | `/discounts/:id` | Get discount details (?include=usage\|analytics) | ✅ | ✅ |
| PATCH | `/discounts/:id` | Manage discount (?action=create\|update\|delete\|activate\|deactivate) | ✅ | ✅ |
| POST | `/discounts/validate` | Validate discount code (?code=discount_code) | ❌ | ✅ |
| POST | `/discounts/apply` | Apply discount (?order_id=id, ?code=discount_code) | ❌ | ✅ |

### Gift Cards & Store Credit
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/gift-cards` | List gift cards (?customer_id=id) | ✅ | ✅ |
| PATCH | `/gift-cards/:id` | Manage gift card (?action=create\|update\|delete\|validate\|use) | ✅/❌ | ✅ |
| GET | `/store-credit/:customer_id` | Get customer store credit | ✅ | ✅ |
| PATCH | `/store-credit/:customer_id` | Update store credit (?action=add\|use\|refund) | ✅ | ✅ |

## Search Module (6 endpoints)

### Search Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/search` | Global search (?q=query, ?type=products\|customers\|orders, ?filters=category\|price\|brand) | ❌/✅ | ✅ |
| GET | `/search/products` | Product search (?q=query, ?category=id, ?price_min=amount, ?price_max=amount, ?sort=relevance\|price\|rating) | ❌ | ✅ |
| GET | `/search/suggestions` | Search suggestions (?q=query, ?type=products\|categories) | ❌ | ✅ |
| GET | `/search/analytics` | Search analytics (?view=popular_terms\|no_results\|trends) | ✅ | ✅ |
| POST | `/search/reindex` | Reindex search data (?type=products\|categories\|all) | ✅ | ✅ |
| GET | `/search/filters` | Get available search filters | ❌ | ✅ |

## Settings Module (5 endpoints)

### Settings Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/settings/store` | Get store settings (?view=public\|private) | ✅ | ✅ |
| PATCH | `/settings/store` | Update store settings (?action=update\|reset) | ✅ | ✅ |
| GET | `/settings/payments` | Get payment settings | ✅ | ✅ |

### Maintenance Mode
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/settings/maintenance` | Get maintenance mode status | ✅ | ✅ |
| PATCH | `/settings/maintenance` | Toggle maintenance mode (?action=enable\|disable, ?message=custom_message) | ✅ | ✅ |

## Public Access Module (3 endpoints)

### Public Access & Customer Registration
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/public/register` | Customer self-registration | ❌ | ✅ |
| GET | `/public/pages/:slug` | Get public pages (about, terms, privacy, etc.) | ❌ | ✅ |
| GET | `/public/content/menus` | Get store navigation menus | ❌ | ✅ |

## Address Module (7 endpoints)

### Customer Address Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/addresses` | List user addresses (?type=billing\|shipping\|all, ?default=true) | ✅ | ✅ |
| POST | `/addresses` | Create new address | ✅ | ✅ |
| GET | `/addresses/:id` | Get address details | ✅ | ✅ |
| PATCH | `/addresses/:id` | Update address (?action=update\|set_default\|delete) | ✅ | ✅ |
| DELETE | `/addresses/:id` | Delete address | ✅ | ✅ |

### Address Validation & Services  
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/addresses/validate` | Validate address (with postal service) | ✅ | ✅ |
| GET | `/addresses/postal-codes/:code` | Get area info by postal code | ❌ | ✅ |

## Returns Module (6 endpoints)

### Customer Returns Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/returns` | Request return/refund | ✅ | ✅ |
| GET | `/returns` | List returns (?status=pending\|approved\|rejected\|completed, ?customer=current) | ✅ | ✅ |
| GET | `/returns/:id` | Get return details (?include=messages\|shipping_label) | ✅ | ✅ |
| PATCH | `/returns/:id` | Update return (?action=approve\|reject\|complete\|cancel\|add_message) | ✅ | ✅ |

### Return Policies & Processing
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/returns/policies` | Get return policies | ❌ | ✅ |
| POST | `/returns/:id/refund` | Process refund (?method=original\|store_credit\|manual) | ✅ | ✅ |

## Finance Module (6 endpoints)

### Financial Management & Reporting
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/finance/reports` | Get financial reports (?type=revenue\|expenses\|profit_loss, ?period=day\|week\|month\|year) | ✅ | ✅ |
| GET | `/finance/transactions` | List transactions (?type=sale\|refund\|fee, ?status=pending\|completed\|failed) | ✅ | ✅ |
| GET | `/finance/reconciliation` | Get payment reconciliation (?date_from=date, ?date_to=date, ?gateway=sslcommerz) | ✅ | ✅ |
| POST | `/finance/payouts` | Request payout (?method=bank_transfer\|mobile_banking) | ✅ | ✅ |
| GET | `/finance/dashboard` | Get financial dashboard overview | ✅ | ✅ |

## Reviews Module (10 endpoints)

### Review Management  
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/reviews` | Create review | ✅ | ✅ |
| GET | `/reviews` | List reviews (?product_id=id, ?status=pending\|approved, ?view=stats\|trends\|top_products) | ✅ | ✅ |
| GET | `/reviews/:id` | Get review details (?include=replies) | ✅ | ✅ |
| PATCH | `/reviews/:id` | Update review (?action=approve\|reject\|spam\|reply\|react) | ✅ | ✅ |
| DELETE | `/reviews/:id` | Delete review | ✅ | ✅ |
| POST | `/reviews/bulk` | Bulk operations (?operation=moderate\|approve\|reject) | ✅ | ✅ |
| GET | `/reviews/invitations` | Manage review invitations (?token=invitation_token for public access) | ❌/✅ | ❌/✅ |
| GET | `/reviews/settings` | Get/update review settings (?action=update) | ✅ | ✅ |

### Public Reviews
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/public/reviews/:product_id` | Get product reviews (public) | ❌ | ✅ |
| POST | `/public/reviews` | Submit product review (customer) | ✅ | ✅ |

## Support & Contact Module (14 endpoints)

### Support Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/support/tickets` | Create support ticket | ✅ | ✅ |
| GET | `/support/tickets` | List tickets (?status=open\|closed, ?assigned_to=user_id, ?view=stats) | ✅ | ✅ |
| GET | `/support/tickets/:id` | Get ticket details (?include=messages\|history) | ✅ | ✅ |
| PATCH | `/support/tickets/:id` | Update ticket (?action=assign\|resolve\|close\|add_message\|delete) | ✅ | ✅ |
| GET | `/support/faqs` | List FAQs (?category=category_name) | ✅ | ✅ |
| PATCH | `/support/faqs/:id` | Manage FAQ (?action=create\|update\|delete) | ✅ | ✅ |
| GET | `/support/knowledge-base` | List articles (?slug=article_slug) | ✅ | ✅ |
| GET | `/support/settings` | Get support settings | ✅ | ✅ |

### Contact Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/contacts` | List contacts (?status=new\|replied, ?view=analytics\|export, ?assigned_to=user_id) | ✅ | ✅ |
| GET | `/contacts/:id` | Get contact details (?include=interactions\|notes\|replies) | ✅ | ✅ |
| PATCH | `/contacts/:id` | Update contact (?action=create\|reply\|assign\|close\|add_note\|delete) | ✅ | ✅ |
| POST | `/contact-forms` | Submit contact form (?form_type=support\|general) | ❌ | ✅ |
| GET | `/contact-templates` | List contact templates | ✅ | ✅ |
| PATCH | `/contact-settings` | Update contact settings | ✅ | ✅ |

## Content Management Module (8 endpoints)

### Content Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/content/pages` | List pages (?status=published\|draft) | ✅ | ✅ |
| PATCH | `/content/pages/:id` | Manage page (?action=create\|update\|delete\|publish\|unpublish) | ✅ | ✅ |
| GET | `/content/posts` | List posts (?status=published\|draft) | ✅ | ✅ |
| PATCH | `/content/posts/:id` | Manage post (?action=create\|update\|delete\|publish\|unpublish) | ✅ | ✅ |
| GET | `/content/media` | List media files | ✅ | ✅ |
| POST | `/content/media` | Upload media | ✅ | ✅ |
| GET | `/content/menus` | List menus | ✅ | ✅ |
| PATCH | `/content/menus/:id` | Manage menu (?action=create\|update\|delete) | ✅ | ✅ |

## Webhook Module (7 endpoints)

### Webhook Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/webhooks/endpoints` | List webhook endpoints | ✅ | ✅ |
| PATCH | `/webhooks/endpoints/:id` | Manage endpoint (?action=create\|update\|delete\|test\|enable\|disable) | ✅ | ✅ |
| GET | `/webhooks/deliveries` | List webhook deliveries (?endpoint_id=id, ?status=success\|failed) | ✅ | ✅ |
| GET | `/webhooks/deliveries/:id` | Get delivery details | ✅ | ✅ |
| PATCH | `/webhooks/deliveries/:id` | Retry webhook delivery | ✅ | ✅ |
| GET | `/webhooks/events` | List webhook events | ✅ | ✅ |
| PATCH | `/webhooks/events/:id` | Manage event (?action=create\|update\|delete) | ✅ | ✅ |

## Billing Module (12 endpoints)

### Billing Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/billing/plans` | List billing plans | ✅ | ❌ |
| PATCH | `/billing/plans/:plan_id` | Manage plan (?action=create\|update\|delete) | ✅ | ❌ |
| GET | `/billing/subscriptions` | Get subscription details | ✅ | ✅ |
| PATCH | `/billing/subscriptions` | Manage subscription (?action=create\|update\|cancel\|upgrade\|downgrade) | ✅ | ✅ |
| GET | `/billing/usage` | Get usage summary (?include=limits) | ✅ | ✅ |
| POST | `/billing/usage` | Record usage | ✅ | ✅ |
| GET | `/billing/invoices` | List invoices (?view=analytics) | ✅ | ✅ |
| GET | `/billing/invoices/:invoice_id` | Get invoice details | ✅ | ✅ |
| PATCH | `/billing/invoices/:invoice_id` | Process invoice (?action=payment\|refund) | ✅ | ✅ |
| GET | `/billing/reports` | Get billing reports (?type=revenue\|churn) | ✅ | ❌ |

### Admin Operations
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/billing/admin/process` | Process billing (?operation=billing\|retry_payments\|dunning) | ✅ | ❌ |
| PATCH | `/billing/admin/tenants/:tenant_id` | Update tenant service status | ✅ | ❌ |

## Tenant Module (7 endpoints)


### Tenant Operations (Tenant Admin APIs)
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/tenants` | Register new tenant (self-registration) | ❌ | ❌ |
| GET | `/tenants/current` | Get current tenant details (tenant admin view) | ✅ | ✅ |
| PATCH | `/tenants/current` | Update current tenant (?action=update_profile\|update_settings\|change_plan) | ✅ | ✅ |

### Tenant Onboarding & Configuration  
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/tenants/check-subdomain/:subdomain` | Check subdomain availability | ❌ | ❌ |
| POST | `/tenants/:id/domain` | Configure custom domain | ✅ | ❌ |
| GET | `/tenants/:id/onboarding-status` | Get onboarding completion status | ✅ | ❌ |
| PATCH | `/tenants/:id/approve` | Approve tenant onboarding (?action=approve\|reject) | ✅ | ❌ |

### Tenant Subscription & Limits
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/tenants/:id/limits` | Get plan limits and usage (?check=products\|users\|storage) | ✅ | ❌ |
| PATCH | `/tenants/:id/migrate` | Migrate tenant (?action=upgrade_plan\|change_database\|export_data) | ✅ | ❌ |

## Observability Module (8 endpoints)

### System Monitoring
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/observability/health` | Get health status (?detailed=true) | ✅ | ❌ |
| GET | `/observability/metrics` | Get metrics (?view=summary\|detailed) | ✅ | ❌ |
| GET | `/observability/logs` | Get log entries | ✅ | ❌ |
| POST | `/observability/logs` | Create log entry | ✅ | ❌ |
| GET | `/observability/traces` | Get traces (?trace_id=id) | ✅ | ❌ |
| GET | `/observability/alerts` | List alerts | ✅ | ❌ |
| POST | `/observability/alerts` | Create alert | ✅ | ❌ |
| GET | `/observability/system` | Get system information (?view=stats) | ✅ | ❌ |

## Cart Module (8 endpoints)

### Cart Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/cart` | Get cart (?include=summary\|shipping_methods\|totals) | ✅ | ✅ |
| POST | `/cart/items` | Add item to cart | ✅ | ✅ |
| PATCH | `/cart/items/:id` | Update cart item (?action=update_quantity\|save_later\|move_to_cart\|remove) | ✅ | ✅ |
| PATCH | `/cart` | Update cart (?action=apply_discount\|remove_discount\|clear\|merge\|validate) | ✅ | ✅ |
| POST | `/cart/estimates` | Get estimates (?type=shipping\|total) | ✅ | ✅ |

### Guest Cart & Checkout (Token-based)
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/cart/guest` | Get guest cart (via guest token) | ❌ | ✅ |
| PATCH | `/cart/guest` | Update guest cart (?action=add_item\|update_item\|remove_item\|apply_discount) | ❌ | ✅ |
| POST | `/checkout/guest` | Process guest checkout | ❌ | ✅ |

## Wishlist Module (6 endpoints)

### Wishlist Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/wishlists` | List wishlists (?view=analytics\|popular_items) | ✅ | ✅ |
| POST | `/wishlists` | Create wishlist | ✅ | ✅ |
| GET | `/wishlists/:id` | Get wishlist details (?token=shared_token for public access) | ❌/✅ | ✅ |
| PATCH | `/wishlists/:id` | Update wishlist (?action=add_item\|remove_item\|move_to_cart\|share\|unshare) | ✅ | ✅ |
| DELETE | `/wishlists/:id` | Delete wishlist | ✅ | ✅ |
| POST | `/wishlists/bulk` | Bulk operations (?operation=add_items\|remove_items\|cleanup) | ✅ | ✅ |

## Real-time Features (WebSocket)

### WebSocket Connection
```
ws://localhost:8080/ws
wss://api.yourplatform.com/ws
```

### Real-time Event Types
- `inventory_updated` - Product inventory changes
- `order_created` - New orders
- `order_status_changed` - Order status updates
- `product_updated` - Product modifications
- `dashboard_metrics_updated` - Real-time dashboard updates
- `system_notification` - System alerts and notifications

## Summary

The e-commerce platform implements **271 documented API endpoints** across **27 optimized modules** with token-based REST design:

- ✅ **User Module** - 23 endpoints (authentication, phone verification, profile management, admin operations, bulk operations)
- ✅ **Admin Dashboard Module** - 8 endpoints (admin dashboard, staff management, roles, activity logs)
- ✅ **Platform Admin Module** - 15 endpoints (platform management, tenant management, system configuration)
- ✅ **Security Module** - 16 endpoints (audit logs, fraud detection, tenant boundaries, 2FA, vulnerability scanning, API key management, rate limiting)
- ✅ **Product Module** - 31 endpoints (products, variants, categories, tags, stock reservation, inventory, stock management, public access, analytics) 
- ✅ **Order Module** - 11 endpoints (order management, disputes, tracking, flexible operations)
- ✅ **Cart Module** - 8 endpoints (cart management, guest cart, checkout)
- ✅ **Wishlist Module** - 6 endpoints (wishlist management with bulk operations)
- ✅ **Address Module** - 7 endpoints (address management, validation, geocoding)
- ✅ **Payment Module** - 6 endpoints (payment processing, methods, webhooks)
- ✅ **Shipping Module** - 11 endpoints (zones, rates, labels, tracking, webhooks)
- ✅ **Notification Module** - 8 endpoints (notifications, templates, preferences)
- ✅ **Analytics Module** - 10 endpoints (tracking, dashboard, insights, reports, customer segmentation)
- ✅ **Marketing Module** - 10 endpoints (campaigns, templates, segments, automation)
- ✅ **Discount Module** - 9 endpoints (discounts, gift cards, store credit)
- ✅ **Search Module** - 6 endpoints (global search, product search, suggestions, analytics, filters)
- ✅ **Settings Module** - 5 endpoints (store settings, maintenance mode, SEO, appearance, integrations)
- ✅ **Public Access Module** - 3 endpoints (public pages, customer registration)
- ✅ **Reviews Module** - 10 endpoints (reviews, moderation, invitations, public reviews)
- ✅ **Support & Contact Module** - 14 endpoints (tickets, FAQ, knowledge base, contact management)
- ✅ **Content Management Module** - 8 endpoints (pages, posts, media, menus)
- ✅ **Webhook Module** - 7 endpoints (endpoint management, deliveries, events)
- ✅ **Billing Module** - 12 endpoints (plans, subscriptions, usage, invoices, admin)
- ✅ **Tenant Module** - 7 endpoints (tenant operations, onboarding, configuration)
- ✅ **Observability Module** - 8 endpoints (health, metrics, logs, alerts)
- ✅ **Returns Module** - 6 endpoints (return management, reasons, processing)
- ✅ **Finance Module** - 6 endpoints (ledger, transactions, reports, payouts)

### Key Optimizations
- **Enhanced Security**: Added comprehensive Security Module with audit trails, fraud detection, and tenant boundary validation
- **Improved Multi-tenancy**: Enhanced Tenant Module with onboarding, approval workflows, and migration capabilities  
- **Eliminated Duplicates**: Merged redundant Customer Module into User/Public Access modules
- **Consolidated Support**: Combined Contact and Support modules for better customer service management
- **Token-based REST**: No session management, pure JWT token authentication
- **Unified Authentication**: Single `/auth` endpoints for all users (admin/customer)
- **Unified Actions**: Single endpoints handle multiple operations via `?action` parameters
- **Flexible Queries**: Rich filtering and view options reduce need for separate endpoints  
- **Consistent Patterns**: PATCH for updates, GET with query params for filtering/views
- **Better Organization**: Logical grouping with consistent naming conventions aligned with USER_FLOWS.md

### Authentication & Security
- JWT-based authentication for all protected endpoints
- Multi-tenant architecture with proper isolation
- Public endpoints for customer-facing features
- WebSocket support for real-time updates

### Technology Stack
- **Backend Framework**: Gin (Go HTTP framework)
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT tokens
- **Real-time**: WebSocket connections
- **Caching**: Redis integration
- **Architecture**: Hexagonal Architecture pattern

---

*Note: All endpoints require appropriate authentication and authorization. Refer to the Authentication section for details on obtaining and using access tokens.*