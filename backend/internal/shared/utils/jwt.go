package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWT claims structure
type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	TenantID *uuid.UUID `json:"tenant_id,omitempty"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	TokenID  string    `json:"token_id"` // For token revocation
	jwt.RegisteredClaims
}

// JWTManager handles JWT operations
type JWTManager struct {
	secretKey []byte
	issuer    string
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secretKey, issuer string) *JWTManager {
	return &JWTManager{
		secretKey: []byte(secretKey),
		issuer:    issuer,
	}
}

// GenerateTokens generates access and refresh tokens
func (j *JWTManager) GenerateTokens(userID uuid.UUID, tenantID *uuid.UUID, email, role string) (string, string, error) {
	tokenID := uuid.New().String()
	
	// Access token (15 minutes)
	accessClaims := &Claims{
		UserID:   userID,
		TenantID: tenantID,
		Email:    email,
		Role:     role,
		TokenID:  tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
			Subject:   userID.String(),
			ID:        tokenID,
		},
	}
	
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(j.secretKey)
	if err != nil {
		return "", "", err
	}
	
	// Refresh token (7 days)
	refreshClaims := &Claims{
		UserID:   userID,
		TenantID: tenantID,
		Email:    email,
		Role:     role,
		TokenID:  tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
			Subject:   userID.String(),
			ID:        tokenID,
		},
	}
	
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(j.secretKey)
	if err != nil {
		return "", "", err
	}
	
	return accessTokenString, refreshTokenString, nil
}

// ValidateToken validates and parses a JWT token
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secretKey, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("invalid token")
}

// RefreshToken generates new tokens from a valid refresh token
func (j *JWTManager) RefreshToken(refreshTokenString string) (string, string, error) {
	claims, err := j.ValidateToken(refreshTokenString)
	if err != nil {
		return "", "", err
	}
	
	// Check if it's actually a refresh token (longer expiry)
	if claims.ExpiresAt.Time.Sub(time.Now()) < 6*24*time.Hour {
		return "", "", errors.New("invalid refresh token")
	}
	
	return j.GenerateTokens(claims.UserID, claims.TenantID, claims.Email, claims.Role)
}

// GenerateResetToken generates a password reset token
func GenerateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GenerateVerificationToken generates an email verification token
func GenerateVerificationToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
