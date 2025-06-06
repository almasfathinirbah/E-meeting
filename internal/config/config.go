package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration settings for the application
type Config struct {
	// Application environment and port settings
	AppEnv  string // Environment (development, production, etc.)
	AppPort string // Port on which the application runs

	// Database connection settings
	DBHost               string // Database host address
	DBPort               int    // Database port number
	DBUser               string // Database username
	DBPassword           string // Database password
	DBName               string // Database name
	DBMaxOpenConnections int    // Maximum number of open database connections
	DBMaxIdleConnections int    // Maximum number of idle database connections

	// JWT authentication settings
	JWT struct {
		SecretKey     string // Secret key for signing JWT tokens
		TokenDuration int    // Duration (in hours) for which JWT tokens are valid
	}

	// SMTP email service configuration
	SMTP struct {
		Host               string // SMTP server host
		Port               int    // SMTP server port
		Username           string // SMTP authentication username
		Password           string // SMTP authentication password
		FromEmail          string // Email address used as sender
		TemplatePath       string // Path to email template files
		TemplateLogoURL    string // URL for logo in email templates
		TimeoutDuration    int    // Timeout duration for SMTP operations
		InsecureSkipVerify bool   // Skip TLS certificate verification
		UseTLS             bool   // Use TLS for SMTP connection
	}

	// Cloudflare R2 storage configuration
	CloudflareR2BucketName string // R2 bucket name
	CloudflareR2APIKey     string // Cloudflare API key
	CloudflareR2APISecret  string // Cloudflare API secret
	CloudflareR2Token      string // R2 access token
	CloudflareR2AccountID  string // Cloudflare account ID
	CloudflareR2PublicURL  string // Public URL for R2 bucket

	// Server configuration
	Server struct {
		Port int // Server port number
	}
}

// findRootDir searches for the root directory of the project by looking for .env file
// Returns:
//   - The root directory path where .env file is found, or "." if not found
func findRootDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".env")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "."
		}
		dir = parent
	}
}

// setDefaults initializes default configuration values using Viper
// This ensures the application has sensible defaults when configuration values are not explicitly set
func setDefaults() {
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("DATABASE_HOST", "localhost")
	viper.SetDefault("DATABASE_PORT", 5432)
	viper.SetDefault("DATABASE_USER", "postgres")
	viper.SetDefault("DATABASE_PASSWORD", "postgres")
	viper.SetDefault("DATABASE_NAME", "e_metting")
	viper.SetDefault("DATABASE_MAX_OPEN_CONNECTIONS", 25)
	viper.SetDefault("DATABASE_MAX_IDLE_CONNECTIONS", 5)
	viper.SetDefault("JWT_SECRET_KEY", "your-secret-key")
	viper.SetDefault("JWT_TOKEN_DURATION", 24)

	viper.SetDefault("SMTP_HOST", "smtp.gmail.com")
	viper.SetDefault("SMTP_PORT", 587)
	viper.SetDefault("SMTP_USERNAME", "your-email@gmail.com")
	viper.SetDefault("SMTP_PASSWORD", "your-password")
	viper.SetDefault("SMTP_FROM_EMAIL", "your-email@gmail.com")
	viper.SetDefault("TEMPLATE_PATH", "templates/reset_password_email.html")
	viper.SetDefault("TEMPLATE_LOGO_URL", "https://example.com/logo.png")
	viper.SetDefault("SMTP_TIMEOUT_DURATION", 10)
	viper.SetDefault("SMTP_INSECURE_SKIP_VERIFY", false)
	viper.SetDefault("SMTP_USE_TLS", true)

	viper.SetDefault("CLOUDFLARE_R2_BUCKET_NAME", "")
	viper.SetDefault("CLOUDFLARE_R2_API_KEY", "")
	viper.SetDefault("CLOUDFLARE_R2_API_SECRET", "")
	viper.SetDefault("CLOUDFLARE_R2_TOKEN", "")
	viper.SetDefault("CLOUDFLARE_R2_ACCOUNT_ID", "")
	viper.SetDefault("CLOUDFLARE_R2_PUBLIC_URL", "")
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("No config file found at %s, using defaults", path)
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config

	config.AppEnv = viper.GetString("APP_ENV")
	config.AppPort = viper.GetString("APP_PORT")

	config.DBHost = viper.GetString("DATABASE_HOST")
	config.DBPort = viper.GetInt("DATABASE_PORT")
	config.DBUser = viper.GetString("DATABASE_USER")
	config.DBPassword = viper.GetString("DATABASE_PASSWORD")
	config.DBName = viper.GetString("DATABASE_NAME")
	config.DBMaxOpenConnections = viper.GetInt("DATABASE_MAX_OPEN_CONNECTIONS")
	config.DBMaxIdleConnections = viper.GetInt("DATABASE_MAX_IDLE_CONNECTIONS")

	config.JWT.SecretKey = viper.GetString("JWT_SECRET_KEY")
	config.JWT.TokenDuration = viper.GetInt("JWT_TOKEN_DURATION")

	config.SMTP.Host = viper.GetString("SMTP_HOST")
	config.SMTP.Port = viper.GetInt("SMTP_PORT")
	config.SMTP.Username = viper.GetString("SMTP_USERNAME")
	config.SMTP.Password = viper.GetString("SMTP_PASSWORD")
	config.SMTP.FromEmail = viper.GetString("SMTP_FROM_EMAIL")
	config.SMTP.TemplatePath = viper.GetString("TEMPLATE_PATH")
	config.SMTP.TemplateLogoURL = viper.GetString("TEMPLATE_LOGO_URL")
	config.SMTP.TimeoutDuration = viper.GetInt("SMTP_TIMEOUT_DURATION")
	config.SMTP.InsecureSkipVerify = viper.GetBool("SMTP_INSECURE_SKIP_VERIFY")
	config.SMTP.UseTLS = viper.GetBool("SMTP_USE_TLS")

	config.CloudflareR2BucketName = viper.GetString("CLOUDFLARE_R2_BUCKET_NAME")
	config.CloudflareR2APIKey = viper.GetString("CLOUDFLARE_R2_API_KEY")
	config.CloudflareR2APISecret = viper.GetString("CLOUDFLARE_R2_API_SECRET")
	config.CloudflareR2Token = viper.GetString("CLOUDFLARE_R2_TOKEN")
	config.CloudflareR2AccountID = viper.GetString("CLOUDFLARE_R2_ACCOUNT_ID")
	config.CloudflareR2PublicURL = viper.GetString("CLOUDFLARE_R2_PUBLIC_URL")

	port, err := strconv.Atoi(strings.TrimSpace(config.AppPort))
	if err != nil {
		return nil, fmt.Errorf("invalid APP_PORT: %v", err)
	}
	config.Server.Port = port

	return &config, nil
}

func NewConfig() *Config {
	rootDir := findRootDir()
	configPath := filepath.Join(rootDir, ".env")

	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	return config
}

func (c *Config) GetAppPort() string {
	return c.AppPort
}

func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}
