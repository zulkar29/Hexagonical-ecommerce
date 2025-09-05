# API Documentation

Comprehensive REST API specification for the e-commerce SaaS platform with **200+ endpoints** covering all business operations across 27 active modules with multi-tenant architecture, authentication, and WebSocket real-time capabilities.

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

## Core Endpoints

### System Health
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/health` | Basic health check | ❌ | ❌ |

## User Module (17 endpoints)

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

### User Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/users/profile` | Get user profile | ✅ | ❌ |
| PUT | `/users/profile` | Update user profile | ✅ | ❌ |
| POST | `/users/change-password` | Change user password | ✅ | ❌ |
| DELETE | `/users/account` | Delete user account | ✅ | ❌ |
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

### Admin Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/admin/dashboard` | Get admin dashboard overview (?period=day\|week\|month, ?metrics=sales\|orders\|customers) | ✅ | ✅ |
| GET | `/admin/quick-stats` | Get quick statistics (orders, revenue, customers, products) | ✅ | ✅ |
| GET | `/admin/staff` | List admin staff (?role=admin\|manager\|staff, ?status=active\|inactive) | ✅ | ✅ |
| PATCH | `/admin/staff/:id` | Manage staff (?action=create\|update\|delete\|assign_roles\|change_status) | ✅ | ✅ |
| GET | `/admin/roles` | List roles & permissions (?include_permissions=true) | ✅ | ✅ |
| PATCH | `/admin/roles/:id` | Manage roles (?action=create\|update\|delete\|assign_permissions) | ✅ | ✅ |
| GET | `/admin/activity-logs` | Get admin activity logs (?user_id=id, ?action=login\|update\|delete, ?date_from=date, ?date_to=date) | ✅ | ✅ |
| GET | `/admin/system-health` | Get system health & performance metrics | ✅ | ✅ |

## Product Module (18 endpoints)

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
| DELETE | `/products/:id/images/:image-id` | Delete product image | ✅ | ✅ |
| GET | `/products/:id/analytics` | Get product analytics (?type=related\|history\|performance) | ✅ | ✅ |

### Product Variants
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/products/:id/variants` | Create product variant | ✅ | ✅ |
| GET | `/products/:id/variants` | Get product variants | ✅ | ✅ |
| PATCH | `/products/:id/variants/:variant-id` | Update variant (?action=update\|delete\|upload_image) | ✅ | ✅ |

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

## Order Module (7 endpoints)

### Orders
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/orders` | Create order | ✅ | ✅ |
| GET | `/orders` | List orders (?customer=current\|id, ?status=pending\|processing\|shipped, ?view=stats\|analytics, ?search=number, pagination) | ✅ | ✅ |
| GET | `/orders/:id` | Get order details (?include=items\|notes\|history\|invoice\|documents) | ✅ | ✅ |
| PATCH | `/orders/:id` | Update order (?action=cancel\|fulfill\|refund\|add_item\|remove_item\|add_note\|update_status) | ✅ | ✅ |
| DELETE | `/orders/:id` | Delete order | ✅ | ✅ |
| GET | `/orders/lookup/:number` | Get order by number (?public=true for customer access) | ❌/✅ | ✅ |
| GET | `/orders/:id/tracking` | Track order (?public=true for customer access) | ❌/✅ | ✅ |

## Payment Module (6 endpoints)

### Payment Processing
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/payments` | Create payment | ✅ | ✅ |
| GET | `/payments` | List payments (?status=pending\|completed\|failed, ?method=card\|bkash\|nagad, ?view=stats) | ✅ | ✅ |
| GET | `/payments/:id` | Get payment details | ✅ | ✅ |
| PATCH | `/payments/:id` | Update payment (?action=process\|refund\|capture\|void) | ✅ | ✅ |
| GET | `/payments/methods` | List payment methods (?active=true) | ✅ | ✅ |
| PATCH | `/payments/methods/:id` | Update payment method (?action=enable\|disable\|update) | ✅ | ✅ |

### Payment Webhooks
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/webhooks/payment/:provider` | Payment provider webhook (sslcommerz, bkash, nagad, stripe) | ❌ | ❌ |

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
| GET | `/shipping/track/:tracking-number` | Track shipment | ✅ | ✅ |
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


## Analytics Module (5 endpoints)

### Analytics
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/analytics/track` | Track events (?type=event\|page_view\|product_view\|purchase) | ❌ | ✅ |
| GET | `/analytics/dashboard` | Get dashboard analytics (?type=traffic\|sales\|realtime\|overview, ?period=day\|week\|month) | ✅ | ✅ |
| GET | `/analytics/insights` | Get insights (?type=top_products\|top_pages\|referrers\|cohorts\|funnel\|retention\|clv) | ✅ | ✅ |
| GET | `/analytics/reports` | List reports (?status=scheduled\|completed) | ✅ | ✅ |
| PATCH | `/analytics/reports/:id` | Manage reports (?action=generate\|schedule\|export\|delete) | ✅ | ✅ |

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
| GET | `/store-credit/:customer-id` | Get customer store credit | ✅ | ✅ |
| PATCH | `/store-credit/:customer-id` | Update store credit (?action=add\|use\|refund) | ✅ | ✅ |

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

## Settings Module (3 endpoints)

### Settings Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/settings` | Get all settings (?section=general\|seo\|appearance\|integrations) | ✅ | ✅ |
| PATCH | `/settings` | Update settings (logo, title, contact, SEO codes, theme, integrations) | ✅ | ✅ |
| GET | `/public/settings` | Get public settings (store name, logo, contact, theme) | ❌ | ✅ |

## Customer Module (3 endpoints)

*Note: Uses unified /auth endpoints and User table. No separate customer authentication needed.*

### Customer-specific Operations
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/customers/orders` | Get customer order history (?status=completed\|pending, ?limit=20) | ✅ | ✅ |

*Note: Customer profile uses `/users/profile` and addresses use `/addresses` endpoints*

### Public Shopping Features
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/public/pages/:slug` | Get public pages (about, terms, privacy, etc.) | ❌ | ✅ |
| GET | `/public/content/menus` | Get store navigation menus | ❌ | ✅ |

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
| GET | `/public/reviews/:product-id` | Get product reviews (public) | ❌ | ✅ |
| POST | `/public/reviews` | Submit product review (customer) | ✅ | ✅ |

## Support Module (8 endpoints)

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

## Contact Module (8 endpoints)

### Contact Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/contacts` | List contacts (?status=new\|replied, ?view=analytics\|export, ?assigned_to=user_id) | ✅ | ✅ |
| GET | `/contacts/:id` | Get contact details (?include=interactions\|notes\|replies) | ✅ | ✅ |
| PATCH | `/contacts/:id` | Update contact (?action=create\|reply\|assign\|close\|add_note\|delete) | ✅ | ✅ |
| POST | `/contacts/bulk` | Bulk operations (?operation=assign\|close\|delete) | ✅ | ✅ |
| GET | `/contact-forms` | List contact forms (?form_type=support\|general) | ❌/✅ | ✅ |
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

## Webhook Module (9 endpoints)

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

### External Webhooks
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/webhooks/payment/:provider` | Payment webhooks (stripe, paypal, bkash, nagad, rocket, upay, sslcommerz) | ❌ | ❌ |
| POST | `/webhooks/shipping/:provider` | Shipping webhooks (pathao, redx, paperfly, dhl, fedex) | ❌ | ❌ |

## Billing Module (13 endpoints)

### Billing Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/billing/plans` | List billing plans | ✅ | ❌ |
| PATCH | `/billing/plans/:plan-id` | Manage plan (?action=create\|update\|delete) | ✅ | ❌ |
| GET | `/billing/subscriptions` | Get subscription details | ✅ | ✅ |
| PATCH | `/billing/subscriptions` | Manage subscription (?action=create\|update\|cancel\|upgrade\|downgrade) | ✅ | ✅ |
| GET | `/billing/usage` | Get usage summary (?include=limits) | ✅ | ✅ |
| POST | `/billing/usage` | Record usage | ✅ | ✅ |
| GET | `/billing/invoices` | List invoices (?view=analytics) | ✅ | ✅ |
| GET | `/billing/invoices/:invoice-id` | Get invoice details | ✅ | ✅ |
| PATCH | `/billing/invoices/:invoice-id` | Process invoice (?action=payment\|refund) | ✅ | ✅ |
| GET | `/billing/reports` | Get billing reports (?type=revenue\|churn) | ✅ | ❌ |

### Admin Operations
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/billing/admin/process` | Process billing (?operation=billing\|retry_payments\|dunning) | ✅ | ❌ |
| PATCH | `/billing/admin/tenants/:tenant-id` | Update tenant service status | ✅ | ❌ |

## Tenant Module (4 endpoints)

### Tenant Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/tenants` | List tenants (?view=stats, ?include=users) | ✅ | ❌ |
| GET | `/tenants/:id` | Get tenant details (?subdomain=name for lookup) | ❌/✅ | ❌ |
| PATCH | `/tenants/:id` | Update tenant (?action=create\|update\|update_plan\|update_status) | ✅ | ❌ |
| GET | `/tenants/check-subdomain/:subdomain` | Check subdomain availability | ❌ | ❌ |

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
| GET | `/cart` | Get cart (?include=summary\|shipping_methods\|taxes) | ✅ | ✅ |
| POST | `/cart/items` | Add item to cart | ✅ | ✅ |
| PATCH | `/cart/items/:id` | Update cart item (?action=update_quantity\|save_later\|move_to_cart\|remove) | ✅ | ✅ |
| PATCH | `/cart` | Update cart (?action=apply_discount\|remove_discount\|clear\|merge\|validate) | ✅ | ✅ |
| POST | `/cart/estimates` | Get estimates (?type=shipping\|tax\|total) | ✅ | ✅ |

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

## Address Module (7 endpoints)

### Address Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/addresses` | List addresses (?customer_id=id, ?view=stats\|validation_trends) | ✅ | ✅ |
| POST | `/addresses` | Create address (?action=validate\|geocode) | ✅ | ✅ |
| GET | `/addresses/:id` | Get address details | ✅ | ✅ |
| PATCH | `/addresses/:id` | Update address (?action=set_default\|validate\|geocode) | ✅ | ✅ |
| DELETE | `/addresses/:id` | Delete address | ✅ | ✅ |
| POST | `/addresses/bulk` | Bulk operations (?operation=create\|update\|delete\|validate\|cleanup) | ✅ | ✅ |
| GET | `/addresses/suggestions` | Get address suggestions (?query=search_term) | ✅ | ✅ |

## Returns Module (6 endpoints)

### Return Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/returns` | Create return request | ✅ | ✅ |
| GET | `/returns` | List returns (?status=pending\|approved\|completed, ?view=stats) | ✅ | ✅ |
| GET | `/returns/:id` | Get return details (?include=label) | ✅ | ✅ |
| PATCH | `/returns/:id` | Update return (?action=approve\|reject\|process\|exchange\|complete) | ✅ | ✅ |
| GET | `/returns/reasons` | List return reasons | ✅ | ✅ |
| PATCH | `/returns/reasons/:id` | Manage return reason (?action=create\|update\|delete) | ✅ | ✅ |

## Loyalty Module (8 endpoints)

### Loyalty Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/loyalty/points/:customer-id` | Get customer points (?include=history) | ✅ | ✅ |
| PATCH | `/loyalty/points/:customer-id` | Update points (?action=earn\|redeem\|adjust) | ✅ | ✅ |
| GET | `/loyalty/tiers` | List loyalty tiers | ✅ | ✅ |
| PATCH | `/loyalty/tiers/:id` | Manage tier (?action=create\|update\|delete) | ✅ | ✅ |
| GET | `/loyalty/rewards` | List rewards catalog | ✅ | ✅ |
| PATCH | `/loyalty/rewards/:id` | Manage reward (?action=create\|update\|delete) | ✅ | ✅ |
| GET | `/loyalty/analytics` | Get loyalty analytics (?type=stats\|leaderboard\|tiers) | ✅ | ✅ |
| PATCH | `/loyalty/campaigns/:id` | Manage loyalty campaign (?action=create\|update\|delete) | ✅ | ✅ |

## Finance Module (6 endpoints)

### Financial Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/finance/ledger` | Get ledger (?type=receivable\|payable\|reconciliation) | ✅ | ✅ |
| GET | `/finance/transactions` | List transactions | ✅ | ✅ |
| POST | `/finance/transactions` | Create transaction | ✅ | ✅ |
| GET | `/finance/reports` | Get reports (?type=profit_loss\|balance_sheet\|cash_flow\|tax\|revenue\|expenses) | ✅ | ✅ |
| GET | `/finance/payouts` | List payouts | ✅ | ✅ |
| PATCH | `/finance/payouts/:id` | Manage payout (?action=create\|process\|update) | ✅ | ✅ |



## Tax Module (9 endpoints)

### Tax Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/tax/rules` | List tax rules (?status=active\|expired) | ✅ | ✅ |
| PATCH | `/tax/rules/:id` | Manage tax rule (?action=create\|update\|delete\|activate\|deactivate) | ✅ | ✅ |
| GET | `/tax/rates` | List tax rates (?location=location_id) | ✅ | ✅ |
| PATCH | `/tax/rates/:id` | Manage tax rate (?action=create\|update\|delete) | ✅ | ✅ |
| POST | `/tax/calculate` | Calculate tax (?preview=true, ?entity=order\|product\|customer, ?entity_id=id) | ✅ | ✅ |
| GET | `/tax/calculations` | Get tax calculations (?order_id=id, ?product_id=id, ?customer_id=id) | ✅ | ✅ |
| GET | `/tax/analytics` | Get tax analytics (?type=stats\|applicable_rules) | ✅ | ✅ |
| POST | `/tax/validate` | Validate tax (?type=location\|rule) | ✅ | ✅ |
| POST | `/tax/maintenance` | Tax maintenance (?operation=cleanup\|refresh) | ✅ | ✅ |

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

The e-commerce platform implements **200+ API endpoints** across **27 active modules** with token-based REST design:

- ✅ **User Module** - 17 endpoints (authentication, profile management, admin operations, bulk operations)
- ✅ **Admin Dashboard Module** - 8 endpoints (admin dashboard, staff management, roles, activity logs)
- ✅ **Product Module** - 18 endpoints (products, variants, categories, inventory, public access, analytics) 
- ✅ **Order Module** - 7 endpoints (order management, tracking, flexible operations)
- ✅ **Cart Module** - 8 endpoints (cart management, guest cart, checkout)
- ✅ **Wishlist Module** - 6 endpoints (wishlist management with bulk operations)
- ✅ **Address Module** - 7 endpoints (address management, validation, geocoding)
- ✅ **Payment Module** - 6 endpoints (payment processing, methods, webhooks)
- ✅ **Shipping Module** - 11 endpoints (zones, rates, labels, tracking, webhooks)
- ✅ **Notification Module** - 8 endpoints (notifications, templates, preferences)
- ✅ **Analytics Module** - 5 endpoints (tracking, dashboard, insights, reports)
- ✅ **Marketing Module** - 10 endpoints (campaigns, templates, segments, automation)
- ✅ **Discount Module** - 9 endpoints (discounts, gift cards, store credit)
- ✅ **Search Module** - 6 endpoints (global search, product search, suggestions, analytics, filters)
- ✅ **Settings Module** - 3 endpoints (store settings, SEO, appearance, integrations)
- ✅ **Customer Module** - 3 endpoints (customer orders, public pages, menus)
- ✅ **Reviews Module** - 10 endpoints (reviews, moderation, invitations, public reviews)
- ✅ **Support Module** - 8 endpoints (tickets, FAQ, knowledge base)
- ✅ **Contact Module** - 8 endpoints (contact management, forms, templates)
- ✅ **Content Management Module** - 8 endpoints (pages, posts, media, menus)
- ✅ **Webhook Module** - 9 endpoints (endpoint management, deliveries, events)
- ✅ **Billing Module** - 13 endpoints (plans, subscriptions, usage, invoices, admin)
- ✅ **Tenant Module** - 4 endpoints (multi-tenancy management)
- ✅ **Observability Module** - 8 endpoints (health, metrics, logs, alerts)
- ✅ **Returns Module** - 6 endpoints (return management, reasons, processing)
- ✅ **Loyalty Module** - 8 endpoints (points, tiers, rewards, analytics)
- ✅ **Finance Module** - 6 endpoints (ledger, transactions, reports, payouts)
- ✅ **Tax Module** - 9 endpoints (rules, rates, calculations, analytics)

### Key Optimizations
- **Token-based REST**: No session management, pure JWT token authentication
- **Unified Authentication**: Single `/auth` endpoints for all users (admin/customer)
- **Eliminated Duplicates**: Removed redundant endpoints, consolidated functionality
- **Unified Actions**: Single endpoints handle multiple operations via `?action` parameters
- **Flexible Queries**: Rich filtering and view options reduce need for separate endpoints  
- **Consistent Patterns**: PATCH for updates, GET with query params for filtering/views
- **Better Organization**: Logical grouping with consistent naming conventions

### Authentication & Security
- JWT-based authentication for all protected endpoints
- Multi-tenant architecture with proper isolation
- Public endpoints for customer-facing features
- WebSocket support for real-time updates

### Technology Stack
- **Backend Framework**: Gin (Go HTTP router)
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT tokens
- **Real-time**: WebSocket connections
- **Caching**: Redis integration
- **Architecture**: Hexagonal (Clean) Architecture pattern

---

*Note: All endpoints require appropriate authentication and authorization. Refer to the Authentication section for details on obtaining and using access tokens.*