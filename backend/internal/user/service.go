package user

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	
	"ecommerce-saas/internal/shared/utils"
)

// Service handles user business logic
type Service struct {
	repo       Repository
	jwtManager *utils.JWTManager
}

// NewService creates a new user service
func NewService(repo Repository, jwtManager *utils.JWTManager) *Service {
	return &Service{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

// RegisterUser registers a new user
func (s *Service) RegisterUser(req RegisterRequest) (*User, error) {
	// Validate input
	if err := s.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	// Check if user already exists
	existingUser, _ := s.repo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists with this email")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &User{
		ID:        uuid.New(),
		TenantID:  req.TenantID,
		Email:     strings.ToLower(strings.TrimSpace(req.Email)),
		Password:  string(hashedPassword),
		FirstName: strings.TrimSpace(req.FirstName),
		LastName:  strings.TrimSpace(req.LastName),
		Phone:     req.Phone,
		Role:      req.Role,
		Status:    StatusPending, // Requires email verification
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set default role if not provided
	if user.Role == "" {
		user.Role = RoleCustomer
	}

	// Save user
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// TODO: Send verification email
	// s.sendVerificationEmail(user)

	return user, nil
}

// LoginUser authenticates a user and returns tokens
func (s *Service) LoginUser(email, password string) (*LoginResponse, error) {
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
	// For now, just mark as verified
	now := time.Now()
	user.EmailVerified = true
	user.EmailVerifiedAt = &now
	user.Status = StatusActive
	user.UpdatedAt = now

	return s.repo.Update(user)
}

// ChangePassword changes user password
func (s *Service) ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error {
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

// ResetPassword initiates password reset
func (s *Service) ResetPassword(email string) error {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		// Don't reveal if email exists
		return nil
	}

	// Generate reset token
	resetToken, err := utils.GenerateResetToken()
	if err != nil {
		return err
	}

	// TODO: Store reset token with expiry
	// TODO: Send reset email

	_ = resetToken
	_ = user

	return nil
}

// GetProfile gets user profile
func (s *Service) GetProfile(userID uuid.UUID) (*User, error) {
	return s.repo.GetByID(userID)
}

// UpdateProfile updates user profile
func (s *Service) UpdateProfile(userID uuid.UUID, req UpdateProfileRequest) (*User, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.FirstName != "" {
		user.FirstName = strings.TrimSpace(req.FirstName)
	}
	if req.LastName != "" {
		user.LastName = strings.TrimSpace(req.LastName)
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	user.UpdatedAt = time.Now()

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Helper methods
func (s *Service) validateRegisterRequest(req RegisterRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if req.FirstName == "" {
		return errors.New("first name is required")
	}
	if req.LastName == "" {
		return errors.New("last name is required")
	}

	// Validate email format
	if !utils.IsValidEmail(req.Email) {
		return errors.New("invalid email format")
	}

	// Validate password strength
	if err := s.validatePassword(req.Password); err != nil {
		return err
	}

	return nil
}

func (s *Service) validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Add more password validation rules as needed
	return nil
}
