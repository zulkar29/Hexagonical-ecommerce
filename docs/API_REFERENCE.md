# API Documentation

Comprehensive REST API specification for the e-commerce SaaS platform with **195+ implemented endpoints** covering all business operations across 13 active modules with multi-tenant architecture, authentication, and WebSocket real-time capabilities.

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

## Product Module (29 endpoints)

### Products
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/products` | Create a new product | ✅ | ✅ |
| GET | `/products` | List products with filtering and pagination | ✅ | ✅ |
| GET | `/products/search` | Search products by query | ✅ | ✅ |
| GET | `/products/stats` | Get product statistics | ✅ | ✅ |
| GET | `/products/low-stock` | Get low stock products | ✅ | ✅ |
| PATCH | `/products/bulk` | Bulk update multiple products | ✅ | ✅ |
| GET | `/products/:id` | Get a specific product | ✅ | ✅ |
| PUT | `/products/:id` | Update a product | ✅ | ✅ |
| DELETE | `/products/:id` | Delete a product | ✅ | ✅ |
| PATCH | `/products/:id/status` | Update product status | ✅ | ✅ |
| PATCH | `/products/:id/inventory` | Update product inventory | ✅ | ✅ |
| POST | `/products/:id/duplicate` | Duplicate a product | ✅ | ✅ |
| GET | `/products/slug/:slug` | Get product by slug (storefront) | ✅ | ✅ |

### Product Variants
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/products/:id/variants` | Create product variant | ✅ | ✅ |
| GET | `/products/:id/variants` | Get product variants | ✅ | ✅ |
| PUT | `/products/:id/variants/:variantId` | Update product variant | ✅ | ✅ |
| DELETE | `/products/:id/variants/:variantId` | Delete product variant | ✅ | ✅ |

### Categories
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/categories` | Create a category | ✅ | ✅ |
| GET | `/categories` | List all categories | ✅ | ✅ |
| GET | `/categories/root` | Get root categories | ✅ | ✅ |
| GET | `/categories/:id` | Get a specific category | ✅ | ✅ |
| PUT | `/categories/:id` | Update a category | ✅ | ✅ |
| DELETE | `/categories/:id` | Delete a category | ✅ | ✅ |
| GET | `/categories/:id/children` | Get category children | ✅ | ✅ |

### Public Product Access (Storefront)
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/public/products` | Browse products (public) | ❌ | ✅ |
| GET | `/public/products/search` | Search products (public) | ❌ | ✅ |
| GET | `/public/products/slug/:slug` | Get product by slug (public) | ❌ | ✅ |
| GET | `/public/products/:id` | Get product details (public) | ❌ | ✅ |
| GET | `/public/products/:id/variants` | Get product variants (public) | ❌ | ✅ |
| GET | `/public/categories` | Browse categories (public) | ❌ | ✅ |
| GET | `/public/categories/root` | Get root categories (public) | ❌ | ✅ |
| GET | `/public/categories/:id` | Get category details (public) | ❌ | ✅ |
| GET | `/public/categories/:id/children` | Get category children (public) | ❌ | ✅ |

## Order Module (15 endpoints)

### Orders
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/orders` | Create a new order | ✅ | ✅ |
| GET | `/orders` | List orders with filtering and pagination | ✅ | ✅ |
| GET | `/orders/stats` | Get order statistics | ✅ | ✅ |
| GET | `/orders/my-orders` | Get customer orders | ✅ | ✅ |
| GET | `/orders/:id` | Get order details | ✅ | ✅ |
| PATCH | `/orders/:id/status` | Update order status | ✅ | ✅ |
| POST | `/orders/:id/cancel` | Cancel an order | ✅ | ✅ |
| POST | `/orders/:id/payment` | Process order payment | ✅ | ✅ |
| POST | `/orders/:id/refund` | Refund an order | ✅ | ✅ |
| GET | `/orders/:id/invoice` | Get order invoice | ✅ | ✅ |
| GET | `/orders/number/:number` | Get order by number | ✅ | ✅ |
| GET | `/public/orders/track/:number` | Track order by number (public) | ❌ | ✅ |
| GET | `/public/orders/number/:number` | Get order by number (public) | ❌ | ✅ |

## User Module (12 endpoints)

### Authentication
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/auth/register` | User registration | ❌ | ❌ |
| POST | `/auth/login` | User login | ❌ | ❌ |
| POST | `/auth/refresh` | Refresh JWT token | ❌ | ❌ |
| POST | `/auth/logout` | User logout | ✅ | ❌ |
| POST | `/auth/forgot-password` | Request password reset | ❌ | ❌ |
| POST | `/auth/reset-password` | Reset password | ❌ | ❌ |
| POST | `/auth/verify-email` | Verify email address | ❌ | ❌ |

### User Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/users/profile` | Get user profile | ✅ | ❌ |
| PUT | `/users/profile` | Update user profile | ✅ | ❌ |
| POST | `/users/change-password` | Change user password | ✅ | ❌ |
| GET | `/users` | List users (admin) | ✅ | ❌ |
| GET | `/users/:id` | Get user by ID (admin) | ✅ | ❌ |

## Tenant Module (11 endpoints)

### Tenant Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/tenants` | Create a new tenant | ✅ | ❌ |
| GET | `/tenants` | List all tenants | ✅ | ❌ |
| GET | `/tenants/:id` | Get tenant details | ✅ | ❌ |
| PUT | `/tenants/:id` | Update tenant | ✅ | ❌ |
| PUT | `/tenants/:id/plan` | Update tenant plan | ✅ | ❌ |
| POST | `/tenants/:id/activate` | Activate tenant | ✅ | ❌ |
| POST | `/tenants/:id/deactivate` | Deactivate tenant | ✅ | ❌ |
| GET | `/tenants/:id/stats` | Get tenant statistics | ✅ | ❌ |
| GET | `/tenants/subdomain/:subdomain` | Get tenant by subdomain | ❌ | ❌ |
| GET | `/tenants/check-subdomain/:subdomain` | Check subdomain availability | ❌ | ❌ |

## Payment Module (6 endpoints)

### Payment Processing
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/payments` | Create a payment | ✅ | ✅ |
| GET | `/payments` | List payments | ✅ | ✅ |
| GET | `/payments/:id` | Get payment details | ✅ | ✅ |
| POST | `/payments/:id/process` | Process a payment | ✅ | ✅ |
| POST | `/payments/:id/refund` | Refund a payment | ✅ | ✅ |

### Payment Webhooks
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/webhooks/sslcommerz` | SSLCommerz payment webhook | ❌ | ❌ |

## Notification Module (14 endpoints)

### Notification Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/notifications` | Send notification | ✅ | ✅ |
| GET | `/notifications` | List notifications | ✅ | ✅ |
| GET | `/notifications/:id` | Get notification details | ✅ | ✅ |
| PUT | `/notifications/:id/read` | Mark notification as read | ✅ | ✅ |
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
| GET | `/notifications/stats` | Get notification statistics | ✅ | ✅ |

## Billing Module (31 endpoints)

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
| GET | `/billing/invoices` | Get invoices | ✅ | ✅ |
| GET | `/billing/invoices/:invoiceId` | Get specific invoice | ✅ | ✅ |
| POST | `/billing/invoices/:invoiceId/payment` | Process invoice payment | ✅ | ✅ |
| POST | `/billing/invoices/:invoiceId/refund` | Refund invoice payment | ✅ | ✅ |

### Analytics & Reports
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/billing/analytics` | Get billing analytics | ✅ | ❌ |
| GET | `/billing/reports/revenue` | Get revenue report | ✅ | ❌ |
| GET | `/billing/reports/churn` | Get churn analysis | ✅ | ❌ |

### Admin Operations
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/billing/admin/process-billing` | Process recurring billing | ✅ | ❌ |
| POST | `/billing/admin/retry-payments` | Retry failed payments | ✅ | ❌ |
| POST | `/billing/admin/process-dunning` | Process dunning | ✅ | ❌ |
| POST | `/billing/admin/tenants/:tenantId/suspend` | Suspend tenant service | ✅ | ❌ |
| POST | `/billing/admin/tenants/:tenantId/reactivate` | Reactivate tenant service | ✅ | ❌ |

## Analytics Module (20 endpoints)

### Event Tracking
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/track/event` | Track custom event | ❌ | ✅ |
| POST | `/track/page-view` | Track page view | ❌ | ✅ |
| POST | `/track/product-view` | Track product view | ❌ | ✅ |
| POST | `/track/purchase` | Track purchase | ❌ | ✅ |

### Dashboard Analytics
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/dashboard` | Get dashboard stats | ✅ | ✅ |
| GET | `/traffic` | Get traffic statistics | ✅ | ✅ |
| GET | `/sales` | Get sales statistics | ✅ | ✅ |
| GET | `/realtime` | Get real-time statistics | ✅ | ✅ |

### Top Performers
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/top/products` | Get top products | ✅ | ✅ |
| GET | `/top/pages` | Get top pages | ✅ | ✅ |
| GET | `/top/referrers` | Get top referrers | ✅ | ✅ |

### Advanced Analytics
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/advanced/cohorts` | Get cohort analysis | ✅ | ✅ |
| GET | `/advanced/funnel` | Get funnel analysis | ✅ | ✅ |
| GET | `/advanced/clv` | Get customer lifetime value | ✅ | ✅ |
| GET | `/advanced/retention` | Get retention rate | ✅ | ✅ |

### Reports
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/reports/generate` | Generate report | ✅ | ✅ |
| POST | `/reports/schedule` | Schedule report | ✅ | ✅ |
| GET | `/reports/scheduled` | Get scheduled reports | ✅ | ✅ |
| DELETE | `/reports/scheduled/:id` | Delete scheduled report | ✅ | ✅ |
| POST | `/export` | Export data | ✅ | ✅ |

## Marketing Module (29 endpoints)

### Campaigns
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/marketing/campaigns` | Create campaign | ✅ | ✅ |
| GET | `/marketing/campaigns` | Get campaigns | ✅ | ✅ |
| GET | `/marketing/campaigns/:id` | Get specific campaign | ✅ | ✅ |
| PUT | `/marketing/campaigns/:id` | Update campaign | ✅ | ✅ |
| DELETE | `/marketing/campaigns/:id` | Delete campaign | ✅ | ✅ |
| POST | `/marketing/campaigns/:id/schedule` | Schedule campaign | ✅ | ✅ |
| POST | `/marketing/campaigns/:id/start` | Start campaign | ✅ | ✅ |
| POST | `/marketing/campaigns/:id/pause` | Pause campaign | ✅ | ✅ |
| POST | `/marketing/campaigns/:id/stop` | Stop campaign | ✅ | ✅ |
| GET | `/marketing/campaigns/:id/emails` | Get campaign emails | ✅ | ✅ |
| GET | `/marketing/campaigns/:id/stats` | Get campaign statistics | ✅ | ✅ |

### Templates
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/marketing/templates` | Create template | ✅ | ✅ |
| GET | `/marketing/templates` | Get templates | ✅ | ✅ |
| GET | `/marketing/templates/:id` | Get specific template | ✅ | ✅ |
| PUT | `/marketing/templates/:id` | Update template | ✅ | ✅ |
| DELETE | `/marketing/templates/:id` | Delete template | ✅ | ✅ |

### Segments
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/marketing/segments` | Create segment | ✅ | ✅ |
| GET | `/marketing/segments` | Get segments | ✅ | ✅ |
| GET | `/marketing/segments/:id` | Get specific segment | ✅ | ✅ |
| PUT | `/marketing/segments/:id` | Update segment | ✅ | ✅ |
| DELETE | `/marketing/segments/:id` | Delete segment | ✅ | ✅ |
| POST | `/marketing/segments/:id/refresh` | Refresh segment | ✅ | ✅ |

### Newsletter
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/marketing/newsletter/subscribe` | Subscribe to newsletter | ❌ | ✅ |
| POST | `/marketing/newsletter/unsubscribe` | Unsubscribe from newsletter | ❌ | ✅ |
| GET | `/marketing/newsletter/subscribers` | Get subscribers | ✅ | ✅ |
| GET | `/marketing/newsletter/subscribers/:email` | Get specific subscriber | ✅ | ✅ |

### Abandoned Carts & Settings
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/marketing/abandoned-carts` | Create abandoned cart record | ✅ | ✅ |
| GET | `/marketing/abandoned-carts` | Get abandoned carts | ✅ | ✅ |
| GET | `/marketing/settings` | Get marketing settings | ✅ | ✅ |
| PUT | `/marketing/settings` | Update marketing settings | ✅ | ✅ |
| GET | `/marketing/overview` | Get marketing overview | ✅ | ✅ |
| GET | `/marketing/track/open/:emailId` | Track email open | ❌ | ✅ |
| GET | `/marketing/track/click/:emailId` | Track email click | ❌ | ✅ |

## Discount Module (22 endpoints)

### Discount Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/discounts` | Create discount/coupon | ✅ | ✅ |
| GET | `/discounts` | Get discounts | ✅ | ✅ |
| GET | `/discounts/:id` | Get specific discount | ✅ | ✅ |
| PUT | `/discounts/:id` | Update discount | ✅ | ✅ |
| DELETE | `/discounts/:id` | Delete discount | ✅ | ✅ |
| GET | `/discounts/:id/usage` | Get discount usage | ✅ | ✅ |
| GET | `/discounts/stats` | Get discount statistics | ✅ | ✅ |
| GET | `/discounts/performance` | Get top discounts | ✅ | ✅ |
| GET | `/discounts/revenue-impact` | Get discount revenue impact | ✅ | ✅ |

### Discount Application (Public)
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/validate-discount` | Validate discount code | ❌ | ✅ |
| POST | `/apply-discount` | Apply discount to order | ❌ | ✅ |
| DELETE | `/remove-discount/:orderId` | Remove discount from order | ❌ | ✅ |

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
| GET | `/store-credit/:customerId` | Get customer store credit | ✅ | ✅ |
| POST | `/store-credit/:customerId/add` | Add store credit | ✅ | ✅ |
| POST | `/store-credit/:customerId/use` | Use store credit | ✅ | ✅ |

## Reviews Module (25 endpoints)

### Review Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/reviews` | Create review | ✅ | ✅ |
| GET | `/reviews` | Get reviews | ✅ | ✅ |
| GET | `/reviews/:id` | Get specific review | ✅ | ✅ |
| PUT | `/reviews/:id` | Update review | ✅ | ✅ |
| DELETE | `/reviews/:id` | Delete review | ✅ | ✅ |

### Review Moderation
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/reviews/:id/approve` | Approve review | ✅ | ✅ |
| POST | `/reviews/:id/reject` | Reject review | ✅ | ✅ |
| POST | `/reviews/:id/spam` | Mark review as spam | ✅ | ✅ |
| POST | `/reviews/bulk-moderate` | Bulk moderate reviews | ✅ | ✅ |
| GET | `/reviews/pending` | Get pending reviews | ✅ | ✅ |

### Review Interactions
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/reviews/:id/replies` | Add review reply | ✅ | ✅ |
| GET | `/reviews/:id/replies` | Get review replies | ✅ | ✅ |
| POST | `/reviews/:id/react` | React to review | ✅ | ✅ |
| DELETE | `/reviews/:id/react` | Remove reaction | ✅ | ✅ |

### Product Reviews
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/products/:productId/reviews` | Get product reviews | ✅ | ✅ |
| GET | `/products/:productId/reviews/summary` | Get product review summary | ✅ | ✅ |

### Review Invitations & Analytics
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/review-invitations` | Create review invitation | ✅ | ✅ |
| GET | `/review-invitations` | Get review invitations | ✅ | ✅ |
| POST | `/review-invitations/:id/send` | Send review invitation | ✅ | ✅ |
| GET | `/review-invite/:token` | Process invitation click | ❌ | ❌ |
| GET | `/reviews/stats` | Get review statistics | ✅ | ✅ |
| GET | `/reviews/trends` | Get review trends | ✅ | ✅ |
| GET | `/reviews/top-products` | Get top rated products | ✅ | ✅ |
| GET | `/reviews/recent` | Get recent reviews | ✅ | ✅ |
| GET | `/reviews/settings` | Get review settings | ✅ | ✅ |
| PUT | `/reviews/settings` | Update review settings | ✅ | ✅ |

## Support Module (15 endpoints)

### Support Tickets
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/support/tickets` | Create support ticket | ✅ | ✅ |
| GET | `/support/tickets` | Get support tickets | ✅ | ✅ |
| GET | `/support/tickets/:id` | Get specific ticket | ✅ | ✅ |
| PUT | `/support/tickets/:id` | Update ticket | ✅ | ✅ |
| DELETE | `/support/tickets/:id` | Delete ticket | ✅ | ✅ |
| POST | `/support/tickets/:id/assign` | Assign ticket | ✅ | ✅ |
| POST | `/support/tickets/:id/resolve` | Resolve ticket | ✅ | ✅ |
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
| GET | `/support/stats` | Get ticket statistics | ✅ | ✅ |

## Contact Module (31 endpoints)

### Contact Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/contacts` | Create contact | ✅ | ✅ |
| GET | `/contacts` | List contacts | ✅ | ✅ |
| GET | `/contacts/:id` | Get contact details | ✅ | ✅ |
| PUT | `/contacts/:id` | Update contact | ✅ | ✅ |
| DELETE | `/contacts/:id` | Delete contact | ✅ | ✅ |
| POST | `/contacts/bulk` | Bulk update contacts | ✅ | ✅ |
| POST | `/contacts/export` | Export contacts | ✅ | ✅ |

### Contact Status Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| PUT | `/contacts/:id/status` | Update contact status | ✅ | ✅ |
| PUT | `/contacts/:id/assign` | Assign contact | ✅ | ✅ |
| PUT | `/contacts/:id/priority` | Update contact priority | ✅ | ✅ |
| POST | `/contacts/:id/tags` | Add contact tags | ✅ | ✅ |
| DELETE | `/contacts/:id/tags` | Remove contact tags | ✅ | ✅ |

### Contact Interactions
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/contacts/:id/replies` | Create contact reply | ✅ | ✅ |
| GET | `/contacts/:id/replies` | List contact replies | ✅ | ✅ |
| POST | `/contacts/:id/notes` | Add contact note | ✅ | ✅ |
| GET | `/contacts/:id/notes` | List contact notes | ✅ | ✅ |

### Contact Forms
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/contact-forms` | Create contact form | ✅ | ✅ |
| GET | `/contact-forms` | List contact forms | ✅ | ✅ |
| GET | `/contact-forms/public/:form_type` | Get public contact form | ❌ | ✅ |
| POST | `/contact-forms/public/:form_type/submit` | Submit public contact form | ❌ | ✅ |

### Contact Templates & Analytics
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/contact-templates` | Create contact template | ✅ | ✅ |
| GET | `/contact-templates` | List contact templates | ✅ | ✅ |
| GET | `/contact-settings` | Get contact settings | ✅ | ✅ |
| PUT | `/contact-settings` | Update contact settings | ✅ | ✅ |
| GET | `/contact-analytics` | Get contact analytics | ✅ | ✅ |
| GET | `/contact-analytics/metrics` | Get contact metrics | ✅ | ✅ |
| GET | `/contact-analytics/performance` | Get agent performance | ✅ | ✅ |
| GET | `/contact-analytics/satisfaction` | Get customer satisfaction | ✅ | ✅ |
| GET | `/contact-analytics/resolution-time` | Get resolution time analytics | ✅ | ✅ |
| GET | `/contact-analytics/response-time` | Get response time analytics | ✅ | ✅ |

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
| GET | `/observability/system/info` | Get system information | ✅ | ❌ |
| GET | `/observability/system/stats` | Get system statistics | ✅ | ❌ |

## Modules in Development

### Shipping Module
**Status:** TODO - Implementation pending
**Planned Features:**
- Shipping rate calculation
- Shipping label creation and management
- Package tracking integration
- Multi-carrier support (Pathao, RedX, Paperfly, DHL, FedEx)

### Content Module  
**Status:** TODO - Implementation pending
**Planned Features:**
- Page and blog management
- Media library management
- SEO management
- Navigation and menu management

### Webhook Module
**Status:** TODO - Implementation pending
**Planned Features:**
- Webhook endpoint management
- Webhook delivery monitoring
- Provider webhooks (Stripe, PayPal, Bkash, Nagad, etc.)

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

The e-commerce platform currently implements **195+ API endpoints** across **13 active modules**:

- ✅ **Product Module** - 29 endpoints (complete CRUD, variants, categories)
- ✅ **Order Module** - 15 endpoints (order management, tracking, payments)
- ✅ **User Module** - 12 endpoints (authentication, user management)
- ✅ **Tenant Module** - 11 endpoints (multi-tenancy management)
- ✅ **Payment Module** - 6 endpoints (payment processing, webhooks)
- ✅ **Notification Module** - 14 endpoints (notifications, templates, preferences)
- ✅ **Billing Module** - 31 endpoints (plans, subscriptions, usage, invoices)
- ✅ **Analytics Module** - 20 endpoints (tracking, dashboard, reports)
- ✅ **Marketing Module** - 29 endpoints (campaigns, templates, segments)
- ✅ **Discount Module** - 22 endpoints (discounts, gift cards, store credit)
- ✅ **Reviews Module** - 25 endpoints (reviews, moderation, invitations)
- ✅ **Support Module** - 15 endpoints (tickets, FAQ, knowledge base)
- ✅ **Contact Module** - 31 endpoints (contact management, forms, templates)
- ✅ **Observability Module** - 12 endpoints (health, metrics, logs, alerts)

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