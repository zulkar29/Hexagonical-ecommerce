# API Documentation

Comprehensive REST API specification for the e-commerce SaaS platform with **460+ endpoints** covering all business operations including advanced features, AI-powered capabilities, enterprise-grade functionality, and comprehensive observability.

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

### Authentication & Users
```
POST   /auth/login                    # User login
POST   /auth/register                 # User registration  
POST   /auth/refresh                  # Refresh JWT token
POST   /auth/logout                   # User logout
POST   /auth/password/forgot          # Send password reset email
POST   /auth/password/reset           # Reset password with token
GET    /auth/me                       # Get current user profile
PUT    /auth/me                       # Update current user profile
PUT    /auth/password                 # Change password
PUT    /auth/mfa/enable               # Enable multi-factor auth
PUT    /auth/mfa/disable              # Disable multi-factor auth
POST   /auth/mfa/verify               # Verify MFA token
GET    /users                         # List all users (admin)
GET    /users/:id                     # Get user details
PUT    /users/:id                     # Update user (admin)
DELETE /users/:id                     # Delete user (admin)
PATCH  /users/:id/status              # Update user status
GET    /users/:id/permissions         # Get user permissions
PUT    /users/:id/permissions         # Update user permissions
GET    /users/:id/activity            # Get user activity log
```

### Tenants & Store Management
```
GET    /tenants                       # List all tenants (super admin)
POST   /tenants                       # Create new tenant
GET    /tenants/:id                   # Get tenant details
PUT    /tenants/:id                   # Update tenant information
DELETE /tenants/:id                   # Delete tenant
PATCH  /tenants/:id/status            # Update tenant status
GET    /tenants/:id/settings          # Get tenant settings
PUT    /tenants/:id/settings          # Update tenant settings
POST   /tenants/:id/domain            # Set custom domain
DELETE /tenants/:id/domain            # Remove custom domain
GET    /tenants/:id/stats             # Get tenant statistics
GET    /tenants/:id/users             # List tenant users
POST   /tenants/:id/users             # Add user to tenant
DELETE /tenants/:id/users/:userId     # Remove user from tenant
GET    /tenants/:id/subscription      # Get tenant subscription
PUT    /tenants/:id/subscription      # Update tenant subscription
POST   /tenants/:id/subscription/upgrade  # Upgrade subscription plan
POST   /tenants/:id/subscription/downgrade  # Downgrade subscription plan
POST   /tenants/:id/subscription/cancel    # Cancel subscription
```

### Subscription Plans
```
GET    /plans                         # List all subscription plans
GET    /plans/:id                     # Get plan details
POST   /plans                         # Create new plan (admin)
PUT    /plans/:id                     # Update plan (admin)
DELETE /plans/:id                     # Delete plan (admin)
PATCH  /plans/:id/status              # Enable/disable plan
GET    /plans/:id/features            # Get plan features
PUT    /plans/:id/features            # Update plan features
```

### Products & Catalog
```
GET    /products                      # List products with filters
POST   /products                      # Create new product
GET    /products/:id                  # Get product by ID
PUT    /products/:id                  # Update product
DELETE /products/:id                  # Delete product
PATCH  /products/:id/status           # Update product status
GET    /products/:id/variants         # Get product variants
POST   /products/:id/variants         # Create product variant
PUT    /products/:id/variants/:variantId  # Update product variant
DELETE /products/:id/variants/:variantId  # Delete product variant
POST   /products/:id/images           # Upload product images
PUT    /products/:id/images/:imageId  # Update product image
DELETE /products/:id/images/:imageId  # Delete product image
GET    /products/:id/reviews          # Get product reviews
POST   /products/:id/reviews          # Create product review
PUT    /products/:id/reviews/:reviewId   # Update product review
DELETE /products/:id/reviews/:reviewId   # Delete product review
GET    /products/:id/attributes       # Get product attributes
PUT    /products/:id/attributes       # Update product attributes
GET    /products/:id/tags             # Get product tags
PUT    /products/:id/tags             # Update product tags
GET    /products/:id/collections      # Get product collections
PUT    /products/:id/collections      # Update product collections
GET    /products/:id/metafields       # Get product metafields
PUT    /products/:id/metafields       # Update product metafields
GET    /products/:id/pricing-rules    # Get product pricing rules
PUT    /products/:id/pricing-rules    # Update product pricing rules
GET    /products/search               # Search products
GET    /products/featured             # Get featured products
GET    /products/bestsellers          # Get bestselling products
GET    /products/new-arrivals         # Get new arrival products
GET    /products/related/:id          # Get related products
GET    /products/recommendations/:customerId  # Get personalized recommendations
PUT    /products/bulk                 # Bulk update products
DELETE /products/bulk                 # Bulk delete products
POST   /products/import               # Import products from CSV
GET    /products/export               # Export products to CSV
POST   /products/:id/duplicate        # Duplicate product
```

### Categories & Collections
```
GET    /categories                    # List all categories
POST   /categories                    # Create new category  
GET    /categories/:id                # Get category by ID
PUT    /categories/:id                # Update category
DELETE /categories/:id                # Delete category
GET    /categories/:id/products       # Get products in category
POST   /categories/:id/products       # Add product to category
DELETE /categories/:id/products/:productId  # Remove product from category
PUT    /categories/reorder            # Reorder categories
GET    /collections                   # List smart collections
POST   /collections                   # Create smart collection
GET    /collections/:id               # Get collection details
PUT    /collections/:id               # Update collection
DELETE /collections/:id               # Delete collection
GET    /collections/:id/products      # Get products in collection
POST   /collections/:id/products      # Add product to collection
DELETE /collections/:id/products/:productId  # Remove product from collection
PUT    /collections/:id/rules         # Update collection rules
GET    /collections/smart             # List smart collections
GET    /collections/manual            # List manual collections
```

### Inventory Management
```
GET    /inventory                     # List inventory items
GET    /inventory/products/:id        # Get product inventory
PUT    /inventory/products/:id        # Update product inventory
GET    /inventory/variants/:id        # Get variant inventory
PUT    /inventory/variants/:id        # Update variant inventory
POST   /inventory/adjustments         # Create inventory adjustment
GET    /inventory/adjustments         # List inventory adjustments
GET    /inventory/adjustments/:id     # Get adjustment details
PUT    /inventory/adjustments/:id     # Update inventory adjustment
DELETE /inventory/adjustments/:id     # Delete inventory adjustment
GET    /inventory/low-stock           # Get low stock items
GET    /inventory/out-of-stock        # Get out of stock items
POST   /inventory/transfers           # Transfer inventory between locations
GET    /inventory/transfers           # List inventory transfers
GET    /inventory/transfers/:id       # Get transfer details
PUT    /inventory/transfers/:id       # Update transfer status
GET    /inventory/locations           # List inventory locations
POST   /inventory/locations           # Create inventory location
PUT    /inventory/locations/:id       # Update inventory location
DELETE /inventory/locations/:id       # Delete inventory location
PUT    /inventory/bulk                # Bulk inventory update
GET    /inventory/history/:id         # Get inventory history
```

### Orders & Fulfillment
```
GET    /orders                        # List orders with filters
POST   /orders                        # Create new order
GET    /orders/:id                    # Get order details
PUT    /orders/:id                    # Update order
DELETE /orders/:id                    # Cancel order
PATCH  /orders/:id/status             # Update order status
GET    /orders/:id/items              # Get order items
POST   /orders/:id/items              # Add item to order
PUT    /orders/:id/items/:itemId      # Update order item
DELETE /orders/:id/items/:itemId      # Remove order item
POST   /orders/:id/fulfill            # Fulfill order
POST   /orders/:id/refund             # Refund order
GET    /orders/:id/payments           # Get order payments
POST   /orders/:id/payments           # Process payment
GET    /orders/:id/shipments          # Get order shipments
POST   /orders/:id/shipments          # Create shipment
PUT    /orders/:id/shipments/:shipId  # Update shipment
GET    /orders/:id/tracking           # Get tracking information
POST   /orders/bulk/export            # Bulk export orders
```

### Customers & Accounts
```
GET    /customers                     # List customers
POST   /customers                     # Create customer
GET    /customers/:id                 # Get customer details
PUT    /customers/:id                 # Update customer
DELETE /customers/:id                 # Delete customer
PATCH  /customers/:id/status          # Update customer status
GET    /customers/:id/orders          # Get customer orders
GET    /customers/:id/addresses       # Get customer addresses
POST   /customers/:id/addresses       # Add customer address
PUT    /customers/:id/addresses/:addressId  # Update customer address
DELETE /customers/:id/addresses/:addressId  # Delete customer address
GET    /customers/:id/wishlist        # Get customer wishlist
POST   /customers/:id/wishlist        # Add item to wishlist
DELETE /customers/:id/wishlist/:productId   # Remove from wishlist
GET    /customers/:id/reviews         # Get customer reviews
GET    /customers/:id/activity        # Get customer activity log
GET    /customers/:id/stats           # Get customer statistics
POST   /customers/search              # Search customers
POST   /customers/import              # Import customers from CSV
GET    /customers/export              # Export customers to CSV
POST   /customers/bulk-update         # Bulk update customers
DELETE /customers/bulk                # Bulk delete customers
```

### Shopping Cart & Checkout
```
GET    /cart                          # Get current cart
POST   /cart/items                    # Add item to cart
PUT    /cart/items/:itemId            # Update cart item quantity
DELETE /cart/items/:itemId            # Remove item from cart
DELETE /cart                          # Clear cart
POST   /cart/discounts                # Apply discount code
DELETE /cart/discounts                # Remove discount code
POST   /cart/save                     # Save cart for later
POST   /cart/restore                  # Restore saved cart
GET    /cart/abandoned                # List abandoned carts (admin)
GET    /checkout/shipping             # Get shipping options
POST   /checkout/shipping             # Set shipping method
GET    /checkout/payment-methods      # Get available payment methods
POST   /checkout/calculate            # Calculate totals
POST   /checkout/validate             # Validate checkout data
POST   /checkout/complete             # Complete checkout
GET    /checkout/success/:orderId     # Checkout success page data
```

### Discounts, Coupons & Gift Cards
```
GET    /discounts                     # List discount codes
POST   /discounts                     # Create discount code
GET    /discounts/:id                 # Get discount details
PUT    /discounts/:id                 # Update discount
DELETE /discounts/:id                 # Delete discount
PATCH  /discounts/:id/status          # Enable/disable discount
GET    /discounts/:id/usage           # Get discount usage stats
POST   /discounts/validate            # Validate discount code
GET    /gift-cards                    # List gift cards
POST   /gift-cards                    # Create gift card
GET    /gift-cards/:id                # Get gift card details
PUT    /gift-cards/:id                # Update gift card
DELETE /gift-cards/:id                # Delete gift card
PATCH  /gift-cards/:id/status         # Enable/disable gift card
POST   /gift-cards/:id/balance        # Check gift card balance
POST   /gift-cards/:id/transactions   # Add gift card transaction
GET    /gift-cards/:id/transactions   # Get gift card transactions
GET    /store-credits                 # List store credits
POST   /store-credits                 # Issue store credit
GET    /store-credits/:customerId     # Get customer store credits
POST   /store-credits/:customerId/apply  # Apply store credit
```

### Payments & Billing
```
GET    /payments                      # List payments
GET    /payments/:id                  # Get payment details
POST   /payments/process              # Process payment
POST   /payments/refund               # Process refund
GET    /payment-methods               # List payment methods
POST   /payment-methods               # Add payment method
PUT    /payment-methods/:id           # Update payment method
DELETE /payment-methods/:id           # Remove payment method
GET    /billing/invoices              # List invoices
GET    /billing/invoices/:id          # Get invoice
POST   /billing/invoices/:id/send     # Send invoice email
```

### Shipping & Logistics
```
GET    /shipping/zones                # List shipping zones
POST   /shipping/zones                # Create shipping zone
PUT    /shipping/zones/:id            # Update shipping zone
DELETE /shipping/zones/:id            # Delete shipping zone
GET    /shipping/rates                # List shipping rates
POST   /shipping/rates                # Create shipping rate
PUT    /shipping/rates/:id            # Update shipping rate
DELETE /shipping/rates/:id            # Delete shipping rate
POST   /shipping/calculate            # Calculate shipping cost
GET    /shipping/carriers             # List shipping carriers
POST   /shipping/labels               # Create shipping label
GET    /shipping/tracking/:number     # Track shipment
```

### Analytics & Reports
```
GET    /analytics/dashboard           # Get dashboard metrics
GET    /analytics/sales               # Sales analytics
GET    /analytics/products            # Product performance
GET    /analytics/customers           # Customer analytics
GET    /analytics/traffic             # Website traffic
GET    /reports/sales                 # Sales reports
GET    /reports/inventory             # Inventory reports
GET    /reports/customers             # Customer reports
GET    /reports/tax                   # Tax reports
POST   /reports/generate              # Generate custom report
GET    /reports/export/:id            # Export report
```

### Content Management
```
GET    /pages                         # List pages
POST   /pages                         # Create page
GET    /pages/:id                     # Get page
PUT    /pages/:id                     # Update page
DELETE /pages/:id                     # Delete page
PATCH  /pages/:id/status              # Update page status
GET    /blogs                         # List blog posts
POST   /blogs                         # Create blog post
GET    /blogs/:id                     # Get blog post
PUT    /blogs/:id                     # Update blog post
DELETE /blogs/:id                     # Delete blog post
PATCH  /blogs/:id/status              # Update blog status
GET    /menus                         # List navigation menus
POST   /menus                         # Create menu
PUT    /menus/:id                     # Update menu
DELETE /menus/:id                     # Delete menu
GET    /redirects                     # List URL redirects
POST   /redirects                     # Create redirect
PUT    /redirects/:id                 # Update redirect
DELETE /redirects/:id                 # Delete redirect
GET    /sitemap                       # Generate sitemap
GET    /robots.txt                    # Get robots.txt
PUT    /robots.txt                    # Update robots.txt
```

### Media & Assets
```
GET    /media                         # List media files
POST   /media/upload                  # Upload media file
DELETE /media/:id                     # Delete media file
POST   /media/bulk-upload             # Bulk upload files
GET    /media/search                  # Search media files
POST   /media/organize                # Organize into folders
```

### Settings & Configuration
```
GET    /settings/store                # Get store settings
PUT    /settings/store                # Update store settings
GET    /settings/theme                # Get theme settings
PUT    /settings/theme                # Update theme
GET    /settings/seo                  # Get SEO settings
PUT    /settings/seo                  # Update SEO settings
GET    /settings/notifications        # Get notification settings
PUT    /settings/notifications        # Update notifications
GET    /settings/integrations         # List integrations
POST   /settings/integrations         # Add integration
PUT    /settings/integrations/:id     # Update integration
DELETE /settings/integrations/:id     # Remove integration
GET    /settings/domains              # Get domain settings
PUT    /settings/domains              # Update domain settings
GET    /settings/email                # Get email settings
PUT    /settings/email                # Update email settings
```

### Themes & Templates
```
GET    /themes                        # List available themes
GET    /themes/:id                    # Get theme details
POST   /themes                        # Upload custom theme (admin)
PUT    /themes/:id                    # Update theme (admin)
DELETE /themes/:id                    # Delete theme (admin)
POST   /themes/:id/install            # Install theme
GET    /themes/current                # Get current active theme
POST   /themes/:id/preview            # Preview theme
GET    /templates                     # List page templates
GET    /templates/:id                 # Get template content
PUT    /templates/:id                 # Update template content
```

### Webhooks & Events
```
GET    /webhooks                      # List webhooks
POST   /webhooks                      # Create webhook
GET    /webhooks/:id                  # Get webhook
PUT    /webhooks/:id                  # Update webhook
DELETE /webhooks/:id                  # Delete webhook
POST   /webhooks/:id/test             # Test webhook
GET    /webhooks/:id/deliveries       # Get webhook deliveries
POST   /webhooks/:id/redeliver        # Redeliver webhook
```

### Reviews & Ratings
```
GET    /reviews                       # List all reviews (admin)
GET    /reviews/:id                   # Get review details
PUT    /reviews/:id                   # Update review
DELETE /reviews/:id                   # Delete review
PATCH  /reviews/:id/status            # Approve/reject review
POST   /reviews/:id/reply             # Reply to review
GET    /reviews/pending               # Get pending reviews
POST   /reviews/bulk-moderate         # Bulk moderate reviews
```

### Taxes & Legal
```
GET    /taxes/rates                   # Get tax rates
POST   /taxes/rates                   # Create tax rate
PUT    /taxes/rates/:id               # Update tax rate
DELETE /taxes/rates/:id               # Delete tax rate
POST   /taxes/calculate               # Calculate tax for order
GET    /legal/terms                   # Get terms of service
PUT    /legal/terms                   # Update terms of service
GET    /legal/privacy                 # Get privacy policy
PUT    /legal/privacy                 # Update privacy policy
GET    /legal/cookies                 # Get cookie policy
PUT    /legal/cookies                 # Update cookie policy
```

### Notifications & Messages
```
GET    /notifications                 # List user notifications
POST   /notifications                 # Create notification
PATCH  /notifications/:id/read        # Mark as read
DELETE /notifications/:id             # Delete notification
PUT    /notifications/read-all        # Mark all as read
GET    /messages                      # List messages
POST   /messages                      # Send message
GET    /messages/:id                  # Get message
DELETE /messages/:id                  # Delete message
GET    /messages/unread               # Get unread count
```

### Advanced Features & Extensions
```
# Customer Segments & Loyalty
GET    /customer-segments             # List customer segments
POST   /customer-segments             # Create customer segment
GET    /customer-segments/:id         # Get segment details
PUT    /customer-segments/:id         # Update segment
DELETE /customer-segments/:id         # Delete segment
GET    /customer-segments/:id/customers # Get customers in segment
GET    /loyalty-programs              # List loyalty programs
POST   /loyalty-programs              # Create loyalty program
GET    /loyalty-programs/:id          # Get program details
PUT    /loyalty-programs/:id          # Update program
DELETE /loyalty-programs/:id          # Delete program
GET    /customers/:id/loyalty-points  # Get customer loyalty points
POST   /customers/:id/loyalty-points  # Add loyalty points
GET    /rewards                       # List available rewards
POST   /rewards                       # Create reward
POST   /rewards/:id/redeem            # Redeem reward

# Abandoned Cart Recovery
POST   /cart/abandoned/recovery       # Send recovery email
GET    /cart/abandoned/stats          # Abandoned cart statistics
POST   /cart/abandoned/:id/restore    # Restore abandoned cart

# Advanced Search & Filters
GET    /search                        # Advanced search across all entities
GET    /search/suggestions            # Get search suggestions
POST   /search/filters                # Create custom search filter
GET    /search/filters                # List saved search filters
GET    /search/popular                # Get popular search terms
POST   /search/index/rebuild          # Rebuild search index (admin)

# Multi-currency & Localization
GET    /currencies                    # List supported currencies
POST   /currencies                    # Add currency (admin)
GET    /currencies/:code/rates        # Get exchange rates
PUT    /currencies/:code/rates        # Update exchange rates
GET    /locales                       # List supported locales
POST   /locales                       # Add locale (admin)
GET    /translations/:locale          # Get translations for locale
PUT    /translations/:locale          # Update translations

# Advanced Analytics & Reporting
GET    /analytics/cohorts             # Cohort analysis
GET    /analytics/ltv                 # Customer lifetime value
GET    /analytics/funnel              # Conversion funnel analysis
GET    /analytics/retention           # Customer retention rates
GET    /analytics/segments            # Segment performance analysis
GET    /analytics/predictions         # AI-powered predictions
GET    /reports/scheduled             # List scheduled reports
POST   /reports/scheduled             # Create scheduled report
GET    /reports/custom                # List custom reports
POST   /reports/custom                # Create custom report
POST   /reports/export                # Export report data

# AI & Personalization
GET    /ai/recommendations/products   # AI product recommendations
GET    /ai/recommendations/customers  # Customer recommendations
POST   /ai/analyze/sentiment          # Sentiment analysis
GET    /ai/insights/trends            # Market trend insights
POST   /ai/optimize/pricing           # AI pricing optimization
GET    /personalization/rules         # List personalization rules
POST   /personalization/rules         # Create personalization rule

# Advanced Marketing
GET    /campaigns                     # List marketing campaigns
POST   /campaigns                     # Create campaign
GET    /campaigns/:id                 # Get campaign details
PUT    /campaigns/:id                 # Update campaign
DELETE /campaigns/:id                 # Delete campaign
GET    /campaigns/:id/performance     # Campaign performance metrics
POST   /email-marketing/automations   # Create email automation
GET    /email-marketing/templates     # List email templates
POST   /email-marketing/send          # Send marketing email
GET    /social-media/integrations     # List social media integrations
POST   /social-media/post             # Create social media post

# Marketplace & Multi-vendor
GET    /vendors                       # List vendors (if multi-vendor)
POST   /vendors                       # Create vendor account
GET    /vendors/:id                   # Get vendor details
PUT    /vendors/:id                   # Update vendor
GET    /vendors/:id/products          # Get vendor products
GET    /vendors/:id/orders            # Get vendor orders
GET    /vendors/:id/earnings          # Get vendor earnings
POST   /vendors/:id/payout            # Process vendor payout
GET    /marketplace/fees              # Get marketplace fee structure
PUT    /marketplace/fees              # Update fee structure

# Advanced Fulfillment
GET    /fulfillment/centers           # List fulfillment centers
POST   /fulfillment/centers           # Add fulfillment center
GET    /fulfillment/rules             # List fulfillment rules
POST   /fulfillment/rules             # Create fulfillment rule
GET    /fulfillment/capacity          # Check fulfillment capacity
POST   /fulfillment/allocate          # Allocate inventory to center
GET    /dropshipping/suppliers        # List dropshipping suppliers
POST   /dropshipping/suppliers        # Add supplier
GET    /dropshipping/products         # List supplier products
POST   /dropshipping/sync             # Sync supplier inventory

# Advanced Security & Compliance
GET    /security/sessions             # List active sessions
DELETE /security/sessions/:id         # Terminate session
GET    /security/audit-logs           # Security audit logs
POST   /security/2fa/setup            # Setup 2FA for user
POST   /security/2fa/verify           # Verify 2FA code
GET    /compliance/gdpr/data/:email   # Get user data (GDPR)
DELETE /compliance/gdpr/delete/:email # Delete user data (GDPR)
POST   /compliance/gdpr/export        # Export user data
GET    /compliance/pci/status         # PCI compliance status

# API Management
GET    /api/keys                      # List API keys
POST   /api/keys                      # Generate API key
DELETE /api/keys/:id                 # Revoke API key
GET    /api/usage                     # API usage statistics
GET    /api/rate-limits               # Current rate limits
GET    /api/webhooks/events           # List webhook event types
```

### System & Health
```
GET    /health                        # Health check endpoint
GET    /version                       # API version information
GET    /status                        # System status
GET    /metrics                       # System metrics (admin only)
PUT    /system/maintenance/enable     # Enable maintenance mode
PUT    /system/maintenance/disable    # Disable maintenance mode
GET    /system/logs                   # System logs (admin only)
POST   /system/backup                 # Create system backup
GET    /system/backups                # List backups
POST   /system/restore/:backupId      # Restore from backup
```

### Observability & Monitoring
```
GET    /observability/health          # Basic health status
GET    /observability/health/detailed # Detailed health with all services
GET    /observability/metrics         # Collected system metrics
GET    /observability/metrics/summary # Metrics summary statistics
GET    /observability/logs            # System log entries with filters
POST   /observability/logs            # Create log entry
GET    /observability/traces          # Distributed traces
GET    /observability/traces/:traceId # Get specific trace details
GET    /observability/alerts          # Current system alerts
POST   /observability/alerts          # Create new alert
PUT    /observability/alerts/:alertId/resolve  # Resolve alert
GET    /observability/system/info     # System information
GET    /observability/system/stats    # Real-time system statistics
```

## Request/Response Format

### Standard Response Structure
```json
{
  "success": true,
  "data": {},
  "message": "Success message",
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "pages": 10
  }
}
```

### Error Response Structure
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "email",
        "message": "Invalid email format"
      }
    ]
  }
}
```

## HTTP Status Codes
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `422` - Validation Error
- `500` - Internal Server Error

## Rate Limiting
- **Anonymous requests**: 100 requests/hour
- **Authenticated requests**: 1000 requests/hour
- **Tenant API calls**: Based on subscription plan

## Pagination
Query parameters for list endpoints:
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10, max: 100)
- `sort` - Sort field (default: created_at)
- `order` - Sort order: asc/desc (default: desc)

## Filtering
Common query parameters:
- `search` - Text search across relevant fields
- `status` - Filter by status
- `created_after` - Filter by creation date
- `created_before` - Filter by creation date

## API Examples

### Create Product
```bash
POST /api/v1/products
Content-Type: application/json
Authorization: Bearer <token>

{
  "name": "Premium T-Shirt",
  "description": "High-quality cotton t-shirt",
  "price": 29.99,
  "sku": "TSHIRT-001",
  "inventory": {
    "quantity": 100,
    "track_quantity": true
  },
  "categories": ["clothing", "t-shirts"]
}
```

### Get Products with Filters
```bash
GET /api/v1/products?search=shirt&status=active&page=1&limit=20
```

## Webhooks
Platform supports webhooks for real-time notifications:

### Supported Events
- `order.created`
- `order.updated`
- `payment.completed`
- `product.updated`
- `customer.created`

### Webhook Payload
```json
{
  "event": "order.created",
  "tenant_id": "uuid",
  "data": {},
  "timestamp": "2025-01-01T00:00:00Z"
}
```

## SDK Support
Official SDKs available for:
- JavaScript/Node.js
- PHP  
- Python
- Go

## Postman Collection
Import our Postman collection for easy API testing:
```
[Download Postman Collection](./postman/ecommerce-SaaS-api.json)
```