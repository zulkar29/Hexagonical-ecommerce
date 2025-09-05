package reviews

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ReviewStatus represents the status of a review
type ReviewStatus string

// ReviewType represents different types of reviews
type ReviewType string

const (
	StatusPending  ReviewStatus = "pending"
	StatusApproved ReviewStatus = "approved"
	StatusRejected ReviewStatus = "rejected"
	StatusSpam     ReviewStatus = "spam"
	StatusDeleted  ReviewStatus = "deleted"
)

const (
	TypeProduct ReviewType = "product"
	TypeStore   ReviewType = "store"
	TypeOrder   ReviewType = "order"
)

// Review represents a customer review
type Review struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	
	// What is being reviewed
	ProductID *uuid.UUID `json:"product_id,omitempty" gorm:"index"`
	OrderID   *uuid.UUID `json:"order_id,omitempty" gorm:"index"`
	Type      ReviewType `json:"type" gorm:"not null;default:product"`
	
	// Reviewer information
	UserID        *uuid.UUID `json:"user_id,omitempty" gorm:"index"`
	CustomerName  string     `json:"customer_name" gorm:"not null"`
	CustomerEmail string     `json:"customer_email" gorm:"not null"`
	IsVerified    bool       `json:"is_verified" gorm:"default:false"` // Verified buyer
	
	// Review content
	Rating  int    `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content" gorm:"type:text"`
	
	// Review attributes
	Pros []string `json:"pros,omitempty" gorm:"serializer:json"`
	Cons []string `json:"cons,omitempty" gorm:"serializer:json"`
	
	// Media attachments
	Images []string `json:"images,omitempty" gorm:"serializer:json"`
	Videos []string `json:"videos,omitempty" gorm:"serializer:json"`
	
	// Review status and moderation
	Status       ReviewStatus `json:"status" gorm:"default:pending"`
	ModerationNote string     `json:"moderation_note,omitempty"`
	ModeratedBy  *uuid.UUID   `json:"moderated_by,omitempty" gorm:"index"`
	ModeratedAt  *time.Time   `json:"moderated_at,omitempty"`
	
	// Engagement metrics
	HelpfulCount   int `json:"helpful_count" gorm:"default:0"`
	UnhelpfulCount int `json:"unhelpful_count" gorm:"default:0"`
	ReportCount    int `json:"report_count" gorm:"default:0"`
	
	// Metadata
	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	Source    string `json:"source,omitempty"` // web, mobile_app, etc.
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relations
	Replies   []ReviewReply    `json:"replies,omitempty" gorm:"foreignKey:ReviewID"`
	Reactions []ReviewReaction `json:"reactions,omitempty" gorm:"foreignKey:ReviewID"`
}

// ReviewReply represents a reply to a review (from merchant or other customers)
type ReviewReply struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	ReviewID uuid.UUID `json:"review_id" gorm:"not null;index"`
	
	// Reply author
	UserID      *uuid.UUID `json:"user_id,omitempty" gorm:"index"` // User (merchant or customer)
	AuthorName  string     `json:"author_name" gorm:"not null"`
	AuthorEmail string     `json:"author_email" gorm:"not null"`
	IsMerchant  bool       `json:"is_merchant" gorm:"default:false"`
	
	// Reply content
	Content string `json:"content" gorm:"type:text;not null"`
	
	// Status
	IsVisible bool `json:"is_visible" gorm:"default:true"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ReviewReaction represents helpful/unhelpful reactions to reviews
type ReviewReaction struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	ReviewID uuid.UUID `json:"review_id" gorm:"not null;index"`
	
	// Reactor information
	UserID        *uuid.UUID `json:"user_id,omitempty" gorm:"index"`
	CustomerEmail string     `json:"customer_email,omitempty"`
	
	// Reaction type
	IsHelpful bool `json:"is_helpful"` // true = helpful, false = unhelpful
	
	// Metadata
	IPAddress string `json:"ip_address,omitempty"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	
	// Unique constraint to prevent duplicate reactions
	// UniqueIndex: tenant_id, review_id, customer_email
}

// ReviewSummary represents aggregated review statistics for a product/store
type ReviewSummary struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	
	// What this summary is for
	ProductID *uuid.UUID `json:"product_id,omitempty" gorm:"uniqueIndex:idx_review_summary_product"`
	Type      ReviewType `json:"type" gorm:"not null;default:product"`
	
	// Aggregate statistics
	TotalReviews    int     `json:"total_reviews" gorm:"default:0"`
	ApprovedReviews int     `json:"approved_reviews" gorm:"default:0"`
	AverageRating   float64 `json:"average_rating" gorm:"default:0"`
	
	// Rating distribution
	Rating1Count int `json:"rating_1_count" gorm:"default:0"`
	Rating2Count int `json:"rating_2_count" gorm:"default:0"`
	Rating3Count int `json:"rating_3_count" gorm:"default:0"`
	Rating4Count int `json:"rating_4_count" gorm:"default:0"`
	Rating5Count int `json:"rating_5_count" gorm:"default:0"`
	
	// Additional metrics
	VerifiedReviews int `json:"verified_reviews" gorm:"default:0"`
	WithPhotos      int `json:"with_photos" gorm:"default:0"`
	WithVideos      int `json:"with_videos" gorm:"default:0"`
	
	// Timestamps
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// ReviewSettings represents review module settings for a tenant
type ReviewSettings struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"unique;not null"`
	
	// Review collection settings
	ReviewsEnabled       bool `json:"reviews_enabled" gorm:"default:true"`
	RequireModeration    bool `json:"require_moderation" gorm:"default:true"`
	RequireVerifiedBuyer bool `json:"require_verified_buyer" gorm:"default:false"`
	AllowAnonymous       bool `json:"allow_anonymous" gorm:"default:false"`
	
	// Review content settings
	AllowImages          bool `json:"allow_images" gorm:"default:true"`
	AllowVideos          bool `json:"allow_videos" gorm:"default:false"`
	MaxImagesPerReview   int  `json:"max_images_per_review" gorm:"default:5"`
	MaxVideosPerReview   int  `json:"max_videos_per_review" gorm:"default:1"`
	MinContentLength     int  `json:"min_content_length" gorm:"default:10"`
	MaxContentLength     int  `json:"max_content_length" gorm:"default:2000"`
	
	// Engagement settings
	AllowReactions       bool `json:"allow_reactions" gorm:"default:true"`
	AllowReplies         bool `json:"allow_replies" gorm:"default:true"`
	AllowCustomerReplies bool `json:"allow_customer_replies" gorm:"default:false"`
	
	// Display settings
	ShowReviewerName     bool `json:"show_reviewer_name" gorm:"default:true"`
	ShowReviewDate       bool `json:"show_review_date" gorm:"default:true"`
	ShowVerifiedBadge    bool `json:"show_verified_badge" gorm:"default:true"`
	ReviewsPerPage       int  `json:"reviews_per_page" gorm:"default:20"`
	
	// Notification settings
	EmailOnNewReview     bool   `json:"email_on_new_review" gorm:"default:true"`
	NotificationEmail    string `json:"notification_email,omitempty"`
	AutoRequestReviews   bool   `json:"auto_request_reviews" gorm:"default:false"`
	RequestReviewsAfter  int    `json:"request_reviews_after" gorm:"default:7"` // days after delivery
	
	// Incentive settings
	RewardForReviews     bool    `json:"reward_for_reviews" gorm:"default:false"`
	RewardPoints         int     `json:"reward_points" gorm:"default:0"`
	RewardDiscount       float64 `json:"reward_discount" gorm:"default:0"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ReviewInvitation represents review invitation/request sent to customers
type ReviewInvitation struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	
	// Order and customer information
	OrderID       uuid.UUID `json:"order_id" gorm:"not null;index"`
	UserID        *uuid.UUID `json:"user_id,omitempty" gorm:"index"`
	CustomerEmail string    `json:"customer_email" gorm:"not null"`
	CustomerName  string    `json:"customer_name,omitempty"`
	
	// Products to review
	ProductIDs []string `json:"product_ids" gorm:"serializer:json"`
	
	// Invitation status
	Status      string     `json:"status" gorm:"default:pending"` // pending, sent, clicked, completed
	SentAt      *time.Time `json:"sent_at,omitempty"`
	ClickedAt   *time.Time `json:"clicked_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	
	// Invitation details
	InvitationToken  string    `json:"invitation_token" gorm:"unique;not null"`
	ExpiresAt        time.Time `json:"expires_at" gorm:"not null"`
	ReminderCount    int       `json:"reminder_count" gorm:"default:0"`
	LastReminderSent *time.Time `json:"last_reminder_sent,omitempty"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Business Logic Methods

// IsExpired checks if a review invitation has expired
func (ri *ReviewInvitation) IsExpired() bool {
	return time.Now().After(ri.ExpiresAt)
}

// CanSendReminder checks if a reminder can be sent for this invitation
func (ri *ReviewInvitation) CanSendReminder(maxReminders int, reminderInterval time.Duration) bool {
	if ri.Status != "sent" || ri.ReminderCount >= maxReminders {
		return false
	}
	
	if ri.LastReminderSent != nil {
		return time.Since(*ri.LastReminderSent) >= reminderInterval
	}
	
	return time.Since(*ri.SentAt) >= reminderInterval
}

// CalculateRatingPercentage calculates percentage for each rating
func (rs *ReviewSummary) CalculateRatingPercentage(rating int) float64 {
	if rs.ApprovedReviews == 0 {
		return 0
	}
	
	var count int
	switch rating {
	case 1:
		count = rs.Rating1Count
	case 2:
		count = rs.Rating2Count
	case 3:
		count = rs.Rating3Count
	case 4:
		count = rs.Rating4Count
	case 5:
		count = rs.Rating5Count
	default:
		return 0
	}
	
	return float64(count) / float64(rs.ApprovedReviews) * 100
}

// UpdateAverageRating recalculates the average rating
func (rs *ReviewSummary) UpdateAverageRating() {
	if rs.ApprovedReviews == 0 {
		rs.AverageRating = 0
		return
	}
	
	totalPoints := (rs.Rating1Count * 1) +
		(rs.Rating2Count * 2) +
		(rs.Rating3Count * 3) +
		(rs.Rating4Count * 4) +
		(rs.Rating5Count * 5)
	
	rs.AverageRating = float64(totalPoints) / float64(rs.ApprovedReviews)
}

// Validation methods

// ValidateReview validates review data
func (r *Review) Validate() error {
	if r.Rating < 1 || r.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}
	
	if r.CustomerName == "" {
		return errors.New("customer name is required")
	}
	
	if r.CustomerEmail == "" {
		return errors.New("customer email is required")
	}
	
	if r.Type == TypeProduct && r.ProductID == nil {
		return errors.New("product ID is required for product reviews")
	}
	
	return nil
}

// ValidateReviewReply validates review reply data
func (rr *ReviewReply) Validate() error {
	if rr.Content == "" {
		return errors.New("reply content is required")
	}
	
	if rr.AuthorName == "" {
		return errors.New("author name is required")
	}
	
	if rr.AuthorEmail == "" {
		return errors.New("author email is required")
	}
	
	// UserID should be set for authenticated users
	if rr.UserID == nil && rr.AuthorEmail == "" {
		return errors.New("reply must have either user ID or author email")
	}
	
	return nil
}