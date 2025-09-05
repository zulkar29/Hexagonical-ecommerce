package user

import (
	"net/http"
	"strconv"

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
		auth.POST("/resend-verification", h.ResendVerification)
	}

	// Protected routes (require authentication)
	users := router.Group("/users")
	// users.Use(middleware.AuthMiddleware()) // Will be enabled when middleware is complete
	{
		// User profile management
		users.GET("/profile", h.GetProfile)
		users.PUT("/profile", h.UpdateProfile)
		users.POST("/change-password", h.ChangePassword)
		users.DELETE("/account", h.DeleteAccount)
		
		// User preferences
		users.GET("/preferences", h.GetPreferences)
		users.PUT("/preferences", h.UpdatePreferences)
		
		// Admin user management
		users.GET("", h.ListUsers)
		users.GET("/:id", h.GetUser)
		users.GET("/:id/activity", h.GetUserActivity)
		users.PATCH("/:id", h.UpdateUser)
		users.POST("/bulk-import", h.BulkImportUsers)
		users.POST("/export", h.ExportUsers)
		
		// User related data
		users.GET("/:id/orders", h.GetUserOrders)
		users.GET("/:id/addresses", h.GetUserAddresses)
	}
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := h.service.RegisterUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully. Please verify your email.",
		"user":    createdUser,
	})
}

// Login handles user login
func (h *Handler) Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.LoginUser(c.Request.Context(), loginData.Email, loginData.Password)
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
	var verifyData struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&verifyData); err != nil {
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

	if err := h.service.VerifyEmail(userID, verifyData.Token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// ForgotPassword handles password reset request
func (h *Handler) ForgotPassword(c *gin.Context) {
	var forgotData struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&forgotData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ForgotPassword(c.Request.Context(), forgotData.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset instructions sent to your email"})
}

// ResetPassword handles password reset
func (h *Handler) ResetPassword(c *gin.Context) {
	var resetData struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&resetData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), resetData.Token, resetData.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	var updates User
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.UpdateProfile(c.Request.Context(), userID, &updates)
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

	var passwordData struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&passwordData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ChangePassword(c.Request.Context(), userID, passwordData.OldPassword, passwordData.NewPassword); err != nil {
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

// ResendVerification resends email verification
func (h *Handler) ResendVerification(c *gin.Context) {
	var resendData struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&resendData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ResendVerification(c.Request.Context(), resendData.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent successfully"})
}

// DeleteAccount handles user account deletion
func (h *Handler) DeleteAccount(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var deleteData struct {
		Password string `json:"password" binding:"required"`
		Confirm  bool   `json:"confirm" binding:"required"`
	}
	if err := c.ShouldBindJSON(&deleteData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !deleteData.Confirm {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account deletion must be confirmed"})
		return
	}

	if err := h.service.DeleteAccount(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

// GetPreferences gets user preferences
func (h *Handler) GetPreferences(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	preferences, err := h.service.GetUserPreferences(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"preferences": preferences})
}

// UpdatePreferences updates user preferences
func (h *Handler) UpdatePreferences(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var preferences map[string]interface{}
	if err := c.ShouldBindJSON(&preferences); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPreferences, err := h.service.UpdateUserPreferences(c.Request.Context(), userID, preferences)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Preferences updated successfully",
		"preferences": updatedPreferences,
	})
}

// GetUserActivity gets user activity logs (admin only)
func (h *Handler) GetUserActivity(c *gin.Context) {
	// TODO: Check admin permissions
	
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	activityType := c.Query("type")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	activity, total, err := h.service.GetUserActivity(userID, activityType, dateFrom, dateTo, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"activity": activity,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// UpdateUser updates user (admin only)
func (h *Handler) UpdateUser(c *gin.Context) {
	// TODO: Check admin permissions
	
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminUserID := h.getUserIDFromContext(c)
	updatedUser, err := h.service.UpdateUserByAdmin(c.Request.Context(), adminUserID, userID, updates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    updatedUser,
	})
}

// BulkImportUsers handles bulk user import (admin only)
func (h *Handler) BulkImportUsers(c *gin.Context) {
	// TODO: Check admin permissions
	
	var importData struct {
		Users []User `json:"users" binding:"required"`
	}
	if err := c.ShouldBindJSON(&importData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminUserID := h.getUserIDFromContext(c)
	result, err := h.service.BulkImportUsers(c.Request.Context(), adminUserID, importData.Users)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bulk import completed",
		"result":  result,
	})
}

// ExportUsers handles user data export (admin only)
func (h *Handler) ExportUsers(c *gin.Context) {
	// TODO: Check admin permissions
	
	format := c.DefaultQuery("format", "csv")
	filters := map[string]string{
		"role":   c.Query("role"),
		"status": c.Query("status"),
		"search": c.Query("search"),
	}

	adminUserID := h.getUserIDFromContext(c)
	exportData, err := h.service.ExportUsers(c.Request.Context(), adminUserID, format, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set appropriate headers for file download
	filename := "users_export." + format
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.Data(http.StatusOK, "application/octet-stream", exportData)
}

// GetUserOrders gets user's orders (admin only)
func (h *Handler) GetUserOrders(c *gin.Context) {
	// TODO: Check admin permissions
	
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	orders, total, err := h.service.GetUserOrders(userID, status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetUserAddresses gets user's addresses (admin only)
func (h *Handler) GetUserAddresses(c *gin.Context) {
	// TODO: Check admin permissions
	
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	addresses, err := h.service.GetUserAddresses(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"addresses": addresses})
}
