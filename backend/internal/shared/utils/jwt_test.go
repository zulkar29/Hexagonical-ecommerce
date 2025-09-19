package utils

import (
	"testing"

	"github.com/google/uuid"
)

func TestJWTManager_GenerateTokens(t *testing.T) {
	// Setup
	jwtManager := NewJWTManager("test-secret-key", "test-issuer")
	userID := uuid.New()
	tenantID := uuid.New()
	email := "test@example.com"
	role := "admin"

	// Test
	accessToken, refreshToken, err := jwtManager.GenerateTokens(userID, &tenantID, email, role)

	// Assertions
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if accessToken == "" {
		t.Error("Expected access token, got empty string")
	}

	if refreshToken == "" {
		t.Error("Expected refresh token, got empty string")
	}

	// Tokens should be different
	if accessToken == refreshToken {
		t.Error("Access token and refresh token should be different")
	}
}

func TestJWTManager_ValidateToken(t *testing.T) {
	// Setup
	jwtManager := NewJWTManager("test-secret-key", "test-issuer")
	userID := uuid.New()
	tenantID := uuid.New()
	email := "test@example.com"
	role := "admin"

	// Generate token
	accessToken, _, err := jwtManager.GenerateTokens(userID, &tenantID, email, role)
	if err != nil {
		t.Fatalf("Failed to generate tokens: %v", err)
	}

	// Test - Validate valid token
	claims, err := jwtManager.ValidateToken(accessToken)

	// Assertions
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected user ID %v, got %v", userID, claims.UserID)
	}

	if claims.TenantID == nil || *claims.TenantID != tenantID {
		t.Errorf("Expected tenant ID %v, got %v", tenantID, claims.TenantID)
	}

	if claims.Email != email {
		t.Errorf("Expected email %s, got %s", email, claims.Email)
	}

	if claims.Role != role {
		t.Errorf("Expected role %s, got %s", role, claims.Role)
	}
}

func TestJWTManager_ValidateToken_InvalidToken(t *testing.T) {
	// Setup
	jwtManager := NewJWTManager("test-secret-key", "test-issuer")

	// Test cases
	testCases := []struct {
		name  string
		token string
	}{
		{"Empty token", ""},
		{"Invalid token", "invalid.token.here"},
		{"Malformed token", "header.payload"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test
			claims, err := jwtManager.ValidateToken(tc.token)

			// Assertions
			if err == nil {
				t.Error("Expected error for invalid token")
			}

			if claims != nil {
				t.Error("Expected nil claims for invalid token")
			}
		})
	}
}

func TestJWTManager_ValidateToken_ExpiredToken(t *testing.T) {
	// This test would require manipulating time or creating a token with past expiration
	// For now, we'll skip it or implement with a time mock
	t.Skip("Expired token test requires time manipulation")
}

func TestJWTManager_GenerateRandomString(t *testing.T) {
	// Test the random string generation
	length := 32
	str1 := GenerateRandomString(length)
	str2 := GenerateRandomString(length)

	// Assertions
	if len(str1) == 0 {
		t.Error("Expected non-empty string")
	}

	if len(str2) == 0 {
		t.Error("Expected non-empty string")
	}

	if str1 == str2 {
		t.Error("Expected different random strings")
	}
}

func TestJWTManager_DifferentSecrets(t *testing.T) {
	// Setup with different secrets
	jwtManager1 := NewJWTManager("secret1", "issuer")
	jwtManager2 := NewJWTManager("secret2", "issuer")

	userID := uuid.New()
	email := "test@example.com"
	role := "admin"

	// Generate token with first manager
	token, _, err := jwtManager1.GenerateTokens(userID, nil, email, role)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Try to validate with second manager (should fail)
	_, err = jwtManager2.ValidateToken(token)
	if err == nil {
		t.Error("Expected error when validating token with different secret")
	}
}

func TestJWTManager_NilTenantID(t *testing.T) {
	// Setup
	jwtManager := NewJWTManager("test-secret-key", "test-issuer")
	userID := uuid.New()
	email := "test@example.com"
	role := "super_admin"

	// Test with nil tenant ID (for super admin)
	accessToken, refreshToken, err := jwtManager.GenerateTokens(userID, nil, email, role)

	// Assertions
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if accessToken == "" {
		t.Error("Expected access token")
	}

	if refreshToken == "" {
		t.Error("Expected refresh token")
	}

	// Validate token
	claims, err := jwtManager.ValidateToken(accessToken)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.TenantID != nil {
		t.Error("Expected nil tenant ID for super admin")
	}
}

// Benchmark tests
func BenchmarkJWTManager_GenerateTokens(b *testing.B) {
	jwtManager := NewJWTManager("test-secret-key", "test-issuer")
	userID := uuid.New()
	tenantID := uuid.New()
	email := "test@example.com"
	role := "admin"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := jwtManager.GenerateTokens(userID, &tenantID, email, role)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJWTManager_ValidateToken(b *testing.B) {
	jwtManager := NewJWTManager("test-secret-key", "test-issuer")
	userID := uuid.New()
	tenantID := uuid.New()
	email := "test@example.com"
	role := "admin"

	token, _, err := jwtManager.GenerateTokens(userID, &tenantID, email, role)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := jwtManager.ValidateToken(token)
		if err != nil {
			b.Fatal(err)
		}
	}
}
