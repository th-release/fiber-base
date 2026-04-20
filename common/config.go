package common

import (
	"os"
	"strconv"
	"strings"
	"time"

	"cth.release/common/utils"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Web      WebConfig      `json:"web"`
	Database DatabaseConfig `json:"database"`
	JWT      JWTConfig      `json:"jwt"`
	SMTP     SMTPConfig     `json:"smtp"`
	FCM      FCMConfig      `json:"fcm"`
	XSS      XSSConfig      `json:"xss"`
}

type ServerConfig struct {
	Port string `json:"PORT"`
}

type WebConfig struct {
	StaticDir    string `json:"STATIC_DIR"`
	SPAIndexFile string `json:"SPA_INDEX_FILE"`
}

type DatabaseConfig struct {
	Host string `json:"DATABASE_HOST"`
	Port string `json:"DATABASE_PORT"`
	User string `json:"DATABASE_USER"`
	Pass string `json:"DATABASE_PASS"`
	Name string `json:"DATABASE_NAME"`
}

type JWTConfig struct {
	Secret string        `json:"JWT_SECRET"`
	Issuer string        `json:"JWT_ISSUER"`
	Expiry time.Duration `json:"JWT_EXPIRY"`
}

type SMTPConfig struct {
	Host               string `json:"SMTP_HOST"`
	Port               int    `json:"SMTP_PORT"`
	User               string `json:"SMTP_USER"`
	Pass               string `json:"SMTP_PASS"`
	From               string `json:"SMTP_FROM"`
	AuthIdentity       string `json:"SMTP_AUTH_IDENTITY"`
	StartTLS           bool   `json:"SMTP_STARTTLS"`
	InsecureSkipVerify bool   `json:"SMTP_INSECURE_SKIP_VERIFY"`
}

type FCMConfig struct {
	ProjectID       string `json:"FCM_PROJECT_ID"`
	ClientEmail     string `json:"FCM_CLIENT_EMAIL"`
	PrivateKey      string `json:"FCM_PRIVATE_KEY"`
	TokenURI        string `json:"FCM_TOKEN_URI"`
	Scope           string `json:"FCM_SCOPE"`
	BaseURL         string `json:"FCM_BASE_URL"`
	CredentialsJSON string `json:"FCM_CREDENTIALS_JSON"`
	CredentialsFile string `json:"FCM_CREDENTIALS_FILE"`
}

type XSSConfig struct {
	AllowElements []string `json:"XSS_ALLOW_ELEMENTS"`
}

func GetConfig() *Config {
	return &Config{
		Server:   loadServerConfig(),
		Web:      loadWebConfig(),
		Database: loadDatabaseConfig(),
		JWT:      loadJWTConfig(),
		SMTP:     loadSMTPConfig(),
		FCM:      loadFCMConfig(),
		XSS:      loadXSSConfig(),
	}
}

func loadServerConfig() ServerConfig {
	return ServerConfig{
		Port: utils.ThreeTermString(len(os.Getenv("PORT")) > 0, os.Getenv("PORT"), "9000"),
	}
}

func loadWebConfig() WebConfig {
	return WebConfig{
		StaticDir:    getEnvString("STATIC_DIR", "./public"),
		SPAIndexFile: getEnvString("SPA_INDEX_FILE", "index.html"),
	}
}

func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host: utils.ThreeTermString(len(os.Getenv("DATABASE_HOST")) > 0, os.Getenv("DATABASE_HOST"), "host.docker.internal"),
		Port: utils.ThreeTermString(len(os.Getenv("DATABASE_PORT")) > 0, os.Getenv("DATABASE_PORT"), "5432"),
		User: utils.ThreeTermString(len(os.Getenv("DATABASE_USER")) > 0, os.Getenv("DATABASE_USER"), "root"),
		Pass: utils.ThreeTermString(len(os.Getenv("DATABASE_PASS")) > 0, os.Getenv("DATABASE_PASS"), "password"),
		Name: utils.ThreeTermString(len(os.Getenv("DATABASE_NAME")) > 0, os.Getenv("DATABASE_NAME"), "database"),
	}
}

func loadJWTConfig() JWTConfig {
	return JWTConfig{
		Secret: os.Getenv("JWT_SECRET"),
		Issuer: getEnvString("JWT_ISSUER", "fiber-base"),
		Expiry: getEnvDuration("JWT_EXPIRY", 8*time.Hour),
	}
}

func loadSMTPConfig() SMTPConfig {
	return SMTPConfig{
		Host:               os.Getenv("SMTP_HOST"),
		Port:               getEnvInt("SMTP_PORT", 587),
		User:               os.Getenv("SMTP_USER"),
		Pass:               os.Getenv("SMTP_PASS"),
		From:               os.Getenv("SMTP_FROM"),
		AuthIdentity:       os.Getenv("SMTP_AUTH_IDENTITY"),
		StartTLS:           getEnvBool("SMTP_STARTTLS", true),
		InsecureSkipVerify: getEnvBool("SMTP_INSECURE_SKIP_VERIFY", false),
	}
}

func loadFCMConfig() FCMConfig {
	return FCMConfig{
		ProjectID:       os.Getenv("FCM_PROJECT_ID"),
		ClientEmail:     os.Getenv("FCM_CLIENT_EMAIL"),
		PrivateKey:      os.Getenv("FCM_PRIVATE_KEY"),
		TokenURI:        getEnvString("FCM_TOKEN_URI", "https://oauth2.googleapis.com/token"),
		Scope:           getEnvString("FCM_SCOPE", "https://www.googleapis.com/auth/firebase.messaging"),
		BaseURL:         getEnvString("FCM_BASE_URL", "https://fcm.googleapis.com"),
		CredentialsJSON: os.Getenv("FCM_CREDENTIALS_JSON"),
		CredentialsFile: os.Getenv("FCM_CREDENTIALS_FILE"),
	}
}

func loadXSSConfig() XSSConfig {
	return XSSConfig{
		AllowElements: getEnvCSV("XSS_ALLOW_ELEMENTS", nil),
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

func getEnvCSV(key string, fallback []string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
	}

	if len(result) == 0 {
		return fallback
	}
	return result
}
