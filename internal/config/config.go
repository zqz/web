package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	// Server
	Port int    `env:"PORT" envDefault:"3000"`
	Host string `env:"HOST" envDefault:"0.0.0.0"`
	Env  string `env:"ENV" envDefault:"development"`

	// Database
	DatabaseURL string `env:"DATABASE_URL,required"`

	// Storage
	FilesPath string `env:"FILES_PATH" envDefault:"./files"`

	// OAuth (Google)
	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
	GoogleCallbackURL  string `env:"GOOGLE_CALLBACK_URL"`

	// Session
	SessionSecret string `env:"SESSION_SECRET,required"`

	// Feature flags
	EnableThumbnails bool `env:"ENABLE_THUMBNAILS" envDefault:"true"`
	ThumbnailSize    int  `env:"THUMBNAIL_SIZE" envDefault:"256"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	var cfg Config

	// Load .env file in development
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	// Parse environment variables
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	if c.SessionSecret == "" {
		return fmt.Errorf("SESSION_SECRET is required")
	}

	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("PORT must be between 1 and 65535")
	}

	return nil
}

// IsProduction returns true if running in production
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

// IsDevelopment returns true if running in development
func (c *Config) IsDevelopment() bool {
	return c.Env == "development" || c.Env == ""
}

// Address returns the server address
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
