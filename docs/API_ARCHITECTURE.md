# API Architecture & Design

## Overview
Technical architecture for the API layer, covering REST, GraphQL, and real-time capabilities to optimize performance and developer experience for the multi-tenant e-commerce platform.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    CLIENT APPLICATIONS                     │
│         (Next.js Storefront + Dashboard)                   │
└─────────┬─────────────┬─────────────┬─────────────┬─────────┘
          │             │             │             │
    ┌─────▼─────┐ ┌─────▼─────┐ ┌─────▼─────┐ ┌─────▼─────┐
    │    REST   │ │  GraphQL  │ │    SSE    │ │ WebSocket │
    │    API    │ │    API    │ │ Real-time │ │   Events  │
    └─────┬─────┘ └─────┬─────┘ └─────┬─────┘ └─────┬─────┘
          └─────────────┼─────────────┼─────────────┘
                        │             │
┌─────────────────────▼─────────────▼─────────────────────────┐
│                   API GATEWAY LAYER                        │
│  ┌─────────────┬─────────────┬─────────────┬─────────────┐  │
│  │   Request   │   Rate      │   Auth &    │  Response   │  │
│  │ Deduplication│  Limiting  │ Validation  │   Caching   │  │
│  └─────────────┴─────────────┴─────────────┴─────────────┘  │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│               APPLICATION SERVICES                         │
│              (Go Modular Monolith)                         │
│  ┌─────────────┬─────────────┬─────────────┬─────────────┐  │
│  │   Product   │   Order     │   User      │   Tenant    │  │
│  │   Module    │   Module    │   Module    │   Module    │  │
│  └─────────────┴─────────────┴─────────────┴─────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## 1. GraphQL for Dashboard Complex Queries

### GraphQL Schema Design

#### Product Management Schema
```graphql
# Core Types
type Product {
  id: ID!
  tenantId: ID!
  name: String!
  description: String
  sku: String!
  price: Money!
  comparePrice: Money
  status: ProductStatus!
  inventory: ProductInventory!
  images: [ProductImage!]!
  variants: [ProductVariant!]!
  categories: [Category!]!
  tags: [String!]!
  seo: ProductSEO!
  createdAt: DateTime!
  updatedAt: DateTime!
  
  # Computed fields
  isInStock: Boolean!
  totalSales: Int!
  conversionRate: Float!
  profitMargin: Float!
}

type ProductInventory {
  trackQuantity: Boolean!
  quantity: Int!
  reservedQuantity: Int!
  availableQuantity: Int! # Computed
  lowStockAlert: Int!
  allowBackorder: Boolean!
}

type ProductImage {
  id: ID!
  url: String!
  altText: String
  position: Int!
  width: Int
  height: Int
}

type ProductVariant {
  id: ID!
  title: String!
  sku: String!
  price: Money
  inventory: ProductInventory!
  options: [VariantOption!]!
}

type Category {
  id: ID!
  name: String!
  slug: String!
  description: String
  parentId: ID
  children: [Category!]!
  productCount: Int!
  isActive: Boolean!
}

# Input Types
input ProductFilter {
  status: [ProductStatus!]
  categories: [ID!]
  priceRange: PriceRange
  inStock: Boolean
  hasVariants: Boolean
  tags: [String!]
  createdAfter: DateTime
  createdBefore: DateTime
}

input ProductSort {
  field: ProductSortField!
  direction: SortDirection!
}

enum ProductSortField {
  NAME
  PRICE
  CREATED_AT
  UPDATED_AT
  TOTAL_SALES
  INVENTORY
}

enum SortDirection {
  ASC
  DESC
}
```

#### Dashboard Analytics Schema
```graphql
type DashboardStats {
  revenue: RevenueStats!
  orders: OrderStats!
  products: ProductStats!
  customers: CustomerStats!
  traffic: TrafficStats!
  conversionFunnel: ConversionFunnel!
}

type RevenueStats {
  total: Money!
  change: Float! # Percentage change from previous period
  trend: [DataPoint!]!
  byPaymentMethod: [PaymentMethodStats!]!
  byCategory: [CategoryStats!]!
  topProducts: [ProductRevenueStats!]!
}

type OrderStats {
  total: Int!
  pending: Int!
  completed: Int!
  cancelled: Int!
  averageOrderValue: Money!
  trend: [DataPoint!]!
  fulfillmentStatus: [FulfillmentStats!]!
}

type DataPoint {
  date: DateTime!
  value: Float!
}

# Complex Queries
type Query {
  # Product Management
  products(
    first: Int = 20
    after: String
    filter: ProductFilter
    sort: ProductSort
  ): ProductConnection!
  
  product(id: ID!): Product
  
  # Dashboard Analytics  
  dashboardStats(period: Period!): DashboardStats!
  
  # Advanced Analytics
  salesReport(
    period: Period!
    groupBy: GroupBy!
    filters: ReportFilters
  ): SalesReport!
  
  inventoryReport(
    lowStockOnly: Boolean = false
    categories: [ID!]
  ): InventoryReport!
  
  customerAnalytics(
    segmentation: CustomerSegmentation!
  ): CustomerAnalytics!
}

# Mutations
type Mutation {
  # Product Management
  createProduct(input: CreateProductInput!): CreateProductResult!
  updateProduct(id: ID!, input: UpdateProductInput!): UpdateProductResult!
  deleteProduct(id: ID!): DeleteProductResult!
  
  # Bulk Operations
  bulkUpdateProducts(
    ids: [ID!]!
    updates: BulkProductUpdates!
  ): BulkUpdateResult!
  
  bulkImportProducts(
    file: Upload!
    mapping: ImportMapping!
  ): ImportResult!
  
  # Inventory Management
  adjustInventory(
    productId: ID!
    adjustment: InventoryAdjustment!
  ): InventoryAdjustmentResult!
}
```

#### Real-time Subscriptions
```graphql
type Subscription {
  # Product Updates
  productUpdated(tenantId: ID!): Product!
  inventoryChanged(tenantId: ID!, productIds: [ID!]): InventoryUpdate!
  
  # Order Updates
  orderStatusChanged(tenantId: ID!): Order!
  newOrder(tenantId: ID!): Order!
  
  # Dashboard Updates
  dashboardStatsUpdated(tenantId: ID!): DashboardStats!
  
  # System Notifications
  systemNotification(tenantId: ID!): SystemNotification!
}
```

### GraphQL Resolver Implementation
```go
// Product resolver with data loader pattern
type ProductResolver struct {
    productLoader *dataloader.Loader
    imageLoader   *dataloader.Loader
    categoryLoader *dataloader.Loader
}

func (r *ProductResolver) Product(ctx context.Context, id string) (*Product, error) {
    tenantID := GetTenantIDFromContext(ctx)
    
    // Use DataLoader to batch database calls
    key := fmt.Sprintf("%s:%s", tenantID, id)
    product, err := r.productLoader.Load(ctx, key)()
    if err != nil {
        return nil, err
    }
    
    return product.(*Product), nil
}

// Efficient N+1 prevention with DataLoader
func (r *ProductResolver) Images(ctx context.Context, product *Product) ([]*ProductImage, error) {
    // Batch load images for all products in current request
    images, err := r.imageLoader.LoadMany(ctx, []string{product.ID.String()})()
    if err != nil {
        return nil, err
    }
    
    return images.([]*ProductImage), nil
}

// Complex dashboard resolver
func (r *QueryResolver) DashboardStats(ctx context.Context, period Period) (*DashboardStats, error) {
    tenantID := GetTenantIDFromContext(ctx)
    
    // Run analytics queries in parallel
    var wg sync.WaitGroup
    results := make(chan interface{}, 5)
    
    // Revenue stats
    wg.Add(1)
    go func() {
        defer wg.Done()
        stats := r.analyticsService.GetRevenueStats(tenantID, period)
        results <- RevenueResult{stats}
    }()
    
    // Order stats
    wg.Add(1)
    go func() {
        defer wg.Done()
        stats := r.analyticsService.GetOrderStats(tenantID, period)
        results <- OrderResult{stats}
    }()
    
    // Collect results
    wg.Wait()
    close(results)
    
    dashboardStats := &DashboardStats{}
    for result := range results {
        switch r := result.(type) {
        case RevenueResult:
            dashboardStats.Revenue = r.Stats
        case OrderResult:
            dashboardStats.Orders = r.Stats
        }
    }
    
    return dashboardStats, nil
}
```

### GraphQL Performance Optimizations

#### Query Complexity Analysis
```go
// Prevent expensive queries
func ComplexityLimitMiddleware(maxComplexity int) gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        requestContext := c.Request.Context()
        
        complexity := graphql.GetRequestContext(requestContext).ComplexityLimit
        if complexity > maxComplexity {
            c.JSON(400, gin.H{
                "error": "Query too complex",
                "max_complexity": maxComplexity,
                "query_complexity": complexity,
            })
            c.Abort()
            return
        }
        
        c.Next()
    })
}

// Field-level complexity calculation
func ProductComplexity(childComplexity int, first *int, filter *ProductFilter) int {
    complexity := childComplexity
    
    if first != nil {
        complexity *= *first
    } else {
        complexity *= 20 // Default limit
    }
    
    // Add complexity for filters
    if filter != nil && filter.Categories != nil {
        complexity += len(*filter.Categories) * 2
    }
    
    return complexity
}
```

#### Query Depth Limiting
```yaml
# GraphQL Configuration
max_query_depth: 15
max_query_complexity: 1000
query_timeout: 30s
enable_introspection: false # Production
enable_playground: false   # Production

# Rate Limiting per Tenant
rate_limits:
  queries_per_minute: 60
  mutations_per_minute: 30
  subscriptions_per_tenant: 10
```

## 2. Server-Sent Events for Real-time Updates

### SSE Implementation Architecture
```go
// SSE Connection Manager
type SSEManager struct {
    connections map[string]map[string]*SSEConnection
    mutex       sync.RWMutex
    eventBus    *EventBus
}

type SSEConnection struct {
    ID       string
    TenantID string
    UserID   string
    Writer   http.ResponseWriter
    Request  *http.Request
    Done     chan bool
    Events   chan SSEEvent
}

type SSEEvent struct {
    Type string      `json:"type"`
    Data interface{} `json:"data"`
    ID   string      `json:"id,omitempty"`
}

// SSE endpoint handler
func (s *SSEManager) HandleSSE(c *gin.Context) {
    tenantID := c.GetString("tenant_id")
    userID := c.GetString("user_id")
    
    // Set SSE headers
    c.Header("Content-Type", "text/event-stream")
    c.Header("Cache-Control", "no-cache")
    c.Header("Connection", "keep-alive")
    c.Header("Access-Control-Allow-Origin", "*")
    
    // Create connection
    conn := &SSEConnection{
        ID:       uuid.New().String(),
        TenantID: tenantID,
        UserID:   userID,
        Writer:   c.Writer,
        Request:  c.Request,
        Done:     make(chan bool),
        Events:   make(chan SSEEvent, 100),
    }
    
    // Register connection
    s.addConnection(conn)
    defer s.removeConnection(conn)
    
    // Handle connection
    s.handleConnection(conn)
}

func (s *SSEManager) handleConnection(conn *SSEConnection) {
    // Send initial connection event
    conn.Events <- SSEEvent{
        Type: "connected",
        Data: map[string]string{
            "connection_id": conn.ID,
            "timestamp": time.Now().Format(time.RFC3339),
        },
    }
    
    for {
        select {
        case event := <-conn.Events:
            // Send event to client
            eventData, _ := json.Marshal(event.Data)
            fmt.Fprintf(conn.Writer, "event: %s\n", event.Type)
            fmt.Fprintf(conn.Writer, "data: %s\n", string(eventData))
            if event.ID != "" {
                fmt.Fprintf(conn.Writer, "id: %s\n", event.ID)
            }
            fmt.Fprintf(conn.Writer, "\n")
            
            // Flush to client
            if flusher, ok := conn.Writer.(http.Flusher); ok {
                flusher.Flush()
            }
            
        case <-conn.Done:
            return
            
        case <-conn.Request.Context().Done():
            return
            
        case <-time.After(30 * time.Second):
            // Send keepalive
            fmt.Fprintf(conn.Writer, ": keepalive\n\n")
            if flusher, ok := conn.Writer.(http.Flusher); ok {
                flusher.Flush()
            }
        }
    }
}
```

### Real-time Event Types
```yaml
# Inventory Updates
inventory_updated:
  product_id: "uuid"
  tenant_id: "uuid"
  old_quantity: 100
  new_quantity: 85
  threshold_alert: true

# Order Events
order_created:
  order_id: "uuid"
  tenant_id: "uuid"  
  customer_email: "customer@example.com"
  total_amount: 99.99
  status: "pending"

order_status_changed:
  order_id: "uuid"
  tenant_id: "uuid"
  old_status: "pending"
  new_status: "processing"
  timestamp: "2024-01-15T10:30:00Z"

# Product Updates  
product_updated:
  product_id: "uuid"
  tenant_id: "uuid"
  changes: ["price", "inventory", "status"]
  updated_fields:
    price: 29.99
    inventory: 150

# Dashboard Metrics
dashboard_metrics_updated:
  tenant_id: "uuid"
  metrics:
    total_revenue: 15420.50
    orders_today: 25
    visitors_online: 12
  updated_at: "2024-01-15T10:30:00Z"

# System Notifications
system_notification:
  tenant_id: "uuid"
  type: "warning" # info, warning, error
  title: "Low Stock Alert"
  message: "5 products are running low on inventory"
  action_url: "/dashboard/inventory"
```

### SSE Event Broadcasting
```go
// Event broadcasting service
func (s *SSEManager) BroadcastToTenant(tenantID string, event SSEEvent) {
    s.mutex.RLock()
    connections, exists := s.connections[tenantID]
    s.mutex.RUnlock()
    
    if !exists {
        return
    }
    
    for _, conn := range connections {
        select {
        case conn.Events <- event:
            // Event sent successfully
        case <-time.After(1 * time.Second):
            // Connection might be dead, remove it
            s.removeConnection(conn)
        }
    }
}

// Integration with business logic
func (s *ProductService) UpdateProduct(product *Product) error {
    oldProduct, _ := s.repo.GetProduct(product.TenantID, product.ID)
    
    // Update product
    err := s.repo.UpdateProduct(product)
    if err != nil {
        return err
    }
    
    // Broadcast real-time update
    changes := s.detectChanges(oldProduct, product)
    if len(changes) > 0 {
        event := SSEEvent{
            Type: "product_updated",
            Data: map[string]interface{}{
                "product_id": product.ID,
                "tenant_id":  product.TenantID,
                "changes":    changes,
                "updated_fields": s.getUpdatedFields(oldProduct, product, changes),
            },
        }
        
        s.sseManager.BroadcastToTenant(product.TenantID.String(), event)
    }
    
    return nil
}
```

## 3. API Pagination with Cursor-based Approach

### Cursor Pagination Implementation
```go
// Cursor-based pagination structure
type CursorPagination struct {
    First  *int    `json:"first,omitempty"`
    After  *string `json:"after,omitempty"`
    Last   *int    `json:"last,omitempty"`
    Before *string `json:"before,omitempty"`
}

type PageInfo struct {
    HasNextPage     bool    `json:"hasNextPage"`
    HasPreviousPage bool    `json:"hasPreviousPage"`
    StartCursor     *string `json:"startCursor,omitempty"`
    EndCursor       *string `json:"endCursor,omitempty"`
}

type ProductConnection struct {
    Edges    []*ProductEdge `json:"edges"`
    PageInfo *PageInfo      `json:"pageInfo"`
    TotalCount int          `json:"totalCount"`
}

type ProductEdge struct {
    Node   *Product `json:"node"`
    Cursor string   `json:"cursor"`
}

// Cursor encoding/decoding
func EncodeCursor(id uuid.UUID, timestamp time.Time) string {
    cursor := fmt.Sprintf("%s|%d", id.String(), timestamp.Unix())
    return base64.StdEncoding.EncodeToString([]byte(cursor))
}

func DecodeCursor(cursor string) (uuid.UUID, time.Time, error) {
    decoded, err := base64.StdEncoding.DecodeString(cursor)
    if err != nil {
        return uuid.Nil, time.Time{}, err
    }
    
    parts := strings.Split(string(decoded), "|")
    if len(parts) != 2 {
        return uuid.Nil, time.Time{}, errors.New("invalid cursor format")
    }
    
    id, err := uuid.Parse(parts[0])
    if err != nil {
        return uuid.Nil, time.Time{}, err
    }
    
    timestamp, err := strconv.ParseInt(parts[1], 10, 64)
    if err != nil {
        return uuid.Nil, time.Time{}, err
    }
    
    return id, time.Unix(timestamp, 0), nil
}
```

### Database Query Implementation
```go
// Optimized cursor-based queries
func (r *ProductRepository) GetProductsWithCursor(
    tenantID uuid.UUID, 
    pagination CursorPagination, 
    filter *ProductFilter,
) (*ProductConnection, error) {
    query := r.db.Model(&Product{}).Where("tenant_id = ?", tenantID)
    
    // Apply filters
    if filter != nil {
        query = r.applyProductFilters(query, filter)
    }
    
    // Handle cursor pagination
    if pagination.After != nil {
        afterID, afterTime, err := DecodeCursor(*pagination.After)
        if err != nil {
            return nil, err
        }
        
        // Use composite cursor for stable pagination
        query = query.Where("(created_at > ? OR (created_at = ? AND id > ?))", 
            afterTime, afterTime, afterID)
    }
    
    if pagination.Before != nil {
        beforeID, beforeTime, err := DecodeCursor(*pagination.Before)
        if err != nil {
            return nil, err
        }
        
        query = query.Where("(created_at < ? OR (created_at = ? AND id < ?))", 
            beforeTime, beforeTime, beforeID)
    }
    
    // Determine limit
    limit := 20 // default
    if pagination.First != nil {
        limit = *pagination.First
    } else if pagination.Last != nil {
        limit = *pagination.Last
    }
    
    // Add safety limit
    if limit > 100 {
        limit = 100
    }
    
    // Execute query with one extra record to determine hasNextPage
    var products []Product
    query = query.Order("created_at ASC, id ASC").Limit(limit + 1)
    
    err := query.Find(&products).Error
    if err != nil {
        return nil, err
    }
    
    // Process results
    hasNextPage := len(products) > limit
    if hasNextPage {
        products = products[:limit] // Remove extra record
    }
    
    // Build edges
    edges := make([]*ProductEdge, len(products))
    for i, product := range products {
        edges[i] = &ProductEdge{
            Node:   &product,
            Cursor: EncodeCursor(product.ID, product.CreatedAt),
        }
    }
    
    // Build page info
    pageInfo := &PageInfo{
        HasNextPage:     hasNextPage,
        HasPreviousPage: pagination.After != nil,
    }
    
    if len(edges) > 0 {
        startCursor := edges[0].Cursor
        endCursor := edges[len(edges)-1].Cursor
        pageInfo.StartCursor = &startCursor
        pageInfo.EndCursor = &endCursor
    }
    
    // Get total count (optional, can be expensive)
    var totalCount int64
    countQuery := r.db.Model(&Product{}).Where("tenant_id = ?", tenantID)
    if filter != nil {
        countQuery = r.applyProductFilters(countQuery, filter)
    }
    countQuery.Count(&totalCount)
    
    return &ProductConnection{
        Edges:      edges,
        PageInfo:   pageInfo,
        TotalCount: int(totalCount),
    }, nil
}
```

### REST API Cursor Pagination
```go
// REST endpoint with cursor pagination
func (h *ProductHandler) GetProducts(c *gin.Context) {
    tenantID := GetTenantID(c)
    
    // Parse pagination parameters
    pagination := CursorPagination{
        First: parseIntQuery(c, "first"),
        After: parseStringQuery(c, "after"),
        Last:  parseIntQuery(c, "last"),
        Before: parseStringQuery(c, "before"),
    }
    
    // Parse filters
    filter := parseProductFilters(c)
    
    // Get products
    connection, err := h.productRepo.GetProductsWithCursor(
        tenantID, pagination, filter)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    // Return paginated response
    c.JSON(200, gin.H{
        "data":       connection.Edges,
        "pageInfo":   connection.PageInfo,
        "totalCount": connection.TotalCount,
    })
}
```

## 4. Request Deduplication for High-Traffic Endpoints

### Deduplication Middleware Implementation
```go
// Request deduplication middleware
type RequestDeduplicator struct {
    cache   *redis.Client
    ttl     time.Duration
    keyFunc func(*gin.Context) string
}

func NewRequestDeduplicator(cache *redis.Client, ttl time.Duration) *RequestDeduplicator {
    return &RequestDeduplicator{
        cache: cache,
        ttl:   ttl,
        keyFunc: defaultKeyFunc,
    }
}

func (rd *RequestDeduplicator) Middleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        // Only deduplicate safe methods for now
        if c.Request.Method != "GET" {
            c.Next()
            return
        }
        
        // Generate deduplication key
        key := rd.keyFunc(c)
        
        // Check if request is in progress
        lockKey := fmt.Sprintf("lock:%s", key)
        acquired, err := rd.cache.SetNX(c.Request.Context(), lockKey, "1", rd.ttl).Result()
        if err != nil {
            // If Redis fails, proceed without deduplication
            c.Next()
            return
        }
        
        if !acquired {
            // Request is already in progress, wait for result
            rd.waitForResult(c, key)
            return
        }
        
        // This request will handle the computation
        defer rd.cache.Del(c.Request.Context(), lockKey)
        
        // Capture response
        responseRecorder := &ResponseRecorder{
            ResponseWriter: c.Writer,
            body:          bytes.NewBuffer(nil),
            statusCode:    200,
            headers:       make(http.Header),
        }
        c.Writer = responseRecorder
        
        // Process request
        c.Next()
        
        // Store result for other waiting requests
        if responseRecorder.statusCode == 200 {
            result := CachedResponse{
                StatusCode: responseRecorder.statusCode,
                Headers:    responseRecorder.headers,
                Body:       responseRecorder.body.Bytes(),
                CreatedAt:  time.Now(),
            }
            
            resultJSON, _ := json.Marshal(result)
            rd.cache.Set(c.Request.Context(), key, resultJSON, rd.ttl)
        }
    })
}

func defaultKeyFunc(c *gin.Context) string {
    // Include tenant context, path, and query parameters
    tenantID := c.GetString("tenant_id")
    h := sha256.New()
    h.Write([]byte(fmt.Sprintf("%s:%s:%s", 
        tenantID, 
        c.Request.URL.Path,
        c.Request.URL.RawQuery,
    )))
    return fmt.Sprintf("req:%x", h.Sum(nil)[:12])
}

func (rd *RequestDeduplicator) waitForResult(c *gin.Context, key string) {
    timeout := time.After(30 * time.Second)
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case <-timeout:
            // Timeout waiting for result
            c.JSON(408, gin.H{"error": "Request timeout"})
            return
            
        case <-ticker.C:
            // Check if result is available
            resultJSON, err := rd.cache.Get(c.Request.Context(), key).Result()
            if err == redis.Nil {
                continue // Not ready yet
            } else if err != nil {
                // Redis error, return error
                c.JSON(500, gin.H{"error": "Internal server error"})
                return
            }
            
            // Parse and return cached result
            var result CachedResponse
            if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
                c.JSON(500, gin.H{"error": "Internal server error"})
                return
            }
            
            // Copy headers
            for key, values := range result.Headers {
                for _, value := range values {
                    c.Header(key, value)
                }
            }
            
            c.Data(result.StatusCode, "application/json", result.Body)
            return
        }
    }
}
```

### Smart Deduplication Strategies
```go
// Endpoint-specific deduplication configuration
type EndpointConfig struct {
    Path        string
    Methods     []string
    TTL         time.Duration
    KeyFunc     func(*gin.Context) string
    ShouldCache func(*gin.Context) bool
}

var DeduplicationConfig = []EndpointConfig{
    {
        Path:    "/api/v1/products",
        Methods: []string{"GET"},
        TTL:     5 * time.Minute,
        KeyFunc: func(c *gin.Context) string {
            return fmt.Sprintf("products:%s:%s", 
                c.GetString("tenant_id"),
                c.Request.URL.RawQuery)
        },
        ShouldCache: func(c *gin.Context) bool {
            // Only cache if no real-time filters
            return c.Query("real_time") != "true"
        },
    },
    {
        Path:    "/api/v1/analytics/dashboard",
        Methods: []string{"GET"},
        TTL:     10 * time.Minute,
        KeyFunc: func(c *gin.Context) string {
            return fmt.Sprintf("dashboard:%s:%s", 
                c.GetString("tenant_id"),
                c.Query("period"))
        },
    },
    {
        Path:    "/api/v1/search",
        Methods: []string{"GET"},
        TTL:     15 * time.Minute,
        KeyFunc: func(c *gin.Context) string {
            h := sha256.New()
            h.Write([]byte(fmt.Sprintf("%s:%s", 
                c.GetString("tenant_id"),
                c.Query("q"))))
            return fmt.Sprintf("search:%x", h.Sum(nil)[:8])
        },
    },
}
```

### Idempotency for Mutations
```go
// Idempotency key middleware for mutations
func IdempotencyMiddleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        // Only apply to mutations
        if c.Request.Method == "GET" || c.Request.Method == "HEAD" {
            c.Next()
            return
        }
        
        // Check for idempotency key
        idempotencyKey := c.GetHeader("Idempotency-Key")
        if idempotencyKey == "" {
            c.Next()
            return
        }
        
        tenantID := c.GetString("tenant_id")
        key := fmt.Sprintf("idempotent:%s:%s", tenantID, idempotencyKey)
        
        // Check if operation already completed
        resultJSON, err := rd.cache.Get(c.Request.Context(), key).Result()
        if err == nil {
            // Return cached result
            var result CachedResponse
            json.Unmarshal([]byte(resultJSON), &result)
            
            for k, v := range result.Headers {
                for _, value := range v {
                    c.Header(k, value)
                }
            }
            c.Header("Idempotency-Replayed", "true")
            c.Data(result.StatusCode, "application/json", result.Body)
            return
        }
        
        // Process request and cache result
        responseRecorder := &ResponseRecorder{
            ResponseWriter: c.Writer,
            body:          bytes.NewBuffer(nil),
            statusCode:    200,
            headers:       make(http.Header),
        }
        c.Writer = responseRecorder
        
        c.Next()
        
        // Cache successful operations
        if responseRecorder.statusCode >= 200 && responseRecorder.statusCode < 300 {
            result := CachedResponse{
                StatusCode: responseRecorder.statusCode,
                Headers:    responseRecorder.headers,
                Body:       responseRecorder.body.Bytes(),
                CreatedAt:  time.Now(),
            }
            
            resultJSON, _ := json.Marshal(result)
            // Cache for 24 hours
            rd.cache.Set(c.Request.Context(), key, resultJSON, 24*time.Hour)
        }
    })
}
```

## 5. Performance Monitoring & Analytics

### API Performance Metrics
```yaml
# Request Metrics
api_requests_total{method, endpoint, tenant_id, status}
api_request_duration_seconds{method, endpoint, tenant_id}
api_request_size_bytes{method, endpoint, tenant_id}
api_response_size_bytes{method, endpoint, tenant_id}

# GraphQL Metrics
graphql_queries_total{operation, tenant_id, complexity}
graphql_query_duration_seconds{operation, tenant_id}
graphql_resolver_duration_seconds{resolver, tenant_id}
graphql_errors_total{error_type, operation, tenant_id}

# SSE Metrics
sse_connections_active{tenant_id}
sse_events_sent_total{event_type, tenant_id}
sse_connection_duration_seconds{tenant_id}

# Deduplication Metrics
request_deduplication_hits_total{endpoint, tenant_id}
request_deduplication_misses_total{endpoint, tenant_id}
request_deduplication_wait_duration_seconds{endpoint, tenant_id}
```

## 6. Implementation Checklist

### Phase 1: GraphQL Foundation
- [ ] Set up GraphQL server with Go (gqlgen)
- [ ] Design core schemas (Product, Order, Analytics)
- [ ] Implement DataLoader pattern for N+1 prevention
- [ ] Add query complexity limiting
- [ ] Create GraphQL playground for development

### Phase 2: Real-time Capabilities
- [ ] Implement Server-Sent Events infrastructure
- [ ] Add GraphQL subscriptions support
- [ ] Create event broadcasting system
- [ ] Integrate real-time updates with business logic
- [ ] Add SSE connection monitoring

### Phase 3: Advanced Pagination
- [ ] Implement cursor-based pagination
- [ ] Add cursor encoding/decoding utilities
- [ ] Update all list endpoints with cursor support
- [ ] Add pagination performance monitoring
- [ ] Create pagination UI components

### Phase 4: Request Optimization
- [ ] Implement request deduplication middleware
- [ ] Add idempotency key support for mutations
- [ ] Configure endpoint-specific deduplication policies
- [ ] Add request deduplication monitoring
- [ ] Optimize cache key strategies

### Phase 5: Performance & Monitoring
- [ ] Add comprehensive API metrics
- [ ] Create performance monitoring dashboards
- [ ] Implement automatic performance alerts
- [ ] Add API response time SLAs
- [ ] Create API usage analytics

This comprehensive API design improvement strategy will significantly enhance performance, user experience, and developer productivity while maintaining scalability and reliability.