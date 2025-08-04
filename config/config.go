package config 

import (
	"os"
)

type Config struct {
	AppName string
	Port    string
	Env   string
}

func Load() *Config {
	return &Config{
		AppName: os.Getenv("APP_NAME"),
		Port:    os.Getenv("PORT"),
		Env:     os.Getenv("ENV"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

