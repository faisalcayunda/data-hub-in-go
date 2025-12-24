package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	MinIO    MinIOConfig
}

// AppConfig contains application metadata
type AppConfig struct {
	Name        string
	Environment string
	Debug       bool
	Version     string
}

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig contains database connection configuration
type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	Database     string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

// RedisConfig contains Redis connection configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// JWTConfig contains JWT token configuration
type JWTConfig struct {
	Secret            string
	AccessTokenExpiry time.Duration
	RefreshTokenExpiry time.Duration
	Issuer            string
}

// MinIOConfig contains MinIO storage configuration
type MinIOConfig struct {
	Endpoint        string
	AccessKey       string
	SecretKey       string
	Bucket          string
	UseSSL          bool
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if exists (optional for production)
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "portal-data-backend"),
			Environment: getEnv("APP_ENV", "development"),
			Debug:       getEnv("APP_DEBUG", "true") == "true",
			Version:     getEnv("APP_VERSION", "1.0.0"),
		},
		Server: ServerConfig{
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnvAsInt("DB_PORT", 5432),
			User:         getEnv("DB_USER", "postgres"),
			Password:     getEnv("DB_PASSWORD", ""),
			Database:     getEnv("DB_NAME", "portal_data"),
			SSLMode:      getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			MaxLifetime:  getEnvAsDuration("DB_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:            getEnv("JWT_SECRET", "change-me-in-production"),
			AccessTokenExpiry: getEnvAsDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
			RefreshTokenExpiry: getEnvAsDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
			Issuer:            getEnv("JWT_ISSUER", "portal-data-backend"),
		},
		MinIO: MinIOConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("MINIO_BUCKET", "portal-data"),
			UseSSL:    getEnv("MINIO_USE_SSL", "false") == "true",
		},
	}

	// Validate required configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if c.Database.Database == "" {
		return fmt.Errorf("database name is required")
	}
	if c.JWT.Secret == "" || c.JWT.Secret == "change-me-in-production" {
		if c.App.Environment == "production" {
			return fmt.Errorf("JWT secret must be set in production")
		}
	}
	return nil
}

// DSN returns the PostgreSQL Data Source Name
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode,
	)
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := getEnvAsIntValue(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsIntValue(value string) (int, error) {
	var intVal int
	if _, err := fmt.Sscanf(value, "%d", &intVal); err != nil {
		return 0, fmt.Errorf("invalid integer value: %s", value)
	}
	return intVal, nil
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
