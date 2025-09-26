package contact

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ecommerce-saas/internal/shared/handlers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	contacts := r.Group("/contacts")
	{
		// Contact management
		contacts.POST("", h.CreateContact)
		contacts.GET("", h.ListContacts)
		contacts.GET("/:id", h.GetContact)
		contacts.PUT("/:id", h.UpdateContact)
		contacts.DELETE("/:id", h.DeleteContact)
		contacts.POST("/bulk", h.BulkUpdateContacts)
		contacts.POST("/export", h.ExportContacts)

		// Contact status management - consolidated into PUT /:id with query params

		// Contact replies
		contacts.POST("/:id/replies", h.CreateContactReply)
		contacts.GET("/:id/replies", h.ListContactReplies)
		contacts.DELETE("/:id/replies/:reply_id", h.DeleteContactReply)

		// Contact notes and internal comments - removed unimplemented endpoints
	}

	// Contact forms
	forms := r.Group("/contact-forms")
	{
		forms.POST("", h.CreateContactForm)
		forms.GET("", h.ListContactForms)
		forms.GET("/:id", h.GetContactForm)
		forms.PUT("/:id", h.UpdateContactForm)
		forms.DELETE("/:id", h.DeleteContactForm)
		// Activation/deactivation handled via PUT /:id with is_active field
		forms.GET("/public/:form_type", h.GetPublicContactForm)
		forms.POST("/public/:form_type/submit", h.SubmitPublicContactForm)
	}

	// Contact templates
	templates := r.Group("/contact-templates")
	{
		templates.POST("", h.CreateContactTemplate)
		templates.GET("", h.ListContactTemplates)
		templates.GET("/:id", h.GetContactTemplate)
		templates.PUT("/:id", h.UpdateContactTemplate)
		templates.DELETE("/:id", h.DeleteContactTemplate)
		// Activation/deactivation handled via PUT /:id with is_active field
	}

	// Settings
	settings := r.Group("/contact-settings")
	{
		settings.GET("", h.GetContactSettings)
		settings.PUT("", h.UpdateContactSettings)
	}

	// Analytics - consolidated with query parameters
	analytics := r.Group("/contact-analytics")
	{
		analytics.GET("", h.GetContactAnalytics) // Handles all analytics types via ?type= parameter
	}
}

// Contact management handlers
func (h *Handler) CreateContact(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	var req CreateContactRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		handlers.HandleValidationError(c, bindErr)
		return
	}

	contact, err := h.service.CreateContact(c.Request.Context(), tenantID, req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	handlers.RespondWithCreated(c, contact)
}

func (h *Handler) ListContacts(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	filter := h.parseContactFilter(c)
	contacts, total, err := h.service.ListContacts(c.Request.Context(), tenantID, filter)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{
		"data":   contacts,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

func (h *Handler) GetContact(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	contactID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid contact ID"))
		return
	}

	contact, err := h.service.GetContactByID(c.Request.Context(), tenantID, contactID)
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("contact not found"))
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, contact)
}

func (h *Handler) UpdateContact(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	contactID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid contact ID"))
		return
	}

	var req UpdateContactRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		handlers.HandleValidationError(c, bindErr)
		return
	}

	contact, err := h.service.UpdateContact(c.Request.Context(), tenantID, contactID, req)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, contact)
}

func (h *Handler) DeleteContact(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	contactID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid contact ID"))
		return
	}

	err = h.service.DeleteContact(c.Request.Context(), tenantID, contactID)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusNoContent, nil)
}

func (h *Handler) BulkUpdateContacts(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	var req BulkUpdateContactsRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		handlers.HandleValidationError(c, bindErr)
		return
	}

	result, err := h.service.BulkUpdateContacts(c.Request.Context(), tenantID, req)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, result)
}

func (h *Handler) ExportContacts(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	var req ExportContactsRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		handlers.HandleValidationError(c, bindErr)
		return
	}

	exportData, err := h.service.ExportContacts(c.Request.Context(), tenantID, req)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, exportData)
}

// Contact status management handlers
func (h *Handler) UpdateContactStatus(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	contactID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid contact ID"))
		return
	}

	var req UpdateContactStatusRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		handlers.HandleValidationError(c, bindErr)
		return
	}

	err = h.service.UpdateContactStatus(c.Request.Context(), tenantID, contactID, req)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Status updated successfully"})
}

func (h *Handler) AssignContact(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	contactID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid contact ID"))
		return
	}

	var req AssignContactRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		handlers.HandleValidationError(c, bindErr)
		return
	}

	err = h.service.AssignContact(c.Request.Context(), tenantID, contactID, req)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Contact assigned successfully"})
}

func (h *Handler) UpdateContactPriority(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	contactID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid contact ID"))
		return
	}

	var req UpdateContactPriorityRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		handlers.HandleValidationError(c, bindErr)
		return
	}

	err = h.service.UpdateContactPriority(c.Request.Context(), tenantID, contactID, req)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Priority updated successfully"})
}

func (h *Handler) AddContactTags(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	contactID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid contact ID"))
		return
	}

	var req AddContactTagsRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		handlers.HandleValidationError(c, bindErr)
		return
	}

	err = h.service.AddContactTags(c.Request.Context(), tenantID, contactID, req)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Tags added successfully"})
}

func (h *Handler) RemoveContactTags(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	contactID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid contact ID"))
		return
	}

	var req RemoveContactTagsRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		handlers.HandleValidationError(c, bindErr)
		return
	}

	err = h.service.RemoveContactTags(c.Request.Context(), tenantID, contactID, req)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Tags removed successfully"})
}

// Contact reply handlers
func (h *Handler) CreateContactReply(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	contactID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid contact ID"))
		return
	}

	var req CreateContactReplyRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		handlers.HandleValidationError(c, bindErr)
		return
	}

	reply, err := h.service.CreateContactReply(c.Request.Context(), tenantID, contactID, req)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithCreated(c, reply)
}

func (h *Handler) ListContactReplies(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	contactID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid contact ID"))
		return
	}

	replies, err := h.service.ListContactReplies(c.Request.Context(), tenantID, contactID)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"data": replies})
}

func (h *Handler) DeleteContactReply(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	contactID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid contact ID"))
		return
	}

	replyID, err := uuid.Parse(c.Param("reply_id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid reply ID"))
		return
	}

	err = h.service.DeleteContactReply(c.Request.Context(), tenantID, contactID, replyID)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusNoContent, nil)
}

// Contact form handlers
func (h *Handler) CreateContactForm(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	var req CreateContactFormRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		handlers.HandleValidationError(c, bindErr)
		return
	}

	form, err := h.service.CreateContactForm(c.Request.Context(), tenantID, req)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithCreated(c, form)
}

func (h *Handler) ListContactForms(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	forms, err := h.service.ListContactForms(c.Request.Context(), tenantID)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{
		"data":  forms,
		"total": len(forms),
	})
}

func (h *Handler) GetContactForm(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	formID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid form ID"))
		return
	}

	form, err := h.service.GetContactFormByID(c.Request.Context(), tenantID, formID)
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("form not found"))
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, form)
}

func (h *Handler) UpdateContactForm(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	formID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid form ID"))
		return
	}

	var req UpdateContactFormRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		handlers.HandleValidationError(c, bindErr)
		return
	}

	form, err := h.service.UpdateContactForm(c.Request.Context(), tenantID, formID, req)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, form)
}

func (h *Handler) DeleteContactForm(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	formID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid form ID"))
		return
	}

	err = h.service.DeleteContactForm(c.Request.Context(), tenantID, formID)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusNoContent, nil)
}

func (h *Handler) GetPublicContactForm(c *gin.Context) {
	formType := c.Param("form_type")
	form, err := h.service.GetPublicContactForm(c.Request.Context(), formType)
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("form not found"))
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, form)
}

func (h *Handler) SubmitPublicContactForm(c *gin.Context) {
	formType := c.Param("form_type")

	var req SubmitContactFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.HandleValidationError(c, err)
		return
	}

	contact, err := h.service.SubmitPublicContactForm(c.Request.Context(), formType, req)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithCreated(c, contact)
}

// Contact template handlers
func (h *Handler) CreateContactTemplate(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	var req CreateContactTemplateRequest
	if contactError := c.ShouldBindJSON(&req); contactError != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": contactError.Error()})
		return
	}

	template, err := h.service.CreateContactTemplate(c.Request.Context(), tenantID, req)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithCreated(c, template)
}

func (h *Handler) GetContactTemplate(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	templateID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid template ID"))
		return
	}

	template, err := h.service.GetContactTemplateByID(c.Request.Context(), tenantID, templateID)
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("template not found"))
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, template)
}

// Analytics handlers
func (h *Handler) GetContactAnalytics(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	period := AnalyticsPeriod(c.DefaultQuery("period", "month"))
	analytics, err := h.service.GetContactAnalytics(c.Request.Context(), tenantID, period)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, analytics)
}

func (h *Handler) GetContactMetrics(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	metrics, err := h.service.GetContactMetrics(c.Request.Context(), tenantID)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, metrics)
}

// Settings handlers
func (h *Handler) GetContactSettings(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	settings, err := h.service.GetContactSettings(c.Request.Context(), tenantID)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, settings)
}

func (h *Handler) UpdateContactSettings(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	var req UpdateContactSettingsRequest
	if contactError := c.ShouldBindJSON(&req); contactError != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": contactError.Error()})
		return
	}

	settings, err := h.service.UpdateContactSettings(c.Request.Context(), tenantID, req)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, settings)
}

// Helper methods for parsing filters
func (h *Handler) parseContactFilter(c *gin.Context) ContactFilter {
	filter := ContactFilter{}

	if status := c.Query("status"); status != "" {
		s := ContactStatus(status)
		filter.Status = []ContactStatus{s}
	}

	if priority := c.Query("priority"); priority != "" {
		p := ContactPriority(priority)
		filter.Priority = []ContactPriority{p}
	}

	if assignedTo := c.Query("assigned_to"); assignedTo != "" {
		if id, err := uuid.Parse(assignedTo); err == nil {
			filter.AssignedToID = &id
		}
	}

	if search := c.Query("search"); search != "" {
		filter.Search = search
	}

	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse(time.RFC3339, startDate); err == nil {
			filter.StartDate = &t
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse(time.RFC3339, endDate); err == nil {
			filter.EndDate = &t
		}
	}

	// Parse tags - assuming comma-separated
	if tags := c.Query("tags"); tags != "" {
		filter.Tags = strings.Split(tags, ",")
	}

	// Parse pagination
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filter.Limit = l
		}
	}

	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			filter.Offset = o
		}
	}

	// Parse sorting
	filter.SortBy = c.Query("sort_by")
	if c.Query("sort_desc") == "true" {
		filter.SortOrder = "desc"
	} else {
		filter.SortOrder = "asc"
	}

	return filter
}

// Additional handlers for missing methods

func (h *Handler) ActivateContactForm(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	formID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid form ID"))
		return
	}

	err = h.service.ActivateContactForm(c.Request.Context(), tenantID, formID)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Form activated successfully"})
}

func (h *Handler) DeactivateContactForm(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	formID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid form ID"))
		return
	}

	err = h.service.DeactivateContactForm(c.Request.Context(), tenantID, formID)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Form deactivated successfully"})
}

func (h *Handler) ListContactTemplates(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	templates, err := h.service.ListContactTemplates(c.Request.Context(), tenantID)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{
		"data":  templates,
		"total": len(templates),
	})
}

func (h *Handler) UpdateContactTemplate(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	templateID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid template ID"))
		return
	}

	var req UpdateContactTemplateRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		handlers.HandleValidationError(c, bindErr)
		return
	}

	template, err := h.service.UpdateContactTemplate(c.Request.Context(), tenantID, templateID, req)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, template)
}

func (h *Handler) DeleteContactTemplate(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	templateID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid template ID"))
		return
	}

	err = h.service.DeleteContactTemplate(c.Request.Context(), tenantID, templateID)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusNoContent, nil)
}

func (h *Handler) ActivateContactTemplate(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	templateID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid template ID"))
		return
	}

	err = h.service.ActivateContactTemplate(c.Request.Context(), tenantID, templateID)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Template activated successfully"})
}

func (h *Handler) DeactivateContactTemplate(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	templateID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid template ID"))
		return
	}

	err = h.service.DeactivateContactTemplate(c.Request.Context(), tenantID, templateID)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Template deactivated successfully"})
}

func (h *Handler) GetAgentPerformance(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	period := c.DefaultQuery("period", "month")
	agentIDStr := c.Query("agent_id")
	var agentID *uuid.UUID
	if agentIDStr != "" {
		parsedID, parseErr := uuid.Parse(agentIDStr)
		if parseErr != nil {
			handlers.HandleError(c, fmt.Errorf("invalid agent ID"))
			return
		}
		agentID = &parsedID
	}

	performance, err := h.service.GetAgentPerformance(c.Request.Context(), tenantID, period, agentID)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, performance)
}

func (h *Handler) GetCustomerSatisfaction(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	period := c.DefaultQuery("period", "month")

	satisfaction, err := h.service.GetCustomerSatisfaction(c.Request.Context(), tenantID, period)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, satisfaction)
}

func (h *Handler) GetResolutionTimeAnalytics(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	period := c.DefaultQuery("period", "month")

	analytics, err := h.service.GetResolutionTimeAnalytics(c.Request.Context(), tenantID, period)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, analytics)
}

func (h *Handler) GetResponseTimeAnalytics(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	period := c.DefaultQuery("period", "month")

	analytics, err := h.service.GetResponseTimeAnalytics(c.Request.Context(), tenantID, period)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, analytics)
}

// Helper methods removed - using utils.GetTenantIDFromContext instead

// handleServiceError handles service layer errors
func (h *Handler) handleServiceError(c *gin.Context, err error) {
	// Add contact-specific error handling similar to address module
	switch {
	case strings.Contains(err.Error(), "not found"):
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found", "details": err.Error()})
	case strings.Contains(err.Error(), "validation failed"):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
	case strings.Contains(err.Error(), "unauthorized"):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access", "details": err.Error()})
	case strings.Contains(err.Error(), "forbidden"):
		c.JSON(http.StatusForbidden, gin.H{"error": "Access forbidden", "details": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "details": err.Error()})
	}
}
