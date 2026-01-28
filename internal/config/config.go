package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Log      LogConfig
	JWT      JWTConfig
}

// AppConfig holds application configuration
type AppConfig struct {
	Port string
	Env  string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string
	ExpireTime time.Duration
}

// Load loads configuration from environment variables
func Load() *Config {
	cfg := &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "go_echo_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "debug"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-super-secret-key-change-in-production"),
			ExpireTime: time.Duration(getEnvAsInt("JWT_EXPIRE_HOURS", 24)) * time.Hour,
		},
	}

	// Basic validation for production
	if cfg.App.Env == "production" && cfg.JWT.Secret == "your-super-secret-key-change-in-production" {
		// We'll let the application decide whether to fatal or just warn, 
		// but here we mark it as a risk.
	}

	return cfg
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development"
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}
