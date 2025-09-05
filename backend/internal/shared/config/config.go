package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	// Server configuration
	Server ServerConfig `mapstructure:"server"`
	
	// Database configuration
	Database DatabaseConfig `mapstructure:"database"`
	
	// Redis configuration
	Redis RedisConfig `mapstructure:"redis"`
	
	// JWT configuration
	JWT JWTConfig `mapstructure:"jwt"`
	
	// Email configuration
	Email EmailConfig `mapstructure:"email"`
	
	// Storage configuration
	Storage StorageConfig `mapstructure:"storage"`
	
	// Payment configuration
	Payment PaymentConfig `mapstructure:"payment"`
	
	// Application configuration
	App AppConfig `mapstructure:"app"`
}

type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	MaxOpenConns int `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret      string        `mapstructure:"secret_key"`
	ExpiresIn   time.Duration `mapstructure:"expiration_time"`
	RefreshTime time.Duration `mapstructure:"refresh_time"`
}

type EmailConfig struct {
	SMTPHost     string `mapstructure:"smtp_host"`
	SMTPPort     int    `mapstructure:"smtp_port"`
	SMTPUser     string `mapstructure:"smtp_user"`
	SMTPPassword string `mapstructure:"smtp_password"`
	FromEmail    string `mapstructure:"from_email"`
	FromName     string `mapstructure:"from_name"`
}

type StorageConfig struct {
	Provider        string `mapstructure:"provider"` // "local" or "s3"
	LocalPath       string `mapstructure:"local_path"`
	S3Bucket        string `mapstructure:"s3_bucket"`
	S3Region        string `mapstructure:"s3_region"`
	S3AccessKey     string `mapstructure:"s3_access_key"`
	S3SecretKey     string `mapstructure:"s3_secret_key"`
	S3Endpoint      string `mapstructure:"s3_endpoint"` // For MinIO or other S3-compatible services
	MaxUploadSize   int64  `mapstructure:"max_upload_size"`
	AllowedFileTypes []string `mapstructure:"allowed_file_types"`
}

type PaymentConfig struct {
	Stripe PaymentProviderConfig `mapstructure:"stripe"`
	BKash  PaymentProviderConfig `mapstructure:"bkash"`
}

type PaymentProviderConfig struct {
	SecretKey    string `mapstructure:"secret_key"`
	PublishableKey string `mapstructure:"publishable_key"`
	WebhookSecret string `mapstructure:"webhook_secret"`
	Enabled      bool   `mapstructure:"enabled"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
	Debug       bool   `mapstructure:"debug"`
	LogLevel    string `mapstructure:"log_level"`
	Domain      string `mapstructure:"domain"`
	CORS        CORSConfig `mapstructure:"cors"`
	RateLimit   RateLimitConfig `mapstructure:"rate_limit"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
	AllowCredentials bool   `mapstructure:"allow_credentials"`
}

type RateLimitConfig struct {
	RequestsPerMinute int           `mapstructure:"requests_per_minute"`
	BurstSize         int           `mapstructure:"burst_size"`
	CleanupInterval   time.Duration `mapstructure:"cleanup_interval"`
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	// Set defaults
	setDefaults()

	// Configure viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/esass")

	// Enable environment variable reading
	viper.AutomaticEnv()
	
	// Set environment variable prefix
	viper.SetEnvPrefix("ESASS")

	// Try to read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %v", err)
	}

	// Override with environment variables
	loadFromEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "120s")

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.dbname", "esass")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_lifetime", "1h")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// JWT defaults
	viper.SetDefault("jwt.secret_key", "your-secret-key-change-in-production")
	viper.SetDefault("jwt.expiration_time", "24h")
	viper.SetDefault("jwt.refresh_time", "168h") // 7 days

	// Storage defaults
	viper.SetDefault("storage.provider", "local")
	viper.SetDefault("storage.local_path", "./uploads")
	viper.SetDefault("storage.max_upload_size", 10485760) // 10MB
	viper.SetDefault("storage.allowed_file_types", []string{"jpg", "jpeg", "png", "gif", "pdf", "doc", "docx"})

	// App defaults
	viper.SetDefault("app.name", "E-commerce SaaS")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.debug", true)
	viper.SetDefault("app.log_level", "debug")
	viper.SetDefault("app.domain", "esass.com")

	// CORS defaults
	viper.SetDefault("app.cors.allowed_origins", []string{"*"})
	viper.SetDefault("app.cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("app.cors.allowed_headers", []string{"*"})
	viper.SetDefault("app.cors.allow_credentials", true)

	// Rate limit defaults
	viper.SetDefault("app.rate_limit.requests_per_minute", 100)
	viper.SetDefault("app.rate_limit.burst_size", 10)
	viper.SetDefault("app.rate_limit.cleanup_interval", "5m")
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv() {
	// Server
	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			viper.Set("server.port", p)
		}
	}

	// Database - Docker environment
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		viper.Set("database.host", dbHost)
	} else {
		// Use postgres service name in Docker
		viper.Set("database.host", "postgres")
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		if p, err := strconv.Atoi(dbPort); err == nil {
			viper.Set("database.port", p)
		}
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		viper.Set("database.user", dbUser)
	} else {
		// Use Docker default
		viper.Set("database.user", "postgres")
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		viper.Set("database.password", dbPassword)
	} else {
		// Use Docker default
		viper.Set("database.password", "postgres123")
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		viper.Set("database.dbname", dbName)
	} else {
		// Use Docker default
		viper.Set("database.dbname", "ecommerce_saas_dev")
	}

	// Redis
	if redisHost := os.Getenv("REDIS_HOST"); redisHost != "" {
		viper.Set("redis.host", redisHost)
	}
	if redisPort := os.Getenv("REDIS_PORT"); redisPort != "" {
		if p, err := strconv.Atoi(redisPort); err == nil {
			viper.Set("redis.port", p)
		}
	}
	if redisPassword := os.Getenv("REDIS_PASSWORD"); redisPassword != "" {
		viper.Set("redis.password", redisPassword)
	}

	// JWT
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		viper.Set("jwt.secret_key", jwtSecret)
	}

	// Environment
	if env := os.Getenv("ENVIRONMENT"); env != "" {
		viper.Set("app.environment", env)
		if env == "production" {
			viper.Set("app.debug", false)
			viper.Set("app.log_level", "info")
		}
	}
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// GetRedisAddr returns the Redis connection address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

// TODO: Add more configuration methods as needed
// - ValidateConfig() error
// - LoadFromFile(path string) error
// - SaveToFile(path string) error
