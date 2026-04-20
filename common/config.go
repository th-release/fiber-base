package common

import (
	"os"
	"strconv"
	"strings"
	"time"

	"cth.release/common/utils"
)

type Config struct {
	Port string `json:"PORT"`

	DatabaseHost string `json:"DATABASE_HOST"`
	DatabasePort string `json:"DATABASE_PORT"`
	DatabaseUser string `json:"DATABASE_USER"`
	DatabasePass string `json:"DATABASE_PASS"`
	DatabaseName string `json:"DATABASE_NAME"`

	JWTSecret        string        `json:"JWT_SECRET"`
	JWTIssuer        string        `json:"JWT_ISSUER"`
	JWTExpiry        time.Duration `json:"JWT_EXPIRY"`
	StaticDir        string        `json:"STATIC_DIR"`
	SPAIndexFile     string        `json:"SPA_INDEX_FILE"`
	SMTPHost         string        `json:"SMTP_HOST"`
	SMTPPort         int           `json:"SMTP_PORT"`
	SMTPUser         string        `json:"SMTP_USER"`
	SMTPPass         string        `json:"SMTP_PASS"`
	SMTPFrom         string        `json:"SMTP_FROM"`
	SMTPAuthIdentity string        `json:"SMTP_AUTH_IDENTITY"`
	SMTPStartTLS     bool          `json:"SMTP_STARTTLS"`
	SMTPInsecureSkip bool          `json:"SMTP_INSECURE_SKIP_VERIFY"`
	FCMProjectID     string        `json:"FCM_PROJECT_ID"`
	FCMClientEmail   string        `json:"FCM_CLIENT_EMAIL"`
	FCMPrivateKey    string        `json:"FCM_PRIVATE_KEY"`
	FCMTokenURI      string        `json:"FCM_TOKEN_URI"`
	FCMScope         string        `json:"FCM_SCOPE"`
	FCMBaseURL       string        `json:"FCM_BASE_URL"`
	FCMCredsJSON     string        `json:"FCM_CREDENTIALS_JSON"`
	FCMCredsFile     string        `json:"FCM_CREDENTIALS_FILE"`
}

func GetConfig() *Config {
	return &Config{
		Port:             utils.ThreeTermString(len(os.Getenv("PORT")) > 0, os.Getenv("PORT"), "9000"),
		DatabaseHost:     utils.ThreeTermString(len(os.Getenv("DATABASE_HOST")) > 0, os.Getenv("DATABASE_HOST"), "host.docker.internal"),
		DatabasePort:     utils.ThreeTermString(len(os.Getenv("DATABASE_PORT")) > 0, os.Getenv("DATABASE_PORT"), "5432"),
		DatabaseUser:     utils.ThreeTermString(len(os.Getenv("DATABASE_USER")) > 0, os.Getenv("DATABASE_USER"), "root"),
		DatabasePass:     utils.ThreeTermString(len(os.Getenv("DATABASE_PASS")) > 0, os.Getenv("DATABASE_PASS"), "password"),
		DatabaseName:     utils.ThreeTermString(len(os.Getenv("DATABASE_NAME")) > 0, os.Getenv("DATABASE_NAME"), "database"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		JWTIssuer:        getEnvString("JWT_ISSUER", "fiber-base"),
		JWTExpiry:        getEnvDuration("JWT_EXPIRY", 8*time.Hour),
		StaticDir:        getEnvString("STATIC_DIR", "./public"),
		SPAIndexFile:     getEnvString("SPA_INDEX_FILE", "index.html"),
		SMTPHost:         os.Getenv("SMTP_HOST"),
		SMTPPort:         getEnvInt("SMTP_PORT", 587),
		SMTPUser:         os.Getenv("SMTP_USER"),
		SMTPPass:         os.Getenv("SMTP_PASS"),
		SMTPFrom:         os.Getenv("SMTP_FROM"),
		SMTPAuthIdentity: os.Getenv("SMTP_AUTH_IDENTITY"),
		SMTPStartTLS:     getEnvBool("SMTP_STARTTLS", true),
		SMTPInsecureSkip: getEnvBool("SMTP_INSECURE_SKIP_VERIFY", false),
		FCMProjectID:     os.Getenv("FCM_PROJECT_ID"),
		FCMClientEmail:   os.Getenv("FCM_CLIENT_EMAIL"),
		FCMPrivateKey:    os.Getenv("FCM_PRIVATE_KEY"),
		FCMTokenURI:      getEnvString("FCM_TOKEN_URI", "https://oauth2.googleapis.com/token"),
		FCMScope:         getEnvString("FCM_SCOPE", "https://www.googleapis.com/auth/firebase.messaging"),
		FCMBaseURL:       getEnvString("FCM_BASE_URL", "https://fcm.googleapis.com"),
		FCMCredsJSON:     os.Getenv("FCM_CREDENTIALS_JSON"),
		FCMCredsFile:     os.Getenv("FCM_CREDENTIALS_FILE"),
	}
}

func getEnvString(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}
