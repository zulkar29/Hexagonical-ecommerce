# Testing Strategy

Comprehensive testing approach for the e-commerce SaaS platform.

## Testing Philosophy

### Testing Pyramid
```
    /\
   /  \    E2E Tests (Few)
  /____\   
 /      \  Integration Tests (Some)
/________\ Unit Tests (Many)
```

### Test-Driven Development (TDD)
1. **Red**: Write failing test
2. **Green**: Write minimal code to pass
3. **Refactor**: Improve code while keeping tests green

## Backend Testing (Go)

### Unit Tests
```go
// Example: Product service unit test
func TestProductService_CreateProduct(t *testing.T) {
    // Arrange
    mockRepo := &mocks.ProductRepository{}
    service := NewProductService(mockRepo)
    
    product := &entities.Product{
        Name:  "Test Product",
        Price: 29.99,
        SKU:   "TEST-001",
    }
    
    mockRepo.On("Create", mock.Anything, product).Return(nil)
    
    // Act
    err := service.CreateProduct(context.Background(), product)
    
    // Assert
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

### Integration Tests
```go
// Example: Database integration test
func TestProductRepository_Integration(t *testing.T) {
    // Skip if not running integration tests
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    repo := database.NewProductRepository(db)
    
    // Test create and retrieve
    product := &entities.Product{
        Name:  "Integration Test Product",
        Price: 19.99,
    }
    
    err := repo.Create(context.Background(), product)
    require.NoError(t, err)
    
    retrieved, err := repo.GetByID(context.Background(), product.ID)
    require.NoError(t, err)
    assert.Equal(t, product.Name, retrieved.Name)
}
```

### API Tests
```go
// Example: HTTP handler test
func TestProductHandler_CreateProduct(t *testing.T) {
    app := fiber.New()
    
    // Mock service
    mockService := &mocks.ProductService{}
    handler := http.NewProductHandler(mockService)
    
    app.Post("/products", handler.CreateProduct)
    
    // Prepare request
    reqBody := `{"name":"Test Product","price":29.99}`
    req := httptest.NewRequest("POST", "/products", strings.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")
    
    // Execute
    resp, err := app.Test(req)
    require.NoError(t, err)
    
    // Assert
    assert.Equal(t, 201, resp.StatusCode)
}
```

### Test Configuration
```go
// testutil/setup.go
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)
    
    // Run migrations
    err = db.AutoMigrate(&entities.Product{}, &entities.Order{})
    require.NoError(t, err)
    
    return db
}
```

## Frontend Testing

### Unit Tests (React/Next.js)
```javascript
// Example: Component unit test
import { render, screen } from '@testing-library/react'
import { ProductCard } from '../ProductCard'

describe('ProductCard', () => {
  const mockProduct = {
    id: '1',
    name: 'Test Product',
    price: 29.99,
    image: '/test-image.jpg'
  }

  it('renders product information', () => {
    render(<ProductCard product={mockProduct} />)
    
    expect(screen.getByText('Test Product')).toBeInTheDocument()
    expect(screen.getByText('৳3,299')).toBeInTheDocument()
  })

  it('calls onAddToCart when button is clicked', () => {
    const mockOnAddToCart = jest.fn()
    
    render(
      <ProductCard 
        product={mockProduct} 
        onAddToCart={mockOnAddToCart} 
      />
    )
    
    const addButton = screen.getByRole('button', { name: /add to cart/i })
    fireEvent.click(addButton)
    
    expect(mockOnAddToCart).toHaveBeenCalledWith(mockProduct.id)
  })
})
```

### Integration Tests (React)
```javascript
// Example: Component integration test
import { render, screen, waitFor } from '@testing-library/react'
import { ProductList } from '../ProductList'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { rest } from 'msw'
import { setupServer } from 'msw/node'

const server = setupServer(
  rest.get('/api/products', (req, res, ctx) => {
    return res(ctx.json({
      data: [
        { id: '1', name: 'Product 1', price: 19.99 },
        { id: '2', name: 'Product 2', price: 29.99 }
      ]
    }))
  })
)

describe('ProductList Integration', () => {
  beforeAll(() => server.listen())
  afterEach(() => server.resetHandlers())
  afterAll(() => server.close())

  it('fetches and displays products', async () => {
    const queryClient = new QueryClient()
    
    render(
      <QueryClientProvider client={queryClient}>
        <ProductList />
      </QueryClientProvider>
    )

    await waitFor(() => {
      expect(screen.getByText('Product 1')).toBeInTheDocument()
      expect(screen.getByText('Product 2')).toBeInTheDocument()
    })
  })
})
```

### E2E Tests (Playwright)
```javascript
// tests/e2e/checkout.spec.js
import { test, expect } from '@playwright/test'

test.describe('Checkout Flow', () => {
  test('complete purchase flow', async ({ page }) => {
    // Navigate to product page
    await page.goto('/products/test-product')
    
    // Add to cart
    await page.click('[data-testid="add-to-cart"]')
    await expect(page.locator('[data-testid="cart-count"]')).toHaveText('1')
    
    // Go to cart
    await page.click('[data-testid="cart-button"]')
    await expect(page).toHaveURL('/cart')
    
    // Proceed to checkout
    await page.click('[data-testid="checkout-button"]')
    await expect(page).toHaveURL('/checkout')
    
    // Fill checkout form
    await page.fill('[name="email"]', 'test@example.com')
    await page.fill('[name="firstName"]', 'John')
    await page.fill('[name="lastName"]', 'Doe')
    
    // Complete purchase
    await page.click('[data-testid="complete-purchase"]')
    await expect(page).toHaveURL('/order-confirmation')
  })
})
```

## Test Data Management

### Test Fixtures
```go
// fixtures/products.go
var TestProducts = []entities.Product{
    {
        ID:    uuid.New(),
        Name:  "Test T-Shirt",
        Price: 29.99,
        SKU:   "TSHIRT-001",
    },
    {
        ID:    uuid.New(), 
        Name:  "Test Jeans",
        Price: 59.99,
        SKU:   "JEANS-001",
    },
}
```

### Database Seeding
```go
// testutil/seed.go
func SeedTestData(db *gorm.DB) error {
    // Create test categories
    categories := fixtures.TestCategories
    for _, category := range categories {
        if err := db.Create(&category).Error; err != nil {
            return err
        }
    }
    
    // Create test products
    products := fixtures.TestProducts
    for _, product := range products {
        if err := db.Create(&product).Error; err != nil {
            return err
        }
    }
    
    return nil
}
```

## Test Environment Setup

### Docker Test Environment
```yaml
# docker-compose.test.yml
version: '3.8'
services:
  postgres-test:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: ecommerce_test
      POSTGRES_USER: test
      POSTGRES_PASSWORD: test
    ports:
      - "5433:5432"
  
  redis-test:
    image: redis:7-alpine
    ports:
      - "6380:6379"
```

### GitHub Actions CI/CD
```yaml
# .github/workflows/test.yml
name: Test Suite

on: [push, pull_request]

jobs:
  test-backend:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.21
      
      - name: Run tests
        run: |
          cd backend
          go test -v -race -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3

  test-frontend:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: 18
      
      - name: Install dependencies
        run: |
          cd storefront && npm ci
          cd ../dashboard && npm ci
      
      - name: Run tests
        run: |
          cd storefront && npm test -- --coverage
          cd ../dashboard && npm test -- --coverage

  e2e-tests:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: 18
      
      - name: Install Playwright
        run: npx playwright install
      
      - name: Start services
        run: make dev-up
      
      - name: Run E2E tests
        run: npx playwright test
```

## Performance Testing

### Load Testing (k6)
```javascript
// scripts/load-test.js
import http from 'k6/http'
import { check } from 'k6'

export let options = {
  stages: [
    { duration: '2m', target: 10 },   // Ramp up
    { duration: '5m', target: 50 },   // Stay at 50 users
    { duration: '2m', target: 0 },    // Ramp down
  ],
}

export default function () {
  // Test product listing
  let response = http.get('http://localhost:8080/api/v1/products')
  check(response, {
    'products endpoint status is 200': (r) => r.status === 200,
    'response time < 200ms': (r) => r.timings.duration < 200,
  })
  
  // Test product creation (authenticated)
  const token = getAuthToken()
  response = http.post('http://localhost:8080/api/v1/products', 
    JSON.stringify({
      name: 'Load Test Product',
      price: 29.99
    }),
    {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    }
  )
  
  check(response, {
    'create product status is 201': (r) => r.status === 201,
  })
}
```

## Test Reporting

### Coverage Reports
```bash
# Backend coverage
cd backend
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Frontend coverage  
cd storefront
npm test -- --coverage --coverageReporters=html

# Combined reporting
make test-coverage
```

### Test Metrics
- **Code Coverage**: Target 80%+ for critical paths
- **Test Execution Time**: < 5 minutes for full suite
- **Flaky Test Rate**: < 1%
- **Bug Escape Rate**: < 5%

## Quality Gates

### Pre-commit Hooks
```yaml
# .pre-commit-config.yaml
repos:
  - repo: local
    hooks:
      - id: go-test
        name: Go tests
        entry: make test-backend
        language: system
        pass_filenames: false
      
      - id: go-lint
        name: Go lint
        entry: golangci-lint run
        language: system
        files: \.go$
      
      - id: js-test
        name: JavaScript tests
        entry: npm test
        language: system
        files: \.(js|jsx|ts|tsx)$
```

### Deployment Gates
1. **Unit Tests**: Must pass 100%
2. **Integration Tests**: Must pass 100%
3. **Coverage**: Must meet minimum threshold
4. **Linting**: No critical issues
5. **Security Scan**: No high/critical vulnerabilities
6. **Performance Tests**: Response times within SLA