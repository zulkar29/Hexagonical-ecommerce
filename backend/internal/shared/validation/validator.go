package validation

import (
	"errors"
	"fmt"
	"net/mail"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Validator interface for entities that can be validated
type Validator interface {
	Validate() error
}

// ValidationRule represents a validation rule function
type ValidationRule func(interface{}) error

// FieldValidator provides field-level validation utilities
type FieldValidator struct {
	validator *validator.Validate
}

// NewFieldValidator creates a new field validator with custom rules
func NewFieldValidator() *FieldValidator {
	v := validator.New()

	// Register custom validators
	v.RegisterValidation("uuid", validateUUID)
	v.RegisterValidation("slug", validateSlug)
	v.RegisterValidation("phone", validatePhone)

	return &FieldValidator{validator: v}
}

// ValidateStruct validates a struct using struct tags
func (fv *FieldValidator) ValidateStruct(s interface{}) error {
	if err := fv.validator.Struct(s); err != nil {
		return fv.formatValidationErrors(err)
	}
	return nil
}

// formatValidationErrors converts validator errors to readable format
func (fv *FieldValidator) formatValidationErrors(err error) error {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var errorMessages []string
		for _, e := range validationErrors {
			errorMessages = append(errorMessages, fv.getErrorMessage(e))
		}
		return errors.New(strings.Join(errorMessages, "; "))
	}
	return err
}

// getErrorMessage returns a user-friendly error message for validation errors
func (fv *FieldValidator) getErrorMessage(e validator.FieldError) string {
	field := strings.ToLower(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, e.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "slug":
		return fmt.Sprintf("%s must be a valid slug (lowercase, alphanumeric, hyphens only)", field)
	case "phone":
		return fmt.Sprintf("%s must be a valid phone number", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, e.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// Common validation functions
func ValidateRequired(value interface{}, fieldName string) error {
	if value == nil {
		return fmt.Errorf("%s is required", fieldName)
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		if strings.TrimSpace(v.String()) == "" {
			return fmt.Errorf("%s is required", fieldName)
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		if v.Len() == 0 {
			return fmt.Errorf("%s is required", fieldName)
		}
	}

	return nil
}

func ValidateEmail(email, fieldName string) error {
	if email == "" {
		return nil // Allow empty, use ValidateRequired for required fields
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("%s must be a valid email address", fieldName)
	}

	return nil
}

func ValidateUUID(id, fieldName string) error {
	if id == "" {
		return nil // Allow empty, use ValidateRequired for required fields
	}

	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%s must be a valid UUID", fieldName)
	}

	return nil
}

func ValidateLength(value string, min, max int, fieldName string) error {
	length := len(strings.TrimSpace(value))

	if min > 0 && length < min {
		return fmt.Errorf("%s must be at least %d characters long", fieldName, min)
	}

	if max > 0 && length > max {
		return fmt.Errorf("%s must be at most %d characters long", fieldName, max)
	}

	return nil
}

func ValidatePositiveNumber(value float64, fieldName string) error {
	if value < 0 {
		return fmt.Errorf("%s must be a positive number", fieldName)
	}
	return nil
}

func ValidateSlug(slug, fieldName string) error {
	if slug == "" {
		return nil // Allow empty, use ValidateRequired for required fields
	}

	// Slug should be lowercase, alphanumeric with hyphens
	matched, _ := regexp.MatchString(`^[a-z0-9]+(?:-[a-z0-9]+)*$`, slug)
	if !matched {
		return fmt.Errorf("%s must be a valid slug (lowercase, alphanumeric, hyphens only)", fieldName)
	}

	return nil
}

func ValidateDateRange(startDate, endDate time.Time, fieldName string) error {
	if !endDate.IsZero() && !startDate.IsZero() && endDate.Before(startDate) {
		return fmt.Errorf("%s end date must be after start date", fieldName)
	}
	return nil
}

func ValidateOneOf(value string, allowedValues []string, fieldName string) error {
	if value == "" {
		return nil // Allow empty, use ValidateRequired for required fields
	}

	for _, allowed := range allowedValues {
		if value == allowed {
			return nil
		}
	}

	return fmt.Errorf("%s must be one of: %s", fieldName, strings.Join(allowedValues, ", "))
}

// Custom validator functions for go-playground/validator
func validateUUID(fl validator.FieldLevel) bool {
	_, err := uuid.Parse(fl.Field().String())
	return err == nil
}

func validateSlug(fl validator.FieldLevel) bool {
	slug := fl.Field().String()
	if slug == "" {
		return true // Allow empty, required tag handles this
	}
	matched, _ := regexp.MatchString(`^[a-z0-9]+(?:-[a-z0-9]+)*$`, slug)
	return matched
}

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true // Allow empty, required tag handles this
	}
	// Basic phone validation - can be enhanced based on requirements
	matched, _ := regexp.MatchString(`^\+?[1-9]\d{1,14}$`, phone)
	return matched
}