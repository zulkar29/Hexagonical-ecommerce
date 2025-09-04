package reviews

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *repository) GetReviewByID(ctx context.Context, tenantID, reviewID uuid.UUID) (*Review, error) {
	var review Review
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", reviewID, tenantID).
		Preload("Replies").
		Preload("Reactions").
		First(&review).Error
	return &review, err
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
	
	if filter.CustomerID != nil {
		query = query.Where("customer_id = ?", *filter.CustomerID)
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
	return reviews, err
}

func (r *repository) UpdateReview(ctx context.Context, tenantID, reviewID uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&Review{}).
		Where("id = ? AND tenant_id = ?", reviewID, tenantID).
		Updates(updates).Error
}

func (r *repository) DeleteReview(ctx context.Context, tenantID, reviewID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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
	return count, err
}

// Review reply operations
func (r *repository) CreateReply(ctx context.Context, reply *ReviewReply) error {
	return r.db.WithContext(ctx).Create(reply).Error
}

func (r *repository) GetRepliesByReviewID(ctx context.Context, tenantID, reviewID uuid.UUID) ([]ReviewReply, error) {
	var replies []ReviewReply
	
	// Verify review belongs to tenant first
	var count int64
	r.db.WithContext(ctx).Model(&Review{}).
		Where("id = ? AND tenant_id = ?", reviewID, tenantID).
		Count(&count)
	
	if count == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	
	err := r.db.WithContext(ctx).
		Where("review_id = ? AND is_visible = ?", reviewID, true).
		Order("created_at ASC").
		Find(&replies).Error
	
	return replies, err
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
	
	return &reply, err
}

func (r *repository) UpdateReply(ctx context.Context, tenantID, replyID uuid.UUID, updates map[string]interface{}) error {
	// Ensure tenant isolation through join
	return r.db.WithContext(ctx).
		Table("review_replies").
		Joins("JOIN reviews ON reviews.id = review_replies.review_id").
		Where("review_replies.id = ? AND reviews.tenant_id = ?", replyID, tenantID).
		Updates(updates).Error
}

func (r *repository) DeleteReply(ctx context.Context, tenantID, replyID uuid.UUID) error {
	// Ensure tenant isolation through join
	return r.db.WithContext(ctx).
		Table("review_replies").
		Joins("JOIN reviews ON reviews.id = review_replies.review_id").
		Where("review_replies.id = ? AND reviews.tenant_id = ?", replyID, tenantID).
		Delete(&ReviewReply{}).Error
}

// Review reaction operations
func (r *repository) CreateReaction(ctx context.Context, reaction *ReviewReaction) error {
	return r.db.WithContext(ctx).Create(reaction).Error
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
	
	return &reaction, err
}

func (r *repository) UpdateReaction(ctx context.Context, tenantID, reviewID uuid.UUID, email string, isHelpful bool) error {
	// Verify tenant isolation and update
	return r.db.WithContext(ctx).
		Table("review_reactions").
		Joins("JOIN reviews ON reviews.id = review_reactions.review_id").
		Where("review_reactions.review_id = ? AND review_reactions.customer_email = ? AND reviews.tenant_id = ?", 
			reviewID, email, tenantID).
		Update("is_helpful", isHelpful).Error
}

func (r *repository) DeleteReaction(ctx context.Context, tenantID, reviewID uuid.UUID, email string) error {
	return r.db.WithContext(ctx).
		Table("review_reactions").
		Joins("JOIN reviews ON reviews.id = review_reactions.review_id").
		Where("review_reactions.review_id = ? AND review_reactions.customer_email = ? AND reviews.tenant_id = ?", 
			reviewID, email, tenantID).
		Delete(&ReviewReaction{}).Error
}

func (r *repository) UpdateReviewReactionCounts(ctx context.Context, reviewID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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
}

// Review summary operations
func (r *repository) CreateReviewSummary(ctx context.Context, summary *ReviewSummary) error {
	return r.db.WithContext(ctx).Create(summary).Error
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
	}
	
	return &summary, err
}

func (r *repository) UpdateReviewSummary(ctx context.Context, tenantID, productID uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&ReviewSummary{}).
		Where("tenant_id = ? AND product_id = ? AND type = ?", tenantID, productID, TypeProduct).
		Updates(updates).Error
}

func (r *repository) RecalculateReviewSummary(ctx context.Context, tenantID, productID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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
		var totalReviews, verifiedCount, withPhotosCount int
		var ratingCounts [6]int // Index 0 unused, 1-5 for ratings
		var totalPoints int
		
		for _, r := range ratings {
			totalReviews += r.Count
			totalPoints += r.Rating * r.Count
			if r.Rating >= 1 && r.Rating <= 5 {
				ratingCounts[r.Rating] = r.Count
			}
		}
		
		// Count verified reviews and reviews with photos
		tx.Model(&Review{}).
			Where("tenant_id = ? AND product_id = ? AND status = ? AND is_verified = ?", 
				tenantID, productID, StatusApproved, true).
			Count((*int64)(&verifiedCount))
		
		tx.Model(&Review{}).
			Where("tenant_id = ? AND product_id = ? AND status = ? AND JSON_LENGTH(images) > 0", 
				tenantID, productID, StatusApproved).
			Count((*int64)(&withPhotosCount))
		
		// Calculate average rating
		var avgRating float64
		if totalReviews > 0 {
			avgRating = float64(totalPoints) / float64(totalReviews)
		}
		
		// Update summary
		updates := map[string]interface{}{
			"total_reviews":     totalReviews,
			"approved_reviews":  totalReviews,
			"average_rating":    avgRating,
			"rating_1_count":    ratingCounts[1],
			"rating_2_count":    ratingCounts[2],
			"rating_3_count":    ratingCounts[3],
			"rating_4_count":    ratingCounts[4],
			"rating_5_count":    ratingCounts[5],
			"verified_reviews":  verifiedCount,
			"with_photos":       withPhotosCount,
			"updated_at":        time.Now(),
		}
		
		return r.UpdateReviewSummary(ctx, tenantID, productID, updates)
	})
}

// Review invitation operations
func (r *repository) CreateInvitation(ctx context.Context, invitation *ReviewInvitation) error {
	return r.db.WithContext(ctx).Create(invitation).Error
}

func (r *repository) GetInvitationByID(ctx context.Context, tenantID, invitationID uuid.UUID) (*ReviewInvitation, error) {
	var invitation ReviewInvitation
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", invitationID, tenantID).
		First(&invitation).Error
	return &invitation, err
}

func (r *repository) GetInvitationByToken(ctx context.Context, token string) (*ReviewInvitation, error) {
	var invitation ReviewInvitation
	err := r.db.WithContext(ctx).
		Where("invitation_token = ?", token).
		First(&invitation).Error
	return &invitation, err
}

func (r *repository) GetInvitationsByStatus(ctx context.Context, tenantID uuid.UUID, status string) ([]ReviewInvitation, error) {
	var invitations []ReviewInvitation
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND status = ?", tenantID, status).
		Order("created_at DESC").
		Find(&invitations).Error
	return invitations, err
}

func (r *repository) UpdateInvitation(ctx context.Context, tenantID, invitationID uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&ReviewInvitation{}).
		Where("id = ? AND tenant_id = ?", invitationID, tenantID).
		Updates(updates).Error
}

func (r *repository) GetExpiredInvitations(ctx context.Context, tenantID uuid.UUID) ([]ReviewInvitation, error) {
	var invitations []ReviewInvitation
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND status = ? AND expires_at < ?", tenantID, "sent", time.Now()).
		Find(&invitations).Error
	return invitations, err
}

// Settings operations
func (r *repository) CreateSettings(ctx context.Context, settings *ReviewSettings) error {
	return r.db.WithContext(ctx).Create(settings).Error
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
	}
	
	return &settings, err
}

func (r *repository) UpdateSettings(ctx context.Context, tenantID uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&ReviewSettings{}).
		Where("tenant_id = ?", tenantID).
		Updates(updates).Error
}

// Analytics operations
func (r *repository) GetReviewStatsByPeriod(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) (*ReviewStats, error) {
	// TODO: Implement comprehensive analytics query
	return nil, fmt.Errorf("TODO: implement GetReviewStatsByPeriod")
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
	
	return ratings, err
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
		return nil, err
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
		return nil, err
	}
	
	counts := make(map[int]int)
	for _, result := range results {
		counts[result.Rating] = result.Count
	}
	
	return counts, nil
}