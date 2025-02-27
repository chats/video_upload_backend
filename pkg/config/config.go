package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server    ServerConfig
	DB        DBConfig
	Storage   StorageConfig
	Transcode TranscodeConfig
	Auth      AuthConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string
	AppEnv       string
	LogLevel     string
	AllowOrigins string
}

// DBConfig holds database configuration
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// StorageConfig holds S3/Minio configuration
type StorageConfig struct {
	AccessKey  string
	SecretKey  string
	Region     string
	BucketName string
	Endpoint   string
	UseSSL     bool
}

// TranscodeConfig holds FFmpeg configuration
type TranscodeConfig struct {
	FFmpegPath        string
	FFprobePath       string
	MaxConcurrentJobs int
	SegmentDuration   int
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret string
	JWTExpiry string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Default values
	config := &Config{
		Server: ServerConfig{
			Port:         getEnvOrDefault("SERVER_PORT", "8080"),
			AppEnv:       getEnvOrDefault("APP_ENV", "development"),
			LogLevel:     getEnvOrDefault("LOG_LEVEL", "info"),
			AllowOrigins: getEnvOrDefault("ALLOW_ORIGINS", "*"),
		},
		DB: DBConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvIntOrDefault("DB_PORT", 5432),
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
			Name:     getEnvOrDefault("DB_NAME", "video_system"),
			SSLMode:  getEnvOrDefault("DB_SSL_MODE", "disable"),
		},
		Storage: StorageConfig{
			AccessKey:  getEnvOrDefault("STORAGE_ACCESS_KEY", "minioadmin"),
			SecretKey:  getEnvOrDefault("STORAGE_SECRET_KEY", "minioadmin"),
			Region:     getEnvOrDefault("STORAGE_REGION", "us-east-1"),
			BucketName: getEnvOrDefault("STORAGE_BUCKET_NAME", "videos"),
			Endpoint:   getEnvOrDefault("STORAGE_ENDPOINT", "http://localhost:9000"),
			UseSSL:     getEnvBoolOrDefault("STORAGE_USE_SSL", false),
		},
		Transcode: TranscodeConfig{
			FFmpegPath:        getEnvOrDefault("FFMPEG_PATH", "ffmpeg"),
			FFprobePath:       getEnvOrDefault("FFPROBE_PATH", "ffprobe"),
			MaxConcurrentJobs: getEnvIntOrDefault("MAX_CONCURRENT_TRANSCODES", 2),
			SegmentDuration:   getEnvIntOrDefault("SEGMENT_DURATION", 10),
		},
		Auth: AuthConfig{
			JWTSecret: getEnvOrDefault("JWT_SECRET", "your-secret-key-change-this-in-production"),
			JWTExpiry: getEnvOrDefault("JWT_EXPIRY", "24h"),
		},
	}

	// Validate required configuration
	if config.Storage.AccessKey == "" || config.Storage.SecretKey == "" {
		return nil, fmt.Errorf("storage access key and secret key are required")
	}

	return config, nil
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvIntOrDefault gets an integer environment variable or returns a default value
func getEnvIntOrDefault(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// getEnvBoolOrDefault gets a boolean environment variable or returns a default value
func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		// Try to parse "yes", "no", etc.
		switch strings.ToLower(valueStr) {
		case "yes", "y", "true", "t", "1":
			return true
		default:
			return false
		}
	}

	return value
}
