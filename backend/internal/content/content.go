package content

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TODO: Implement content management entities
// This will handle:
// - CMS pages and blog posts
// - SEO optimization
// - Media library management
// - Navigation and menu management

type ContentType string
type PageStatus string
type MediaType string

const (
	TypePage        ContentType = "page"
	TypeBlogPost    ContentType = "blog_post"
	TypeProduct     ContentType = "product"
	TypeCategory    ContentType = "category"
	TypeLanding     ContentType = "landing"
	TypeEmail       ContentType = "email_template"
)

const (
	StatusDraft     PageStatus = "draft"
	StatusPublished PageStatus = "published"
	StatusArchived  PageStatus = "archived"
	StatusScheduled PageStatus = "scheduled"
)

const (
	MediaImage MediaType = "image"
	MediaVideo MediaType = "video"
	MediaAudio MediaType = "audio"
	MediaDocument MediaType = "document"
)

type Page struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID    uuid.UUID   `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Type        ContentType `json:"type" gorm:"size:50;not null;default:'page'"`
	Title       string      `json:"title" gorm:"size:255;not null"`
	Slug        string      `json:"slug" gorm:"size:255;not null;index"`
	Content     string      `json:"content" gorm:"type:text"`
	Excerpt     string      `json:"excerpt" gorm:"type:text"`
	Status      PageStatus  `json:"status" gorm:"size:20;not null;default:'draft'"`
	
	// SEO fields
	MetaTitle       string `json:"meta_title" gorm:"size:255"`
	MetaDescription string `json:"meta_description" gorm:"type:text"`
	MetaKeywords    string `json:"meta_keywords" gorm:"type:text"`
	OGTitle         string `json:"og_title" gorm:"size:255"`
	OGDescription   string `json:"og_description" gorm:"type:text"`
	OGImage         string `json:"og_image" gorm:"size:500"`
	
	// Featured image
	FeaturedImageID *uuid.UUID `json:"featured_image_id" gorm:"type:uuid"`
	FeaturedImage   *Media     `json:"featured_image" gorm:"foreignKey:FeaturedImageID"`
	
	// Author and publishing
	AuthorID     uuid.UUID  `json:"author_id" gorm:"type:uuid;not null"`
	PublishedAt  *time.Time `json:"published_at"`
	ScheduledAt  *time.Time `json:"scheduled_at"`
	
	// Template and layout
	Template     string `json:"template" gorm:"size:100;default:'default'"`
	LayoutData   string `json:"layout_data" gorm:"type:json"` // For page builder data
	
	// Analytics
	ViewCount    int `json:"view_count" gorm:"default:0"`
	
	// Hierarchy (for pages)
	ParentID     *uuid.UUID `json:"parent_id" gorm:"type:uuid"`
	Parent       *Page      `json:"parent" gorm:"foreignKey:ParentID"`
	Children     []Page     `json:"children" gorm:"foreignKey:ParentID"`
	MenuOrder    int        `json:"menu_order" gorm:"default:0"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	
	// Relations
	Tags       []Tag       `json:"tags" gorm:"many2many:page_tags"`
	Categories []Category  `json:"categories" gorm:"many2many:page_categories"`
}

type Media struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID    uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Type        MediaType `json:"type" gorm:"size:20;not null"`
	Title       string    `json:"title" gorm:"size:255;not null"`
	Description string    `json:"description" gorm:"type:text"`
	FileName    string    `json:"file_name" gorm:"size:255;not null"`
	FilePath    string    `json:"file_path" gorm:"size:500;not null"`
	FileSize    int64     `json:"file_size" gorm:"not null"`
	MimeType    string    `json:"mime_type" gorm:"size:100;not null"`
	
	// Image specific
	Width       int    `json:"width" gorm:"default:0"`
	Height      int    `json:"height" gorm:"default:0"`
	AltText     string `json:"alt_text" gorm:"size:255"`
	
	// Video/Audio specific
	Duration    int    `json:"duration" gorm:"default:0"` // seconds
	
	// SEO
	MetaTitle       string `json:"meta_title" gorm:"size:255"`
	MetaDescription string `json:"meta_description" gorm:"type:text"`
	
	UploadedBy uuid.UUID `json:"uploaded_by" gorm:"type:uuid;not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Menu struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID    uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Location    string    `json:"location" gorm:"size:50;not null"` // header, footer, sidebar
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	
	Items []MenuItem `json:"items" gorm:"foreignKey:MenuID"`
}

type MenuItem struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MenuID      uuid.UUID  `json:"menu_id" gorm:"type:uuid;not null;index"`
	ParentID    *uuid.UUID `json:"parent_id" gorm:"type:uuid"`
	Title       string     `json:"title" gorm:"size:255;not null"`
	URL         string     `json:"url" gorm:"size:500"`
	PageID      *uuid.UUID `json:"page_id" gorm:"type:uuid"` // Link to internal page
	Target      string     `json:"target" gorm:"size:20;default:'_self'"`
	CSSClass    string     `json:"css_class" gorm:"size:100"`
	IconClass   string     `json:"icon_class" gorm:"size:100"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	SortOrder   int        `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	
	// Relations
	Parent   *MenuItem  `json:"parent" gorm:"foreignKey:ParentID"`
	Children []MenuItem `json:"children" gorm:"foreignKey:ParentID"`
	Page     *Page      `json:"page" gorm:"foreignKey:PageID"`
}

type Tag struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID    uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Slug        string    `json:"slug" gorm:"size:100;not null;index"`
	Description string    `json:"description" gorm:"type:text"`
	Color       string    `json:"color" gorm:"size:7"` // Hex color
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Category struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID    uuid.UUID  `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Name        string     `json:"name" gorm:"size:100;not null"`
	Slug        string     `json:"slug" gorm:"size:100;not null;index"`
	Description string     `json:"description" gorm:"type:text"`
	ParentID    *uuid.UUID `json:"parent_id" gorm:"type:uuid"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	SortOrder   int        `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	
	// Relations
	Parent   *Category  `json:"parent" gorm:"foreignKey:ParentID"`
	Children []Category `json:"children" gorm:"foreignKey:ParentID"`
}

type SEOSettings struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID    uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	
	// Global SEO settings
	SiteTitle         string `json:"site_title" gorm:"size:255"`
	SiteDescription   string `json:"site_description" gorm:"type:text"`
	DefaultMetaTitle  string `json:"default_meta_title" gorm:"size:255"`
	DefaultMetaDesc   string `json:"default_meta_description" gorm:"type:text"`
	DefaultOGImage    string `json:"default_og_image" gorm:"size:500"`
	
	// Analytics
	GoogleAnalyticsID string `json:"google_analytics_id" gorm:"size:50"`
	GoogleTagManagerID string `json:"google_tag_manager_id" gorm:"size:50"`
	FacebookPixelID   string `json:"facebook_pixel_id" gorm:"size:50"`
	
	// Verification codes
	GoogleVerification string `json:"google_verification" gorm:"size:100"`
	BingVerification   string `json:"bing_verification" gorm:"size:100"`
	
	// Robots and sitemaps
	RobotsTxt        string `json:"robots_txt" gorm:"type:text"`
	EnableSitemap    bool   `json:"enable_sitemap" gorm:"default:true"`
	SitemapFrequency string `json:"sitemap_frequency" gorm:"size:20;default:'weekly'"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Business Logic Methods

// IsPublished checks if page is published and live
func (p *Page) IsPublished() bool {
	return p.Status == StatusPublished && 
		   (p.PublishedAt == nil || p.PublishedAt.Before(time.Now()))
}

// IsScheduled checks if page is scheduled for future publishing
func (p *Page) IsScheduled() bool {
	return p.Status == StatusScheduled && 
		   p.ScheduledAt != nil && 
		   p.ScheduledAt.After(time.Now())
}

// GetURL generates the full URL for the page
func (p *Page) GetURL() string {
	if p.Type == TypeBlogPost {
		return "/blog/" + p.Slug
	}
	return "/" + p.Slug
}

// GetFullTitle returns meta title or falls back to title
func (p *Page) GetFullTitle() string {
	if p.MetaTitle != "" {
		return p.MetaTitle
	}
	return p.Title
}

// GetFullDescription returns meta description or falls back to excerpt
func (p *Page) GetFullDescription() string {
	if p.MetaDescription != "" {
		return p.MetaDescription
	}
	return p.Excerpt
}

// IsImage checks if media is an image
func (m *Media) IsImage() bool {
	return m.Type == MediaImage
}

// IsVideo checks if media is a video
func (m *Media) IsVideo() bool {
	return m.Type == MediaVideo
}

// GetThumbnailURL returns thumbnail URL for media
func (m *Media) GetThumbnailURL() string {
	if m.IsImage() {
		return m.FilePath + "?w=300&h=300"
	}
	return m.FilePath
}

// TODO: Add more business logic methods
// - GenerateSlug(title string) string
// - ValidateSlugUniqueness(slug string) bool
// - GenerateSitemap() string
// - OptimizeImageSizes() error
// - ValidateSEOScore() int
