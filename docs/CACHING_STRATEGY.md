# Caching Strategy Enhancement

## Overview
Comprehensive caching strategy to optimize performance, reduce database load, and improve user experience across the multi-tenant e-commerce platform.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    CLIENT APPLICATIONS                     │
│              (Storefront + Dashboard)                      │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                     CDN LAYER                              │
│        (CloudFront/CloudFlare for Static Assets)          │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                  API GATEWAY                               │
│            (Response Caching Layer)                        │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                 APPLICATION LAYER                          │
│              (Go Modular Monolith)                         │
│  ┌─────────────────┬─────────────────┬─────────────────┐    │
│  │ Application     │ Query Result    │ Session         │    │
│  │ Cache (L1)      │ Cache (L2)      │ Cache (L3)      │    │
│  └─────────────────┴─────────────────┴─────────────────┘    │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                   REDIS CLUSTER                            │
│        (Centralized Caching Infrastructure)               │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│               POSTGRESQL DATABASE                          │
│              (Primary Data Store)                          │
└─────────────────────────────────────────────────────────────┘
```

## 1. Application-Level Caching with Redis

### Redis Cluster Configuration
```yaml
# Redis Cluster Setup
cluster_nodes:
  - redis-node-1:7000
  - redis-node-2:7000  
  - redis-node-3:7000
  - redis-node-4:7000
  - redis-node-5:7000
  - redis-node-6:7000

# Memory Configuration
maxmemory: 4gb
maxmemory-policy: allkeys-lru
timeout: 300

# Persistence
save: 900 1    # Save if at least 1 key changed in 900 seconds
appendonly: yes
appendfsync: everysec
```

### Cache Key Patterns
```yaml
# Tenant-scoped Keys
tenant_data: "tenant:{tenant_id}:data"
tenant_settings: "tenant:{tenant_id}:settings"
tenant_plan: "tenant:{tenant_id}:plan"

# Product Catalog Caching
product_detail: "product:{tenant_id}:{product_id}"
product_list: "products:{tenant_id}:{category_id}:{page}:{sort}"
product_search: "search:{tenant_id}:{query}:{filters}:{page}"
category_tree: "categories:{tenant_id}:tree"

# User & Session Caching  
user_session: "session:{session_id}"
user_profile: "user:{tenant_id}:{user_id}"
user_cart: "cart:{tenant_id}:{user_id}"
user_wishlist: "wishlist:{tenant_id}:{user_id}"

# Order & Inventory Caching
order_summary: "order:{tenant_id}:{order_id}"
inventory_level: "inventory:{tenant_id}:{product_id}"
stock_alerts: "alerts:stock:{tenant_id}"

# Analytics & Metrics
analytics_daily: "analytics:{tenant_id}:{date}"
metrics_hourly: "metrics:{tenant_id}:{hour}"
dashboard_stats: "dashboard:{tenant_id}:{period}"
```

### Cache Implementation Patterns

#### Cache-Aside Pattern (Lazy Loading)
```go
func (s *ProductService) GetProduct(tenantID, productID uuid.UUID) (*Product, error) {
    // Try cache first
    cacheKey := fmt.Sprintf("product:%s:%s", tenantID, productID)
    if cached, err := s.cache.Get(cacheKey); err == nil {
        var product Product
        json.Unmarshal([]byte(cached), &product)
        return &product, nil
    }
    
    // Cache miss - get from database
    product, err := s.repo.GetProduct(tenantID, productID)
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    productJSON, _ := json.Marshal(product)
    s.cache.Set(cacheKey, string(productJSON), 1*time.Hour)
    
    return product, nil
}
```

#### Write-Through Pattern (Critical Data)
```go
func (s *ProductService) UpdateProduct(tenantID uuid.UUID, product *Product) error {
    // Update database first
    err := s.repo.UpdateProduct(tenantID, product)
    if err != nil {
        return err
    }
    
    // Update cache immediately
    cacheKey := fmt.Sprintf("product:%s:%s", tenantID, product.ID)
    productJSON, _ := json.Marshal(product)
    s.cache.Set(cacheKey, string(productJSON), 1*time.Hour)
    
    // Invalidate related caches
    s.invalidateProductListCaches(tenantID)
    
    return nil
}
```

#### Write-Behind Pattern (Analytics Data)
```go
func (s *AnalyticsService) TrackEvent(tenantID uuid.UUID, event Event) {
    // Store in cache immediately
    cacheKey := fmt.Sprintf("events:%s:%s", tenantID, time.Now().Format("2006-01-02"))
    s.cache.LPush(cacheKey, event.ToJSON())
    
    // Batch write to database (background job)
    s.eventQueue.Enqueue(event)
}
```

### Cache TTL Strategy
```yaml
# Short-lived Caches (High Change Frequency)
inventory_levels: 5m
cart_contents: 15m
user_sessions: 30m
stock_alerts: 10m

# Medium-lived Caches (Moderate Change Frequency)  
product_details: 1h
category_trees: 2h
user_profiles: 1h
order_summaries: 4h

# Long-lived Caches (Low Change Frequency)
tenant_settings: 24h
product_catalogs: 6h
analytics_reports: 12h
dashboard_configs: 24h

# Static Data Caches
plan_features: 7d
country_lists: 30d
currency_rates: 6h
```

## 2. Query Result Caching for Product Catalogs

### Database Query Caching Strategy

#### Product Listing Queries
```go
// Cached product listing with faceted search
func (s *ProductService) GetProductList(params ProductListParams) (*ProductList, error) {
    // Generate cache key from parameters
    cacheKey := generateListCacheKey(params)
    
    // Try cache first
    if cached, found := s.cache.Get(cacheKey); found {
        return cached.(*ProductList), nil
    }
    
    // Execute database query
    products, total, facets, err := s.repo.GetProductsWithFacets(params)
    if err != nil {
        return nil, err
    }
    
    result := &ProductList{
        Products: products,
        Total: total,
        Facets: facets,
        Page: params.Page,
        PerPage: params.PerPage,
    }
    
    // Cache result for 1 hour
    s.cache.Set(cacheKey, result, 1*time.Hour)
    return result, nil
}

func generateListCacheKey(params ProductListParams) string {
    h := sha256.New()
    h.Write([]byte(fmt.Sprintf(
        "products:%s:cat:%s:q:%s:sort:%s:filters:%s:page:%d:limit:%d",
        params.TenantID,
        params.CategoryID, 
        params.Query,
        params.SortBy,
        params.FiltersHash(),
        params.Page,
        params.PerPage,
    )))
    return fmt.Sprintf("list:%x", h.Sum(nil)[:8])
}
```

#### Complex Aggregation Queries
```go
// Cached analytics queries
func (s *AnalyticsService) GetDashboardStats(tenantID uuid.UUID, period string) (*DashboardStats, error) {
    cacheKey := fmt.Sprintf("dashboard:%s:%s", tenantID, period)
    
    if cached, found := s.cache.Get(cacheKey); found {
        return cached.(*DashboardStats), nil
    }
    
    // Complex aggregation query
    stats := &DashboardStats{}
    
    // Revenue metrics
    stats.Revenue = s.calculateRevenue(tenantID, period)
    stats.Orders = s.calculateOrders(tenantID, period)
    stats.Visitors = s.calculateVisitors(tenantID, period)
    stats.ConversionRate = s.calculateConversion(tenantID, period)
    
    // Cache for 1 hour (dashboard data)
    s.cache.Set(cacheKey, stats, 1*time.Hour)
    return stats, nil
}
```

#### Search Result Caching
```go
// Elasticsearch/Product search with caching
func (s *SearchService) SearchProducts(tenantID uuid.UUID, query SearchQuery) (*SearchResults, error) {
    // Generate search cache key
    searchKey := fmt.Sprintf("search:%s:%s", tenantID, query.Hash())
    
    if cached, found := s.cache.Get(searchKey); found {
        // Update search analytics asynchronously
        go s.trackSearchHit(tenantID, query.Query)
        return cached.(*SearchResults), nil
    }
    
    // Execute search
    results, err := s.searchEngine.Search(tenantID, query)
    if err != nil {
        return nil, err
    }
    
    // Cache popular searches longer
    ttl := 30 * time.Minute
    if s.isPopularSearch(query.Query) {
        ttl = 2 * time.Hour
    }
    
    s.cache.Set(searchKey, results, ttl)
    
    // Track search analytics
    go s.trackSearch(tenantID, query.Query, results.Total)
    
    return results, nil
}
```

### Cache Invalidation Strategies

#### Smart Cache Invalidation
```go
// Product update invalidates related caches
func (s *ProductService) invalidateProductCaches(tenantID uuid.UUID, productID uuid.UUID) {
    product, _ := s.repo.GetProduct(tenantID, productID)
    
    // Invalidate direct product cache
    s.cache.Delete(fmt.Sprintf("product:%s:%s", tenantID, productID))
    
    // Invalidate category listings
    for _, categoryID := range product.CategoryIDs {
        pattern := fmt.Sprintf("products:%s:%s:*", tenantID, categoryID)
        s.cache.DeletePattern(pattern)
    }
    
    // Invalidate search results containing this product
    s.cache.DeletePattern(fmt.Sprintf("search:%s:*", tenantID))
    
    // Invalidate dashboard stats
    s.cache.DeletePattern(fmt.Sprintf("dashboard:%s:*", tenantID))
}
```

#### Time-based Cache Warming
```go
// Background cache warming
func (s *CacheService) WarmProductCaches(tenantID uuid.UUID) {
    // Popular products
    popularProducts := s.getPopularProducts(tenantID, 100)
    for _, product := range popularProducts {
        go s.productService.GetProduct(tenantID, product.ID)
    }
    
    // Category trees
    go s.categoryService.GetCategoryTree(tenantID)
    
    // Dashboard stats for common periods
    periods := []string{"today", "week", "month"}
    for _, period := range periods {
        go s.analyticsService.GetDashboardStats(tenantID, period)
    }
}
```

## 3. CDN Integration for Static Assets

### CloudFront Configuration
```yaml
# CDN Distribution Settings
origin_domain: api.yourstore.com
behaviors:
  # Static Assets (Images, CSS, JS)
  - path_pattern: "/static/*"
    origin: S3 Bucket
    cache_policy: CachingOptimized
    ttl: 31536000  # 1 year
    compress: true
    
  # Product Images  
  - path_pattern: "/images/*"
    origin: S3 Bucket
    cache_policy: CachingOptimized
    ttl: 86400     # 24 hours
    compress: true
    viewer_protocol_policy: redirect-to-https
    
  # API Responses (GET only)
  - path_pattern: "/api/v1/products/*"
    origin: ALB
    cache_policy: CachingDisabled  # Let API handle caching
    allowed_methods: [GET, HEAD, OPTIONS, PUT, POST, PATCH, DELETE]
    
  # Storefront Pages (SSR)
  - path_pattern: "/*"
    origin: Next.js Server
    cache_policy: CachingOptimizedForUncompressedObjects
    ttl: 3600      # 1 hour
    compress: true

# Custom Headers
response_headers:
  - header: "Cache-Control"
    value: "public, max-age=31536000, immutable"
    paths: ["/static/*", "/images/*"]
  - header: "Cache-Control"  
    value: "public, max-age=3600"
    paths: ["/api/v1/products/*"]
```

### Image Optimization Strategy
```yaml
# Multi-format Image Delivery
image_formats:
  webp: "Modern browsers (Chrome, Firefox, Edge)"
  avif: "Cutting-edge browsers (Chrome 85+)"
  jpeg: "Legacy browser fallback"
  
# Responsive Image Sizes
image_sizes:
  thumbnail: [150x150, 300x300]
  product_grid: [300x300, 600x600, 1200x1200]
  product_detail: [800x800, 1600x1600]
  hero_banner: [1920x600, 2560x800]

# CDN Image Transformations
transformations:
  - resize: "width=300,height=300,fit=crop"
  - quality: "auto"
  - format: "auto"
  - progressive: true
```

### Static Asset Optimization
```go
// Asset URL generation with CDN
func (s *AssetService) GetAssetURL(path string, options ...AssetOption) string {
    baseURL := s.config.CDNBaseURL
    
    // Add cache busting
    if s.config.Environment == "production" {
        path = s.addCacheBusting(path)
    }
    
    // Add transformations for images
    if s.isImage(path) {
        path = s.addImageTransformations(path, options...)
    }
    
    return fmt.Sprintf("%s%s", baseURL, path)
}

// Automatic WebP/AVIF format selection
func (s *AssetService) GetOptimizedImageURL(imagePath string, width, height int) string {
    params := url.Values{}
    params.Set("width", strconv.Itoa(width))
    params.Set("height", strconv.Itoa(height))
    params.Set("format", "auto")
    params.Set("quality", "auto")
    
    return fmt.Sprintf("%s%s?%s", s.config.CDNBaseURL, imagePath, params.Encode())
}
```

## 4. API Response Caching for Read-Heavy Endpoints

### HTTP Response Caching Middleware
```go
// Cache middleware for GET endpoints
func CacheMiddleware(cache CacheService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Only cache GET requests
        if c.Method() != "GET" {
            return c.Next()
        }
        
        // Generate cache key from request
        cacheKey := generateRequestCacheKey(c)
        
        // Try to get cached response
        if cached, found := cache.Get(cacheKey); found {
            cachedResponse := cached.(CachedResponse)
            
            // Set cached headers
            c.Set("Content-Type", cachedResponse.ContentType)
            c.Set("X-Cache", "HIT")
            c.Set("Cache-Control", cachedResponse.CacheControl)
            
            return c.Send(cachedResponse.Body)
        }
        
        // Continue to handler
        err := c.Next()
        if err != nil {
            return err
        }
        
        // Cache successful responses
        if c.Response().StatusCode() == 200 {
            response := CachedResponse{
                Body:         c.Response().Body(),
                ContentType:  string(c.Response().Header.ContentType()),
                CacheControl: determineCacheControl(c.Path()),
            }
            
            ttl := determineCacheTTL(c.Path())
            cache.Set(cacheKey, response, ttl)
            
            c.Set("X-Cache", "MISS")
        }
        
        return nil
    }
}

func generateRequestCacheKey(c *fiber.Ctx) string {
    // Include tenant context
    tenantID := c.Locals("tenant_id").(string)
    
    // Create key from path, query params, and tenant
    h := sha256.New()
    h.Write([]byte(fmt.Sprintf("%s:%s:%s", 
        tenantID, 
        c.Path(), 
        c.Request().URI().QueryString(),
    )))
    
    return fmt.Sprintf("api:%x", h.Sum(nil)[:12])
}
```

### Endpoint-Specific Caching Strategy
```yaml
# Product Endpoints
"/api/v1/products":
  cache_ttl: 5m
  cache_key: "tenant_id + query_params"
  invalidate_on: ["product_update", "inventory_change"]
  
"/api/v1/products/{id}":
  cache_ttl: 1h  
  cache_key: "tenant_id + product_id"
  invalidate_on: ["product_update"]

# Category Endpoints  
"/api/v1/categories":
  cache_ttl: 2h
  cache_key: "tenant_id"
  invalidate_on: ["category_update", "category_tree_change"]
  
"/api/v1/categories/{id}/products":
  cache_ttl: 30m
  cache_key: "tenant_id + category_id + query_params" 
  invalidate_on: ["product_update", "category_assignment"]

# Search Endpoints
"/api/v1/search":
  cache_ttl: 15m
  cache_key: "tenant_id + search_query + filters"
  invalidate_on: ["product_update", "search_index_update"]

# Analytics Endpoints (Read-heavy)
"/api/v1/analytics/dashboard":
  cache_ttl: 1h
  cache_key: "tenant_id + date_range"
  invalidate_on: ["order_update", "daily_analytics_job"]
  
"/api/v1/analytics/reports/*":
  cache_ttl: 6h
  cache_key: "tenant_id + report_type + parameters"
  invalidate_on: ["report_generation_job"]
```

### Conditional Caching with ETags
```go
// ETag-based conditional caching
func ETagMiddleware(cache CacheService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Generate ETag from tenant + resource version
        resourceKey := fmt.Sprintf("version:%s:%s", 
            c.Locals("tenant_id"), 
            c.Path())
            
        version, _ := cache.Get(resourceKey)
        etag := fmt.Sprintf(`"%s"`, version)
        
        // Check If-None-Match header
        if c.Get("If-None-Match") == etag {
            return c.SendStatus(304) // Not Modified
        }
        
        // Set ETag header
        c.Set("ETag", etag)
        
        return c.Next()
    }
}

// Update resource version on changes
func (s *ProductService) UpdateProduct(product *Product) error {
    err := s.repo.UpdateProduct(product)
    if err != nil {
        return err
    }
    
    // Update resource version for ETag
    versionKey := fmt.Sprintf("version:%s:/api/v1/products/%s", 
        product.TenantID, product.ID)
    newVersion := generateNewVersion()
    s.cache.Set(versionKey, newVersion, 24*time.Hour)
    
    return nil
}
```

## 5. Cache Performance Monitoring

### Redis Monitoring Metrics
```yaml
# Redis Performance Metrics
redis_commands_processed_total: "Total commands processed"
redis_keyspace_hits_total: "Cache hits"
redis_keyspace_misses_total: "Cache misses"  
redis_memory_used_bytes: "Memory usage"
redis_connected_clients: "Active connections"
redis_blocked_clients: "Blocked clients"
redis_evicted_keys_total: "Evicted keys"

# Custom Application Metrics
cache_hit_ratio{endpoint, tenant_id}: "Hit ratio per endpoint"
cache_response_time{operation}: "Cache operation latency"
cache_key_distribution{pattern}: "Key pattern distribution"
cache_ttl_effectiveness{key_type}: "TTL optimization metrics"
```

### Cache Performance Alerts
```yaml
- alert: LowCacheHitRatio
  expr: cache_hit_ratio < 0.7
  for: 5m
  severity: warning
  
- alert: RedisMemoryHigh
  expr: redis_memory_used_bytes / redis_memory_max_bytes > 0.9
  for: 2m
  severity: critical
  
- alert: RedisConnectionPoolExhaustion
  expr: redis_connected_clients > redis_max_clients * 0.8
  for: 1m
  severity: warning
  
- alert: CacheOperationLatency
  expr: cache_response_time > 100ms
  for: 3m
  severity: warning
```

## 6. Implementation Checklist

### Phase 1: Core Caching Infrastructure
- [ ] Set up Redis cluster with persistence
- [ ] Implement cache service wrapper with tenant isolation
- [ ] Add basic cache-aside pattern for product data
- [ ] Configure cache TTL policies
- [ ] Set up Redis monitoring and alerts

### Phase 2: Query Result Caching
- [ ] Implement query result caching for product catalogs
- [ ] Add search result caching
- [ ] Implement cache warming strategies
- [ ] Add smart cache invalidation logic
- [ ] Create cache performance monitoring

### Phase 3: CDN Integration
- [ ] Configure CloudFront distribution
- [ ] Set up S3 bucket for static assets
- [ ] Implement image optimization pipeline
- [ ] Add responsive image delivery
- [ ] Configure cache headers optimization

### Phase 4: API Response Caching
- [ ] Implement HTTP response caching middleware
- [ ] Add ETag-based conditional caching
- [ ] Configure endpoint-specific caching policies  
- [ ] Add cache invalidation webhooks
- [ ] Implement cache analytics dashboard

### Phase 5: Optimization & Monitoring
- [ ] Fine-tune cache TTL values based on usage patterns
- [ ] Implement cache preloading for popular content
- [ ] Add A/B testing for cache strategies
- [ ] Create cache performance reports
- [ ] Implement automated cache tuning

This comprehensive caching strategy will significantly improve application performance, reduce database load, and enhance user experience across the platform.