# Contributing Guidelines

Development standards and contribution guidelines for the e-commerce SaaS platform.

## Development Workflow

### Git Workflow (GitHub Flow)
1. **Create Branch**: Create feature branch from `main`
2. **Develop**: Make changes with atomic commits
3. **Test**: Ensure all tests pass locally
4. **Pull Request**: Create PR with detailed description
5. **Review**: Code review and approval
6. **Merge**: Squash merge to `main`
7. **Deploy**: Automatic deployment to staging

### Branch Naming
```
feature/add-product-variants
bugfix/fix-cart-calculation  
hotfix/security-patch-auth
refactor/improve-database-schema
docs/update-api-documentation
```

### Commit Messages
Follow [Conventional Commits](https://conventionalcommits.org/):
```
feat: add product variant support
fix: resolve cart total calculation bug
docs: update API documentation
test: add integration tests for orders
refactor: improve product service structure
perf: optimize database queries
chore: update dependencies
```

## Code Standards

### Backend (Go)

#### Project Structure
```go
// Follow hexagonal architecture
internal/
├── domain/          // Business entities and logic
├── application/     // Use cases and services  
├── ports/          // Interface definitions
├── adapters/       // External integrations
└── infrastructure/ // Framework and config
```

#### Coding Standards
```go
// Good: Clear, descriptive names
func CreateProductWithVariants(ctx context.Context, product *Product, variants []Variant) error {
    if err := p.validateProduct(product); err != nil {
        return fmt.Errorf("invalid product: %w", err)
    }
    
    return p.repo.Create(ctx, product, variants)
}

// Bad: Unclear names, no error handling
func Create(p *Product, v []Variant) {
    p.repo.Save(p, v)
}
```

#### Error Handling
```go
// Wrap errors with context
func (s *ProductService) GetProduct(ctx context.Context, id string) (*Product, error) {
    product, err := s.repo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, repository.ErrNotFound) {
            return nil, ErrProductNotFound
        }
        return nil, fmt.Errorf("failed to get product %s: %w", id, err)
    }
    
    return product, nil
}
```

#### Testing Requirements
```go
// Always write tests for new functions
func TestProductService_CreateProduct(t *testing.T) {
    tests := []struct {
        name     string
        product  *Product
        wantErr  bool
        errType  error
    }{
        {
            name: "valid product",
            product: &Product{Name: "Test", Price: 29.99},
            wantErr: false,
        },
        {
            name: "invalid price",
            product: &Product{Name: "Test", Price: -1},
            wantErr: true,
            errType: ErrInvalidPrice,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

### Frontend (React/Next.js)

#### Component Structure
```jsx
// Good: Clear component structure
import { useState, useCallback } from 'react'
import { useAtom } from 'jotai'
import { cartAtom } from '@/stores/cartStore'

export function ProductCard({ product, onAddToCart }) {
  const [cart, setCart] = useAtom(cartAtom)
  const [isLoading, setIsLoading] = useState(false)
  
  const handleAddToCart = useCallback(async () => {
    setIsLoading(true)
    try {
      await onAddToCart(product.id)
      // Update cart state
    } catch (error) {
      console.error('Failed to add to cart:', error)
    } finally {
      setIsLoading(false)
    }
  }, [product.id, onAddToCart])
  
  return (
    <div className="product-card">
      <img src={product.image} alt={product.name} />
      <h3>{product.name}</h3>
      <p>${product.price}</p>
      <button 
        onClick={handleAddToCart}
        disabled={isLoading}
        data-testid="add-to-cart"
      >
        {isLoading ? 'Adding...' : 'Add to Cart'}
      </button>
    </div>
  )
}
```

#### State Management (Jotai)
```jsx
// atoms/productAtoms.js
import { atom } from 'jotai'

export const productsAtom = atom([])
export const selectedProductAtom = atom(null)
export const productFiltersAtom = atom({
  category: null,
  priceRange: [0, 1000],
  search: ''
})

// Derived atom
export const filteredProductsAtom = atom((get) => {
  const products = get(productsAtom)
  const filters = get(productFiltersAtom)
  
  return products.filter(product => {
    if (filters.category && product.category !== filters.category) {
      return false
    }
    // Apply other filters
    return true
  })
})
```

## Code Review Process

### Pull Request Requirements
- [ ] **Description**: Clear description of changes
- [ ] **Tests**: New tests for new functionality  
- [ ] **Documentation**: Updated relevant docs
- [ ] **Breaking Changes**: Documented if any
- [ ] **Performance**: No performance regressions
- [ ] **Security**: No security vulnerabilities

### Review Checklist

#### Functionality
- [ ] Code solves the stated problem
- [ ] Edge cases are handled
- [ ] Error conditions are managed
- [ ] Business logic is correct

#### Code Quality
- [ ] Code is readable and maintainable
- [ ] Functions are single-purpose
- [ ] Variables and functions are well-named
- [ ] No code duplication
- [ ] Appropriate abstractions

#### Testing
- [ ] Unit tests cover new functionality
- [ ] Integration tests for complex flows
- [ ] Test names clearly describe what's being tested
- [ ] Tests are reliable and not flaky

#### Security
- [ ] Input validation is present
- [ ] No hardcoded secrets
- [ ] Proper authorization checks
- [ ] SQL injection prevention

## Development Environment

### Required Tools
```bash
# Backend
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/cosmtrek/air@latest

# Frontend
npm install -g prettier eslint

# Database
brew install postgresql
brew install redis
```

### IDE Setup

#### VS Code Extensions
```json
{
  "recommendations": [
    "golang.go",
    "bradlc.vscode-tailwindcss", 
    "esbenp.prettier-vscode",
    "ms-vscode.vscode-eslint",
    "ms-vscode.thunder-client"
  ]
}
```

#### Settings
```json
{
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "package",
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  }
}
```

## Documentation Standards

### API Documentation
- Use OpenAPI/Swagger specifications
- Include request/response examples
- Document error codes and messages
- Provide SDK examples

### Code Documentation
```go
// ProductService handles product-related business logic.
// It orchestrates between the domain layer and external adapters.
type ProductService struct {
    repo ProductRepository
    validator ProductValidator
}

// CreateProduct creates a new product with validation.
// It returns ErrInvalidProduct if validation fails.
// It returns ErrDuplicateSKU if SKU already exists.
func (s *ProductService) CreateProduct(ctx context.Context, product *Product) error {
    // implementation
}
```

### README Files
Every package/module should have a README with:
- Purpose and responsibilities
- Setup instructions
- Usage examples
- Testing instructions

## Performance Guidelines

### Backend Performance
```go
// Use context for cancellation
func (s *Service) ProcessOrder(ctx context.Context, orderID string) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // process order
    }
}

// Use database transactions for consistency
func (s *Service) CreateOrderWithItems(ctx context.Context, order Order) error {
    return s.db.WithTx(ctx, func(tx *gorm.DB) error {
        if err := tx.Create(&order).Error; err != nil {
            return err
        }
        
        for _, item := range order.Items {
            if err := tx.Create(&item).Error; err != nil {
                return err
            }
        }
        
        return nil
    })
}
```

### Frontend Performance
```jsx
// Use React.memo for expensive components
export const ProductList = React.memo(({ products }) => {
  return (
    <div>
      {products.map(product => (
        <ProductCard key={product.id} product={product} />
      ))}
    </div>
  )
})

// Use useMemo for expensive calculations
function ProductSummary({ products }) {
  const totalValue = useMemo(() => {
    return products.reduce((sum, product) => sum + product.price, 0)
  }, [products])
  
  return <div>Total: ${totalValue}</div>
}
```

## Security Guidelines

### Input Validation
```go
// Always validate input
type CreateProductRequest struct {
    Name        string  `json:"name" validate:"required,min=1,max=255"`
    Description string  `json:"description" validate:"max=10000"`
    Price       float64 `json:"price" validate:"required,min=0"`
    SKU         string  `json:"sku" validate:"required,alphanum,max=50"`
}
```

### Authentication
```jsx
// Use secure authentication patterns
function useAuthenticatedRequest() {
  const token = useAuthToken()
  
  return useCallback(async (url, options = {}) => {
    return fetch(url, {
      ...options,
      headers: {
        ...options.headers,
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    })
  }, [token])
}
```

## Deployment

### Environment Configuration
```bash
# Development
export NODE_ENV=development
export API_URL=http://localhost:8080

# Staging
export NODE_ENV=staging  
export API_URL=https://staging-api.example.com

# Production
export NODE_ENV=production
export API_URL=https://api.example.com
```

### Database Migrations
```go
// migrations/001_create_products.up.sql
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

// migrations/001_create_products.down.sql
DROP TABLE products;
```

## Getting Help

### Resources
- **Slack**: #dev-team channel
- **Documentation**: Check `/docs` folder first
- **API Issues**: Use Thunder Client or Postman
- **Database Issues**: Check logs in Docker

### Issue Reporting
1. Search existing issues first
2. Provide detailed reproduction steps
3. Include environment information
4. Add relevant logs and screenshots
5. Label appropriately (bug, feature, etc.)

## Recognition

Contributors who consistently follow these guidelines and help improve the codebase will be recognized in:
- Monthly team meetings
- Quarterly reviews
- Project documentation
- Public acknowledgments