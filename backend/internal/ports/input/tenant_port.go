package input

// TODO: Implement simple input ports during development
// These define what our application can do (use cases)

// TenantPort defines tenant-related operations
type TenantPort interface {
	// TODO: Define simple interface methods
	// CreateTenant(ctx context.Context, req CreateTenantRequest) (*TenantResponse, error)
	// GetTenant(ctx context.Context, id string) (*TenantResponse, error)
	// UpdateTenant(ctx context.Context, id string, req UpdateTenantRequest) (*TenantResponse, error)
	// DeleteTenant(ctx context.Context, id string) error
}

// Request/Response DTOs (keep simple)
// type CreateTenantRequest struct {
//     Name      string `json:"name" validate:"required"`
//     Email     string `json:"email" validate:"required,email"`
//     Subdomain string `json:"subdomain" validate:"required"`
// }

// type TenantResponse struct {
//     ID        string `json:"id"`
//     Name      string `json:"name"`
//     Subdomain string `json:"subdomain"`
//     CreatedAt string `json:"created_at"`
// }