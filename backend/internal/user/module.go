package user

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"ecommerce-saas/internal/security"
	"ecommerce-saas/internal/shared/email"
	"ecommerce-saas/internal/shared/utils"
)

// userRepositoryAdapter adapts the user repository to work with the security service
type userRepositoryAdapter struct {
	repo Repository
}

func (u *userRepositoryAdapter) UpdateUser2FA(ctx context.Context, userID uuid.UUID, enabled bool) error {
	return u.repo.UpdateUser2FA(ctx, userID, enabled)
}

func (u *userRepositoryAdapter) GetUserByID(ctx context.Context, userID uuid.UUID) (security.User, error) {
	user, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &userAdapter{user: user}, nil
}

func (u *userRepositoryAdapter) GetByID(ctx context.Context, userID uuid.UUID) (security.User, error) {
	return u.GetUserByID(ctx, userID)
}

// userAdapter adapts the User model to implement the security.User interface
type userAdapter struct {
	user *User
}

func (u *userAdapter) GetID() uuid.UUID {
	return u.user.ID
}

func (u *userAdapter) GetTenantID() *uuid.UUID {
	return u.user.TenantID
}

func (u *userAdapter) IsTwoFactorEnabled() bool {
	return u.user.TwoFactorEnabled
}

// Module represents the user module
type Module struct {
	Repository Repository
	Service    *Service
	Handler    *Handler
}

// NewModule creates a new user module instance
func NewModule(db *gorm.DB, jwtManager *utils.JWTManager) *Module {
	repo := NewRepository(db)
	tokenRepo := NewTokenRepository(db)
	securityService := security.NewSecurityService(db, &userRepositoryAdapter{repo: repo})
	emailService := email.NewService(nil) // TODO: Pass proper config
	svc := NewService(repo, tokenRepo, jwtManager, securityService, emailService)
	handler := NewHandler(svc)

	return &Module{
		Repository: repo,
		Service:    svc,
		Handler:    handler,
	}
}

// RegisterRoutes registers all user routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.Handler.RegisterRoutes(router)
}

// GetHandler returns the user handler
func (m *Module) GetHandler() *Handler {
	return m.Handler
}

// GetService returns the user service
func (m *Module) GetService() *Service {
	return m.Service
}

// GetRepository returns the user repository
func (m *Module) GetRepository() Repository {
	return m.Repository
}

// GetSecurityRepository returns the security-compatible repository adapter
func (m *Module) GetSecurityRepository() security.UserRepository {
	return &userRepositoryAdapter{repo: m.Repository}
}
