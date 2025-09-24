package reviews

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	sharedErrors "ecommerce-saas/internal/shared/errors"
)

// Repository defines the reviews repository interface
type Repository interface {
	// Review operations
	CreateReview(ctx context.Context, review *Review) error
	GetReviewByID(ctx context.Context, tenantID, reviewID uuid.UUID) (*Review, error)
	GetReviews(ctx context.Context, tenantID uuid.UUID, filter ReviewFilter) ([]Review, error)
	UpdateReview(ctx context.Context, tenantID, reviewID uuid.UUID, updates map[string]interface{}) error
	DeleteReview(ctx context.Context, tenantID, reviewID uuid.UUID) error
	GetReviewCount(ctx context.Context, tenantID uuid.UUID, filter ReviewFilter) (int64, error)

	// Review reply operations
	CreateReply(ctx context.Context, reply *ReviewReply) error
	GetRepliesByReviewID(ctx context.Context, tenantID, reviewID uuid.UUID) ([]ReviewReply, error)
	GetReplyByID(ctx context.Context, tenantID, replyID uuid.UUID) (*ReviewReply, error)
	UpdateReply(ctx context.Context, tenantID, replyID uuid.UUID, updates map[string]interface{}) error
	DeleteReply(ctx context.Context, tenantID, replyID uuid.UUID) error

	// Review reaction operations
	CreateReaction(ctx context.Context, reaction *ReviewReaction) error
	GetReactionByReviewAndEmail(ctx context.Context, tenantID, reviewID uuid.UUID, email string) (*ReviewReaction, error)
	UpdateReaction(ctx context.Context, tenantID, reviewID uuid.UUID, email string, isHelpful bool) error
	DeleteReaction(ctx context.Context, tenantID, reviewID uuid.UUID, email string) error
	UpdateReviewReactionCounts(ctx context.Context, reviewID uuid.UUID) error

	// Review summary operations
	CreateReviewSummary(ctx context.Context, summary *ReviewSummary) error
	GetReviewSummary(ctx context.Context, tenantID, productID uuid.UUID) (*ReviewSummary, error)
	UpdateReviewSummary(ctx context.Context, tenantID, productID uuid.UUID, updates map[string]interface{}) error
	RecalculateReviewSummary(ctx context.Context, tenantID, productID uuid.UUID) error

	// Review invitation operations
	CreateInvitation(ctx context.Context, invitation *ReviewInvitation) error
	GetInvitationByID(ctx context.Context, tenantID, invitationID uuid.UUID) (*ReviewInvitation, error)
	GetInvitationByToken(ctx context.Context, token string) (*ReviewInvitation, error)
	GetInvitationsByStatus(ctx context.Context, tenantID uuid.UUID, status string) ([]ReviewInvitation, error)
	UpdateInvitation(ctx context.Context, tenantID, invitationID uuid.UUID, updates map[string]interface{}) error
	DeleteInvitation(ctx context.Context, tenantID, invitationID uuid.UUID) error
	GetExpiredInvitations(ctx context.Context, tenantID uuid.UUID) ([]ReviewInvitation, error)

	// Settings operations
	CreateSettings(ctx context.Context, settings *ReviewSettings) error
	GetSettings(ctx context.Context, tenantID uuid.UUID) (*ReviewSettings, error)
	UpdateSettings(ctx context.Context, tenantID uuid.UUID, updates map[string]interface{}) error

	// Analytics operations
	GetReviewStatsByPeriod(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) (*ReviewStats, error)
	GetTopRatedProducts(ctx context.Context, tenantID uuid.UUID, limit int) ([]ProductRating, error)
	GetReviewCountByStatus(ctx context.Context, tenantID uuid.UUID) (map[ReviewStatus]int, error)
	GetReviewCountByRating(ctx context.Context, tenantID uuid.UUID) (map[int]int, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new reviews repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Review operations
func (r *repository) CreateReview(ctx context.Context, review *Review) error {
	if err := r.db.WithContext(ctx).Create(review).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to create review", err)
	}
	return nil
}

func (r *repository) GetReviewByID(ctx context.Context, tenantID, reviewID uuid.UUID) (*Review, error) {
	var review Review
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", reviewID, tenantID).
		Preload("Replies").
		Preload("Reactions").
		First(&review).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("Review not found")
		}
		return nil, sharedErrors.NewInternalError("Failed to get review", err)
	}
	return &review, nil
}

func (r *repository) GetReviews(ctx context.Context, tenantID uuid.UUID, filter ReviewFilter) ([]Review, error) {
	var reviews []Review
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)

	// Apply filters
	if filter.ProductID != nil {
		query = query.Where("product_id = ?", *filter.ProductID)
	}

	if filter.OrderID != nil {
		query = query.Where("order_id = ?", *filter.OrderID)
	}

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	if len(filter.Type) > 0 {
		query = query.Where("type IN ?", filter.Type)
	}

	if len(filter.Status) > 0 {
		query = query.Where("status IN ?", filter.Status)
	}

	if len(filter.Rating) > 0 {
		query = query.Where("rating IN ?", filter.Rating)
	}

	if filter.IsVerified != nil {
		query = query.Where("is_verified = ?", *filter.IsVerified)
	}

	if filter.HasImages != nil {
		if *filter.HasImages {
			query = query.Where("JSON_LENGTH(images) > 0")
		} else {
			query = query.Where("JSON_LENGTH(images) = 0 OR images IS NULL")
		}
	}

	if filter.HasVideos != nil {
		if *filter.HasVideos {
			query = query.Where("JSON_LENGTH(videos) > 0")
		} else {
			query = query.Where("JSON_LENGTH(videos) = 0 OR videos IS NULL")
		}
	}

	if filter.Search != "" {
		query = query.Where("title ILIKE ? OR content ILIKE ? OR customer_name ILIKE ?",
			"%"+filter.Search+"%", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", filter.StartDate)
	}

	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", filter.EndDate)
	}

	// Sorting
	sortBy := "created_at"
	if filter.SortBy != "" {
		switch filter.SortBy {
		case "rating", "helpful_count", "created_at", "updated_at":
			sortBy = filter.SortBy
		}
	}

	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	query = query.Order(sortBy + " " + sortOrder)

	// Pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
		if filter.Page > 0 {
			query = query.Offset((filter.Page - 1) * filter.Limit)
		}
	}

	// Preload relations
	query = query.Preload("Replies").Preload("Reactions")

	err := query.Find(&reviews).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to get reviews", err)
	}
	return reviews, nil
}

func (r *repository) UpdateReview(ctx context.Context, tenantID, reviewID uuid.UUID, updates map[string]interface{}) error {
	result := r.db.WithContext(ctx).
		Model(&Review{}).
		Where("id = ? AND tenant_id = ?", reviewID, tenantID).
		Updates(updates)
	if result.Error != nil {
		return sharedErrors.NewInternalError("Failed to update review", result.Error)
	}
	if result.RowsAffected == 0 {
		return sharedErrors.NewNotFoundError("Review not found")
	}
	return nil
}

func (r *repository) DeleteReview(ctx context.Context, tenantID, reviewID uuid.UUID) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete related reactions
		if err := tx.Where("review_id = ?", reviewID).Delete(&ReviewReaction{}).Error; err != nil {
			return err
		}

		// Delete related replies
		if err := tx.Where("review_id = ?", reviewID).Delete(&ReviewReply{}).Error; err != nil {
			return err
		}

		// Delete the review
		return tx.Where("id = ? AND tenant_id = ?", reviewID, tenantID).Delete(&Review{}).Error
	})
	if err != nil {
		return sharedErrors.NewInternalError("Failed to delete review", err)
	}
	return nil
}

func (r *repository) GetReviewCount(ctx context.Context, tenantID uuid.UUID, filter ReviewFilter) (int64, error) {
	query := r.db.WithContext(ctx).Model(&Review{}).Where("tenant_id = ?", tenantID)

	// Apply same filters as GetReviews
	if filter.ProductID != nil {
		query = query.Where("product_id = ?", *filter.ProductID)
	}

	if len(filter.Status) > 0 {
		query = query.Where("status IN ?", filter.Status)
	}

	// Add other filter conditions as needed

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return 0, sharedErrors.NewInternalError("Failed to get review count", err)
	}
	return count, nil
}

// Review reply operations
func (r *repository) CreateReply(ctx context.Context, reply *ReviewReply) error {
	if err := r.db.WithContext(ctx).Create(reply).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to create review reply", err)
	}
	return nil
}

func (r *repository) GetRepliesByReviewID(ctx context.Context, tenantID, reviewID uuid.UUID) ([]ReviewReply, error) {
	var replies []ReviewReply

	// Verify review belongs to tenant first
	var count int64
	err := r.db.WithContext(ctx).Model(&Review{}).
		Where("id = ? AND tenant_id = ?", reviewID, tenantID).
		Count(&count).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to verify review ownership", err)
	}

	if count == 0 {
		return nil, sharedErrors.NewNotFoundError("Review not found")
	}

	err = r.db.WithContext(ctx).
		Where("review_id = ? AND is_visible = ?", reviewID, true).
		Order("created_at ASC").
		Find(&replies).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to get review replies", err)
	}

	return replies, nil
}

func (r *repository) GetReplyByID(ctx context.Context, tenantID, replyID uuid.UUID) (*ReviewReply, error) {
	var reply ReviewReply

	// Join with reviews table to ensure tenant isolation
	err := r.db.WithContext(ctx).
		Table("review_replies").
		Select("review_replies.*").
		Joins("JOIN reviews ON reviews.id = review_replies.review_id").
		Where("review_replies.id = ? AND reviews.tenant_id = ?", replyID, tenantID).
		First(&reply).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("Review reply not found")
		}
		return nil, sharedErrors.NewInternalError("Failed to get review reply", err)
	}

	return &reply, nil
}

func (r *repository) UpdateReply(ctx context.Context, tenantID, replyID uuid.UUID, updates map[string]interface{}) error {
	// Ensure tenant isolation through join
	result := r.db.WithContext(ctx).
		Table("review_replies").
		Joins("JOIN reviews ON reviews.id = review_replies.review_id").
		Where("review_replies.id = ? AND reviews.tenant_id = ?", replyID, tenantID).
		Updates(updates)
	if result.Error != nil {
		return sharedErrors.NewInternalError("Failed to update review reply", result.Error)
	}
	if result.RowsAffected == 0 {
		return sharedErrors.NewNotFoundError("Review reply not found")
	}
	return nil
}

func (r *repository) DeleteReply(ctx context.Context, tenantID, replyID uuid.UUID) error {
	// Ensure tenant isolation through join
	result := r.db.WithContext(ctx).
		Table("review_replies").
		Joins("JOIN reviews ON reviews.id = review_replies.review_id").
		Where("review_replies.id = ? AND reviews.tenant_id = ?", replyID, tenantID).
		Delete(&ReviewReply{})
	if result.Error != nil {
		return sharedErrors.NewInternalError("Failed to delete review reply", result.Error)
	}
	if result.RowsAffected == 0 {
		return sharedErrors.NewNotFoundError("Review reply not found")
	}
	return nil
}

// Review reaction operations
func (r *repository) CreateReaction(ctx context.Context, reaction *ReviewReaction) error {
	if err := r.db.WithContext(ctx).Create(reaction).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to create review reaction", err)
	}
	return nil
}

func (r *repository) GetReactionByReviewAndEmail(ctx context.Context, tenantID, reviewID uuid.UUID, email string) (*ReviewReaction, error) {
	var reaction ReviewReaction

	// Verify review belongs to tenant and get reaction
	err := r.db.WithContext(ctx).
		Table("review_reactions").
		Select("review_reactions.*").
		Joins("JOIN reviews ON reviews.id = review_reactions.review_id").
		Where("review_reactions.review_id = ? AND review_reactions.customer_email = ? AND reviews.tenant_id = ?",
			reviewID, email, tenantID).
		First(&reaction).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("Review reaction not found")
		}
		return nil, sharedErrors.NewInternalError("Failed to get review reaction", err)
	}
	return &reaction, nil
}

func (r *repository) UpdateReaction(ctx context.Context, tenantID, reviewID uuid.UUID, email string, isHelpful bool) error {
	// Verify tenant isolation and update
	result := r.db.WithContext(ctx).
		Table("review_reactions").
		Joins("JOIN reviews ON reviews.id = review_reactions.review_id").
		Where("review_reactions.review_id = ? AND review_reactions.customer_email = ? AND reviews.tenant_id = ?",
			reviewID, email, tenantID).
		Update("is_helpful", isHelpful)
	if result.Error != nil {
		return sharedErrors.NewInternalError("Failed to update review reaction", result.Error)
	}
	if result.RowsAffected == 0 {
		return sharedErrors.NewNotFoundError("Review reaction not found")
	}
	return nil
}

func (r *repository) DeleteReaction(ctx context.Context, tenantID, reviewID uuid.UUID, email string) error {
	result := r.db.WithContext(ctx).
		Table("review_reactions").
		Joins("JOIN reviews ON reviews.id = review_reactions.review_id").
		Where("review_reactions.review_id = ? AND review_reactions.customer_email = ? AND reviews.tenant_id = ?",
			reviewID, email, tenantID).
		Delete(&ReviewReaction{})
	if result.Error != nil {
		return sharedErrors.NewInternalError("Failed to delete review reaction", result.Error)
	}
	if result.RowsAffected == 0 {
		return sharedErrors.NewNotFoundError("Review reaction not found")
	}
	return nil
}

func (r *repository) UpdateReviewReactionCounts(ctx context.Context, reviewID uuid.UUID) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Count helpful reactions
		var helpfulCount int64
		if err := tx.Model(&ReviewReaction{}).
			Where("review_id = ? AND is_helpful = ?", reviewID, true).
			Count(&helpfulCount).Error; err != nil {
			return err
		}

		// Count unhelpful reactions
		var unhelpfulCount int64
		if err := tx.Model(&ReviewReaction{}).
			Where("review_id = ? AND is_helpful = ?", reviewID, false).
			Count(&unhelpfulCount).Error; err != nil {
			return err
		}

		// Update review counts
		return tx.Model(&Review{}).
			Where("id = ?", reviewID).
			Updates(map[string]interface{}{
				"helpful_count":   helpfulCount,
				"unhelpful_count": unhelpfulCount,
			}).Error
	})
	if err != nil {
		return sharedErrors.NewInternalError("Failed to update review reaction counts", err)
	}
	return nil
}

// Review summary operations
func (r *repository) CreateReviewSummary(ctx context.Context, summary *ReviewSummary) error {
	if err := r.db.WithContext(ctx).Create(summary).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to create review summary", err)
	}
	return nil
}

func (r *repository) GetReviewSummary(ctx context.Context, tenantID, productID uuid.UUID) (*ReviewSummary, error) {
	var summary ReviewSummary
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND product_id = ? AND type = ?", tenantID, productID, TypeProduct).
		First(&summary).Error

	if err == gorm.ErrRecordNotFound {
		// Create default summary if not exists
		summary = ReviewSummary{
			ID:        uuid.New(),
			TenantID:  tenantID,
			ProductID: &productID,
			Type:      TypeProduct,
		}
		if createErr := r.CreateReviewSummary(ctx, &summary); createErr != nil {
			return nil, createErr
		}
		return &summary, nil
	}

	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to get review summary", err)
	}
	return &summary, nil
}

func (r *repository) UpdateReviewSummary(ctx context.Context, tenantID, productID uuid.UUID, updates map[string]interface{}) error {
	result := r.db.WithContext(ctx).
		Model(&ReviewSummary{}).
		Where("tenant_id = ? AND product_id = ? AND type = ?", tenantID, productID, TypeProduct).
		Updates(updates)
	if result.Error != nil {
		return sharedErrors.NewInternalError("Failed to update review summary", result.Error)
	}
	if result.RowsAffected == 0 {
		return sharedErrors.NewNotFoundError("Review summary not found")
	}
	return nil
}

func (r *repository) RecalculateReviewSummary(ctx context.Context, tenantID, productID uuid.UUID) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		type ratingCount struct {
			Rating int `json:"rating"`
			Count  int `json:"count"`
		}

		var ratings []ratingCount
		err := tx.Model(&Review{}).
			Select("rating, COUNT(*) as count").
			Where("tenant_id = ? AND product_id = ? AND status = ?", tenantID, productID, StatusApproved).
			Group("rating").
			Find(&ratings).Error

		if err != nil {
			return err
		}

		// Calculate summary statistics
		var totalReviews, verifiedCount, withPhotosCount int64
		var ratingCounts [6]int // Index 0 unused, 1-5 for ratings
		var totalPoints int

		for _, r := range ratings {
			totalReviews += int64(r.Count)
			totalPoints += r.Rating * r.Count
			if r.Rating >= 1 && r.Rating <= 5 {
				ratingCounts[r.Rating] = r.Count
			}
		}

		// Count verified reviews and reviews with photos
		tx.Model(&Review{}).
			Where("tenant_id = ? AND product_id = ? AND status = ? AND is_verified = ?",
				tenantID, productID, StatusApproved, true).
			Count(&verifiedCount)

		tx.Model(&Review{}).
			Where("tenant_id = ? AND product_id = ? AND status = ? AND JSON_LENGTH(images) > 0",
				tenantID, productID, StatusApproved).
			Count(&withPhotosCount)

		// Calculate average rating
		var avgRating float64
		if totalReviews > 0 {
			avgRating = float64(totalPoints) / float64(totalReviews)
		}

		// Update summary
		updates := map[string]interface{}{
			"total_reviews":    totalReviews,
			"approved_reviews": totalReviews,
			"average_rating":   avgRating,
			"rating_1_count":   ratingCounts[1],
			"rating_2_count":   ratingCounts[2],
			"rating_3_count":   ratingCounts[3],
			"rating_4_count":   ratingCounts[4],
			"rating_5_count":   ratingCounts[5],
			"verified_reviews": verifiedCount,
			"with_photos":      withPhotosCount,
			"updated_at":       time.Now(),
		}

		return r.UpdateReviewSummary(ctx, tenantID, productID, updates)
	})
	if err != nil {
		return sharedErrors.NewInternalError("Failed to recalculate review summary", err)
	}
	return nil
}

// Review invitation operations
func (r *repository) CreateInvitation(ctx context.Context, invitation *ReviewInvitation) error {
	if err := r.db.WithContext(ctx).Create(invitation).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to create review invitation", err)
	}
	return nil
}

func (r *repository) GetInvitationByID(ctx context.Context, tenantID, invitationID uuid.UUID) (*ReviewInvitation, error) {
	var invitation ReviewInvitation
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", invitationID, tenantID).
		First(&invitation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("Review invitation not found")
		}
		return nil, sharedErrors.NewInternalError("Failed to get review invitation", err)
	}
	return &invitation, nil
}

func (r *repository) GetInvitationByToken(ctx context.Context, token string) (*ReviewInvitation, error) {
	var invitation ReviewInvitation
	err := r.db.WithContext(ctx).
		Where("invitation_token = ?", token).
		First(&invitation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("Review invitation not found")
		}
		return nil, sharedErrors.NewInternalError("Failed to get review invitation", err)
	}
	return &invitation, nil
}

func (r *repository) GetInvitationsByStatus(ctx context.Context, tenantID uuid.UUID, status string) ([]ReviewInvitation, error) {
	var invitations []ReviewInvitation
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND status = ?", tenantID, status).
		Order("created_at DESC").
		Find(&invitations).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to get review invitations by status", err)
	}
	return invitations, nil
}

func (r *repository) UpdateInvitation(ctx context.Context, tenantID, invitationID uuid.UUID, updates map[string]interface{}) error {
	result := r.db.WithContext(ctx).
		Model(&ReviewInvitation{}).
		Where("id = ? AND tenant_id = ?", invitationID, tenantID).
		Updates(updates)
	if result.Error != nil {
		return sharedErrors.NewInternalError("Failed to update review invitation", result.Error)
	}
	if result.RowsAffected == 0 {
		return sharedErrors.NewNotFoundError("Review invitation not found")
	}
	return nil
}

func (r *repository) DeleteInvitation(ctx context.Context, tenantID, invitationID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", invitationID, tenantID).
		Delete(&ReviewInvitation{})
	if result.Error != nil {
		return sharedErrors.NewInternalError("Failed to delete review invitation", result.Error)
	}
	if result.RowsAffected == 0 {
		return sharedErrors.NewNotFoundError("Review invitation not found")
	}
	return nil
}

func (r *repository) GetExpiredInvitations(ctx context.Context, tenantID uuid.UUID) ([]ReviewInvitation, error) {
	var invitations []ReviewInvitation
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND status = ? AND expires_at < ?", tenantID, "sent", time.Now()).
		Find(&invitations).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to get expired review invitations", err)
	}
	return invitations, nil
}

// Settings operations
func (r *repository) CreateSettings(ctx context.Context, settings *ReviewSettings) error {
	if err := r.db.WithContext(ctx).Create(settings).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to create review settings", err)
	}
	return nil
}

func (r *repository) GetSettings(ctx context.Context, tenantID uuid.UUID) (*ReviewSettings, error) {
	var settings ReviewSettings
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		First(&settings).Error

	if err == gorm.ErrRecordNotFound {
		// Create default settings
		settings = ReviewSettings{
			ID:       uuid.New(),
			TenantID: tenantID,
			// Default values are set via gorm tags
		}
		if createErr := r.CreateSettings(ctx, &settings); createErr != nil {
			return nil, createErr
		}
		return &settings, nil
	}

	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to get review settings", err)
	}
	return &settings, nil
}

func (r *repository) UpdateSettings(ctx context.Context, tenantID uuid.UUID, updates map[string]interface{}) error {
	result := r.db.WithContext(ctx).
		Model(&ReviewSettings{}).
		Where("tenant_id = ?", tenantID).
		Updates(updates)
	if result.Error != nil {
		return sharedErrors.NewInternalError("Failed to update review settings", result.Error)
	}
	if result.RowsAffected == 0 {
		return sharedErrors.NewNotFoundError("Review settings not found")
	}
	return nil
}

// Analytics operations
func (r *repository) GetReviewStatsByPeriod(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) (*ReviewStats, error) {
	var stats ReviewStats

	// Get total reviews in period
	var totalReviews int64
	err := r.db.WithContext(ctx).
		Model(&Review{}).
		Where("tenant_id = ? AND created_at >= ? AND created_at <= ?", tenantID, startDate, endDate).
		Count(&totalReviews).Error
	stats.TotalReviews = int(totalReviews)
	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to get total reviews count", err)
	}

	// Get approved reviews in period
	var approvedReviews int64
	err = r.db.WithContext(ctx).
		Model(&Review{}).
		Where("tenant_id = ? AND status = ? AND created_at >= ? AND created_at <= ?", tenantID, StatusApproved, startDate, endDate).
		Count(&approvedReviews).Error
	stats.ApprovedReviews = int(approvedReviews)
	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to get approved reviews count", err)
	}

	// Get average rating in period
	type avgResult struct {
		AvgRating *float64 `json:"avg_rating"`
	}
	var result avgResult
	err = r.db.WithContext(ctx).
		Model(&Review{}).
		Select("AVG(rating) as avg_rating").
		Where("tenant_id = ? AND status = ? AND created_at >= ? AND created_at <= ?", tenantID, StatusApproved, startDate, endDate).
		Scan(&result).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to get average rating", err)
	}

	if result.AvgRating != nil {
		stats.AverageRating = *result.AvgRating
	}

	return &stats, nil
}

func (r *repository) GetTopRatedProducts(ctx context.Context, tenantID uuid.UUID, limit int) ([]ProductRating, error) {
	var ratings []ProductRating

	// This would need to join with products table to get product names
	// For now, returning the structure expected
	err := r.db.WithContext(ctx).
		Table("review_summaries rs").
		Select("rs.product_id, rs.average_rating, rs.approved_reviews as total_reviews").
		Where("rs.tenant_id = ? AND rs.type = ? AND rs.approved_reviews > 0", tenantID, TypeProduct).
		Order("rs.average_rating DESC, rs.approved_reviews DESC").
		Limit(limit).
		Find(&ratings).Error

	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to get top rated products", err)
	}
	return ratings, nil
}

func (r *repository) GetReviewCountByStatus(ctx context.Context, tenantID uuid.UUID) (map[ReviewStatus]int, error) {
	type statusCount struct {
		Status ReviewStatus `json:"status"`
		Count  int          `json:"count"`
	}

	var results []statusCount
	err := r.db.WithContext(ctx).
		Model(&Review{}).
		Select("status, COUNT(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("status").
		Find(&results).Error

	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to get review count by status", err)
	}

	counts := make(map[ReviewStatus]int)
	for _, result := range results {
		counts[result.Status] = result.Count
	}

	return counts, nil
}

func (r *repository) GetReviewCountByRating(ctx context.Context, tenantID uuid.UUID) (map[int]int, error) {
	type ratingCount struct {
		Rating int `json:"rating"`
		Count  int `json:"count"`
	}

	var results []ratingCount
	err := r.db.WithContext(ctx).
		Model(&Review{}).
		Select("rating, COUNT(*) as count").
		Where("tenant_id = ? AND status = ?", tenantID, StatusApproved).
		Group("rating").
		Find(&results).Error

	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to get review count by rating", err)
	}

	counts := make(map[int]int)
	for _, result := range results {
		counts[result.Rating] = result.Count
	}

	return counts, nil
}
