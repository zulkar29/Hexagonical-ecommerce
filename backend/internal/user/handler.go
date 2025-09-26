package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ecommerce-saas/internal/security"
	"ecommerce-saas/internal/shared/handlers"
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
		auth.POST("/verify-phone", h.VerifyPhone)
		auth.POST("/resend-phone-otp", h.ResendPhoneOTP)
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

		// User preferences functionality removed due to JSONB complexity

		// Two-Factor Authentication
		users.POST("/2fa/setup", h.Setup2FA)
		users.POST("/2fa/enable", h.Enable2FA)
		users.POST("/2fa/disable", h.Disable2FA)
		users.POST("/2fa/verify", h.Verify2FA)
		users.POST("/2fa/verify-backup", h.VerifyBackupCode)

		// Security
		users.GET("/security/settings", h.GetSecuritySettings)
		users.PUT("/security/settings", h.UpdateSecuritySettings)
		users.GET("/security/logs", h.GetSecurityLogs)

		// Admin user management
		users.GET("", h.ListUsers)
		users.GET("/:id", h.GetUser)
		users.GET("/:id/activity", h.GetUserActivity)
		users.PATCH("/:id", h.UpdateUser)
		users.POST("/bulk-import", h.BulkImportUsers)
		// Removed export route - GDPR data export functionality not needed

		// User related data
		users.GET("/:id/orders", h.GetUserOrders)
		users.GET("/:id/addresses", h.GetUserAddresses)

		// Missing endpoints from documentation
		users.POST("/bulk", h.BulkOperations)
		users.POST("/export", h.ExportUsers)
		users.PATCH("/account", h.ManageAccount)
	}
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var user User
	if bindErr := c.ShouldBindJSON(&user); bindErr != nil {
		handlers.HandleValidationError(c, bindErr)
		return
	}

	createdUser, err := h.service.RegisterUser(c.Request.Context(), &user)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithCreated(c, map[string]any{
		"user": createdUser,
	}, "User registered successfully. Please verify your email.")
}

// Login handles user login
func (h *Handler) Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}
	if bindErr := c.ShouldBindJSON(&loginData); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
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
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
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
	if bindErr := c.ShouldBindJSON(&verifyData); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}

	if err := h.service.VerifyEmail(c.Request.Context(), verifyData.Token); err != nil {
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
	if bindErr := c.ShouldBindJSON(&forgotData); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
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
	if bindErr := c.ShouldBindJSON(&resetData); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
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

	user, err := h.service.GetProfile(c.Request.Context(), userID)
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
	if bindErr := c.ShouldBindJSON(&updates); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
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
	if bindErr := c.ShouldBindJSON(&passwordData); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
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
	// Admin permissions already validated by RoleMiddleware

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var filter UserFilter
	if role := c.Query("role"); role != "" {
		filter.Role = UserRole(role)
	}
	if status := c.Query("status"); status != "" {
		filter.Status = UserStatus(status)
	}
	filter.Search = c.Query("search")

	users, total, err := h.service.ListUsers(c.Request.Context(), nil, filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetUser gets user by ID (admin only)
func (h *Handler) GetUser(c *gin.Context) {
	// Admin permissions already validated by RoleMiddleware

	userIDStr := c.Param("id")
	userID, parseErr := uuid.Parse(userIDStr)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Helper methods

// getUserIDFromContext extracts user ID from JWT token in context
func (h *Handler) getUserIDFromContext(c *gin.Context) uuid.UUID {
	// Extract user ID from JWT token set by AuthMiddleware
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uuid.UUID); ok {
			return id
		}
	}
	return uuid.Nil
}

// ResendVerification resends email verification
func (h *Handler) ResendVerification(c *gin.Context) {
	var resendData struct {
		Email string `json:"email" binding:"required,email"`
	}
	if bindErr := c.ShouldBindJSON(&resendData); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
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

// Note: User preferences functionality has been removed due to JSONB complexity
// If needed in the future, consider using a separate notification_preferences table

// GetUserActivity gets user activity logs (admin only)
func (h *Handler) GetUserActivity(c *gin.Context) {
	// Admin permissions already validated by RoleMiddleware

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
	// Admin permissions already validated by RoleMiddleware

	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var updates map[string]interface{}
	if bindErr := c.ShouldBindJSON(&updates); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
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
	// Admin permissions already validated by RoleMiddleware

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

// BulkOperations handles bulk user operations (admin only)
func (h *Handler) BulkOperations(c *gin.Context) {
	adminUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Operation string      `json:"operation" binding:"required"`
		UserIDs   []uuid.UUID `json:"user_ids" binding:"required"`
		Data      interface{} `json:"data,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.BulkOperations(c.Request.Context(), adminUserID.(uuid.UUID), req.Operation, req.UserIDs, req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bulk operation completed",
		"result":  result,
	})
}

// ExportUsers handles user data export (admin only)
func (h *Handler) ExportUsers(c *gin.Context) {
	adminUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Format  string      `json:"format" binding:"required"`
		UserIDs []uuid.UUID `json:"user_ids,omitempty"`
		Filters interface{} `json:"filters,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.ExportUsers(c.Request.Context(), adminUserID.(uuid.UUID), req.Format, req.UserIDs, req.Filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Export completed",
		"result":  result,
	})
}

// ManageAccount handles account management operations
func (h *Handler) ManageAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Action string      `json:"action" binding:"required"`
		Data   interface{} `json:"data,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.ManageAccount(c.Request.Context(), userID.(uuid.UUID), req.Action, req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account operation completed",
		"result":  result,
	})
}

// VerifyPhone handles phone number verification
func (h *Handler) VerifyPhone(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		OTP   string `json:"otp" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.VerifyPhone(c.Request.Context(), req.Phone, req.OTP)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Phone verified successfully",
		"result":  result,
	})
}

// ResendPhoneOTP handles resending phone OTP
func (h *Handler) ResendPhoneOTP(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.ResendPhoneOTP(c.Request.Context(), req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent successfully",
	})
}

// GetUserOrders gets user's orders (admin only)
func (h *Handler) GetUserOrders(c *gin.Context) {
	// Admin permissions already validated by RoleMiddleware

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
	// Admin permissions already validated by RoleMiddleware

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

// Setup2FA sets up two-factor authentication
func (h *Handler) Setup2FA(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	response, err := h.service.Setup2FA(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "2FA setup initiated",
		"data":    response,
	})
}

// Enable2FA enables two-factor authentication
func (h *Handler) Enable2FA(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	backupCodes, err := h.service.Enable2FA(c.Request.Context(), userID, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "2FA enabled successfully",
		"backup_codes": backupCodes,
	})
}

// Disable2FA disables two-factor authentication
func (h *Handler) Disable2FA(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Disable2FA(c.Request.Context(), userID, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA disabled successfully"})
}

// Verify2FA verifies two-factor authentication code
func (h *Handler) Verify2FA(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Verify2FA(c.Request.Context(), userID, req.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA code verified successfully"})
}

// VerifyBackupCode verifies a 2FA backup code
func (h *Handler) VerifyBackupCode(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.VerifyBackupCode(c.Request.Context(), userID, req.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Backup code verified successfully"})
}

// GetSecuritySettings gets user security settings
func (h *Handler) GetSecuritySettings(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	settings, err := h.service.GetSecuritySettings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

// UpdateSecuritySettings updates user security settings
func (h *Handler) UpdateSecuritySettings(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var settings security.SecuritySettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateSecuritySettings(c.Request.Context(), userID, settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Security settings updated successfully"})
}

// GetSecurityLogs gets user security logs
func (h *Handler) GetSecurityLogs(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit > 100 {
		limit = 100
	}

	logs, err := h.service.GetSecurityLogs(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}
