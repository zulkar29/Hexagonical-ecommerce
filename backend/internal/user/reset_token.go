package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ResetToken represents a password reset token stored in database
type ResetToken struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Token     string    `json:"token" gorm:"unique;not null;index"`
	Type      TokenType `json:"type" gorm:"not null;index"` // password_reset, email_verification
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`
	UsedAt    *time.Time `json:"used_at" gorm:"index"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// TokenType represents the type of token
type TokenType string

const (
	TokenTypePasswordReset     TokenType = "password_reset"
	TokenTypeEmailVerification TokenType = "email_verification"
)

// TokenRepository interface for reset token operations
type TokenRepository interface {
	CreateToken(ctx context.Context, token *ResetToken) error
	GetValidToken(ctx context.Context, tokenValue string, tokenType TokenType) (*ResetToken, error)
	MarkTokenAsUsed(ctx context.Context, tokenID uuid.UUID) error
	DeleteExpiredTokens(ctx context.Context) error
	DeleteUserTokens(ctx context.Context, userID uuid.UUID, tokenType TokenType) error
}

// tokenRepository implements TokenRepository
type tokenRepository struct {
	db *gorm.DB
}

// NewTokenRepository creates a new token repository
func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

// CreateToken creates a new reset token
func (r *tokenRepository) CreateToken(ctx context.Context, token *ResetToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// GetValidToken retrieves a valid (unused and not expired) token
func (r *tokenRepository) GetValidToken(ctx context.Context, tokenValue string, tokenType TokenType) (*ResetToken, error) {
	var token ResetToken
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("token = ? AND type = ? AND used_at IS NULL AND expires_at > ?",
			tokenValue, tokenType, time.Now()).
		First(&token).Error

	if err != nil {
		return nil, err
	}

	return &token, nil
}

// MarkTokenAsUsed marks a token as used
func (r *tokenRepository) MarkTokenAsUsed(ctx context.Context, tokenID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&ResetToken{}).
		Where("id = ?", tokenID).
		Update("used_at", now).Error
}

// DeleteExpiredTokens removes all expired tokens
func (r *tokenRepository) DeleteExpiredTokens(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&ResetToken{}).Error
}

// DeleteUserTokens deletes all tokens of a specific type for a user
func (r *tokenRepository) DeleteUserTokens(ctx context.Context, userID uuid.UUID, tokenType TokenType) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, tokenType).
		Delete(&ResetToken{}).Error
}

// GenerateSecureToken generates a cryptographically secure token
func GenerateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// TokenExpirationDuration returns the expiration duration for different token types
func TokenExpirationDuration(tokenType TokenType) time.Duration {
	switch tokenType {
	case TokenTypePasswordReset:
		return 1 * time.Hour // 1 hour for password reset
	case TokenTypeEmailVerification:
		return 24 * time.Hour // 24 hours for email verification
	default:
		return 1 * time.Hour
	}
}