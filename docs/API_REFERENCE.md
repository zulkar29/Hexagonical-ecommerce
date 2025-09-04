# API Documentation

Comprehensive REST API specification for the e-commerce SaaS platform with **280+ implemented endpoints** covering all business operations across 13 active modules with multi-tenant architecture, authentication, and WebSocket real-time capabilities.

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

## Authentication Module (8 endpoints)

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

## User Module (12 endpoints)

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

## Product Module (35 endpoints)

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

## Inventory Module (9 endpoints)

### Inventory Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/inventory` | List inventory items | ✅ | ✅ |
| GET | `/inventory/:id` | Get inventory item | ✅ | ✅ |
| POST | `/inventory` | Create inventory item | ✅ | ✅ |
| PUT | `/inventory/:id` | Update inventory item | ✅ | ✅ |
| DELETE | `/inventory/:id` | Delete inventory item | ✅ | ✅ |
| POST | `/inventory/adjust` | Adjust inventory levels | ✅ | ✅ |
| GET | `/inventory/movements` | Get inventory movements | ✅ | ✅ |
| GET | `/inventory/low-stock` | Get low stock items | ✅ | ✅ |
| POST | `/inventory/bulk-update` | Bulk update inventory | ✅ | ✅ |

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

## Payment Module (11 endpoints)

### Payment Processing
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/payments` | Create a payment | ✅ | ✅ |
| GET | `/payments` | List payments | ✅ | ✅ |
| GET | `/payments/:id` | Get payment details | ✅ | ✅ |
| POST | `/payments/:id/process` | Process a payment | ✅ | ✅ |
| POST | `/payments/:id/refund` | Refund a payment | ✅ | ✅ |
| POST | `/payments/:id/capture` | Capture authorized payment | ✅ | ✅ |
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
| GET | `/shipping/zones` | Get shipping zones | ✅ | ✅ |
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
| GET | `/shipping/track/:trackingNumber` | Track shipment | ✅ | ✅ |
| GET | `/shipping/providers` | Get shipping providers | ✅ | ✅ |
| POST | `/shipping/providers` | Create shipping provider | ✅ | ✅ |
| PUT | `/shipping/providers/:id` | Update shipping provider | ✅ | ✅ |
| DELETE | `/shipping/providers/:id` | Delete shipping provider | ✅ | ✅ |
| GET | `/shipping/stats` | Get shipping statistics | ✅ | ✅ |

### Shipping Webhooks
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/shipping/webhooks/pathao` | Pathao shipping webhook | ❌ | ❌ |
| POST | `/shipping/webhooks/redx` | RedX shipping webhook | ❌ | ❌ |
| POST | `/shipping/webhooks/paperfly` | Paperfly shipping webhook | ❌ | ❌ |
| POST | `/shipping/webhooks/dhl` | DHL shipping webhook | ❌ | ❌ |
| POST | `/shipping/webhooks/fedex` | FedEx shipping webhook | ❌ | ❌ |

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

## Content Management Module (16 endpoints)

### Pages & Posts
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/content/pages` | List pages | ✅ | ✅ |
| POST | `/content/pages` | Create page | ✅ | ✅ |
| GET | `/content/pages/:id` | Get page | ✅ | ✅ |
| PUT | `/content/pages/:id` | Update page | ✅ | ✅ |
| DELETE | `/content/pages/:id` | Delete page | ✅ | ✅ |
| POST | `/content/pages/:id/publish` | Publish page | ✅ | ✅ |
| POST | `/content/pages/:id/unpublish` | Unpublish page | ✅ | ✅ |
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
| POST | `/webhooks/endpoints/:id/test` | Test webhook endpoint | ✅ | ✅ |

### Webhook Deliveries
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/webhooks/deliveries` | List webhook deliveries | ✅ | ✅ |
| GET | `/webhooks/deliveries/:id` | Get webhook delivery | ✅ | ✅ |
| POST | `/webhooks/deliveries/:id/retry` | Retry webhook delivery | ✅ | ✅ |

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
| GET | `/observability/system/info` | Get system information | ✅ | ❌ |
| GET | `/observability/system/stats` | Get system statistics | ✅ | ❌ |

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

## Category Module (15 endpoints)

### Category Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/categories` | List all categories | ❌ | ✅ |
| POST | `/categories` | Create category | ✅ | ✅ |
| GET | `/categories/:id` | Get category details | ❌ | ✅ |
| PUT | `/categories/:id` | Update category | ✅ | ✅ |
| DELETE | `/categories/:id` | Delete category | ✅ | ✅ |
| GET | `/categories/:id/children` | Get child categories | ❌ | ✅ |
| GET | `/categories/:id/products` | Get category products | ❌ | ✅ |
| POST | `/categories/:id/move` | Move category | ✅ | ✅ |

### Category Operations
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/categories/bulk` | Bulk create categories | ✅ | ✅ |
| PUT | `/categories/bulk` | Bulk update categories | ✅ | ✅ |
| DELETE | `/categories/bulk` | Bulk delete categories | ✅ | ✅ |
| GET | `/categories/tree` | Get category tree | ❌ | ✅ |
| GET | `/categories/stats` | Get category statistics | ✅ | ✅ |
| POST | `/categories/reorder` | Reorder categories | ✅ | ✅ |
| POST | `/categories/cleanup` | Cleanup empty categories | ✅ | ✅ |

## Cart Module (12 endpoints)

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
| POST | `/cart/apply-discount` | Apply discount code | ✅ | ✅ |
| DELETE | `/cart/remove-discount` | Remove discount | ✅ | ✅ |
| GET | `/cart/summary` | Get cart summary | ✅ | ✅ |
| POST | `/cart/save-for-later/:id` | Save item for later | ✅ | ✅ |
| POST | `/cart/move-to-cart/:id` | Move saved item to cart | ✅ | ✅ |
| POST | `/cart/cleanup` | Cleanup abandoned carts | ✅ | ✅ |

## Wishlist Module (15 endpoints)

### Wishlist Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/wishlists` | Get user wishlists | ✅ | ✅ |
| POST | `/wishlists` | Create wishlist | ✅ | ✅ |
| GET | `/wishlists/:id` | Get wishlist details | ✅ | ✅ |
| PUT | `/wishlists/:id` | Update wishlist | ✅ | ✅ |
| DELETE | `/wishlists/:id` | Delete wishlist | ✅ | ✅ |
| POST | `/wishlists/:id/items` | Add item to wishlist | ✅ | ✅ |
| DELETE | `/wishlists/:id/items/:itemId` | Remove wishlist item | ✅ | ✅ |

### Wishlist Operations
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| POST | `/wishlists/:id/share` | Share wishlist | ✅ | ✅ |
| GET | `/wishlists/shared/:token` | View shared wishlist | ❌ | ✅ |
| POST | `/wishlists/:id/move-to-cart` | Move items to cart | ✅ | ✅ |
| POST | `/wishlists/bulk/items` | Bulk add items | ✅ | ✅ |
| DELETE | `/wishlists/bulk/items` | Bulk remove items | ✅ | ✅ |
| GET | `/wishlists/analytics` | Get wishlist analytics | ✅ | ✅ |
| POST | `/wishlists/cleanup` | Cleanup old wishlists | ✅ | ✅ |
| GET | `/wishlists/popular-items` | Get popular wishlist items | ✅ | ✅ |

## Address Module (18 endpoints)

### Address Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/addresses` | Get user addresses | ✅ | ✅ |
| POST | `/addresses` | Create address | ✅ | ✅ |
| GET | `/addresses/:id` | Get address details | ✅ | ✅ |
| PUT | `/addresses/:id` | Update address | ✅ | ✅ |
| DELETE | `/addresses/:id` | Delete address | ✅ | ✅ |
| POST | `/addresses/:id/set-default` | Set default address | ✅ | ✅ |
| POST | `/addresses/:id/validate` | Validate address | ✅ | ✅ |

### Address Operations
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/addresses/customers/:customerId` | Get customer addresses | ✅ | ✅ |
| POST | `/addresses/bulk` | Bulk create addresses | ✅ | ✅ |
| PUT | `/addresses/bulk` | Bulk update addresses | ✅ | ✅ |
| DELETE | `/addresses/bulk` | Bulk delete addresses | ✅ | ✅ |
| POST | `/addresses/validate-bulk` | Bulk validate addresses | ✅ | ✅ |
| GET | `/addresses/stats` | Get address statistics | ✅ | ✅ |
| GET | `/addresses/validation-stats` | Get validation statistics | ✅ | ✅ |
| POST | `/addresses/cleanup` | Cleanup unvalidated addresses | ✅ | ✅ |
| GET | `/addresses/validation-trends` | Get validation trends | ✅ | ✅ |
| POST | `/addresses/geocode` | Geocode address | ✅ | ✅ |
| GET | `/addresses/suggestions` | Get address suggestions | ✅ | ✅ |

## Tax Module (22 endpoints)

### Tax Rules Management
| Method | URL | Description | Auth | Tenant |
|--------|-----|-------------|------|--------|
| GET | `/tax/rules` | List tax rules | ✅ | ✅ |
| POST | `/tax/rules` | Create tax rule | ✅ | ✅ |
| GET | `/tax/rules/:id` | Get tax rule | ✅ | ✅ |
| PUT | `/tax/rules/:id` | Update tax rule | ✅ | ✅ |
| DELETE | `/tax/rules/:id` | Delete tax rule | ✅ | ✅ |
| POST | `/tax/rules/:id/activate` | Activate tax rule | ✅ | ✅ |
| POST | `/tax/rules/:id/deactivate` | Deactivate tax rule | ✅ | ✅ |

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

The e-commerce platform currently implements **360+ API endpoints** across **19 active modules**:

- ✅ **Authentication Module** - 8 endpoints (user authentication, password management)
- ✅ **User Module** - 12 endpoints (profile management, preferences)
- ✅ **Product Module** - 35 endpoints (products, variants, categories, public access)
- ✅ **Category Module** - 15 endpoints (hierarchical category management)
- ✅ **Inventory Module** - 9 endpoints (inventory management, tracking)
- ✅ **Cart Module** - 12 endpoints (shopping cart management, operations)
- ✅ **Wishlist Module** - 15 endpoints (wishlist management, sharing, analytics)
- ✅ **Address Module** - 18 endpoints (address management, validation, geocoding)
- ✅ **Tax Module** - 22 endpoints (tax rules, rates, calculations)
- ✅ **Order Module** - 15 endpoints (order management, tracking, payments)
- ✅ **Payment Module** - 11 endpoints (payment processing, methods, webhooks)
- ✅ **Shipping Module** - 25 endpoints (zones, rates, labels, tracking, webhooks)
- ✅ **Notification Module** - 14 endpoints (notifications, templates, preferences)
- ✅ **Analytics Module** - 20 endpoints (tracking, dashboard, reports)
- ✅ **Marketing Module** - 29 endpoints (campaigns, templates, segments)
- ✅ **Discount Module** - 22 endpoints (discounts, gift cards, store credit)
- ✅ **Reviews Module** - 25 endpoints (reviews, moderation, invitations)
- ✅ **Support Module** - 15 endpoints (tickets, FAQ, knowledge base)
- ✅ **Contact Module** - 31 endpoints (contact management, forms, templates)
- ✅ **Content Management Module** - 16 endpoints (pages, posts, media, menus)
- ✅ **Webhook Module** - 25 endpoints (endpoint management, deliveries, events)
- ✅ **Billing Module** - 31 endpoints (plans, subscriptions, usage, invoices)
- ✅ **Tenant Module** - 11 endpoints (multi-tenancy management)
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

---

*Note: All endpoints require appropriate authentication and authorization. Refer to the Authentication section for details on obtaining and using access tokens.*