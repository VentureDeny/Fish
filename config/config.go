// config/config.go
package config

import (
	"log"
	"os"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	ServerAddr string
}

func LoadConfig() *Config {
	cfg := &Config{
		DBUser:     getEnv("DB_USER", "yourusername"),
		DBPassword: getEnv("DB_PASSWORD", "yourpassword"),
		DBName:     getEnv("DB_NAME", "websocket_db"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		ServerAddr: getEnv("SERVER_ADDR", ":8080"),
	}
	return cfg
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("环境变量 %s 未设置，使用默认值: %s", key, defaultValue)
		return defaultValue
	}
	return value
}
