package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"

	sharedErrors "ecommerce-saas/internal/shared/errors"
	"ecommerce-saas/internal/shared/utils"
)

// ResetTokenData represents password reset token data
type ResetTokenData struct {
	UserID    uuid.UUID
	ExpiresAt time.Time
}

// Service handles user business logic
type Service struct {
	repo           Repository
	jwtManager     *utils.JWTManager
	securityService *SecurityService
	resetTokens    map[string]ResetTokenData
}

// NewService creates a new user service
func NewService(repo Repository, jwtManager *utils.JWTManager, securityService *SecurityService) *Service {
	return &Service{
		repo:           repo,
		jwtManager:     jwtManager,
		securityService: securityService,
		resetTokens:    make(map[string]ResetTokenData),
	}
}

// RegisterUser creates a new user account
func (s *Service) RegisterUser(ctx context.Context, user *User) (*User, error) {
	// Validate input
	if err := s.validateUser(user); err != nil {
		return nil, err
	}

	// Validate password with security policies
	if err := s.securityService.ValidatePassword(ctx, user.TenantID, user.Password); err != nil {
		return nil, err
	}

	// Check if user already exists
	existingUser, _ := s.repo.GetUserByEmail(ctx, user.Email)
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
	user, err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, sharedErrors.Wrap(err, sharedErrors.CodeInternal, "Failed to create user", 500)
	}

	// Save initial password to history
	if err := s.securityService.SavePasswordHistory(ctx, user.ID, user.Password); err != nil {
		log.Printf("Failed to save initial password history: %v", err)
	}

	// Log security event
	s.securityService.LogSecurityEvent(ctx, user.ID, user.TenantID, SecurityEventAccountCreated, "", "", map[string]interface{}{
		"message": "User account created",
		"email":   user.Email,
	})

	// Send verification email
	s.sendVerificationEmail(user)

	return user, nil
}

// LoginUser authenticates a user and returns tokens
func (s *Service) LoginUser(ctx context.Context, email, password string) (*LoginResponse, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	
	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		// Log failed login attempt even if user doesn't exist
		s.securityService.LogSecurityEvent(ctx, uuid.Nil, nil, SecurityEventFailedLogin, "", "", map[string]interface{}{
			"message": "User not found",
			"email":   email,
		})
		return nil, sharedErrors.ErrInvalidCredentials
	}

	// Check if account is locked
	if locked, err := s.securityService.IsAccountLocked(ctx, user.ID, user.TenantID); err != nil {
		return nil, sharedErrors.NewInternalError("Security check failed", err)
	} else if locked {
		return nil, sharedErrors.NewForbiddenError("Account is temporarily locked due to multiple failed login attempts")
	}

	// Check password
	if tempErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); tempErr != nil {
		// Record failed login attempt
		if err := s.securityService.RecordFailedLogin(ctx, user.ID, user.TenantID, "", ""); err != nil {
			log.Printf("Failed to record failed login: %v", err)
		}
		return nil, sharedErrors.ErrInvalidCredentials
	}

	// Check user status
	if user.Status != StatusActive {
		return nil, sharedErrors.ErrAccountInactive
	}

	// Reset failed login attempts on successful login
	if err := s.securityService.ResetFailedLogins(ctx, user.ID); err != nil {
		log.Printf("Failed to reset failed login attempts: %v", err)
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
	if _, err := s.repo.UpdateUser(ctx, user); err != nil {
		log.Printf("Failed to update user last login: %v", err)
	}

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
		return nil, sharedErrors.ErrInvalidToken
	}

	// Check if session exists and is active
	session, err := s.repo.GetSessionByToken(refreshToken)
	if err != nil || !session.IsActive {
		return nil, sharedErrors.NewUnauthorizedError("Session not found or inactive")
	}

	// Get user
	user, err := s.repo.GetUserByID(context.Background(), claims.UserID)
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
	if err := s.repo.UpdateSession(session); err != nil {
		log.Printf("Failed to update session: %v", err)
	}

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
	user, err := s.repo.GetUserByID(context.Background(), userID)
	if err != nil {
		return err
	}

	// Validate verification token (email verification service integration needed)
	// For now, just mark as verified if token is provided
	if token == "" {
		return errors.New("verification token is required")
	}

	now := time.Now()
	user.EmailVerified = true
	user.EmailVerifiedAt = &now
	user.Status = StatusActive
	user.UpdatedAt = now

	if _, err := s.repo.UpdateUser(context.Background(), user); err != nil {
		return sharedErrors.Wrap(err, sharedErrors.CodeInternal, "Failed to update user", 500)
	}
	return nil
}

// ChangePassword changes user password
func (s *Service) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if tempErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); tempErr != nil {
		return errors.New("current password is incorrect")
	}

	// Validate new password with security policies
	if err := s.securityService.ValidatePassword(ctx, user.TenantID, newPassword); err != nil {
		return err
	}

	// Check password history
	if reused, err := s.securityService.IsPasswordReused(ctx, userID, user.TenantID, newPassword); err != nil {
		return errors.New("failed to check password history")
	} else if reused {
		return errors.New("password has been used recently, please choose a different password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Save password to history
	if err := s.securityService.SavePasswordHistory(ctx, userID, string(hashedPassword)); err != nil {
		log.Printf("Failed to save password history: %v", err)
	}

	// Update password
	now := time.Now()
	user.Password = string(hashedPassword)
	user.PasswordChangedAt = &now
	user.UpdatedAt = now

	// Log security event
		s.securityService.LogSecurityEvent(ctx, userID, user.TenantID, SecurityEventPasswordChange, "", "", map[string]interface{}{
			"message": "Password changed successfully",
			"email":   user.Email,
		})

	// Invalidate all existing sessions
	if err := s.repo.InvalidateUserSessions(ctx, userID); err != nil {
		log.Printf("Failed to invalidate user sessions: %v", err)
	}

	_, err = s.repo.UpdateUser(ctx, user)
	return err
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
	user, err := s.repo.GetUserByID(ctx, tokenData.UserID)
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
	if _, err := s.repo.UpdateUser(ctx, user); err != nil {
		return err
	}

	// Remove used token
	delete(s.resetTokens, token)

	return nil
}

// GetProfile gets user profile
func (s *Service) GetProfile(ctx context.Context, userID uuid.UUID) (*User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

// UpdateProfile updates user profile information
func (s *Service) UpdateProfile(ctx context.Context, userID uuid.UUID, updates *User) (*User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
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

	if _, err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, sharedErrors.Wrap(err, sharedErrors.CodeInternal, "Failed to update user", 500)
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

	return s.repo.GetUserByID(context.Background(), claims.UserID)
}

// ListUsers returns paginated list of users
func (s *Service) ListUsers(ctx context.Context, tenantID *uuid.UUID, filter UserFilter, page, limit int) ([]*User, int64, error) {
	offset := (page - 1) * limit
	return s.repo.ListUsers(ctx, tenantID, filter, offset, limit)
}

// UpdateUserRole updates user role (admin only)
func (s *Service) UpdateUserRole(adminUserID, targetUserID uuid.UUID, newRole UserRole) error {
	// Check if admin has permission
	admin, err := s.repo.GetUserByID(context.Background(), adminUserID)
	if err != nil {
		return err
	}

	if !admin.IsAdmin() {
		return errors.New("insufficient permissions")
	}

	// Get target user
	targetUser, err := s.repo.GetUserByID(context.Background(), targetUserID)
	if err != nil {
		return err
	}

	// Update role
	targetUser.Role = newRole
	targetUser.UpdatedAt = time.Now()

	if _, err := s.repo.UpdateUser(context.Background(), targetUser); err != nil {
		return sharedErrors.Wrap(err, sharedErrors.CodeInternal, "Failed to update user role", 500)
	}
	return nil
}

// UpdateUserStatus updates user status (admin only)
func (s *Service) UpdateUserStatus(adminUserID, targetUserID uuid.UUID, newStatus UserStatus) error {
	// Check if admin has permission
	admin, err := s.repo.GetUserByID(context.Background(), adminUserID)
	if err != nil {
		return err
	}

	if !admin.IsAdmin() {
		return errors.New("insufficient permissions")
	}

	// Get target user
	targetUser, err := s.repo.GetUserByID(context.Background(), targetUserID)
	if err != nil {
		return err
	}

	// Update status
	targetUser.Status = newStatus
	targetUser.UpdatedAt = time.Now()

	if _, err := s.repo.UpdateUser(context.Background(), targetUser); err != nil {
		return sharedErrors.Wrap(err, sharedErrors.CodeInternal, "Failed to update user status", 500)
	}
	return nil
}

// DeleteUser soft deletes a user (admin only)
func (s *Service) DeleteUser(adminUserID, targetUserID uuid.UUID) error {
	// Check if admin has permission
	admin, err := s.repo.GetUserByID(context.Background(), adminUserID)
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
	if err := s.repo.InvalidateUserSessions(context.Background(), targetUserID); err != nil {
		log.Printf("Failed to invalidate user sessions: %v", err)
	}

	// Delete user
	if err := s.repo.DeleteUser(context.Background(), targetUserID); err != nil {
		return sharedErrors.Wrap(err, sharedErrors.CodeInternal, "Failed to delete user", 500)
	}
	return nil
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
func (s *Service) sendVerificationEmail(user *User) {
	// Send verification email (email service integration needed)
	// For now, just log that verification email would be sent
	log.Printf("Verification email would be sent to: %s", user.Email)
}

// ForgotPassword initiates password reset process
func (s *Service) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
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

	// Send reset email with token (email service integration needed)
	log.Printf("Password reset email would be sent to: %s with token: %s", user.Email, resetToken)

	return nil
}

// SendPasswordResetEmail sends password reset email
func (s *Service) SendPasswordResetEmail(email, token string) error {
	// Send password reset email (email service integration needed)
	return nil
}

// ResendVerificationEmail resends verification email
func (s *Service) ResendVerificationEmail(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	if user.EmailVerified {
		return errors.New("email already verified")
	}

	// Generate new verification token and send email (email service integration needed)
	return nil
}

// ResendVerification resends verification email (alias for handler compatibility)
func (s *Service) ResendVerification(ctx context.Context, email string) error {
	return s.ResendVerificationEmail(ctx, email)
}

// DeleteAccount deletes user account (alias for handler compatibility)
func (s *Service) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	return s.DeleteUserAccount(ctx, userID)
}

// DeleteUserAccount deletes user account
func (s *Service) DeleteUserAccount(ctx context.Context, userID uuid.UUID) error {
	// Invalidate all user sessions
	if err := s.repo.InvalidateUserSessions(ctx, userID); err != nil {
		return err
	}

	// Delete user account
	return s.repo.DeleteUser(ctx, userID)
}

// GetUserPreferences gets user preferences
func (s *Service) GetUserPreferences(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	return s.repo.GetUserPreferences(ctx, userID)
}

// UpdateUserPreferences updates user preferences
func (s *Service) UpdateUserPreferences(ctx context.Context, userID uuid.UUID, preferences map[string]interface{}) (map[string]interface{}, error) {
	result, err := s.repo.UpdateUserPreferences(ctx, userID, preferences)
	if err != nil {
		return nil, sharedErrors.Wrap(err, sharedErrors.CodeInternal, "Failed to save user preferences", 500)
	}
	return result, nil
}

// GetUserActivity gets user activity logs
func (s *Service) GetUserActivity(userID uuid.UUID, activityType, dateFrom, dateTo string, page, limit int) ([]interface{}, int64, error) {
	// Get user activity logs (activity service integration needed)
	return []interface{}{}, 0, nil
}

// UpdateUserByAdmin updates user by admin
func (s *Service) UpdateUserByAdmin(ctx context.Context, adminUserID, userID uuid.UUID, updates map[string]interface{}) (*User, error) {
	// Admin permissions validated by middleware
	return s.repo.UpdateUserByAdmin(ctx, userID, updates)
}

// BulkImportUsers handles bulk user import
func (s *Service) BulkImportUsers(ctx context.Context, adminUserID uuid.UUID, users []User) (map[string]interface{}, error) {
	// Admin permissions validated by middleware
	successCount := 0
	failedCount := 0
	errors := []string{}

	for _, user := range users {
		if err := s.validateUser(&user); err != nil {
			failedCount++
			errors = append(errors, err.Error())
			continue
		}

		if _, err := s.repo.CreateUser(ctx, &user); err != nil {
			failedCount++
			errors = append(errors, err.Error())
			continue
		}

		successCount++
	}

	return map[string]interface{}{
		"success_count": successCount,
		"failed_count":  failedCount,
		"errors":        errors,
	}, nil
}

// BulkOperations handles bulk user operations
func (s *Service) BulkOperations(ctx context.Context, adminUserID uuid.UUID, operation string, userIDs []uuid.UUID, data interface{}) (map[string]interface{}, error) {
	// Validate admin permissions
	if err := s.validateAdminPermissions(ctx, adminUserID); err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	successCount := 0
	errorCount := 0
	errors := make([]string, 0)

	switch operation {
	case "activate":
		for _, userID := range userIDs {
			if err := s.repo.UpdateUserStatus(ctx, userID, "active"); err != nil {
				errorCount++
				errors = append(errors, fmt.Sprintf("Failed to activate user %s: %v", userID, err))
			} else {
				successCount++
			}
		}
	case "deactivate":
		for _, userID := range userIDs {
			if err := s.repo.UpdateUserStatus(ctx, userID, "inactive"); err != nil {
				errorCount++
				errors = append(errors, fmt.Sprintf("Failed to deactivate user %s: %v", userID, err))
			} else {
				successCount++
			}
		}
	case "delete":
		for _, userID := range userIDs {
			if err := s.repo.DeleteUser(ctx, userID); err != nil {
				errorCount++
				errors = append(errors, fmt.Sprintf("Failed to delete user %s: %v", userID, err))
			} else {
				successCount++
			}
		}
	default:
		return nil, fmt.Errorf("unsupported operation: %s", operation)
	}

	result["success_count"] = successCount
	result["error_count"] = errorCount
	result["errors"] = errors

	// Log activity
	s.logActivity(ctx, adminUserID, "bulk_operation", fmt.Sprintf("Operation: %s, Users: %d, Success: %d, Errors: %d", operation, len(userIDs), successCount, errorCount))

	return result, nil
}

// ExportUsers handles user data export
func (s *Service) ExportUsers(ctx context.Context, adminUserID uuid.UUID, format string, userIDs []uuid.UUID, filters interface{}) (map[string]interface{}, error) {
	// Validate admin permissions
	if err := s.validateAdminPermissions(ctx, adminUserID); err != nil {
		return nil, err
	}

	// Get users based on filters or IDs
	var users []User
	var err error

	if len(userIDs) > 0 {
		users, err = s.repo.GetUsersByIDs(ctx, userIDs)
	} else {
		users, err = s.repo.GetUsersWithFilters(ctx, filters)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	result := make(map[string]interface{})

	switch format {
	case "json":
		result["data"] = users
		result["format"] = "json"
	case "csv":
		// Convert to CSV format
		csvData := s.convertUsersToCSV(users)
		result["data"] = csvData
		result["format"] = "csv"
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	result["count"] = len(users)
	result["exported_at"] = time.Now()

	// Log activity
	s.logActivity(ctx, adminUserID, "export_users", fmt.Sprintf("Format: %s, Count: %d", format, len(users)))

	return result, nil
}

// ManageAccount handles account management operations
func (s *Service) ManageAccount(ctx context.Context, userID uuid.UUID, action string, data interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	switch action {
	case "close":
		// Close account
		if err := s.repo.UpdateUserStatus(ctx, userID, "closed"); err != nil {
			return nil, fmt.Errorf("failed to close account: %w", err)
		}
		result["status"] = "closed"
	case "suspend":
		// Suspend account
		if err := s.repo.UpdateUserStatus(ctx, userID, "suspended"); err != nil {
			return nil, fmt.Errorf("failed to suspend account: %w", err)
		}
		result["status"] = "suspended"
	case "reactivate":
		// Reactivate account
		if err := s.repo.UpdateUserStatus(ctx, userID, "active"); err != nil {
			return nil, fmt.Errorf("failed to reactivate account: %w", err)
		}
		result["status"] = "active"
	default:
		return nil, fmt.Errorf("unsupported action: %s", action)
	}

	// Log activity
	s.logActivity(ctx, userID, "account_management", fmt.Sprintf("Action: %s", action))

	return result, nil
}

// VerifyPhone handles phone number verification
func (s *Service) VerifyPhone(ctx context.Context, phone, otp string) (map[string]interface{}, error) {
	// Verify OTP
	if err := s.verifyOTP(ctx, phone, otp, "phone"); err != nil {
		return nil, fmt.Errorf("invalid OTP: %w", err)
	}

	// Update phone verification status
	if err := s.repo.UpdatePhoneVerification(ctx, phone, true); err != nil {
		return nil, fmt.Errorf("failed to update phone verification: %w", err)
	}

	result := map[string]interface{}{
		"phone_verified": true,
		"verified_at":    time.Now(),
	}

	return result, nil
}

// ResendPhoneOTP handles resending phone OTP
func (s *Service) ResendPhoneOTP(ctx context.Context, phone string) error {
	// Generate and send OTP
	otp := s.generateOTP()
	if err := s.sendSMS(ctx, phone, fmt.Sprintf("Your verification code is: %s", otp)); err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	// Store OTP for verification
	if err := s.storeOTP(ctx, phone, otp, "phone", 10*time.Minute); err != nil {
		return fmt.Errorf("failed to store OTP: %w", err)
	}

	return nil
}

// Helper methods
func (s *Service) convertUsersToCSV(users []User) string {
	// Simple CSV conversion - in production, use proper CSV library
	var csv strings.Builder
	csv.WriteString("ID,Email,FirstName,LastName,Status,CreatedAt\n")
	for _, user := range users {
		csv.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,%s\n",
			user.ID, user.Email, user.FirstName, user.LastName, user.Status, user.CreatedAt.Format(time.RFC3339)))
	}
	return csv.String()
}

func (s *Service) sendSMS(ctx context.Context, phone, message string) error {
	// Implement SMS sending logic
	// This is a placeholder - integrate with SMS provider
	return nil
}

// GetUserOrders gets user's orders
func (s *Service) GetUserOrders(userID uuid.UUID, status string, page, limit int) ([]interface{}, int64, error) {
	// Get user orders (order service integration needed)
	return []interface{}{}, 0, nil
}

// GetUserAddresses gets user's addresses
func (s *Service) GetUserAddresses(userID uuid.UUID) ([]interface{}, error) {
	// Get user addresses (address service integration needed)
	return []interface{}{}, nil
}

// Setup2FA sets up two-factor authentication for a user
func (s *Service) Setup2FA(ctx context.Context, userID uuid.UUID) (*TwoFactorSetupResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.TwoFactorEnabled {
		return nil, errors.New("2FA is already enabled")
	}

	twoFA, qrURL, err := s.securityService.Setup2FA(ctx, userID, user.TenantID, "E-commerce App", user.Email)
	if err != nil {
		return nil, err
	}

	return &TwoFactorSetupResponse{
		Secret:      twoFA.Secret,
		QRCodeURL:   qrURL,
		BackupCodes: twoFA.BackupCodes,
	}, nil
}

// Verify2FA verifies two-factor authentication code
func (s *Service) Verify2FA(ctx context.Context, userID uuid.UUID, code string) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if !user.TwoFactorEnabled {
		return errors.New("2FA is not enabled")
	}

	valid, err := s.securityService.Verify2FA(ctx, userID, user.TenantID, code)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("invalid 2FA code")
	}

	// Log successful 2FA verification
	s.securityService.LogSecurityEvent(ctx, userID, user.TenantID, SecurityEvent2FAVerified, "", "", map[string]interface{}{
			"message": "2FA code verified successfully",
			"email":   user.Email,
		})

	return nil
}

// Enable2FA enables two-factor authentication after verification
func (s *Service) Enable2FA(ctx context.Context, userID uuid.UUID, code string) ([]string, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.TwoFactorEnabled {
		return nil, errors.New("2FA is already enabled")
	}

	// Verify the setup code first
	valid, err := s.securityService.Verify2FA(ctx, userID, user.TenantID, code)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, errors.New("invalid 2FA code")
	}

	// Enable 2FA
	err = s.securityService.Enable2FA(ctx, userID, user.TenantID, code)
	if err != nil {
		return nil, err
	}

	// Update user record
	now := time.Now()
	user.TwoFactorEnabled = true
	user.UpdatedAt = now
	if _, err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	// Get backup codes from 2FA record
	var twoFA TwoFactorAuth
	query := s.repo.(*repository).db.WithContext(ctx).Where("user_id = ?", userID)
	if user.TenantID != nil {
		query = query.Where("tenant_id = ?", *user.TenantID)
	}
	if err := query.First(&twoFA).Error; err != nil {
		return nil, err
	}

	// Log security event
	s.securityService.LogSecurityEvent(ctx, userID, user.TenantID, SecurityEvent2FAEnabled, "", "", map[string]interface{}{
		"message": "Two-factor authentication enabled",
		"email":   user.Email,
	})

	return twoFA.BackupCodes, nil
}

// Disable2FA disables two-factor authentication
func (s *Service) Disable2FA(ctx context.Context, userID uuid.UUID, password string) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if !user.TwoFactorEnabled {
		return errors.New("2FA is not enabled")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return errors.New("invalid password")
	}

	// Disable 2FA
	if err := s.securityService.Disable2FA(ctx, userID); err != nil {
		return err
	}

	// Update user record
	now := time.Now()
	user.TwoFactorEnabled = false
	user.UpdatedAt = now
	if _, err := s.repo.UpdateUser(ctx, user); err != nil {
		return err
	}

	// Log security event
		s.securityService.LogSecurityEvent(ctx, userID, user.TenantID, SecurityEvent2FADisabled, "", "", map[string]interface{}{
			"message": "Two-factor authentication disabled",
			"email":   user.Email,
		})

	// Invalidate all sessions for security
	if err := s.repo.InvalidateUserSessions(ctx, userID); err != nil {
		log.Printf("Failed to invalidate user sessions: %v", err)
	}

	return nil
}

// VerifyBackupCode verifies a 2FA backup code
func (s *Service) VerifyBackupCode(ctx context.Context, userID uuid.UUID, code string) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if !user.TwoFactorEnabled {
		return errors.New("2FA is not enabled")
	}

	valid, err := s.securityService.VerifyBackupCode(ctx, userID, code)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("invalid backup code")
	}

	// Log backup code usage
	s.securityService.LogSecurityEvent(ctx, userID, user.TenantID, SecurityEventBackupCodeUsed, "", "", map[string]interface{}{
		"message": "2FA backup code used",
		"email":   user.Email,
	})

	return nil
}

// GetSecuritySettings gets user security settings
func (s *Service) GetSecuritySettings(ctx context.Context, userID uuid.UUID) (SecuritySettings, error) {
	return s.securityService.GetSecuritySettings(ctx, userID)
}

// UpdateSecuritySettings updates user security settings
func (s *Service) UpdateSecuritySettings(ctx context.Context, userID uuid.UUID, settings SecuritySettings) error {
	return s.securityService.UpdateSecuritySettings(ctx, userID, settings)
}

// GetSecurityLogs gets user security logs
func (s *Service) GetSecurityLogs(ctx context.Context, userID uuid.UUID, limit int) ([]UserSecurityLog, error) {
	return s.securityService.GetSecurityLogs(ctx, userID, limit)
}

// Helper methods

// validateAdminPermissions validates admin permissions
func (s *Service) validateAdminPermissions(ctx context.Context, adminUserID uuid.UUID) error {
	admin, err := s.repo.GetUserByID(ctx, adminUserID)
	if err != nil {
		return err
	}

	if !admin.IsAdmin() {
		return errors.New("insufficient permissions")
	}

	return nil
}

// logActivity logs user activity
func (s *Service) logActivity(ctx context.Context, userID uuid.UUID, activityType, description string) {
	// Log user activity (activity service integration needed)
	log.Printf("User %s activity: %s - %s", userID, activityType, description)
}

// generateOTP generates a random OTP
func (s *Service) generateOTP() string {
	// Generate 6-digit OTP
	return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
}

// verifyOTP verifies an OTP
func (s *Service) verifyOTP(ctx context.Context, identifier, otp, otpType string) error {
	// Verify OTP (OTP service integration needed)
	// This is a placeholder implementation
	if otp == "" {
		return errors.New("OTP is required")
	}
	return nil
}

// storeOTP stores an OTP for verification
func (s *Service) storeOTP(ctx context.Context, identifier, otp, otpType string, expiry time.Duration) error {
	// Store OTP (OTP service integration needed)
	// This is a placeholder implementation
	return nil
}
