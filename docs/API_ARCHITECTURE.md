# API Architecture & Design

## Overview
Technical architecture for the API layer, covering REST and WebSocket real-time capabilities to optimize performance and developer experience for the multi-tenant e-commerce platform.

> **Note**: For overall system architecture and technology decisions, see [ARCHITECTURE.md](./ARCHITECTURE.md)

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    CLIENT APPLICATIONS                     │
│    (Next.js Storefront + React Dashboard)                  │
└─────────────────────┬───────────────────┬───────────────────┘
                      │                   │
              ┌─────▼─────┐       ┌─────▼─────┐
              │    REST   │       │ WebSocket │
              │    API    │       │Real-time  │
              └─────┬─────┘       └─────┬─────┘
                    └─────────┬─────────┘
                              │
┌─────────────────────────────▼───────────────────────────────┐
│                   API GATEWAY LAYER                        │
│  ┌─────────────┬─────────────┬─────────────┬─────────────┐  │
│  │   Request   │   Rate      │   Auth &    │  Response   │  │
│  │ Deduplication│  Limiting  │ Validation  │   Caching   │  │
│  └─────────────┴─────────────┴─────────────┴─────────────┘  │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│               APPLICATION SERVICES                          │
│              (Go Modular Monolith)                          │
│  ┌─────────────┬─────────────┬─────────────┬─────────────┐  │
│  │   Product   │   Order     │   User      │   Tenant    │  │
│  │   Module    │   Module    │   Module    │   Module    │  │
│  └─────────────┴─────────────┴─────────────┴─────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## 1. REST API Design for Dashboard Operations

### API Endpoint Structure

#### Product Management Endpoints
```go
// Product CRUD Operations
GET    /api/v1/products              // List products with pagination and filters
POST   /api/v1/products              // Create new product
GET    /api/v1/products/{id}         // Get single product with full details
PUT    /api/v1/products/{id}         // Update product
DELETE /api/v1/products/{id}         // Delete product

// Product Variants
GET    /api/v1/products/{id}/variants       // List product variants
POST   /api/v1/products/{id}/variants       // Create variant
PUT    /api/v1/products/{id}/variants/{vid} // Update variant
DELETE /api/v1/products/{id}/variants/{vid} // Delete variant

// Bulk Operations
POST   /api/v1/products/bulk-import        // Bulk import products
PUT    /api/v1/products/bulk-update        // Bulk update products
DELETE /api/v1/products/bulk-delete        // Bulk delete products

```

#### REST API Response Models
```go
// Product Response Model
type ProductResponse struct {
    ID              string               `json:"id"`
    TenantID        string               `json:"tenant_id"`
    Name            string               `json:"name"`
    Description     string               `json:"description,omitempty"`
    SKU             string               `json:"sku"`
    Price           float64              `json:"price"`
    ComparePrice    float64              `json:"compare_price,omitempty"`
    Status          string               `json:"status"` // draft, active, archived
    Inventory       ProductInventory     `json:"inventory"`
    Images          []ProductImage       `json:"images"`
    Variants        []ProductVariant     `json:"variants"`
    Categories      []string             `json:"categories"`
    Tags            []string             `json:"tags"`
    SEO             ProductSEO           `json:"seo"`
    CreatedAt       time.Time            `json:"created_at"`
    UpdatedAt       time.Time            `json:"updated_at"`
    
    // Computed fields
    IsInStock       bool                 `json:"is_in_stock"`
    TotalSales      int                  `json:"total_sales"`
    ConversionRate  float64              `json:"conversion_rate"`
    ProfitMargin    float64              `json:"profit_margin"`
}

type ProductInventory struct {
    TrackQuantity     bool  `json:"track_quantity"`
    Quantity          int   `json:"quantity"`
    ReservedQuantity  int   `json:"reserved_quantity"`
    AvailableQuantity int   `json:"available_quantity"` // computed
    LowStockAlert     int   `json:"low_stock_alert"`
    AllowBackorder    bool  `json:"allow_backorder"`
}

type ProductImage struct {
    ID       string `json:"id"`
    URL      string `json:"url"`
    AltText  string `json:"alt_text,omitempty"`
    Position int    `json:"position"`
    Width    int    `json:"width,omitempty"`
    Height   int    `json:"height,omitempty"`
}
```

#### Query Parameters for Filtering
```go
// GET /api/v1/products?status=active&category=electronics&sort=name&order=asc&limit=20&after=cursor123

type ProductListQuery struct {
    // Filters
    Status       []string `query:"status"`          // draft, active, archived
    Categories   []string `query:"category"`        // category IDs or names
    MinPrice     float64  `query:"min_price"`       // minimum price
    MaxPrice     float64  `query:"max_price"`       // maximum price
    InStock      bool     `query:"in_stock"`        // only in-stock items
    HasVariants  bool     `query:"has_variants"`    // products with variants
    Tags         []string `query:"tags"`            // product tags
    CreatedAfter string   `query:"created_after"`   // ISO date
    CreatedBefore string  `query:"created_before"`  // ISO date
    Search       string   `query:"search"`          // search in name/description
    
    // Sorting
    Sort         string   `query:"sort"`            // name, price, created_at, updated_at
    Order        string   `query:"order"`           // asc, desc
    
    // Pagination
    Limit        int      `query:"limit"`           // default 20, max 100
    After        string   `query:"after"`           // cursor for pagination
    Before       string   `query:"before"`          // cursor for pagination
}
```

#### Dashboard Analytics Endpoints
```go
// Analytics Endpoints
GET    /api/v1/analytics/dashboard         // Main dashboard stats
GET    /api/v1/analytics/revenue           // Revenue breakdown
GET    /api/v1/analytics/orders            // Order analytics  
GET    /api/v1/analytics/products          // Product performance
GET    /api/v1/analytics/customers         // Customer insights
GET    /api/v1/analytics/reports/sales     // Sales reports
GET    /api/v1/analytics/reports/inventory // Inventory reports
```

#### Analytics Response Models
```go
// Dashboard Stats Response
type DashboardStats struct {
    Revenue    RevenueStats    `json:"revenue"`
    Orders     OrderStats      `json:"orders"`
    Products   ProductStats    `json:"products"`
    Customers  CustomerStats   `json:"customers"`
    Traffic    TrafficStats    `json:"traffic"`
}

type RevenueStats struct {
    Total             float64                `json:"total"`
    Change            float64                `json:"change"` // percentage change
    Trend             []DataPoint            `json:"trend"`
    ByPaymentMethod   []PaymentMethodStats   `json:"by_payment_method"`
    ByCategory        []CategoryStats        `json:"by_category"`
    TopProducts       []ProductRevenueStats  `json:"top_products"`
}

type OrderStats struct {
    Total            int              `json:"total"`
    Pending          int              `json:"pending"`
    Completed        int              `json:"completed"`
    Cancelled        int              `json:"cancelled"`
    AverageValue     float64          `json:"average_order_value"`
    Trend            []DataPoint      `json:"trend"`
    FulfillmentStats []FulfillmentStat `json:"fulfillment_status"`
}

type DataPoint struct {
    Date  time.Time `json:"date"`
    Value float64   `json:"value"`
}

```

### REST API Handler Implementation
```go
// Product Handler with repository pattern
type ProductHandler struct {
    productService *ProductService
    imageService   *ImageService
    categoryService *CategoryService
}

// Get single product by ID
func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    productID := c.Params("id")
    
    product, err := h.productService.GetProductByID(tenantID, productID)
    if err != nil {
        if errors.Is(err, ErrProductNotFound) {
            return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
        }
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(product)
}

// List products with filtering and pagination
func (h *ProductHandler) ListProducts(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    
    // Parse query parameters
    var query ProductListQuery
    if err := c.QueryParser(&query); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid query parameters"})
    }
    
    // Apply defaults
    if query.Limit == 0 || query.Limit > 100 {
        query.Limit = 20
    }
    if query.Sort == "" {
        query.Sort = "created_at"
    }
    if query.Order == "" {
        query.Order = "desc"
    }
    
    products, pageInfo, err := h.productService.ListProducts(tenantID, query)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(fiber.Map{
        "data":      products,
        "page_info": pageInfo,
        "total":     len(products),
    })
}

// Dashboard analytics handler
func (h *AnalyticsHandler) GetDashboardStats(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    period := c.Query("period", "30d") // default 30 days
    
    // Run analytics queries in parallel
    var wg sync.WaitGroup
    var revenue RevenueStats
    var orders OrderStats
    var mu sync.Mutex
    
    // Revenue stats
    wg.Add(1)
    go func() {
        defer wg.Done()
        stats, err := h.analyticsService.GetRevenueStats(tenantID, period)
        if err == nil {
            mu.Lock()
            revenue = stats
            mu.Unlock()
        }
    }()
    
    // Order stats
    wg.Add(1)
    go func() {
        defer wg.Done()
        stats, err := h.analyticsService.GetOrderStats(tenantID, period)
        if err == nil {
            mu.Lock()
            orders = stats
            mu.Unlock()
        }
    }()
    
    wg.Wait()
    
    dashboardStats := DashboardStats{
        Revenue: revenue,
        Orders:  orders,
    }
    
    return c.JSON(dashboardStats)
}
```

### REST API Performance Optimizations

#### Request Validation Middleware
```go
// Validate request parameters and limits
func RequestValidationMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Limit query parameter counts
        queryParams := c.Context().QueryArgs()
        if queryParams.Len() > 50 {
            return c.Status(400).JSON(fiber.Map{
                "error": "Too many query parameters",
                "max_allowed": 50,
            })
        }
        
        // Validate pagination limits
        if limit := c.QueryInt("limit", 20); limit > 100 {
            return c.Status(400).JSON(fiber.Map{
                "error": "Limit too large",
                "max_limit": 100,
            })
        }
        
        return c.Next()
    }
}

// Request timeout middleware
func RequestTimeoutMiddleware(timeout time.Duration) fiber.Handler {
    return func(c *fiber.Ctx) error {
        ctx, cancel := context.WithTimeout(c.Context(), timeout)
        defer cancel()
        
        c.SetUserContext(ctx)
        return c.Next()
    }
}
```

#### API Configuration
```yaml
# REST API Configuration
max_request_size: 10MB
request_timeout: 30s
max_query_params: 50
enable_cors: true
enable_compression: true

# Rate Limiting per Tenant
rate_limits:
  requests_per_minute: 1000
  burst_capacity: 100
  dashboard_requests_per_minute: 500
```

## 2. WebSocket for Real-time Updates

### WebSocket Implementation Architecture
```go
// WebSocket Connection Manager
type WSManager struct {
    connections map[string]map[string]*WSConnection
    mutex       sync.RWMutex
    eventBus    *EventBus
    upgrader    websocket.Upgrader
}

type WSConnection struct {
    ID       string
    TenantID string
    UserID   string
    Conn     *websocket.Conn
    Send     chan WSMessage
    Done     chan bool
}

type WSMessage struct {
    Type      string      `json:"type"`
    Data      interface{} `json:"data"`
    Timestamp time.Time   `json:"timestamp"`
}

// WebSocket upgrade and connection handler
func (ws *WSManager) HandleWebSocket(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    userID := c.Locals("user_id").(string)
    
    // Upgrade HTTP connection to WebSocket
    conn, err := ws.upgrader.Upgrade(c.Response(), c.Request(), nil)
    if err != nil {
        return err
    }
    defer conn.Close()
    
    // Create WebSocket connection
    wsConn := &WSConnection{
        ID:       uuid.New().String(),
        TenantID: tenantID,
        UserID:   userID,
        Conn:     conn,
        Send:     make(chan WSMessage, 256),
        Done:     make(chan bool),
    }
    
    // Register connection
    ws.addConnection(wsConn)
    defer ws.removeConnection(wsConn)
    
    // Start goroutines for reading and writing
    go ws.handleRead(wsConn)
    go ws.handleWrite(wsConn)
    
    // Send welcome message
    wsConn.Send <- WSMessage{
        Type: "connected",
        Data: map[string]string{
            "connection_id": wsConn.ID,
            "tenant_id":     wsConn.TenantID,
        },
        Timestamp: time.Now(),
    }
    
    // Wait for connection to close
    <-wsConn.Done
    
    return nil
}

// Handle reading messages from client
func (ws *WSManager) handleRead(conn *WSConnection) {
    defer func() {
        conn.Done <- true
    }()
    
    conn.Conn.SetReadLimit(512)
    conn.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
    conn.Conn.SetPongHandler(func(string) error {
        conn.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        return nil
    })
    
    for {
        var message WSMessage
        err := conn.Conn.ReadJSON(&message)
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("WebSocket error: %v", err)
            }
            break
        }
        
        // Handle incoming messages (ping, subscribe, etc.)
        ws.handleMessage(conn, message)
    }
}

// Handle writing messages to client
func (ws *WSManager) handleWrite(conn *WSConnection) {
    ticker := time.NewTicker(54 * time.Second)
    defer func() {
        ticker.Stop()
        conn.Conn.Close()
    }()
    
    for {
        select {
        case message, ok := <-conn.Send:
            conn.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if !ok {
                conn.Conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            
            if err := conn.Conn.WriteJSON(message); err != nil {
                return
            }
            
        case <-ticker.C:
            conn.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if err := conn.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
            
        case <-conn.Done:
            return
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

### WebSocket Event Broadcasting
```go
// Event broadcasting service
func (ws *WSManager) BroadcastToTenant(tenantID string, message WSMessage) {
    ws.mutex.RLock()
    connections, exists := ws.connections[tenantID]
    ws.mutex.RUnlock()
    
    if !exists {
        return
    }
    
    for _, conn := range connections {
        select {
        case conn.Send <- message:
            // Message sent successfully
        case <-time.After(1 * time.Second):
            // Connection might be blocked, remove it
            close(conn.Send)
            ws.removeConnection(conn)
        }
    }
}

// Broadcast to specific user
func (ws *WSManager) BroadcastToUser(tenantID, userID string, message WSMessage) {
    ws.mutex.RLock()
    tenantConnections, exists := ws.connections[tenantID]
    ws.mutex.RUnlock()
    
    if !exists {
        return
    }
    
    for _, conn := range tenantConnections {
        if conn.UserID == userID {
            select {
            case conn.Send <- message:
                // Message sent successfully
            default:
                // Connection is blocked, remove it
                close(conn.Send)
                ws.removeConnection(conn)
            }
            break
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
    
    // Broadcast real-time update via WebSocket
    changes := s.detectChanges(oldProduct, product)
    if len(changes) > 0 {
        message := WSMessage{
            Type: "product_updated",
            Data: map[string]interface{}{
                "product_id": product.ID,
                "tenant_id":  product.TenantID,
                "changes":    changes,
                "updated_fields": s.getUpdatedFields(oldProduct, product, changes),
            },
            Timestamp: time.Now(),
        }
        
        s.wsManager.BroadcastToTenant(product.TenantID.String(), message)
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
func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
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
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    // Return paginated response
    return c.JSON(fiber.Map{
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

func (rd *RequestDeduplicator) Middleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Only deduplicate safe methods for now
        if c.Method() != "GET" {
            return c.Next()
        }
        
        // Generate deduplication key
        key := rd.keyFunc(c)
        
        // Check if request is in progress
        lockKey := fmt.Sprintf("lock:%s", key)
        acquired, err := rd.cache.SetNX(c.Context(), lockKey, "1", rd.ttl).Result()
        if err != nil {
            // If Redis fails, proceed without deduplication
            return c.Next()
        }
        
        if !acquired {
            // Request is already in progress, wait for result
            return rd.waitForResult(c, key)
        }
        
        // This request will handle the computation
        defer rd.cache.Del(c.Context(), lockKey)
        
        // Process request and capture response
        err = c.Next()
        if err != nil {
            return err
        }
        
        // Store result for other waiting requests if successful
        if c.Response().StatusCode() == 200 {
            result := CachedResponse{
                StatusCode: c.Response().StatusCode(),
                Headers:    make(http.Header),
                Body:       c.Response().Body(),
                CreatedAt:  time.Now(),
            }
            
            resultJSON, _ := json.Marshal(result)
            rd.cache.Set(c.Context(), key, resultJSON, rd.ttl)
        }
        
        return nil
    }
}

func defaultKeyFunc(c *fiber.Ctx) string {
    // Include tenant context, path, and query parameters
    tenantID := c.Locals("tenant_id").(string)
    h := sha256.New()
    h.Write([]byte(fmt.Sprintf("%s:%s:%s", 
        tenantID, 
        c.Path(),
        c.Request().URI().QueryArgs().String(),
    )))
    return fmt.Sprintf("req:%x", h.Sum(nil)[:12])
}

func (rd *RequestDeduplicator) waitForResult(c *fiber.Ctx, key string) error {
    timeout := time.After(30 * time.Second)
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case <-timeout:
            // Timeout waiting for result
            return c.Status(408).JSON(fiber.Map{"error": "Request timeout"})
            
        case <-ticker.C:
            // Check if result is available
            resultJSON, err := rd.cache.Get(c.Context(), key).Result()
            if err == redis.Nil {
                continue // Not ready yet
            } else if err != nil {
                // Redis error, return error
                return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
            }
            
            // Parse and return cached result
            var result CachedResponse
            if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
                return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
            }
            
            // Copy headers
            for key, values := range result.Headers {
                for _, value := range values {
                    c.Set(key, value)
                }
            }
            
            return c.Status(result.StatusCode).Send(result.Body)
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
    KeyFunc     func(*fiber.Ctx) string
    ShouldCache func(*fiber.Ctx) bool
}

var DeduplicationConfig = []EndpointConfig{
    {
        Path:    "/api/v1/products",
        Methods: []string{"GET"},
        TTL:     5 * time.Minute,
        KeyFunc: func(c *fiber.Ctx) string {
            return fmt.Sprintf("products:%s:%s", 
                c.Locals("tenant_id").(string),
                c.Request().URI().QueryArgs().String())
        },
        ShouldCache: func(c *fiber.Ctx) bool {
            // Only cache if no real-time filters
            return c.Query("real_time") != "true"
        },
    },
    {
        Path:    "/api/v1/analytics/dashboard",
        Methods: []string{"GET"},
        TTL:     10 * time.Minute,
        KeyFunc: func(c *fiber.Ctx) string {
            return fmt.Sprintf("dashboard:%s:%s", 
                c.Locals("tenant_id").(string),
                c.Query("period"))
        },
    },
    {
        Path:    "/api/v1/search",
        Methods: []string{"GET"},
        TTL:     15 * time.Minute,
        KeyFunc: func(c *fiber.Ctx) string {
            h := sha256.New()
            h.Write([]byte(fmt.Sprintf("%s:%s", 
                c.Locals("tenant_id").(string),
                c.Query("q"))))
            return fmt.Sprintf("search:%x", h.Sum(nil)[:8])
        },
    },
}
```

### Idempotency for Mutations
```go
// Idempotency key middleware for mutations
func IdempotencyMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Only apply to mutations
        if c.Method() == "GET" || c.Method() == "HEAD" {
            return c.Next()
        }
        
        // Check for idempotency key
        idempotencyKey := c.Get("Idempotency-Key")
        if idempotencyKey == "" {
            return c.Next()
        }
        
        tenantID := c.Locals("tenant_id").(string)
        key := fmt.Sprintf("idempotent:%s:%s", tenantID, idempotencyKey)
        
        // Check if operation already completed
        resultJSON, err := rd.cache.Get(c.Context(), key).Result()
        if err == nil {
            // Return cached result
            var result CachedResponse
            json.Unmarshal([]byte(resultJSON), &result)
            
            for k, v := range result.Headers {
                for _, value := range v {
                    c.Set(k, value)
                }
            }
            c.Set("Idempotency-Replayed", "true")
            return c.Status(result.StatusCode).Send(result.Body)
        }
        
        // Process request
        err = c.Next()
        if err != nil {
            return err
        }
        
        // Cache successful operations
        if c.Response().StatusCode() >= 200 && c.Response().StatusCode() < 300 {
            result := CachedResponse{
                StatusCode: c.Response().StatusCode(),
                Headers:    make(http.Header),
                Body:       c.Response().Body(),
                CreatedAt:  time.Now(),
            }
            
            resultJSON, _ := json.Marshal(result)
            // Cache for 24 hours
            rd.cache.Set(c.Context(), key, resultJSON, 24*time.Hour)
        }
        
        return nil
    }
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

# WebSocket Metrics
websocket_connections_active{tenant_id}
websocket_messages_sent_total{message_type, tenant_id}
websocket_connection_duration_seconds{tenant_id}
websocket_errors_total{error_type, tenant_id}

# Deduplication Metrics
request_deduplication_hits_total{endpoint, tenant_id}
request_deduplication_misses_total{endpoint, tenant_id}
request_deduplication_wait_duration_seconds{endpoint, tenant_id}
```

## 6. Implementation Checklist

### Phase 1: REST API Foundation
- [ ] Set up REST API with Go Fiber
- [ ] Design endpoint structure and response models
- [ ] Implement request validation and filtering
- [ ] Add cursor-based pagination
- [ ] Create API documentation with OpenAPI

### Phase 2: Real-time Capabilities
- [ ] Implement WebSocket infrastructure
- [ ] Add WebSocket connection management
- [ ] Create message broadcasting system
- [ ] Integrate real-time updates with business logic
- [ ] Add WebSocket connection monitoring

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