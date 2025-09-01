package config

// TODO: Implement simple configuration during development

// Config holds application configuration
type Config struct {
	// TODO: Add basic config fields
	// Port         string `mapstructure:"PORT"`
	// DatabaseURL  string `mapstructure:"DATABASE_URL"`
	// RedisURL     string `mapstructure:"REDIS_URL"`
	// JWTSecret    string `mapstructure:"JWT_SECRET"`
	// Environment  string `mapstructure:"ENVIRONMENT"`
}

// Load loads configuration from environment
func Load() *Config {
	// TODO: Implement config loading using Viper
	// - Load from .env file
	// - Override with environment variables
	// - Set defaults
	return &Config{}
}

// IsDevelopment checks if running in development mode
func (c *Config) IsDevelopment() bool {
	// TODO: Implement environment check
	return true // placeholder
}

// IsProduction checks if running in production mode
func (c *Config) IsProduction() bool {
	// TODO: Implement environment check
	return false // placeholder
}