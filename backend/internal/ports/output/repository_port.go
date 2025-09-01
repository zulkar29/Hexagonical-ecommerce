package output

// TODO: Implement simple output ports during development
// These define what we need from external services

// TenantRepository defines data persistence operations
type TenantRepository interface {
	// TODO: Define simple CRUD operations
	// Create(ctx context.Context, tenant *entities.Tenant) error
	// GetByID(ctx context.Context, id string) (*entities.Tenant, error)
	// GetBySubdomain(ctx context.Context, subdomain string) (*entities.Tenant, error)
	// Update(ctx context.Context, tenant *entities.Tenant) error
	// Delete(ctx context.Context, id string) error
}

// ProductRepository defines product data operations
type ProductRepository interface {
	// TODO: Define simple CRUD operations
	// Create(ctx context.Context, product *entities.Product) error
	// GetByID(ctx context.Context, id string) (*entities.Product, error)
	// GetByTenant(ctx context.Context, tenantID string) ([]*entities.Product, error)
	// Update(ctx context.Context, product *entities.Product) error
	// Delete(ctx context.Context, id string) error
}

// EmailService defines email operations
type EmailService interface {
	// TODO: Define simple email operations
	// SendWelcomeEmail(ctx context.Context, email, name string) error
	// SendOrderConfirmation(ctx context.Context, email string, order *entities.Order) error
}