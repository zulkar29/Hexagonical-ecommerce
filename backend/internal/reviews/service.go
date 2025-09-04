package reviews

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Service defines the reviews service interface
type Service interface {
	// Review operations
	CreateReview(ctx context.Context, req CreateReviewRequest) (*Review, error)
	GetReview(ctx context.Context, tenantID, reviewID uuid.UUID) (*Review, error)
	GetReviews(ctx context.Context, tenantID uuid.UUID, filter ReviewFilter) ([]Review, error)
	GetProductReviews(ctx context.Context, tenantID, productID uuid.UUID, filter ProductReviewFilter) ([]Review, error)
	UpdateReview(ctx context.Context, tenantID, reviewID uuid.UUID, req UpdateReviewRequest) (*Review, error)
	DeleteReview(ctx context.Context, tenantID, reviewID uuid.UUID) error
	
	// Review moderation
	ApproveReview(ctx context.Context, tenantID, reviewID uuid.UUID, moderatorID uuid.UUID) error
	RejectReview(ctx context.Context, tenantID, reviewID uuid.UUID, moderatorID uuid.UUID, reason string) error
	MarkAsSpam(ctx context.Context, tenantID, reviewID uuid.UUID, moderatorID uuid.UUID) error
	BulkModerateReviews(ctx context.Context, tenantID uuid.UUID, req BulkModerationRequest) error
	
	// Review replies
	AddReply(ctx context.Context, req AddReplyRequest) (*ReviewReply, error)
	GetReplies(ctx context.Context, tenantID, reviewID uuid.UUID) ([]ReviewReply, error)
	UpdateReply(ctx context.Context, tenantID, replyID uuid.UUID, req UpdateReplyRequest) (*ReviewReply, error)
	DeleteReply(ctx context.Context, tenantID, replyID uuid.UUID) error
	
	// Review reactions (helpful/unhelpful)
	ReactToReview(ctx context.Context, req ReviewReactionRequest) error
	RemoveReaction(ctx context.Context, tenantID, reviewID uuid.UUID, customerEmail string) error
	
	// Review statistics and summaries
	GetReviewSummary(ctx context.Context, tenantID, productID uuid.UUID) (*ReviewSummary, error)
	RefreshReviewSummary(ctx context.Context, tenantID, productID uuid.UUID) (*ReviewSummary, error)
	GetReviewStats(ctx context.Context, tenantID uuid.UUID, period string) (*ReviewStats, error)
	
	// Review invitations
	CreateReviewInvitation(ctx context.Context, req CreateInvitationRequest) (*ReviewInvitation, error)
	SendReviewInvitation(ctx context.Context, tenantID, invitationID uuid.UUID) error
	SendReviewReminder(ctx context.Context, tenantID, invitationID uuid.UUID) error
	ProcessInvitationClick(ctx context.Context, token string) (*ReviewInvitation, error)
	GetPendingInvitations(ctx context.Context, tenantID uuid.UUID) ([]ReviewInvitation, error)
	
	// Settings
	GetSettings(ctx context.Context, tenantID uuid.UUID) (*ReviewSettings, error)
	UpdateSettings(ctx context.Context, tenantID uuid.UUID, req UpdateSettingsRequest) (*ReviewSettings, error)
	
	// Analytics and reporting
	GetTopRatedProducts(ctx context.Context, tenantID uuid.UUID, limit int) ([]ProductRating, error)
	GetRecentReviews(ctx context.Context, tenantID uuid.UUID, limit int) ([]Review, error)
	GetReviewTrends(ctx context.Context, tenantID uuid.UUID, period string) (*ReviewTrends, error)
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new reviews service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Request/Response DTOs
type CreateReviewRequest struct {
	ProductID     *uuid.UUID `json:"product_id"`
	OrderID       *uuid.UUID `json:"order_id"`
	Type          ReviewType `json:"type"`
	CustomerID    *uuid.UUID `json:"customer_id"`
	CustomerName  string     `json:"customer_name" validate:"required"`
	CustomerEmail string     `json:"customer_email" validate:"required,email"`
	Rating        int        `json:"rating" validate:"required,min=1,max=5"`
	Title         string     `json:"title"`
	Content       string     `json:"content" validate:"required"`
	Pros          []string   `json:"pros"`
	Cons          []string   `json:"cons"`
	Images        []string   `json:"images"`
	Videos        []string   `json:"videos"`
	Source        string     `json:"source"`
	IPAddress     string     `json:"ip_address"`
	UserAgent     string     `json:"user_agent"`
}

type UpdateReviewRequest struct {
	Rating  *int     `json:"rating"`
	Title   *string  `json:"title"`
	Content *string  `json:"content"`
	Pros    []string `json:"pros"`
	Cons    []string `json:"cons"`
	Images  []string `json:"images"`
	Videos  []string `json:"videos"`
}

type ReviewFilter struct {
	ProductID   *uuid.UUID     `json:"product_id"`
	OrderID     *uuid.UUID     `json:"order_id"`
	CustomerID  *uuid.UUID     `json:"customer_id"`
	Type        []ReviewType   `json:"type"`
	Status      []ReviewStatus `json:"status"`
	Rating      []int          `json:"rating"`
	IsVerified  *bool          `json:"is_verified"`
	HasImages   *bool          `json:"has_images"`
	HasVideos   *bool          `json:"has_videos"`
	Search      string         `json:"search"`
	StartDate   *time.Time     `json:"start_date"`
	EndDate     *time.Time     `json:"end_date"`
	SortBy      string         `json:"sort_by"` // created_at, rating, helpful_count
	SortOrder   string         `json:"sort_order"` // asc, desc
	Page        int            `json:"page"`
	Limit       int            `json:"limit"`
}

type ProductReviewFilter struct {
	Status      []ReviewStatus `json:"status"`
	Rating      []int          `json:"rating"`
	IsVerified  *bool          `json:"is_verified"`
	HasImages   *bool          `json:"has_images"`
	HasVideos   *bool          `json:"has_videos"`
	Search      string         `json:"search"`
	SortBy      string         `json:"sort_by"`
	SortOrder   string         `json:"sort_order"`
	Page        int            `json:"page"`
	Limit       int            `json:"limit"`
}

type BulkModerationRequest struct {
	ReviewIDs   []uuid.UUID  `json:"review_ids" validate:"required"`
	Action      string       `json:"action" validate:"required"` // approve, reject, spam
	ModeratorID uuid.UUID    `json:"moderator_id" validate:"required"`
	Reason      string       `json:"reason"` // For reject/spam actions
}

type AddReplyRequest struct {
	ReviewID    uuid.UUID  `json:"review_id" validate:"required"`
	UserID      *uuid.UUID `json:"user_id"` // Merchant user
	CustomerID  *uuid.UUID `json:"customer_id"` // Customer reply
	AuthorName  string     `json:"author_name" validate:"required"`
	AuthorEmail string     `json:"author_email" validate:"required,email"`
	Content     string     `json:"content" validate:"required"`
	IsMerchant  bool       `json:"is_merchant"`
}

type UpdateReplyRequest struct {
	Content   *string `json:"content"`
	IsVisible *bool   `json:"is_visible"`
}

type ReviewReactionRequest struct {
	ReviewID      uuid.UUID  `json:"review_id" validate:"required"`
	CustomerID    *uuid.UUID `json:"customer_id"`
	CustomerEmail string     `json:"customer_email" validate:"required,email"`
	IsHelpful     bool       `json:"is_helpful"`
	IPAddress     string     `json:"ip_address"`
}

type CreateInvitationRequest struct {
	OrderID       uuid.UUID   `json:"order_id" validate:"required"`
	CustomerID    *uuid.UUID  `json:"customer_id"`
	CustomerEmail string      `json:"customer_email" validate:"required,email"`
	CustomerName  string      `json:"customer_name"`
	ProductIDs    []string    `json:"product_ids" validate:"required"`
	ExpiresIn     int         `json:"expires_in"` // Days from now
}

type UpdateSettingsRequest struct {
	ReviewsEnabled         *bool   `json:"reviews_enabled"`
	RequireModeration      *bool   `json:"require_moderation"`
	RequireVerifiedBuyer   *bool   `json:"require_verified_buyer"`
	AllowAnonymous         *bool   `json:"allow_anonymous"`
	AllowImages            *bool   `json:"allow_images"`
	AllowVideos            *bool   `json:"allow_videos"`
	MaxImagesPerReview     *int    `json:"max_images_per_review"`
	MaxVideosPerReview     *int    `json:"max_videos_per_review"`
	MinContentLength       *int    `json:"min_content_length"`
	MaxContentLength       *int    `json:"max_content_length"`
	AllowReactions         *bool   `json:"allow_reactions"`
	AllowReplies           *bool   `json:"allow_replies"`
	AllowCustomerReplies   *bool   `json:"allow_customer_replies"`
	ShowReviewerName       *bool   `json:"show_reviewer_name"`
	ShowReviewDate         *bool   `json:"show_review_date"`
	ShowVerifiedBadge      *bool   `json:"show_verified_badge"`
	ReviewsPerPage         *int    `json:"reviews_per_page"`
	EmailOnNewReview       *bool   `json:"email_on_new_review"`
	NotificationEmail      *string `json:"notification_email"`
	AutoRequestReviews     *bool   `json:"auto_request_reviews"`
	RequestReviewsAfter    *int    `json:"request_reviews_after"`
	RewardForReviews       *bool   `json:"reward_for_reviews"`
	RewardPoints           *int    `json:"reward_points"`
	RewardDiscount         *float64 `json:"reward_discount"`
}

// Analytics DTOs
type ReviewStats struct {
	TotalReviews        int                    `json:"total_reviews"`
	ApprovedReviews     int                    `json:"approved_reviews"`
	PendingReviews      int                    `json:"pending_reviews"`
	RejectedReviews     int                    `json:"rejected_reviews"`
	AverageRating       float64                `json:"average_rating"`
	ReviewsByRating     map[int]int            `json:"reviews_by_rating"`
	ReviewsByStatus     map[ReviewStatus]int   `json:"reviews_by_status"`
	VerifiedReviews     int                    `json:"verified_reviews"`
	ReviewsWithMedia    int                    `json:"reviews_with_media"`
	ResponseRate        float64                `json:"response_rate"`
	AverageResponseTime string                 `json:"average_response_time"`
}

type ProductRating struct {
	ProductID     uuid.UUID `json:"product_id"`
	ProductName   string    `json:"product_name"`
	AverageRating float64   `json:"average_rating"`
	TotalReviews  int       `json:"total_reviews"`
}

type ReviewTrends struct {
	Period          string             `json:"period"`
	TotalReviews    int                `json:"total_reviews"`
	AverageRating   float64            `json:"average_rating"`
	DailyReviews    []DailyReviewCount `json:"daily_reviews"`
	TopProducts     []ProductRating    `json:"top_products"`
	RecentReviews   []Review           `json:"recent_reviews"`
}

type DailyReviewCount struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// Implementation methods (TODO: implement business logic)
func (s *service) CreateReview(ctx context.Context, req CreateReviewRequest) (*Review, error) {
	// Validate rating
	if req.Rating < 1 || req.Rating > 5 {
		return nil, fmt.Errorf("rating must be between 1 and 5")
	}

	// Create review entity
	review := &Review{
		ID:            uuid.New(),
		TenantID:      *req.ProductID, // This should be passed from context
		ProductID:     req.ProductID,
		OrderID:       req.OrderID,
		Type:          req.Type,
		CustomerID:    req.CustomerID,
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
		Rating:        req.Rating,
		Title:         req.Title,
		Content:       req.Content,
		Pros:          req.Pros,
		Cons:          req.Cons,
		Images:        req.Images,
		Videos:        req.Videos,
		Source:        req.Source,
		IPAddress:     req.IPAddress,
		UserAgent:     req.UserAgent,
		Status:        StatusPending, // Default to pending for moderation
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Check if customer is verified buyer
	if req.OrderID != nil {
		review.IsVerified = true
	}

	return s.repo.CreateReview(ctx, review)
}

func (s *service) GetReview(ctx context.Context, tenantID, reviewID uuid.UUID) (*Review, error) {
	return s.repo.GetReviewByID(ctx, tenantID, reviewID)
}

func (s *service) GetReviews(ctx context.Context, tenantID uuid.UUID, filter ReviewFilter) ([]Review, error) {
	return s.repo.GetReviews(ctx, tenantID, filter)
}

func (s *service) GetProductReviews(ctx context.Context, tenantID, productID uuid.UUID, filter ProductReviewFilter) ([]Review, error) {
	// Convert ProductReviewFilter to ReviewFilter
	reviewFilter := ReviewFilter{
		ProductID:  &productID,
		Status:     filter.Status,
		Rating:     filter.Rating,
		IsVerified: filter.IsVerified,
		HasImages:  filter.HasImages,
		HasVideos:  filter.HasVideos,
		Search:     filter.Search,
		SortBy:     filter.SortBy,
		SortOrder:  filter.SortOrder,
		Page:       filter.Page,
		Limit:      filter.Limit,
	}
	
	return s.repo.GetReviews(ctx, tenantID, reviewFilter)
}

func (s *service) UpdateReview(ctx context.Context, tenantID, reviewID uuid.UUID, req UpdateReviewRequest) (*Review, error) {
	updates := make(map[string]interface{})

	if req.Rating != nil {
		if *req.Rating < 1 || *req.Rating > 5 {
			return nil, fmt.Errorf("rating must be between 1 and 5")
		}
		updates["rating"] = *req.Rating
	}

	if req.Title != nil {
		updates["title"] = *req.Title
	}

	if req.Content != nil {
		updates["content"] = *req.Content
	}

	if req.Pros != nil {
		updates["pros"] = req.Pros
	}

	if req.Cons != nil {
		updates["cons"] = req.Cons
	}

	if req.Images != nil {
		updates["images"] = req.Images
	}

	if req.Videos != nil {
		updates["videos"] = req.Videos
	}

	updates["updated_at"] = time.Now()

	err := s.repo.UpdateReview(ctx, tenantID, reviewID, updates)
	if err != nil {
		return nil, err
	}

	return s.repo.GetReviewByID(ctx, tenantID, reviewID)
}

func (s *service) DeleteReview(ctx context.Context, tenantID, reviewID uuid.UUID) error {
	// Soft delete by updating status
	updates := map[string]interface{}{
		"status":     StatusDeleted,
		"deleted_at": time.Now(),
		"updated_at": time.Now(),
	}

	err := s.repo.UpdateReview(ctx, tenantID, reviewID, updates)
	if err != nil {
		return err
	}

	// TODO: Update product review summary after deletion
	return nil
}

func (s *service) ApproveReview(ctx context.Context, tenantID, reviewID uuid.UUID, moderatorID uuid.UUID) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":        StatusApproved,
		"moderated_by":  moderatorID,
		"moderated_at":  &now,
	}
	
	err := s.repo.UpdateReview(ctx, tenantID, reviewID, updates)
	if err != nil {
		return err
	}
	
	// TODO: Update review summary and send notifications
	return nil
}

func (s *service) RejectReview(ctx context.Context, tenantID, reviewID uuid.UUID, moderatorID uuid.UUID, reason string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":          StatusRejected,
		"moderated_by":    moderatorID,
		"moderated_at":    &now,
		"moderation_note": reason,
	}
	
	return s.repo.UpdateReview(ctx, tenantID, reviewID, updates)
}

func (s *service) MarkAsSpam(ctx context.Context, tenantID, reviewID uuid.UUID, moderatorID uuid.UUID) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":       StatusSpam,
		"moderated_by": moderatorID,
		"moderated_at": &now,
	}
	
	return s.repo.UpdateReview(ctx, tenantID, reviewID, updates)
}

func (s *service) BulkModerateReviews(ctx context.Context, tenantID uuid.UUID, req BulkModerationRequest) error {
	now := time.Now()
	var updates map[string]interface{}

	switch req.Action {
	case "approve":
		updates = map[string]interface{}{
			"status":       StatusApproved,
			"moderated_by": req.ModeratorID,
			"moderated_at": &now,
		}
	case "reject":
		updates = map[string]interface{}{
			"status":          StatusRejected,
			"moderated_by":    req.ModeratorID,
			"moderated_at":    &now,
			"moderation_note": req.Reason,
		}
	case "spam":
		updates = map[string]interface{}{
			"status":       StatusSpam,
			"moderated_by": req.ModeratorID,
			"moderated_at": &now,
		}
	default:
		return fmt.Errorf("invalid action: %s", req.Action)
	}

	// Update all reviews in batch
	for _, reviewID := range req.ReviewIDs {
		if err := s.repo.UpdateReview(ctx, tenantID, reviewID, updates); err != nil {
			return fmt.Errorf("failed to update review %s: %w", reviewID, err)
		}
	}

	return nil
}

func (s *service) AddReply(ctx context.Context, req AddReplyRequest) (*ReviewReply, error) {
	reply := &ReviewReply{
		ID:          uuid.New(),
		ReviewID:    req.ReviewID,
		UserID:      req.UserID,
		CustomerID:  req.CustomerID,
		AuthorName:  req.AuthorName,
		AuthorEmail: req.AuthorEmail,
		Content:     req.Content,
		IsMerchant:  req.IsMerchant,
		IsVisible:   true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.repo.CreateReply(ctx, reply)
}

func (s *service) GetReplies(ctx context.Context, tenantID, reviewID uuid.UUID) ([]ReviewReply, error) {
	return s.repo.GetRepliesByReviewID(ctx, tenantID, reviewID)
}

func (s *service) UpdateReply(ctx context.Context, tenantID, replyID uuid.UUID, req UpdateReplyRequest) (*ReviewReply, error) {
	// TODO: Implement reply update logic
	return nil, fmt.Errorf("TODO: implement UpdateReply")
}

func (s *service) DeleteReply(ctx context.Context, tenantID, replyID uuid.UUID) error {
	return s.repo.DeleteReply(ctx, tenantID, replyID)
}

func (s *service) ReactToReview(ctx context.Context, req ReviewReactionRequest) error {
	// Check if user already reacted
	existingReaction, err := s.repo.GetReaction(ctx, req.ReviewID, req.CustomerEmail)
	if err == nil && existingReaction != nil {
		// Update existing reaction
		return s.repo.UpdateReaction(ctx, existingReaction.ID, req.IsHelpful)
	}

	// Create new reaction
	reaction := &ReviewReaction{
		ID:            uuid.New(),
		ReviewID:      req.ReviewID,
		CustomerID:    req.CustomerID,
		CustomerEmail: req.CustomerEmail,
		IsHelpful:     req.IsHelpful,
		IPAddress:     req.IPAddress,
		CreatedAt:     time.Now(),
	}

	return s.repo.CreateReaction(ctx, reaction)
}

func (s *service) RemoveReaction(ctx context.Context, tenantID, reviewID uuid.UUID, customerEmail string) error {
	return s.repo.DeleteReaction(ctx, tenantID, reviewID, customerEmail)
}

func (s *service) GetReviewSummary(ctx context.Context, tenantID, productID uuid.UUID) (*ReviewSummary, error) {
	return s.repo.GetReviewSummary(ctx, tenantID, productID)
}

func (s *service) RefreshReviewSummary(ctx context.Context, tenantID, productID uuid.UUID) (*ReviewSummary, error) {
	// Get all approved reviews for the product
	filter := ReviewFilter{
		ProductID: &productID,
		Status:    []ReviewStatus{StatusApproved},
		Limit:     1000, // Get all reviews
	}
	
	reviews, err := s.repo.GetReviews(ctx, tenantID, filter)
	if err != nil {
		return nil, err
	}

	if len(reviews) == 0 {
		// No reviews, create empty summary
		summary := &ReviewSummary{
			ProductID:     productID,
			TotalReviews:  0,
			AverageRating: 0,
			UpdatedAt:     time.Now(),
		}
		return s.repo.UpsertReviewSummary(ctx, tenantID, summary)
	}

	// Calculate statistics
	var totalRating float64
	ratingCounts := make(map[int]int)

	for _, review := range reviews {
		totalRating += float64(review.Rating)
		ratingCounts[review.Rating]++
	}

	avgRating := totalRating / float64(len(reviews))

	summary := &ReviewSummary{
		ProductID:      productID,
		TotalReviews:   len(reviews),
		AverageRating:  avgRating,
		Rating1Count:   ratingCounts[1],
		Rating2Count:   ratingCounts[2],
		Rating3Count:   ratingCounts[3],
		Rating4Count:   ratingCounts[4],
		Rating5Count:   ratingCounts[5],
		UpdatedAt:      time.Now(),
	}

	return s.repo.UpsertReviewSummary(ctx, tenantID, summary)
}

func (s *service) GetReviewStats(ctx context.Context, tenantID uuid.UUID, period string) (*ReviewStats, error) {
	stats, err := s.repo.GetReviewStats(ctx, tenantID, period)
	if err != nil {
		return nil, err
	}

	// Calculate additional metrics
	if stats.TotalReviews > 0 {
		stats.ResponseRate = float64(stats.ApprovedReviews) / float64(stats.TotalReviews) * 100
	}

	return stats, nil
}

func (s *service) CreateReviewInvitation(ctx context.Context, req CreateInvitationRequest) (*ReviewInvitation, error) {
	// Validate email format
	if req.CustomerEmail == "" {
		return nil, fmt.Errorf("customer email is required")
	}

	// Check if invitation already exists
	existing, err := s.repo.GetInvitationByOrderAndEmail(ctx, req.TenantID, req.OrderID, req.CustomerEmail)
	if err == nil && existing != nil {
		return existing, nil // Return existing invitation
	}

	invitation := &ReviewInvitation{
		ID:            uuid.New(),
		TenantID:      req.TenantID,
		OrderID:       req.OrderID,
		ProductID:     req.ProductID,
		CustomerID:    req.CustomerID,
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
		Token:         generateInvitationToken(),
		Status:        InvitationStatusPending,
		ExpiresAt:     time.Now().AddDate(0, 0, 30), // 30 days expiry
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return s.repo.CreateInvitation(ctx, invitation)
}

func generateInvitationToken() string {
	return uuid.New().String()
}

func (s *service) SendReviewInvitation(ctx context.Context, tenantID, invitationID uuid.UUID) error {
	invitation, err := s.repo.GetInvitationByID(ctx, tenantID, invitationID)
	if err != nil {
		return err
	}

	if invitation.Status != InvitationStatusPending {
		return fmt.Errorf("invitation is not in pending status")
	}

	// TODO: Integrate with email service to send invitation
	// For now, just update the status
	updates := map[string]interface{}{
		"status":   InvitationStatusSent,
		"sent_at":  time.Now(),
		"updated_at": time.Now(),
	}

	return s.repo.UpdateInvitation(ctx, tenantID, invitationID, updates)
}

func (s *service) SendReviewReminder(ctx context.Context, tenantID, invitationID uuid.UUID) error {
	// TODO: Implement reminder sending logic
	return fmt.Errorf("TODO: implement SendReviewReminder")
}

func (s *service) ProcessInvitationClick(ctx context.Context, token string) (*ReviewInvitation, error) {
	invitation, err := s.repo.GetInvitationByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if invitation.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("invitation has expired")
	}

	// Update click tracking
	updates := map[string]interface{}{
		"clicked_at": time.Now(),
		"updated_at": time.Now(),
	}

	if invitation.Status == InvitationStatusSent {
		updates["status"] = InvitationStatusClicked
	}

	err = s.repo.UpdateInvitation(ctx, invitation.TenantID, invitation.ID, updates)
	if err != nil {
		return nil, err
	}

	return invitation, nil
}

func (s *service) GetPendingInvitations(ctx context.Context, tenantID uuid.UUID) ([]ReviewInvitation, error) {
	return s.repo.GetInvitationsByStatus(ctx, tenantID, "pending")
}

func (s *service) GetSettings(ctx context.Context, tenantID uuid.UUID) (*ReviewSettings, error) {
	return s.repo.GetSettings(ctx, tenantID)
}

func (s *service) UpdateSettings(ctx context.Context, tenantID uuid.UUID, req UpdateSettingsRequest) (*ReviewSettings, error) {
	updates := make(map[string]interface{})

	if req.AutoApprove != nil {
		updates["auto_approve"] = *req.AutoApprove
	}

	if req.RequireApproval != nil {
		updates["require_approval"] = *req.RequireApproval
	}

	if req.AllowAnonymous != nil {
		updates["allow_anonymous"] = *req.AllowAnonymous
	}

	if req.RequireVerifiedPurchase != nil {
		updates["require_verified_purchase"] = *req.RequireVerifiedPurchase
	}

	if req.EnablePhotos != nil {
		updates["enable_photos"] = *req.EnablePhotos
	}

	if req.EnableVideos != nil {
		updates["enable_videos"] = *req.EnableVideos
	}

	if req.MaxPhotos != nil {
		updates["max_photos"] = *req.MaxPhotos
	}

	if req.MaxVideos != nil {
		updates["max_videos"] = *req.MaxVideos
	}

	if req.AutoInviteAfterDays != nil {
		updates["auto_invite_after_days"] = *req.AutoInviteAfterDays
	}

	if req.ReminderAfterDays != nil {
		updates["reminder_after_days"] = *req.ReminderAfterDays
	}

	updates["updated_at"] = time.Now()

	err := s.repo.UpdateSettings(ctx, tenantID, updates)
	if err != nil {
		return nil, err
	}

	return s.repo.GetSettings(ctx, tenantID)
}

func (s *service) GetTopRatedProducts(ctx context.Context, tenantID uuid.UUID, limit int) ([]ProductRating, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	return s.repo.GetTopRatedProducts(ctx, tenantID, limit)
}

func (s *service) GetRecentReviews(ctx context.Context, tenantID uuid.UUID, limit int) ([]Review, error) {
	filter := ReviewFilter{
		Status:    []ReviewStatus{StatusApproved},
		SortBy:    "created_at",
		SortOrder: "desc",
		Limit:     limit,
		Page:      1,
	}
	
	return s.repo.GetReviews(ctx, tenantID, filter)
}

func (s *service) GetReviewTrends(ctx context.Context, tenantID uuid.UUID, period string) (*ReviewTrends, error) {
	if period == "" {
		period = "30d" // Default to 30 days
	}

	trends, err := s.repo.GetReviewTrends(ctx, tenantID, period)
	if err != nil {
		return nil, err
	}

	// Calculate growth rates if we have previous period data
	if len(trends.DailyReviews) > 1 {
		currentPeriod := trends.DailyReviews[len(trends.DailyReviews)-7:] // Last 7 days
		previousPeriod := trends.DailyReviews[len(trends.DailyReviews)-14:len(trends.DailyReviews)-7] // Previous 7 days

		if len(currentPeriod) > 0 && len(previousPeriod) > 0 {
			currentTotal := 0
			previousTotal := 0

			for _, stat := range currentPeriod {
				currentTotal += stat.Count
			}

			for _, stat := range previousPeriod {
				previousTotal += stat.Count
			}

			if previousTotal > 0 {
				// Add growth rate calculation to trends if needed
			}
		}
	}

	return trends, nil
}