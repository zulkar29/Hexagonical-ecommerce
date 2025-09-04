package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
)

// String utilities

// TrimAndLower trims whitespace and converts to lowercase
func TrimAndLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// ToSlug converts a string to URL-friendly slug
func ToSlug(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)
	
	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")
	
	// Remove leading/trailing hyphens
	s = strings.Trim(s, "-")
	
	return s
}

// Truncate truncates a string to specified length
func Truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		randomBytes := make([]byte, 1)
		rand.Read(randomBytes)
		b[i] = charset[randomBytes[0]%byte(len(charset))]
	}
	return string(b)
}

// GenerateRandomNumbers generates a random numeric string
func GenerateRandomNumbers(length int) string {
	const charset = "0123456789"
	b := make([]byte, length)
	for i := range b {
		randomBytes := make([]byte, 1)
		rand.Read(randomBytes)
		b[i] = charset[randomBytes[0]%byte(len(charset))]
	}
	return string(b)
}

// Password utilities

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword compares a password with its hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	
	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false
	
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	
	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}
	
	return nil
}

// Validation utilities

// IsValidEmail validates email format
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidPhone validates phone number format (Bangladesh)
func IsValidPhone(phone string) bool {
	// Remove spaces and special characters
	phone = regexp.MustCompile(`[^\d]`).ReplaceAllString(phone, "")
	
	// Check Bangladesh phone number patterns
	if strings.HasPrefix(phone, "880") {
		return len(phone) == 13 // 880xxxxxxxxxx
	}
	if strings.HasPrefix(phone, "01") {
		return len(phone) == 11 // 01xxxxxxxxx
	}
	
	return false
}

// IsValidURL validates URL format
func IsValidURL(url string) bool {
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return urlRegex.MatchString(url)
}

// Formatting utilities

// FormatCurrency formats amount as currency (BDT)
func FormatCurrency(amount float64, currency string) string {
	if currency == "" {
		currency = "BDT"
	}
	
	// Format with 2 decimal places
	formatted := fmt.Sprintf("%.2f", amount)
	
	// Add thousand separators
	parts := strings.Split(formatted, ".")
	intPart := parts[0]
	decPart := parts[1]
	
	// Add commas for thousands
	if len(intPart) > 3 {
		var result []string
		for i, char := range intPart {
			if i > 0 && (len(intPart)-i)%3 == 0 {
				result = append(result, ",")
			}
			result = append(result, string(char))
		}
		intPart = strings.Join(result, "")
	}
	
	return fmt.Sprintf("%s %s.%s", currency, intPart, decPart)
}

// FormatPhone formats phone number for display
func FormatPhone(phone string) string {
	// Remove all non-digits
	phone = regexp.MustCompile(`[^\d]`).ReplaceAllString(phone, "")
	
	if len(phone) == 11 && strings.HasPrefix(phone, "01") {
		// Format as: 01XXX-XXXXXX
		return fmt.Sprintf("%s-%s", phone[:5], phone[5:])
	}
	if len(phone) == 13 && strings.HasPrefix(phone, "880") {
		// Format as: +880-1XXX-XXXXXX
		return fmt.Sprintf("+%s-%s-%s", phone[:3], phone[3:7], phone[7:])
	}
	
	return phone
}

// Time utilities

// FormatDate formats date for Bangladesh timezone
func FormatDate(t time.Time) string {
	// Convert to Dhaka timezone
	dhaka, _ := time.LoadLocation("Asia/Dhaka")
	t = t.In(dhaka)
	return t.Format("02-01-2006")
}

// FormatDateTime formats datetime for Bangladesh timezone
func FormatDateTime(t time.Time) string {
	dhaka, _ := time.LoadLocation("Asia/Dhaka")
	t = t.In(dhaka)
	return t.Format("02-01-2006 15:04:05")
}

// TimeAgo returns human-readable time difference
func TimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)
	
	switch {
	case diff < time.Minute:
		return "Just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case diff < 30*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		return FormatDate(t)
	}
}

// Hash utilities

// GenerateHash generates SHA256 hash of input
func GenerateHash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

// GenerateUUID generates a new UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// Number utilities

// ParseFloat safely parses string to float64
func ParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}

// ParseInt safely parses string to int
func ParseInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}

// RoundToTwoDecimals rounds float to 2 decimal places
func RoundToTwoDecimals(f float64) float64 {
	return float64(int(f*100+0.5)) / 100
}

// Array utilities

// Contains checks if slice contains an element
func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Unique returns unique elements from slice
func Unique[T comparable](slice []T) []T {
	keys := make(map[T]bool)
	var result []T
	
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// TODO: Add more utility functions
// - Image processing utilities
// - File upload utilities
// - Cache key generators
// - API response helpers
// - Pagination helpers
