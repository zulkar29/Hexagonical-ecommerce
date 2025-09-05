package user

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	
	"ecommerce-saas/internal/shared/utils"
)

// ResetTokenData represents password reset token data
type ResetTokenData struct {
	UserID    uuid.UUID
	ExpiresAt time.Time
}

// Service handles user business logic
type Service struct {
	repo        Repository
	jwtManager  *utils.JWTManager
	resetTokens map[string]ResetTokenData
}

// NewService creates a new user service
func NewService(repo Repository, jwtManager *utils.JWTManager) *Service {
	return &Service{
		repo:        repo,
		jwtManager:  jwtManager,
		resetTokens: make(map[string]ResetTokenData),
	}
}

// RegisterUser creates a new user account
func (s *Service) RegisterUser(ctx context.Context, user *User) (*User, error) {
	// Validate input
	if err := s.validateUser(user); err != nil {
		return nil, err
	}

	// Check if user already exists
	existingUser, _ := s.repo.GetByEmail(user.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Set user fields
	user.ID = uuid.New()
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Password = string(hashedPassword)
	user.Status = StatusActive
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Set default role if not provided
	if user.Role == "" {
		user.Role = RoleCustomer
	}

	// Save user
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// Send verification email
	if err := s.sendVerificationEmail(user); err != nil {
		// Log error but don't fail registration
		log.Printf("Failed to send verification email: %v", err)
	}

	return user, nil
}

// LoginUser authenticates a user and returns tokens
func (s *Service) LoginUser(ctx context.Context, email, password string) (*LoginResponse, error) {
	// Get user by email
	user, err := s.repo.GetByEmail(strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check user status
	if user.Status != StatusActive {
		return nil, errors.New("account is not active")
	}

	// Generate tokens
	accessToken, refreshToken, err := s.jwtManager.GenerateTokens(
		user.ID,
		user.TenantID,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	// Create session
	session := &UserSession{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateSession(session); err != nil {
		return nil, err
	}

	// Update last login
	user.LastLoginAt = &[]time.Time{time.Now()}[0]
	s.repo.Update(user)

	return &LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    15 * 60, // 15 minutes in seconds
	}, nil
}

// RefreshToken refreshes access token using refresh token
func (s *Service) RefreshToken(refreshToken string) (*TokenResponse, error) {
	// Validate refresh token
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Check if session exists and is active
	session, err := s.repo.GetSessionByToken(refreshToken)
	if err != nil || !session.IsActive {
		return nil, errors.New("session not found or inactive")
	}

	// Get user
	user, err := s.repo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generate new tokens
	accessToken, newRefreshToken, err := s.jwtManager.GenerateTokens(
		user.ID,
		user.TenantID,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	// Update session with new refresh token
	session.Token = newRefreshToken
	session.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
	session.UpdatedAt = time.Now()
	s.repo.UpdateSession(session)

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    15 * 60,
	}, nil
}

// LogoutUser logs out a user by invalidating session
func (s *Service) LogoutUser(userID uuid.UUID, refreshToken string) error {
	session, err := s.repo.GetSessionByToken(refreshToken)
	if err != nil {
		return nil // Already logged out
	}

	if session.UserID != userID {
		return errors.New("unauthorized")
	}

	session.IsActive = false
	session.UpdatedAt = time.Now()
	return s.repo.UpdateSession(session)
}

// VerifyEmail verifies user email with token
func (s *Service) VerifyEmail(userID uuid.UUID, token string) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return err
	}

	// TODO: Validate verification token against stored token
	// For now, just mark as verified if token is provided
	if token == "" {
		return errors.New("verification token is required")
	}

	now := time.Now()
	user.EmailVerified = true
	user.EmailVerifiedAt = &now
	user.Status = StatusActive
	user.UpdatedAt = now

	return s.repo.Update(user)
}

// ChangePassword changes user password
func (s *Service) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Validate new password
	if err := s.validatePassword(newPassword); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	now := time.Now()
	user.Password = string(hashedPassword)
	user.PasswordChangedAt = &now
	user.UpdatedAt = now

	// Invalidate all existing sessions
	s.repo.InvalidateUserSessions(userID)

	return s.repo.Update(user)
}

// ResetPassword resets user password with token
func (s *Service) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Validate token
	tokenData, exists := s.resetTokens[token]
	if !exists {
		return errors.New("invalid reset token")
	}

	// Check if token is expired
	if time.Now().After(tokenData.ExpiresAt) {
		delete(s.resetTokens, token)
		return errors.New("reset token has expired")
	}

	// Get user
	user, err := s.repo.GetByID(tokenData.UserID)
	if err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	user.Password = string(hashedPassword)
	user.PasswordChangedAt = &time.Time{}
	*user.PasswordChangedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Save user
	if err := s.repo.Update(user); err != nil {
		return err
	}

	// Remove used token
	delete(s.resetTokens, token)

	return nil
}

// GetProfile gets user profile
func (s *Service) GetProfile(userID uuid.UUID) (*User, error) {
	return s.repo.GetByID(userID)
}

// UpdateProfile updates user profile information
func (s *Service) UpdateProfile(ctx context.Context, userID uuid.UUID, updates *User) (*User, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if updates.FirstName != "" {
		user.FirstName = strings.TrimSpace(updates.FirstName)
	}
	if updates.LastName != "" {
		user.LastName = strings.TrimSpace(updates.LastName)
	}
	if updates.Phone != "" {
		user.Phone = updates.Phone
	}
	if updates.Avatar != "" {
		user.Avatar = updates.Avatar
	}

	user.UpdatedAt = time.Now()

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Helper methods
func (s *Service) validateUser(user *User) error {
	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}
	if user.FirstName == "" {
		return errors.New("first name is required")
	}
	if user.LastName == "" {
		return errors.New("last name is required")
	}

	// Validate email format
	if !utils.IsValidEmail(user.Email) {
		return errors.New("invalid email format")
	}

	// Validate password strength
	if err := s.validatePassword(user.Password); err != nil {
		return err
	}

	return nil
}

func (s *Service) validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Use the utility function for comprehensive password validation
	return utils.ValidatePassword(password)
}

// GetUserFromToken extracts user from JWT token
func (s *Service) GetUserFromToken(tokenString string) (*User, error) {
	claims, err := s.jwtManager.ValidateToken(tokenString)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	return s.repo.GetByID(claims.UserID)
}

// ListUsers returns paginated list of users
func (s *Service) ListUsers(tenantID *uuid.UUID, filter UserFilter, page, limit int) ([]*User, int64, error) {
	offset := (page - 1) * limit
	return s.repo.List(tenantID, filter, offset, limit)
}

// UpdateUserRole updates user role (admin only)
func (s *Service) UpdateUserRole(adminUserID, targetUserID uuid.UUID, newRole UserRole) error {
	// Check if admin has permission
	admin, err := s.repo.GetByID(adminUserID)
	if err != nil {
		return err
	}

	if !admin.IsAdmin() {
		return errors.New("insufficient permissions")
	}

	// Get target user
	targetUser, err := s.repo.GetByID(targetUserID)
	if err != nil {
		return err
	}

	// Update role
	targetUser.Role = newRole
	targetUser.UpdatedAt = time.Now()

	return s.repo.Update(targetUser)
}

// UpdateUserStatus updates user status (admin only)
func (s *Service) UpdateUserStatus(adminUserID, targetUserID uuid.UUID, newStatus UserStatus) error {
	// Check if admin has permission
	admin, err := s.repo.GetByID(adminUserID)
	if err != nil {
		return err
	}

	if !admin.IsAdmin() {
		return errors.New("insufficient permissions")
	}

	// Get target user
	targetUser, err := s.repo.GetByID(targetUserID)
	if err != nil {
		return err
	}

	// Update status
	targetUser.Status = newStatus
	targetUser.UpdatedAt = time.Now()

	return s.repo.Update(targetUser)
}

// DeleteUser soft deletes a user (admin only)
func (s *Service) DeleteUser(adminUserID, targetUserID uuid.UUID) error {
	// Check if admin has permission
	admin, err := s.repo.GetByID(adminUserID)
	if err != nil {
		return err
	}

	if !admin.IsAdmin() {
		return errors.New("insufficient permissions")
	}

	// Cannot delete self
	if adminUserID == targetUserID {
		return errors.New("cannot delete your own account")
	}

	// Invalidate all sessions first
	s.repo.InvalidateUserSessions(targetUserID)

	// Delete user
	return s.repo.Delete(targetUserID)
}

// GetUserPermissions returns user permissions
func (s *Service) GetUserPermissions(userID uuid.UUID) ([]*Permission, error) {
	return s.repo.GetUserPermissions(userID)
}

// CheckUserPermission checks if user has specific permission
func (s *Service) CheckUserPermission(userID uuid.UUID, resource, action string) (bool, error) {
	return s.repo.CheckUserPermission(userID, resource, action)
}

// CleanupExpiredSessions removes expired sessions
func (s *Service) CleanupExpiredSessions() error {
	return s.repo.CleanupExpiredSessions()
}

// sendVerificationEmail sends email verification email
func (s *Service) sendVerificationEmail(user *User) error {
	// TODO: Implement email service integration
	// For now, just log that verification email would be sent
	log.Printf("Verification email would be sent to: %s", user.Email)
	return nil
}

// ForgotPassword initiates password reset process
func (s *Service) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		// Don't reveal if email exists
		return nil
	}

	// Generate reset token
	resetToken := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour) // 24 hour expiry

	// Store reset token
	s.resetTokens[resetToken] = ResetTokenData{
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}

	// TODO: Send reset email with resetToken
	log.Printf("Password reset email would be sent to: %s with token: %s", user.Email, resetToken)

	return nil
}
