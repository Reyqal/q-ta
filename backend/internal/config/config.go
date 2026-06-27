package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config menyimpan seluruh konfigurasi aplikasi yang dibaca dari environment variables.
type Config struct {
	AppPort  string
	AppEnv   string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret      string
	JWTExpiryHours int

	TaxPercentage float64

	FrontendURL string

	MidtransServerKey    string
	MidtransIsProduction bool

	WhatsAppAPIKey string
	WhatsAppAPIURL string
}

// LoadConfig membaca file .env dan mengembalikan struct Config.
func LoadConfig() (*Config, error) {
	// Load .env file, ignore error if file doesn't exist (e.g., in production with real env vars)
	_ = godotenv.Load()

	jwtExpiry, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		jwtExpiry = 24
	}

	taxPercentage, err := strconv.ParseFloat(getEnv("TAX_PERCENTAGE", "0.5"), 64)
	if err != nil {
		taxPercentage = 0.5
	}

	midtransIsProd, _ := strconv.ParseBool(getEnv("MIDTRANS_IS_PRODUCTION", "false"))

	cfg := &Config{
		AppPort:  getEnv("APP_PORT", "3000"),
		AppEnv:   getEnv("APP_ENV", "development"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "qta_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret:      getEnv("JWT_SECRET", "super-secret-key-change-in-production"),
		JWTExpiryHours: jwtExpiry,

		TaxPercentage: taxPercentage,

		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),

		MidtransServerKey:    getEnv("MIDTRANS_SERVER_KEY", ""),
		MidtransIsProduction: midtransIsProd,

		WhatsAppAPIKey: getEnv("WHATSAPP_API_KEY", ""),
		WhatsAppAPIURL: getEnv("WHATSAPP_API_URL", ""),
	}

	return cfg, nil
}

// DSN mengembalikan PostgreSQL connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

// MigrationDSN mengembalikan connection string format URL untuk golang-migrate.
func (c *Config) MigrationDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode,
	)
}

// IsDevelopment mengecek apakah aplikasi berjalan di mode development.
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

// getEnv membaca environment variable dengan fallback ke default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}
