# API Documentation

Comprehensive REST API specification for the e-commerce SaaS platform with **270+ implemented endpoints** covering all business operations across 27 active modules with multi-tenant architecture, authentication, and WebSocket real-time capabilities.

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
  - `:productId`, `:customerId` for related resource IDs
  - `:variantId`, `:addressId` for nested resources

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

## User Module (25 endpoints)

### User Authentication
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/auth/register` | User registration | ❌ | ❌ |
| POST | `/auth/login` | User login | ❌ | ❌ |
| POST | `/auth/refresh` | Refresh JWT token | ❌ | ❌ |
| POST | `/auth/logout` | User logout | ✅ | ❌ |
| POST | `/auth/forgot-password` | Request password reset | ❌ | ❌ |
| POST | `/auth/reset-password` | Reset password | ❌ | ❌ |
| POST | `/auth/verify-email` | Verify email address | ❌ | ❌ |
| POST | `/auth/resend-verification` | Resend verification email | ❌ | ❌ |

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

## Product Module (31 endpoints) - **CONSOLIDATED**

*Note: Stock management is handled within product endpoints for this single-vendor platform*

### Products
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/products` | Create a new product | ✅ | ✅ |
| GET | `/products` | List products with analytics support (?type=stats\|low-stock\|search, ?q=query, filtering, pagination) | ✅ | ✅ |
| POST | `/products/operations` | Bulk operations (?operation=bulk_update\|import\|export) | ✅ | ✅ |
| GET | `/products/:id` | Get a specific product | ✅ | ✅ |
| PUT | `/products/:id` | Update product with actions (?action=update_inventory\|update_status\|duplicate) | ✅ | ✅ |
| DELETE | `/products/:id` | Delete a product | ✅ | ✅ |
| GET | `/products/slug/:slug` | Get product by slug (storefront) | ✅ | ✅ |
| POST | `/products/:id/images` | Upload product images | ✅ | ✅ |
| DELETE | `/products/:id/images/:image-id` | Delete product image | ✅ | ✅ |
| GET | `/products/:id/related` | Get related products | ✅ | ✅ |
| GET | `/products/:id/inventory-history` | Get inventory movement history | ✅ | ✅ |
| POST | `/products/:id/variants/:variant-id/images` | Upload variant images | ✅ | ✅ |

### Product Variants
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/products/:id/variants` | Create product variant | ✅ | ✅ |
| GET | `/products/:id/variants` | Get product variants | ✅ | ✅ |
| PUT | `/products/:id/variants/:variant-id` | Update product variant | ✅ | ✅ |
| DELETE | `/products/:id/variants/:variant-id` | Delete product variant | ✅ | ✅ |

### Categories
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/categories` | Create a category | ✅ | ✅ |
| GET | `/categories` | List categories with hierarchy support (?hierarchy=root\|children, ?parent_id=id) | ✅ | ✅ |
| GET | `/categories/:id` | Get a specific category | ✅ | ✅ |
| PUT | `/categories/:id` | Update a category | ✅ | ✅ |
| DELETE | `/categories/:id` | Delete a category | ✅ | ✅ |
| POST | `/categories/:id/image` | Upload category image | ✅ | ✅ |
| DELETE | `/categories/:id/image` | Delete category image | ✅ | ✅ |

### Public Product Access (Storefront)
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/public/products` | Browse products (public) | ❌ | ✅ |
| GET | `/public/products/slug/:slug` | Get product by slug (public) | ❌ | ✅ |
| GET | `/public/products/:id` | Get product details (public) | ❌ | ✅ |
| GET | `/public/products/:id/variants` | Get product variants (public) | ❌ | ✅ |
| GET | `/public/categories` | Browse categories (public) | ❌ | ✅ |
| GET | `/public/categories/root` | Get root categories (public) | ❌ | ✅ |
| GET | `/public/categories/:id` | Get category details (public) | ❌ | ✅ |
| GET | `/public/categories/:id/children` | Get category children (public) | ❌ | ✅ |

### Product Inventory Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/products/:id/stock/adjust` | Adjust product stock levels | ✅ | ✅ |
| GET | `/products/:id/stock/movements` | Get stock movement history | ✅ | ✅ |
| POST | `/products/:id/stock/reserve` | Reserve stock for order | ✅ | ✅ |
| POST | `/products/:id/stock/release` | Release reserved stock | ✅ | ✅ |
| GET | `/products/low-stock` | Get low stock products | ✅ | ✅ |
| POST | `/products/stock/bulk-update` | Bulk update stock levels | ✅ | ✅ |

## Order Module (22 endpoints)

### Orders
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/orders` | Create a new order | ✅ | ✅ |
| GET | `/orders` | List orders (supports: ?customer=current, ?view=stats, filtering, pagination) | ✅ | ✅ |
| GET | `/orders/:id` | Get order details (supports ?type=history) | ✅ | ✅ |

| GET | `/orders/:id/invoice` | Get order invoice | ✅ | ✅ |
| GET | `/orders/number/:number` | Get order by number | ✅ | ✅ |
| POST | `/orders/:id/items` | Add item to existing order | ✅ | ✅ |
| PUT | `/orders/:id/items/:item-id` | Update order item | ✅ | ✅ |
| DELETE | `/orders/:id/items/:item-id` | Remove item from order | ✅ | ✅ |

| POST | `/orders/:id/notes` | Add internal order notes | ✅ | ✅ |
| GET | `/orders/:id/notes` | Get internal order notes | ✅ | ✅ |
| GET | `/orders/:id/documents` | Get order documents | ✅ | ✅ |
| PATCH | `/orders/:id` | Update order (status, shipping_status, details, actions: cancel/payment/refund/split/merge) | ✅ | ✅ |
| GET | `/public/orders/track/:number` | Track order by number (public) | ❌ | ✅ |
| GET | `/public/orders/number/:number` | Get order by number (public) | ❌ | ✅ |

## Payment Module (11 endpoints)

### Payment Processing
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/payments` | Create a payment | ✅ | ✅ |
| GET | `/payments` | List payments | ✅ | ✅ |
| GET | `/payments/:id` | Get payment details | ✅ | ✅ |
| PATCH | `/payments/:id` | Update payment (actions: process/refund/capture) | ✅ | ✅ |

| GET | `/payments/methods` | Get payment methods | ✅ | ✅ |
| POST | `/payments/methods` | Create payment method | ✅ | ✅ |
| PUT | `/payments/methods/:id` | Update payment method | ✅ | ✅ |
| DELETE | `/payments/methods/:id` | Delete payment method | ✅ | ✅ |

### Payment Webhooks
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/webhooks/sslcommerz` | SSLCommerz payment webhook | ❌ | ❌ |

## Shipping Module (25 endpoints)

### Shipping Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/shipping/zones` | Get shipping zones (supports ?type=stats) | ✅ | ✅ |
| POST | `/shipping/zones` | Create shipping zone | ✅ | ✅ |
| PUT | `/shipping/zones/:id` | Update shipping zone | ✅ | ✅ |
| DELETE | `/shipping/zones/:id` | Delete shipping zone | ✅ | ✅ |
| GET | `/shipping/rates` | Get shipping rates | ✅ | ✅ |
| POST | `/shipping/rates` | Create shipping rate | ✅ | ✅ |
| PUT | `/shipping/rates/:id` | Update shipping rate | ✅ | ✅ |
| DELETE | `/shipping/rates/:id` | Delete shipping rate | ✅ | ✅ |
| POST | `/shipping/calculate` | Calculate shipping cost | ✅ | ✅ |
| GET | `/shipping/labels` | Get shipping labels | ✅ | ✅ |
| POST | `/shipping/labels` | Create shipping label | ✅ | ✅ |
| GET | `/shipping/labels/:id` | Get shipping label | ✅ | ✅ |
| PUT | `/shipping/labels/:id` | Update shipping label | ✅ | ✅ |
| DELETE | `/shipping/labels/:id` | Delete shipping label | ✅ | ✅ |
| GET | `/shipping/track/:tracking-number` | Track shipment | ✅ | ✅ |
| GET | `/shipping/providers` | Get shipping providers | ✅ | ✅ |
| POST | `/shipping/providers` | Create shipping provider | ✅ | ✅ |
| PUT | `/shipping/providers/:id` | Update shipping provider | ✅ | ✅ |
| DELETE | `/shipping/providers/:id` | Delete shipping provider | ✅ | ✅ |

### Shipping Webhooks
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/shipping/webhooks/pathao` | Pathao shipping webhook | ❌ | ❌ |
| POST | `/shipping/webhooks/redx` | RedX shipping webhook | ❌ | ❌ |
| POST | `/shipping/webhooks/paperfly` | Paperfly shipping webhook | ❌ | ❌ |
| POST | `/shipping/webhooks/dhl` | DHL shipping webhook | ❌ | ❌ |
| POST | `/shipping/webhooks/fedex` | FedEx shipping webhook | ❌ | ❌ |

## Notification Module (13 endpoints)

### Notification Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/notifications` | Send notification | ✅ | ✅ |
| GET | `/notifications` | List notifications with stats support (?type=stats for notification statistics) | ✅ | ✅ |
| GET | `/notifications/:id` | Get notification details | ✅ | ✅ |
| PATCH | `/notifications/:id` | Update notification (mark as read, etc.) | ✅ | ✅ |

| POST | `/notifications/email` | Send email notification | ✅ | ✅ |
| POST | `/notifications/sms` | Send SMS notification | ✅ | ✅ |

### Notification Templates
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/notifications/templates` | Create notification template | ✅ | ✅ |
| GET | `/notifications/templates` | List notification templates | ✅ | ✅ |
| GET | `/notifications/templates/:id` | Get notification template | ✅ | ✅ |
| PUT | `/notifications/templates/:id` | Update notification template | ✅ | ✅ |

### User Preferences
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/notifications/preferences` | Get notification preferences | ✅ | ✅ |
| PUT | `/notifications/preferences` | Update notification preferences | ✅ | ✅ |


## Analytics Module (12 endpoints)

### Event Tracking
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/analytics/track` | Track events (supports ?type=event|page-view|product-view|purchase) | ❌ | ✅ |

### Dashboard Analytics
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/analytics/dashboard` | Get dashboard analytics (supports ?type=traffic|sales|realtime|overview) | ✅ | ✅ |
| GET | `/analytics/top` | Get top performers (supports ?type=products|pages|referrers) | ✅ | ✅ |

### Advanced Analytics
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/analytics/advanced` | Get advanced analytics (supports ?type=cohorts|funnel|retention|clv) | ✅ | ✅ |

### Reports
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/analytics/reports` | Generate or schedule reports (supports ?action=generate|schedule) | ✅ | ✅ |
| GET | `/analytics/reports` | Get scheduled reports | ✅ | ✅ |
| DELETE | `/analytics/reports/:id` | Delete scheduled report | ✅ | ✅ |
| POST | `/analytics/export` | Export analytics data (supports ?type=dashboard|advanced|reports) | ✅ | ✅ |

## Marketing Module (16 endpoints)

### Campaigns
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/marketing/campaigns` | Create campaign | ✅ | ✅ |
| GET | `/marketing/campaigns` | Get campaigns (supports ?type=stats|emails) | ✅ | ✅ |
| GET | `/marketing/campaigns/:id` | Get specific campaign (supports ?type=stats|emails) | ✅ | ✅ |
| PATCH | `/marketing/campaigns/:id` | Update campaign (status, details, schedule, etc.) | ✅ | ✅ |
| DELETE | `/marketing/campaigns/:id` | Delete campaign | ✅ | ✅ |

### Templates & Segments
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/marketing/templates` | Create template | ✅ | ✅ |
| GET | `/marketing/templates` | Get templates | ✅ | ✅ |
| PATCH | `/marketing/templates/:id` | Update template (content, etc.) | ✅ | ✅ |
| DELETE | `/marketing/templates/:id` | Delete template | ✅ | ✅ |
| POST | `/marketing/segments` | Create customer segment | ✅ | ✅ |
| GET | `/marketing/segments` | Get customer segments | ✅ | ✅ |
| PATCH | `/marketing/segments/:id` | Update segment (details, refresh, etc.) | ✅ | ✅ |
| DELETE | `/marketing/segments/:id` | Delete segment | ✅ | ✅ |

### Newsletter & Automation
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/marketing/newsletter` | Subscribe/unsubscribe newsletter (supports ?action=subscribe|unsubscribe) | ❌ | ✅ |
| GET | `/marketing/newsletter/subscribers` | Get newsletter subscribers | ✅ | ✅ |
| GET | `/marketing/abandoned-carts` | Get abandoned carts for marketing automation | ✅ | ✅ |

## Discount Module (19 endpoints)

### Discount Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/discounts` | Create discount/coupon | ✅ | ✅ |
| GET | `/discounts` | Get discounts with analytics support (?type=stats\|performance\|revenue-impact, ?discount_type=percentage\|fixed, filtering, pagination) | ✅ | ✅ |
| GET | `/discounts/:id` | Get specific discount | ✅ | ✅ |
| PUT | `/discounts/:id` | Update discount | ✅ | ✅ |
| DELETE | `/discounts/:id` | Delete discount | ✅ | ✅ |
| GET | `/discounts/:id/usage` | Get discount usage | ✅ | ✅ |

### Discount Application (Public)
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/validate-discount` | Validate discount code | ❌ | ✅ |
| POST | `/apply-discount` | Apply discount to order | ❌ | ✅ |
| DELETE | `/remove-discount/:order-id` | Remove discount from order | ❌ | ✅ |

### Gift Cards
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/gift-cards` | Create gift card | ✅ | ✅ |
| GET | `/gift-cards` | Get gift cards | ✅ | ✅ |
| GET | `/gift-cards/:code` | Get gift card by code | ✅ | ✅ |
| PUT | `/gift-cards/:id` | Update gift card | ✅ | ✅ |
| DELETE | `/gift-cards/:id` | Delete gift card | ✅ | ✅ |
| POST | `/validate-gift-card` | Validate gift card | ❌ | ✅ |
| POST | `/use-gift-card` | Use gift card | ❌ | ✅ |

### Store Credit
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/store-credit/:customer-id` | Get customer store credit | ✅ | ✅ |
| PATCH | `/store-credit/:customer-id` | Update store credit (add/use) | ✅ | ✅ |


## Reviews Module (19 endpoints)

### Review Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/reviews` | Create review | ✅ | ✅ |
| GET | `/reviews` | Get reviews | ✅ | ✅ |
| GET | `/reviews/:id` | Get specific review | ✅ | ✅ |
| PATCH | `/reviews/:id` | Update review (status, content, mark as spam, etc.) | ✅ | ✅ |
| DELETE | `/reviews/:id` | Delete review | ✅ | ✅ |

### Review Moderation
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|

| POST | `/reviews/bulk-moderate` | Bulk moderate reviews | ✅ | ✅ |

### Review Interactions
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/reviews/:id/replies` | Add review reply | ✅ | ✅ |
| GET | `/reviews/:id/replies` | Get review replies | ✅ | ✅ |
| PATCH | `/reviews/:id` | Update review (react, etc.) | ✅ | ✅ |

### Product Reviews
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/products/:product-id/reviews` | Get product reviews | ✅ | ✅ |
| GET | `/products/:product-id/reviews/summary` | Get product review summary | ✅ | ✅ |

### Review Invitations & Analytics
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/review-invitations` | Create review invitation | ✅ | ✅ |
| GET | `/review-invitations` | Get review invitations | ✅ | ✅ |
| PATCH | `/review-invitations/:id` | Update invitation (send, etc.) | ✅ | ✅ |
| GET | `/review-invitations/:token` | Process invitation click | ❌ | ❌ |
| GET | `/reviews` | List reviews (supports: ?status=pending, ?view=stats, ?view=trends, ?view=top_products, ?sort=recent) | ✅ | ✅ |
| GET | `/reviews/settings` | Get review settings | ✅ | ✅ |
| PUT | `/reviews/settings` | Update review settings | ✅ | ✅ |

## Support Module (15 endpoints)

### Support Tickets
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/support/tickets` | Create support ticket | ✅ | ✅ |
| GET | `/support/tickets` | Get support tickets (supports ?type=stats) | ✅ | ✅ |
| GET | `/support/tickets/:id` | Get specific ticket | ✅ | ✅ |
| PATCH | `/support/tickets/:id` | Update ticket (assign, resolve, status, etc.) | ✅ | ✅ |
| DELETE | `/support/tickets/:id` | Delete ticket | ✅ | ✅ |

| GET | `/support/tickets/:id/messages` | Get ticket messages | ✅ | ✅ |
| POST | `/support/tickets/:id/messages` | Add ticket message | ✅ | ✅ |

### FAQ & Knowledge Base
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/support/faqs` | Create FAQ | ✅ | ✅ |
| GET | `/support/faqs` | Get FAQs | ✅ | ✅ |
| GET | `/support/knowledge-base` | Get knowledge base articles | ✅ | ✅ |
| GET | `/support/knowledge-base/:slug` | Get article by slug | ✅ | ✅ |

### Support Settings & Analytics
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/support/settings` | Get support settings | ✅ | ✅ |

## Contact Module (14 endpoints)

### Contact Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/contacts` | Create contact | ✅ | ✅ |
| GET | `/contacts` | List contacts (supports ?type=analytics|export) | ✅ | ✅ |
| GET | `/contacts/:id` | Get contact details (supports ?type=replies|notes) | ✅ | ✅ |
| PATCH | `/contacts/:id` | Update contact (status, assignment, priority, details, tags, etc.) | ✅ | ✅ |
| DELETE | `/contacts/:id` | Delete contact | ✅ | ✅ |
| POST | `/contacts/bulk` | Bulk operations on contacts | ✅ | ✅ |

### Contact Interactions
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/contacts/:id/interactions` | Create contact interaction (supports ?type=reply|note) | ✅ | ✅ |
| GET | `/contacts/:id/interactions` | List contact interactions (supports ?type=replies|notes) | ✅ | ✅ |

### Contact Forms & Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/contact-forms` | Create contact form | ✅ | ✅ |
| GET | `/contact-forms` | List contact forms | ✅ | ✅ |
| GET | `/contact-forms/public/:form_type` | Get public contact form | ❌ | ✅ |
| POST | `/contact-forms/public/:form_type/submit` | Submit public contact form | ❌ | ✅ |
| GET | `/contact-templates` | List contact templates | ✅ | ✅ |
| PATCH | `/contact-settings` | Update contact settings (templates, etc.) | ✅ | ✅ |

## Content Management Module (13 endpoints)

### Pages & Posts
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/content/pages` | List pages | ✅ | ✅ |
| POST | `/content/pages` | Create page | ✅ | ✅ |
| GET | `/content/pages/:id` | Get page | ✅ | ✅ |
| PUT | `/content/pages/:id` | Update page | ✅ | ✅ |
| DELETE | `/content/pages/:id` | Delete page | ✅ | ✅ |
| PATCH | `/content/pages/:id/status` | Update page status (publish/unpublish) | ✅ | ✅ |
| GET | `/content/posts` | List posts | ✅ | ✅ |
| POST | `/content/posts` | Create post | ✅ | ✅ |
| GET | `/content/posts/:id` | Get post | ✅ | ✅ |
| PUT | `/content/posts/:id` | Update post | ✅ | ✅ |
| DELETE | `/content/posts/:id` | Delete post | ✅ | ✅ |

### Media & Menus
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/content/media` | List media | ✅ | ✅ |
| POST | `/content/media` | Upload media | ✅ | ✅ |
| GET | `/content/menus` | List menus | ✅ | ✅ |
| POST | `/content/menus` | Create menu | ✅ | ✅ |

## Webhook Module (25 endpoints)

### Webhook Endpoints
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/webhooks/endpoints` | List webhook endpoints | ✅ | ✅ |
| POST | `/webhooks/endpoints` | Create webhook endpoint | ✅ | ✅ |
| GET | `/webhooks/endpoints/:id` | Get webhook endpoint | ✅ | ✅ |
| PUT | `/webhooks/endpoints/:id` | Update webhook endpoint | ✅ | ✅ |
| DELETE | `/webhooks/endpoints/:id` | Delete webhook endpoint | ✅ | ✅ |
| PATCH | `/webhooks/endpoints/:id` | Update endpoint (test, etc.) | ✅ | ✅ |

### Webhook Deliveries
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/webhooks/deliveries` | List webhook deliveries | ✅ | ✅ |
| GET | `/webhooks/deliveries/:id` | Get webhook delivery | ✅ | ✅ |
| PATCH | `/webhooks/deliveries/:id` | Update delivery (retry, etc.) | ✅ | ✅ |

### Webhook Events
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/webhooks/events` | List webhook events | ✅ | ✅ |
| POST | `/webhooks/events` | Create webhook event | ✅ | ✅ |
| GET | `/webhooks/events/:id` | Get webhook event | ✅ | ✅ |
| PUT | `/webhooks/events/:id` | Update webhook event | ✅ | ✅ |
| DELETE | `/webhooks/events/:id` | Delete webhook event | ✅ | ✅ |

### Payment Provider Webhooks
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/webhooks/stripe` | Stripe payment webhook | ❌ | ❌ |
| POST | `/webhooks/paypal` | PayPal payment webhook | ❌ | ❌ |
| POST | `/webhooks/bkash` | bKash payment webhook | ❌ | ❌ |
| POST | `/webhooks/nagad` | Nagad payment webhook | ❌ | ❌ |
| POST | `/webhooks/rocket` | Rocket payment webhook | ❌ | ❌ |
| POST | `/webhooks/upay` | Upay payment webhook | ❌ | ❌ |

## Billing Module (30 endpoints)

### Billing Plans
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/billing/plans` | Get billing plans | ✅ | ❌ |
| GET | `/billing/plans/:planId` | Get specific billing plan | ✅ | ❌ |
| POST | `/billing/plans` | Create billing plan | ✅ | ❌ |
| PUT | `/billing/plans/:planId` | Update billing plan | ✅ | ❌ |
| DELETE | `/billing/plans/:planId` | Delete billing plan | ✅ | ❌ |

### Subscriptions
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/billing/subscriptions` | Create subscription | ✅ | ✅ |
| GET | `/billing/subscriptions` | Get subscription | ✅ | ✅ |
| PUT | `/billing/subscriptions` | Update subscription | ✅ | ✅ |
| DELETE | `/billing/subscriptions` | Cancel subscription | ✅ | ✅ |
| POST | `/billing/subscriptions/upgrade` | Upgrade plan | ✅ | ✅ |
| POST | `/billing/subscriptions/downgrade` | Downgrade plan | ✅ | ✅ |

### Usage Tracking
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/billing/usage` | Record usage | ✅ | ✅ |
| GET | `/billing/usage` | Get usage summary | ✅ | ✅ |
| GET | `/billing/usage/limits` | Check usage limits | ✅ | ✅ |

### Invoices
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/billing/invoices` | Get invoices (supports ?type=analytics) | ✅ | ✅ |
| GET | `/billing/invoices/:invoiceId` | Get specific invoice | ✅ | ✅ |
| POST | `/billing/invoices/:invoiceId/payment` | Process invoice payment | ✅ | ✅ |
| POST | `/billing/invoices/:invoiceId/refund` | Refund invoice payment | ✅ | ✅ |

### Reports
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/billing/reports/revenue` | Get revenue report | ✅ | ❌ |
| GET | `/billing/reports/churn` | Get churn analysis | ✅ | ❌ |

### Admin Operations
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/billing/admin/process-billing` | Process recurring billing | ✅ | ❌ |
| POST | `/billing/admin/retry-payments` | Retry failed payments | ✅ | ❌ |
| POST | `/billing/admin/process-dunning` | Process dunning | ✅ | ❌ |
| PATCH | `/billing/admin/tenants/:tenant-id/status` | Update tenant service status | ✅ | ❌ |

## Tenant Module (8 endpoints)

### Tenant Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/tenants` | Create a new tenant | ✅ | ❌ |
| GET | `/tenants` | List all tenants (supports ?type=stats) | ✅ | ❌ |
| GET | `/tenants/:id` | Get tenant details | ✅ | ❌ |
| PUT | `/tenants/:id` | Update tenant | ✅ | ❌ |
| PUT | `/tenants/:id/plan` | Update tenant plan | ✅ | ❌ |
| PATCH | `/tenants/:id/status` | Update tenant status | ✅ | ❌ |

| GET | `/tenants/subdomain/:subdomain` | Get tenant by subdomain | ❌ | ❌ |
| GET | `/tenants/check-subdomain/:subdomain` | Check subdomain availability | ❌ | ❌ |
| GET | `/tenants/:id/users` | Get tenant users | ✅ | ❌ |

## Observability Module (12 endpoints)

### Health Monitoring
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/observability/health` | Get basic health status | ✅ | ❌ |
| GET | `/observability/health/detailed` | Get detailed health status | ✅ | ❌ |

### Metrics & Monitoring
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/observability/metrics` | Get system metrics | ✅ | ❌ |
| GET | `/observability/metrics/summary` | Get metrics summary | ✅ | ❌ |

### Logging & Tracing
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/observability/logs` | Get log entries | ✅ | ❌ |
| POST | `/observability/logs` | Create log entry | ✅ | ❌ |
| GET | `/observability/traces` | Get traces | ✅ | ❌ |
| GET | `/observability/traces/:traceId` | Get specific trace | ✅ | ❌ |

### Alerting & System Info
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/observability/alerts` | Get alerts | ✅ | ❌ |
| POST | `/observability/alerts` | Create alert | ✅ | ❌ |
| GET | `/observability/system/info` | Get system information (supports ?type=stats) | ✅ | ❌ |

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

## Category Module (17 endpoints)

### Category Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/categories` | List all categories (supports ?type=stats) | ❌ | ✅ |
| POST | `/categories` | Create category | ✅ | ✅ |
| GET | `/categories/:id` | Get category details | ❌ | ✅ |
| PUT | `/categories/:id` | Update category | ✅ | ✅ |
| DELETE | `/categories/:id` | Delete category | ✅ | ✅ |
| GET | `/categories/:id/children` | Get child categories | ❌ | ✅ |
| GET | `/categories/:id/products` | Get category products | ❌ | ✅ |
| PATCH | `/categories/:id` | Update category (move, reorder, etc.) | ✅ | ✅ |

### Category Operations
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/categories/bulk` | Bulk create categories | ✅ | ✅ |
| PUT | `/categories/bulk` | Bulk update categories | ✅ | ✅ |
| DELETE | `/categories/bulk` | Bulk delete categories | ✅ | ✅ |
| GET | `/categories/tree` | Get category tree | ❌ | ✅ |
| POST | `/categories/cleanup` | Cleanup empty categories | ✅ | ✅ |

## Cart Module (15 endpoints)

### Cart Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/cart` | Get user cart | ✅ | ✅ |
| POST | `/cart/items` | Add item to cart | ✅ | ✅ |
| PUT | `/cart/items/:id` | Update cart item | ✅ | ✅ |
| DELETE | `/cart/items/:id` | Remove cart item | ✅ | ✅ |
| DELETE | `/cart/clear` | Clear cart | ✅ | ✅ |
| POST | `/cart/merge` | Merge guest cart | ✅ | ✅ |

### Cart Operations
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| PATCH | `/cart` | Update cart (apply/remove discount, validate, etc.) | ✅ | ✅ |
| GET | `/cart/summary` | Get cart summary | ✅ | ✅ |
| PATCH | `/cart/items/:id` | Update cart item (save for later, move to cart, etc.) | ✅ | ✅ |
| GET | `/cart/shipping-methods` | Get available shipping methods | ✅ | ✅ |
| POST | `/cart/estimate-taxes` | Estimate taxes for cart | ✅ | ✅ |
| POST | `/cart/cleanup` | Cleanup abandoned carts | ✅ | ✅ |

## Wishlist Module (15 endpoints)

### Wishlist Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/wishlists` | Get user wishlists (supports ?type=analytics) | ✅ | ✅ |
| POST | `/wishlists` | Create wishlist | ✅ | ✅ |
| GET | `/wishlists/:id` | Get wishlist details | ✅ | ✅ |
| PUT | `/wishlists/:id` | Update wishlist | ✅ | ✅ |
| DELETE | `/wishlists/:id` | Delete wishlist | ✅ | ✅ |
| POST | `/wishlists/:id/items` | Add item to wishlist | ✅ | ✅ |
| DELETE | `/wishlists/:id/items/:itemId` | Remove wishlist item | ✅ | ✅ |

### Wishlist Operations
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| PATCH | `/wishlists/:id` | Update wishlist (share, move to cart, etc.) | ✅ | ✅ |
| GET | `/wishlists/shared/:token` | View shared wishlist | ❌ | ✅ |
| POST | `/wishlists/bulk/items` | Bulk add items | ✅ | ✅ |
| DELETE | `/wishlists/bulk/items` | Bulk remove items | ✅ | ✅ |
| POST | `/wishlists/cleanup` | Cleanup old wishlists | ✅ | ✅ |
| GET | `/wishlists/popular-items` | Get popular wishlist items | ✅ | ✅ |

## Address Module (18 endpoints)

### Address Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/addresses` | Get user addresses (supports ?type=stats) | ✅ | ✅ |
| POST | `/addresses` | Create address | ✅ | ✅ |
| GET | `/addresses/:id` | Get address details | ✅ | ✅ |
| PUT | `/addresses/:id` | Update address | ✅ | ✅ |
| DELETE | `/addresses/:id` | Delete address | ✅ | ✅ |
| PATCH | `/addresses/:id` | Update address (set default, validate, etc.) | ✅ | ✅ |

### Address Operations
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/addresses/customers/:customerId` | Get customer addresses | ✅ | ✅ |
| POST | `/addresses/bulk` | Bulk create addresses | ✅ | ✅ |
| PUT | `/addresses/bulk` | Bulk update addresses | ✅ | ✅ |
| DELETE | `/addresses/bulk` | Bulk delete addresses | ✅ | ✅ |
| POST | `/addresses/validate-bulk` | Bulk validate addresses | ✅ | ✅ |
| GET | `/addresses/validation-stats` | Get validation statistics | ✅ | ✅ |
| POST | `/addresses/cleanup` | Cleanup unvalidated addresses | ✅ | ✅ |
| GET | `/addresses/validation-trends` | Get validation trends | ✅ | ✅ |
| POST | `/addresses/geocode` | Geocode address | ✅ | ✅ |
| GET | `/addresses/suggestions` | Get address suggestions | ✅ | ✅ |



## Returns Module (10 endpoints)

### Return Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/returns` | Create return request | ✅ | ✅ |
| GET | `/returns` | List returns | ✅ | ✅ |
| GET | `/returns/:id` | Get return details | ✅ | ✅ |
| PUT | `/returns/:id` | Update return status | ✅ | ✅ |
| PATCH | `/returns/:id` | Update return (status, process, exchange, etc.) | ✅ | ✅ |
| GET | `/returns/reasons` | Get return reasons | ✅ | ✅ |
| POST | `/returns/reasons` | Create return reason | ✅ | ✅ |
| GET | `/returns/stats` | Get return statistics | ✅ | ✅ |
| GET | `/returns/:id/label` | Get return shipping label | ✅ | ✅ |

## Loyalty Module (15 endpoints)

### Points Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/loyalty/points/:customer-id` | Get customer points balance (supports ?type=history) | ✅ | ✅ |
| PATCH | `/loyalty/points/:customer-id` | Update points (earn, redeem, etc.) | ✅ | ✅ |


### Loyalty Program
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/loyalty/tiers` | Get loyalty tiers | ✅ | ✅ |
| POST | `/loyalty/tiers` | Create loyalty tier | ✅ | ✅ |
| PUT | `/loyalty/tiers/:id` | Update loyalty tier | ✅ | ✅ |
| GET | `/loyalty/rewards` | Get reward catalog | ✅ | ✅ |
| POST | `/loyalty/rewards` | Create reward | ✅ | ✅ |
| PUT | `/loyalty/rewards/:id` | Update reward | ✅ | ✅ |
| GET | `/loyalty/customers/:customer-id/tier` | Get customer tier | ✅ | ✅ |
| PUT | `/loyalty/customers/:customer-id/tier` | Update customer tier | ✅ | ✅ |
| GET | `/loyalty/stats` | Get loyalty program statistics | ✅ | ✅ |
| GET | `/loyalty/leaderboard` | Get points leaderboard | ✅ | ✅ |
| POST | `/loyalty/campaigns` | Create loyalty campaign | ✅ | ✅ |

## Finance Module (11 endpoints)

### Accounting & Transactions
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/finance/ledger` | Get general ledger (supports ?type=receivable|payable) | ✅ | ✅ |
| POST | `/finance/transactions` | Create financial transaction | ✅ | ✅ |
| GET | `/finance/transactions` | List financial transactions | ✅ | ✅ |
| GET | `/finance/reconciliation` | Get account reconciliation | ✅ | ✅ |

### Financial Reports
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/finance/reports` | Get financial reports (supports ?type=profit-loss|balance-sheet|cash-flow|tax|revenue|expenses) | ✅ | ✅ |

### Payouts
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/finance/payouts` | Create payout | ✅ | ✅ |
| GET | `/finance/payouts` | List payouts | ✅ | ✅ |
| GET | `/finance/payouts/:id` | Get payout details | ✅ | ✅ |
| PATCH | `/finance/payouts/:id` | Update payout (process, etc.) | ✅ | ✅ |



## Tax Module (20 endpoints)

### Tax Rules Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/tax/rules` | List tax rules | ✅ | ✅ |
| POST | `/tax/rules` | Create tax rule | ✅ | ✅ |
| GET | `/tax/rules/:id` | Get tax rule | ✅ | ✅ |
| PUT | `/tax/rules/:id` | Update tax rule | ✅ | ✅ |
| DELETE | `/tax/rules/:id` | Delete tax rule | ✅ | ✅ |
| PATCH | `/tax/rules/:id` | Update tax rule (status, etc.) | ✅ | ✅ |

### Tax Rates Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/tax/rates` | List tax rates | ✅ | ✅ |
| POST | `/tax/rates` | Create tax rate | ✅ | ✅ |
| GET | `/tax/rates/:id` | Get tax rate | ✅ | ✅ |
| PUT | `/tax/rates/:id` | Update tax rate | ✅ | ✅ |
| DELETE | `/tax/rates/:id` | Delete tax rate | ✅ | ✅ |

### Tax Calculation
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/tax/calculate` | Calculate tax | ✅ | ✅ |
| POST | `/tax/preview` | Preview tax calculation | ✅ | ✅ |
| GET | `/tax/calculations/order/:orderId` | Get order tax calculations | ✅ | ✅ |
| GET | `/tax/calculations/product/:productId` | Get product tax calculations | ✅ | ✅ |
| GET | `/tax/calculations/customer/:customerId` | Get customer tax calculations | ✅ | ✅ |

### Tax Operations
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/tax/stats` | Get tax statistics | ✅ | ✅ |
| GET | `/tax/applicable-rules` | Get applicable tax rules | ✅ | ✅ |
| POST | `/tax/validate-location` | Validate tax location | ✅ | ✅ |
| POST | `/tax/cleanup` | Cleanup expired rules | ✅ | ✅ |

## Summary

The e-commerce platform currently implements **220+ API endpoints** across **25 active modules**:

- ✅ **User Module** - 25 endpoints (user authentication, password management, profile management, preferences, admin operations, bulk operations)
- ✅ **Product Module** - 31 endpoints (products, variants, categories, stock management, public access, image management, bulk operations)
- ✅ **Category Module** - 15 endpoints (hierarchical category management, image management)
- ✅ **Order Module** - 19 endpoints (order management, tracking, payments, items management, notes, order operations)
- ✅ **Cart Module** - 12 endpoints (shopping cart management, operations, validation, tax estimation)
- ✅ **Returns Module** - 10 endpoints (return management, approval workflow, exchange processing)
- ✅ **Loyalty Module** - 15 endpoints (points management, tiers, rewards, campaigns)
- ✅ **Finance Module** - 11 endpoints (accounting, financial reports, payouts, reconciliation)
- ✅ **Wishlist Module** - 15 endpoints (wishlist management, sharing, analytics)
- ✅ **Address Module** - 18 endpoints (address management, validation, geocoding)
- ✅ **Tax Module** - 20 endpoints (tax rules, rates, calculations)
- ✅ **Payment Module** - 9 endpoints (payment processing, methods, webhooks)
- ✅ **Shipping Module** - 25 endpoints (zones, rates, labels, tracking, webhooks)
- ✅ **Notification Module** - 12 endpoints (notifications, templates, preferences)
- ✅ **Analytics Module** - 12 endpoints (tracking, dashboard, reports)
- ✅ **Marketing Module** - 16 endpoints (campaigns, templates, segments)
- ✅ **Discount Module** - 19 endpoints (discounts, gift cards, store credit)
- ✅ **Reviews Module** - 16 endpoints (reviews, moderation, invitations)
- ✅ **Support Module** - 13 endpoints (tickets, FAQ, knowledge base)
- ✅ **Contact Module** - 14 endpoints (contact management, forms, templates)
- ✅ **Content Management Module** - 13 endpoints (pages, posts, media, menus)
- ✅ **Webhook Module** - 23 endpoints (endpoint management, deliveries, events)
- ✅ **Billing Module** - 28 endpoints (plans, subscriptions, usage, invoices)
- ✅ **Tenant Module** - 8 endpoints (multi-tenancy management)
- ✅ **Observability Module** - 11 endpoints (health, metrics, logs, alerts)

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