package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APP_PORT             string
	JWT_SECRET           string
	DB_HOST              string
	DB_PORT              string
	DB_USER              string
	DB_PASSWORD          string
	DB_NAME              string
	GOOGLE_CLIENT_ID     string
	GOOGLE_CLIENT_SECRET string
	GOOGLE_REDIRECT_URL  string
	OTP_APP_KEY          string
	OTP_APP_SECRET       string
	OTP_APP_URL_REQUEST  string
	OTP_APP_URL_VERIFY   string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		APP_PORT:             getEnv("APP_PORT", "3000"),
		JWT_SECRET:           getEnv("JWT_SECRET", "secret"),
		DB_HOST:              getEnv("DB_HOST", "localhost"),
		DB_PORT:              getEnv("DB_PORT", "5432"),
		DB_USER:              getEnv("DB_USER", "postgres"),
		DB_PASSWORD:          getEnv("DB_PASSWORD", ""),
		DB_NAME:              getEnv("DB_NAME", "google_auth_db"),
		GOOGLE_CLIENT_ID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GOOGLE_CLIENT_SECRET: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GOOGLE_REDIRECT_URL:  getEnv("GOOGLE_REDIRECT_URL", ""),
		OTP_APP_KEY:          getEnv("OTP_APP_KEY", ""),
		OTP_APP_SECRET:       getEnv("OTP_APP_SECRET", ""),
		OTP_APP_URL_REQUEST:  getEnv("OTP_APP_URL_REQUEST", "https://otp.thaibulksms.com/v2/otp/request"),
		OTP_APP_URL_VERIFY:   getEnv("OTP_APP_URL_VERIFY", "https://otp.thaibulksms.com/v2/otp/verify"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
