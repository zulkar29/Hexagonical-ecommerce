package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"ecommerce-saas/internal/shared/config"
	"ecommerce-saas/internal/shared/database"
	"ecommerce-saas/internal/notification"
	"ecommerce-saas/internal/webhook"
	"gorm.io/gorm"
)

// Worker represents a background worker
type Worker struct {
	db                  *gorm.DB
	config              *config.Config
	ctx                 context.Context
	cancel              context.CancelFunc
	wg                  sync.WaitGroup
	notificationService notification.Service
	notificationRepo    notification.Repository
	webhookService      *webhook.Service
}

// NewWorker creates a new worker instance
func NewWorker(db *gorm.DB, cfg *config.Config, notificationService notification.Service, notificationRepo notification.Repository, webhookService *webhook.Service) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker{
		db:                  db,
		config:              cfg,
		ctx:                 ctx,
		cancel:              cancel,
		notificationService: notificationService,
		notificationRepo:    notificationRepo,
		webhookService:      webhookService,
	}
}

// Start begins worker processing
func (w *Worker) Start() {
	log.Println("Starting background workers...")

	// Start email queue processor
	w.wg.Add(1)
	go w.processEmailQueue()

	// Start notification dispatcher
	w.wg.Add(1)
	go w.processNotifications()

	// Start cleanup tasks
	w.wg.Add(1)
	go w.processCleanupTasks()

	// Start webhook processor
	w.wg.Add(1)
	go w.processWebhooks()

	log.Println("All workers started successfully")
}

// Stop gracefully shuts down all workers
func (w *Worker) Stop() {
	log.Println("Stopping workers...")
	w.cancel()
	w.wg.Wait()
	log.Println("All workers stopped")
}

// processEmailQueue handles email sending queue
func (w *Worker) processEmailQueue() {
	defer w.wg.Done()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			// Process pending emails
			log.Println("Processing email queue...")
			
			// Get pending email notifications by status
			notifications, _, err := w.notificationRepo.ListByStatus(uuid.Nil, "pending", 0, 100)
			if err != nil {
				log.Printf("Failed to get pending email notifications: %v", err)
				return
			}
			
			for _, notif := range notifications {
				// Send email using notification service
			emailReq := &notification.SendEmailRequest{
				To:      []string{notif.Recipient},
				Subject: notif.Subject,
				Content: notif.Content,
			}
				
				err := w.notificationService.SendEmail(notif.TenantID, emailReq)
				if err != nil {
					log.Printf("Failed to send email notification %s: %v", notif.ID, err)
					// Update notification status to failed
					notif.Status = "failed"
					notif.FailureReason = err.Error()
				} else {
					// Update notification status to delivered
					notif.Status = "delivered"
					now := time.Now()
					notif.DeliveredAt = &now
				}
				
				// Update notification in database
				if updateErr := w.notificationRepo.Update(notif); updateErr != nil {
					log.Printf("Failed to update notification status: %v", updateErr)
				}
			}
		}
	}
}

// processNotifications handles notification dispatch
func (w *Worker) processNotifications() {
	defer w.wg.Done()
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			// Process pending notifications
			log.Println("Processing notifications...")
			
			// Process push notifications
			pushNotifications, _, err := w.notificationRepo.ListByStatus(uuid.Nil, "pending", 0, 100)
			if err != nil {
				log.Printf("Failed to get pending push notifications: %v", err)
			} else {
				for _, notif := range pushNotifications {
					// Send SMS notification using notification service
			smsReq := &notification.SendSMSRequest{
				To:      []string{notif.Recipient},
				Message: notif.Content,
			}
					
					err := w.notificationService.SendSMS(notif.TenantID, smsReq)
					if err != nil {
						log.Printf("Failed to send push notification %s: %v", notif.ID, err)
						notif.Status = "failed"
						notif.FailureReason = err.Error()
					} else {
						notif.Status = "delivered"
						now := time.Now()
						notif.DeliveredAt = &now
					}
					
					if updateErr := w.notificationRepo.Update(notif); updateErr != nil {
						log.Printf("Failed to update notification status: %v", updateErr)
					}
				}
			}
			
			// Process in-app notifications (just mark as delivered since they're stored in DB)
			inAppNotifications, _, err := w.notificationRepo.ListByStatus(uuid.Nil, "pending", 0, 100)
			if err != nil {
				log.Printf("Failed to get pending in-app notifications: %v", err)
			} else {
				for _, notif := range inAppNotifications {
					// In-app notifications are delivered by storing in database
					notif.Status = "delivered"
					now := time.Now()
					notif.DeliveredAt = &now
					
					if updateErr := w.notificationRepo.Update(notif); updateErr != nil {
						log.Printf("Failed to update in-app notification status: %v", updateErr)
					}
				}
			}
		}
	}
}

// processCleanupTasks handles data cleanup
func (w *Worker) processCleanupTasks() {
	defer w.wg.Done()
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			// Perform cleanup tasks
			w.cleanup()
		}
	}
}

func (w *Worker) cleanup() {
	log.Println("Running cleanup tasks...")
	
	// Clean up expired sessions (older than 7 days)
	cutoffTime := time.Now().AddDate(0, 0, -7)
	if err := w.cleanupExpiredSessions(cutoffTime); err != nil {
		log.Printf("Failed to cleanup expired sessions: %v", err)
	}
	
	// Clean up old notifications (older than 30 days)
	notificationCutoff := time.Now().AddDate(0, 0, -30)
	if err := w.cleanupOldNotifications(notificationCutoff); err != nil {
		log.Printf("Failed to cleanup old notifications: %v", err)
	}
	
	// Clean up old audit logs (older than 90 days)
	auditCutoff := time.Now().AddDate(0, 0, -90)
	if err := w.cleanupOldAuditLogs(auditCutoff); err != nil {
		log.Printf("Failed to cleanup old audit logs: %v", err)
	}
	
	// Clean up temporary files and cache
	if err := w.cleanupTempFiles(); err != nil {
		log.Printf("Failed to cleanup temporary files: %v", err)
	}
	
	log.Println("Cleanup tasks completed")
}

func (w *Worker) cleanupExpiredSessions(cutoffTime time.Time) error {
	// TODO: Implement session cleanup when session management is added
	log.Printf("Cleaning up sessions older than %v", cutoffTime)
	return nil
}

func (w *Worker) cleanupOldNotifications(cutoffTime time.Time) error {
	log.Printf("Cleaning up notifications older than %v", cutoffTime)
	// Use raw SQL to delete old notifications
	result := w.db.Exec("DELETE FROM notifications WHERE created_at < ? AND status IN ('delivered', 'failed')", cutoffTime)
	if result.Error != nil {
		return result.Error
	}
	log.Printf("Deleted %d old notifications", result.RowsAffected)
	return nil
}

func (w *Worker) cleanupOldAuditLogs(cutoffTime time.Time) error {
	log.Printf("Cleaning up audit logs older than %v", cutoffTime)
	// Use raw SQL to delete old audit logs
	result := w.db.Exec("DELETE FROM audit_logs WHERE created_at < ?", cutoffTime)
	if result.Error != nil {
		return result.Error
	}
	log.Printf("Deleted %d old audit logs", result.RowsAffected)
	return nil
}

func (w *Worker) cleanupTempFiles() error {
	log.Println("Cleaning up temporary files")
	// TODO: Implement temp file cleanup based on application needs
	// This could include cleaning up uploaded files, cache files, etc.
	return nil
}

// processWebhooks handles webhook processing
func (w *Worker) processWebhooks() {
	defer w.wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			// Process pending webhooks
			log.Println("Processing webhooks...")
			
			// Get pending webhook deliveries
			deliveries, err := w.webhookService.GetPendingRetries(100)
			if err != nil {
				log.Printf("Error getting pending webhook deliveries: %v", err)
				return
			}
			
			// Process webhook deliveries
			for _, delivery := range deliveries {
				log.Printf("Processing webhook delivery: %s", delivery.ID)
				// TODO: Implement webhook retry logic
			}
			
			// Additional webhook processing can be added here
			// The ProcessBackgroundTasks method handles:
			// - Retry failed webhook deliveries
			// - Clean up old webhook data (older than 30 days)
			// - Disable failing endpoints with high failure rates
		}
	}
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	dbConn, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize repositories and services
	notificationRepo := notification.NewRepository(dbConn)
	notificationService := notification.NewService(notificationRepo)

	// Initialize webhook repository and service
	webhookRepo := webhook.NewRepository(dbConn)
	// Get webhook signing key from environment
	webhookSigningKey := []byte(os.Getenv("WEBHOOK_SIGNING_KEY"))
	if len(webhookSigningKey) == 0 {
		webhookSigningKey = []byte("default-signing-key") // fallback for development
	}
	webhookService := webhook.NewService(webhookRepo, webhookSigningKey)

	// Create worker
	worker := NewWorker(dbConn, cfg, notificationService, notificationRepo, webhookService)

	// Start worker
	worker.Start()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Stop worker
	worker.Stop()
	log.Println("Worker stopped")
}
