package common

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server   ServerConfig
	Web      WebConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
}

type WebConfig struct {
	StaticDir    string
	SPAIndexFile string
}

type DatabaseConfig struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

type JWTConfig struct {
	Secret string
	Issuer string
	Expiry time.Duration
}

func GetConfig() *Config {
	return &Config{
		Server:   loadServerConfig(),
		Web:      loadWebConfig(),
		Database: loadDatabaseConfig(),
		JWT:      loadJWTConfig(),
	}
}

func loadServerConfig() ServerConfig {
	return ServerConfig{
		Port: getEnvString("PORT", "9000"),
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
		Host: getEnvString("DATABASE_HOST", "host.docker.internal"),
		Port: getEnvString("DATABASE_PORT", "5432"),
		User: getEnvString("DATABASE_USER", "root"),
		Pass: getEnvString("DATABASE_PASS", "password"),
		Name: getEnvString("DATABASE_NAME", "database"),
	}
}

func loadJWTConfig() JWTConfig {
	return JWTConfig{
		Secret: os.Getenv("JWT_SECRET"),
		Issuer: getEnvString("JWT_ISSUER", "fiber-base"),
		Expiry: getEnvDuration("JWT_EXPIRY", 8*time.Hour),
	}
}

func getEnvString(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	if parsed, err := time.ParseDuration(value); err == nil {
		return parsed
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	if parsed, err := strconv.Atoi(value); err == nil {
		return parsed
	}
	return fallback
}
