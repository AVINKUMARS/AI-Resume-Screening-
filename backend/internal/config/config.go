package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration for the service, sourced from the
// environment (see .env.example).
type Config struct {
	Port          string
	DatabaseURL   string
	GeminiAPIKey  string
	GeminiModel   string
	AllowedOrigin string
}

// Load reads configuration from a local .env file (if present) and the process
// environment. Missing critical values are logged so misconfiguration surfaces
// early instead of at the first request.
func Load() *Config {
	// .env is optional in deployed environments where real env vars are set.
	if err := godotenv.Load(); err != nil {
		log.Println("config: no .env file found, relying on process environment")
	}

	cfg := &Config{
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/resume_screening?sslmode=disable"),
		GeminiAPIKey:  os.Getenv("GEMINI_API_KEY"),
		GeminiModel:   getEnv("GEMINI_MODEL", "gemini-2.5-flash"),
		AllowedOrigin: getEnv("ALLOWED_ORIGIN", "http://localhost:3000"),
	}

	if cfg.GeminiAPIKey == "" {
		log.Println("config: WARNING GEMINI_API_KEY is empty — AI parsing and matching will fail")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
