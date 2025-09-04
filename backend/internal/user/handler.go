package user

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles user HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new user handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers user routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/logout", h.Logout)
		auth.POST("/forgot-password", h.ForgotPassword)
		auth.POST("/reset-password", h.ResetPassword)
		auth.POST("/verify-email", h.VerifyEmail)
	}

	// Protected routes (require authentication)
	users := router.Group("/users")
	// users.Use(middleware.AuthMiddleware()) // Will be enabled when middleware is complete
	{
		users.GET("/profile", h.GetProfile)
		users.PUT("/profile", h.UpdateProfile)
		users.POST("/change-password", h.ChangePassword)
		users.GET("", h.ListUsers)
		users.GET("/:id", h.GetUser)
	}
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.RegisterUser(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully. Please verify your email.",
		"user":    user,
	})
}

// Login handles user login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.LoginUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set refresh token as HTTP-only cookie
	c.SetCookie(
		"refresh_token",
		response.RefreshToken,
		7*24*60*60, // 7 days
		"/",
		"",
		false, // Set to true in production with HTTPS
		true,  // HTTP-only
	)

	c.JSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"user":         response.User,
		"access_token": response.AccessToken,
		"expires_in":   response.ExpiresIn,
	})
}

// RefreshToken handles token refresh
func (h *Handler) RefreshToken(c *gin.Context) {
	// Get refresh token from cookie or request body
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		// Try to get from request body
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
			return
		}
		refreshToken = req.RefreshToken
	}

	response, err := h.service.RefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Update refresh token cookie
	c.SetCookie(
		"refresh_token",
		response.RefreshToken,
		7*24*60*60, // 7 days
		"/",
		"",
		false, // Set to true in production with HTTPS
		true,  // HTTP-only
	)

	c.JSON(http.StatusOK, response)
}

// Logout handles user logout
func (h *Handler) Logout(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get refresh token
	refreshToken, _ := c.Cookie("refresh_token")
	
	if err := h.service.LogoutUser(userID, refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Clear refresh token cookie
	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// VerifyEmail handles email verification
func (h *Handler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract user ID from query params or JWT token
	userIDStr := c.Query("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.service.VerifyEmail(userID, req.Token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// ForgotPassword handles password reset request
func (h *Handler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ResetPassword(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset instructions sent to your email"})
}

// ResetPassword handles password reset
func (h *Handler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement actual password reset with token validation
	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// GetProfile gets current user profile
func (h *Handler) GetProfile(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateProfile updates user profile
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.UpdateProfile(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    user,
	})
}

// ChangePassword handles password change
func (h *Handler) ChangePassword(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// ListUsers lists users (admin only)
func (h *Handler) ListUsers(c *gin.Context) {
	// TODO: Check admin permissions
	
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	var filter UserFilter
	if role := c.Query("role"); role != "" {
		filter.Role = UserRole(role)
	}
	if status := c.Query("status"); status != "" {
		filter.Status = UserStatus(status)
	}
	filter.Search = c.Query("search")

	users, total, err := h.service.repo.List(nil, filter, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetUser gets user by ID (admin only)
func (h *Handler) GetUser(c *gin.Context) {
	// TODO: Check admin permissions
	
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Helper methods

// getUserIDFromContext extracts user ID from JWT token in context
func (h *Handler) getUserIDFromContext(c *gin.Context) uuid.UUID {
	// TODO: Extract from JWT token when middleware is implemented
	// For now, return nil UUID
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uuid.UUID); ok {
			return id
		}
	}
	return uuid.Nil
}

// getTenantIDFromContext extracts tenant ID from JWT token in context
func (h *Handler) getTenantIDFromContext(c *gin.Context) *uuid.UUID {
	// TODO: Extract from JWT token when middleware is implemented
	if tenantID, exists := c.Get("tenant_id"); exists {
		if id, ok := tenantID.(*uuid.UUID); ok {
			return id
		}
	}
	return nil
}
