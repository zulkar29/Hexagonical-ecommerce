package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service defines the platform service interface
type Service interface {
	// Dashboard & Stats
	GetPlatformStats(ctx context.Context, period string) (*PlatformStats, error)
	GetSystemStatus(ctx context.Context) (*SystemStatus, error)

	// Platform Admin Management
	ListPlatformAdmins(ctx context.Context, role string) ([]*PlatformAdmin, error)
	GetPlatformAdmin(ctx context.Context, id uuid.UUID) (*PlatformAdmin, error)
	CreatePlatformAdmin(ctx context.Context, req *CreatePlatformAdminRequest) (*PlatformAdmin, error)
	UpdatePlatformAdmin(ctx context.Context, id uuid.UUID, req *UpdatePlatformAdminRequest) (*PlatformAdmin, error)
	DeletePlatformAdmin(ctx context.Context, id uuid.UUID) error

	// Platform Role Management
	ListPlatformRoles(ctx context.Context) ([]*PlatformRole, error)
	GetPlatformRole(ctx context.Context, id uuid.UUID) (*PlatformRole, error)
	CreatePlatformRole(ctx context.Context, req *CreatePlatformRoleRequest) (*PlatformRole, error)
	UpdatePlatformRole(ctx context.Context, id uuid.UUID, req *UpdatePlatformRoleRequest) (*PlatformRole, error)
	DeletePlatformRole(ctx context.Context, id uuid.UUID) error

	// Platform Settings
	GetPlatformSettings(ctx context.Context, category string) ([]*PlatformSettings, error)
	UpdatePlatformSettings(ctx context.Context, req *UpdatePlatformSettingsRequest) (*PlatformSettings, error)

	// Platform Tenant Management
	ListAllTenants(ctx context.Context, req *ListTenantsRequest) (*ListTenantsResponse, error)
	GetTenantDetails(ctx context.Context, id uuid.UUID, include []string) (map[string]interface{}, error)
	UpdateTenant(ctx context.Context, id uuid.UUID, req *UpdateTenantRequest) error
	DeleteTenant(ctx context.Context, id uuid.UUID) error

	// Audit Logs
	GetPlatformAuditLogs(ctx context.Context, req *GetAuditLogsRequest) (*GetAuditLogsResponse, error)
	LogPlatformActivity(ctx context.Context, userID *uuid.UUID, tenantID *uuid.UUID, action, resource string, details map[string]interface{}) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new platform service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// GetPlatformStats retrieves platform statistics
func (s *service) GetPlatformStats(ctx context.Context, period string) (*PlatformStats, error) {
	// Validate period
	validPeriods := []string{"day", "week", "month", "quarter", "year"}
	if period == "" {
		period = "month" // default
	}
	
	valid := false
	for _, p := range validPeriods {
		if period == p {
			valid = true
			break
		}
	}
	
	if !valid {
		return nil, fmt.Errorf("invalid period: %s. Valid periods are: day, week, month, quarter, year", period)
	}

	return s.repo.GetPlatformStats(ctx, period)
}

// GetSystemStatus retrieves system health status
func (s *service) GetSystemStatus(ctx context.Context) (*SystemStatus, error) {
	return s.repo.GetSystemStatus(ctx)
}

// ListPlatformAdmins retrieves platform administrators
func (s *service) ListPlatformAdmins(ctx context.Context, role string) ([]*PlatformAdmin, error) {
	return s.repo.ListPlatformAdmins(ctx, role)
}

// GetPlatformAdmin retrieves a platform administrator by ID
func (s *service) GetPlatformAdmin(ctx context.Context, id uuid.UUID) (*PlatformAdmin, error) {
	return s.repo.GetPlatformAdmin(ctx, id)
}

// CreatePlatformAdmin creates a new platform administrator
func (s *service) CreatePlatformAdmin(ctx context.Context, req *CreatePlatformAdminRequest) (*PlatformAdmin, error) {
	// Validate request
	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Password == "" {
		return nil, fmt.Errorf("password is required")
	}
	if req.Role == "" {
		return nil, fmt.Errorf("role is required")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create admin
	admin := &PlatformAdmin{
		ID:        uuid.New(),
		Email:     req.Email,
		Name:      req.Name,
		Password:  string(hashedPassword),
		Role:      req.Role,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if req.Permissions != nil {
		admin.Permissions = req.Permissions
	} else {
		admin.Permissions = []string{} // Initialize empty permissions array
	}

	createdAdmin, err := s.repo.CreatePlatformAdmin(ctx, admin)
	if err != nil {
		return nil, fmt.Errorf("failed to create platform admin: %w", err)
	}

	// Log activity
	s.LogPlatformActivity(ctx, nil, nil, "create", "platform_admin", map[string]interface{}{
		"admin_id": createdAdmin.ID,
		"email":    createdAdmin.Email,
		"role":     createdAdmin.Role,
	})

	return createdAdmin, nil
}

// UpdatePlatformAdmin updates a platform administrator
func (s *service) UpdatePlatformAdmin(ctx context.Context, id uuid.UUID, req *UpdatePlatformAdminRequest) (*PlatformAdmin, error) {
	// Get existing admin
	existingAdmin, err := s.repo.GetPlatformAdmin(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("admin not found: %w", err)
	}

	// Update fields
	updatedAdmin := *existingAdmin
	updatedAdmin.UpdatedAt = time.Now()

	if req.Name != nil {
		updatedAdmin.Name = *req.Name
	}
	if req.Email != nil {
		updatedAdmin.Email = *req.Email
	}
	if req.Role != nil {
		updatedAdmin.Role = *req.Role
	}
	if req.Status != nil {
		updatedAdmin.Status = *req.Status
	}
	if req.Permissions != nil {
		updatedAdmin.Permissions = req.Permissions
	}
	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		updatedAdmin.Password = string(hashedPassword)
	}

	result, err := s.repo.UpdatePlatformAdmin(ctx, id, &updatedAdmin)
	if err != nil {
		return nil, fmt.Errorf("failed to update platform admin: %w", err)
	}

	// Log activity
	s.LogPlatformActivity(ctx, nil, nil, "update", "platform_admin", map[string]interface{}{
		"admin_id": id,
		"changes":  req,
	})

	return result, nil
}

// DeletePlatformAdmin deletes a platform administrator
func (s *service) DeletePlatformAdmin(ctx context.Context, id uuid.UUID) error {
	// Check if admin exists
	_, err := s.repo.GetPlatformAdmin(ctx, id)
	if err != nil {
		return fmt.Errorf("admin not found: %w", err)
	}

	err = s.repo.DeletePlatformAdmin(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete platform admin: %w", err)
	}

	// Log activity
	s.LogPlatformActivity(ctx, nil, nil, "delete", "platform_admin", map[string]interface{}{
		"admin_id": id,
	})

	return nil
}

// ListPlatformRoles retrieves platform roles
func (s *service) ListPlatformRoles(ctx context.Context) ([]*PlatformRole, error) {
	return s.repo.ListPlatformRoles(ctx)
}

// GetPlatformRole retrieves a platform role by ID
func (s *service) GetPlatformRole(ctx context.Context, id uuid.UUID) (*PlatformRole, error) {
	return s.repo.GetPlatformRole(ctx, id)
}

// CreatePlatformRole creates a new platform role
func (s *service) CreatePlatformRole(ctx context.Context, req *CreatePlatformRoleRequest) (*PlatformRole, error) {
	// Validate request
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Permissions == nil || len(req.Permissions) == 0 {
		return nil, fmt.Errorf("permissions are required")
	}

	// Create role
	role := &PlatformRole{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permissions,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdRole, err := s.repo.CreatePlatformRole(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("failed to create platform role: %w", err)
	}

	// Log activity
	s.LogPlatformActivity(ctx, nil, nil, "create", "platform_role", map[string]interface{}{
		"role_id": createdRole.ID,
		"name":    createdRole.Name,
	})

	return createdRole, nil
}

// UpdatePlatformRole updates a platform role
func (s *service) UpdatePlatformRole(ctx context.Context, id uuid.UUID, req *UpdatePlatformRoleRequest) (*PlatformRole, error) {
	// Get existing role
	existingRole, err := s.repo.GetPlatformRole(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}

	// Update fields
	updatedRole := *existingRole
	updatedRole.UpdatedAt = time.Now()

	if req.Name != nil {
		updatedRole.Name = *req.Name
	}
	if req.Description != nil {
		updatedRole.Description = *req.Description
	}
	if req.Permissions != nil {
		updatedRole.Permissions = req.Permissions
	}

	result, err := s.repo.UpdatePlatformRole(ctx, id, &updatedRole)
	if err != nil {
		return nil, fmt.Errorf("failed to update platform role: %w", err)
	}

	// Log activity
	s.LogPlatformActivity(ctx, nil, nil, "update", "platform_role", map[string]interface{}{
		"role_id": id,
		"changes": req,
	})

	return result, nil
}

// DeletePlatformRole deletes a platform role
func (s *service) DeletePlatformRole(ctx context.Context, id uuid.UUID) error {
	// Check if role exists
	_, err := s.repo.GetPlatformRole(ctx, id)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	err = s.repo.DeletePlatformRole(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete platform role: %w", err)
	}

	// Log activity
	s.LogPlatformActivity(ctx, nil, nil, "delete", "platform_role", map[string]interface{}{
		"role_id": id,
	})

	return nil
}

// GetPlatformSettings retrieves platform settings
func (s *service) GetPlatformSettings(ctx context.Context, category string) ([]*PlatformSettings, error) {
	return s.repo.GetPlatformSettings(ctx, category)
}

// UpdatePlatformSettings updates platform settings
func (s *service) UpdatePlatformSettings(ctx context.Context, req *UpdatePlatformSettingsRequest) (*PlatformSettings, error) {
	// Validate request
	if req.Category == "" {
		return nil, fmt.Errorf("category is required")
	}
	if req.Key == "" {
		return nil, fmt.Errorf("key is required")
	}

	// Create or update settings
	settings := &PlatformSettings{
		ID:        uuid.New(),
		Category:  req.Category,
		Key:       req.Key,
		Value:     req.Value,
		UpdatedAt: time.Now(),
	}

	if req.Description != nil {
		settings.Description = *req.Description
	}

	result, err := s.repo.UpdatePlatformSettings(ctx, settings)
	if err != nil {
		return nil, fmt.Errorf("failed to update platform settings: %w", err)
	}

	// Log activity
	s.LogPlatformActivity(ctx, nil, nil, "update", "platform_settings", map[string]interface{}{
		"category": req.Category,
		"key":      req.Key,
		"value":    req.Value,
	})

	return result, nil
}

// ListAllTenants retrieves all tenants for platform admin
func (s *service) ListAllTenants(ctx context.Context, req *ListTenantsRequest) (*ListTenantsResponse, error) {
	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	tenants, err := s.repo.ListAllTenants(ctx, req.Status, req.Include)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}

	// Apply pagination
	start := req.Offset
	end := start + req.Limit
	if end > len(tenants) {
		end = len(tenants)
	}
	if start > len(tenants) {
		start = len(tenants)
	}

	paginatedTenants := tenants[start:end]

	return &ListTenantsResponse{
		Tenants: paginatedTenants,
		Total:   len(tenants),
		Limit:   req.Limit,
		Offset:  req.Offset,
	}, nil
}

// GetTenantDetails retrieves detailed tenant information
func (s *service) GetTenantDetails(ctx context.Context, id uuid.UUID, include []string) (map[string]interface{}, error) {
	return s.repo.GetTenantDetails(ctx, id, include)
}

// UpdateTenant updates tenant information
func (s *service) UpdateTenant(ctx context.Context, id uuid.UUID, req *UpdateTenantRequest) error {
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Plan != nil {
		updates["plan"] = *req.Plan
	}
	if req.Domain != nil {
		updates["domain"] = *req.Domain
	}

	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	updates["updated_at"] = time.Now()

	err := s.repo.UpdateTenant(ctx, id, updates)
	if err != nil {
		return fmt.Errorf("failed to update tenant: %w", err)
	}

	// Log activity
	s.LogPlatformActivity(ctx, nil, &id, "update", "tenant", map[string]interface{}{
		"tenant_id": id,
		"changes":   updates,
	})

	return nil
}

// DeleteTenant deletes a tenant
func (s *service) DeleteTenant(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteTenant(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}

	// Log activity
	s.LogPlatformActivity(ctx, nil, &id, "delete", "tenant", map[string]interface{}{
		"tenant_id": id,
	})

	return nil
}

// GetPlatformAuditLogs retrieves platform audit logs
func (s *service) GetPlatformAuditLogs(ctx context.Context, req *GetAuditLogsRequest) (*GetAuditLogsResponse, error) {
	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 50
	}
	if req.Limit > 200 {
		req.Limit = 200
	}

	filter := AuditLogFilter{
		UserID:   req.UserID,
		TenantID: req.TenantID,
		Action:   req.Action,
		Resource: req.Resource,
		DateFrom: req.DateFrom,
		DateTo:   req.DateTo,
		Limit:    req.Limit,
		Offset:   req.Offset,
	}

	logs, err := s.repo.GetPlatformAuditLogs(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}

	return &GetAuditLogsResponse{
		Logs:   logs,
		Total:  len(logs),
		Limit:  req.Limit,
		Offset: req.Offset,
	}, nil
}

// LogPlatformActivity logs platform activity
func (s *service) LogPlatformActivity(ctx context.Context, userID *uuid.UUID, tenantID *uuid.UUID, action, resource string, details map[string]interface{}) error {
	log := &PlatformAuditLog{
		ID:        uuid.New(),
		UserID:    userID,
		TenantID:  tenantID,
		Action:    action,
		Resource:  resource,
		Details:   details,
		CreatedAt: time.Now(),
	}

	return s.repo.CreateAuditLog(ctx, log)
}