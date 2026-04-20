package common

import (
	"os"

	"cth.release/common/utils"
)

type Config struct {
	Port string `json:"PORT"`

	DatabaseHost string `json:"DATABASE_HOST"`
	DatabasePort string `json:"DATABASE_PORT"`
	DatabaseUser string `json:"DATABASE_USER"`
	DatabasePass string `json:"DATABASE_PASS"`
	DatabaseName string `json:"DATABASE_NAME"`
}

func GetConfig() *Config {
	return &Config{
		Port:         utils.ThreeTermString(len(os.Getenv("PORT")) > 0, os.Getenv("PORT"), "9000"),
		DatabaseHost: utils.ThreeTermString(len(os.Getenv("DATABASE_HOST")) > 0, os.Getenv("DATABASE_HOST"), "host.docker.internal"),
		DatabasePort: utils.ThreeTermString(len(os.Getenv("DATABASE_PORT")) > 0, os.Getenv("DATABASE_PORT"), "5432"),
		DatabaseUser: utils.ThreeTermString(len(os.Getenv("DATABASE_USER")) > 0, os.Getenv("DATABASE_USER"), "root"),
		DatabasePass: utils.ThreeTermString(len(os.Getenv("DATABASE_PASS")) > 0, os.Getenv("DATABASE_PASS"), "password"),
		DatabaseName: utils.ThreeTermString(len(os.Getenv("DATABASE_NAME")) > 0, os.Getenv("DATABASE_NAME"), "database"),
	}
}
