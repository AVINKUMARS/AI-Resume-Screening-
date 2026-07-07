package main

import (
	"log"

	"github.com/scalingwolf/ai-resume-screening/backend/internal/config"
	"github.com/scalingwolf/ai-resume-screening/backend/internal/database"
	"github.com/scalingwolf/ai-resume-screening/backend/internal/gemini"
	"github.com/scalingwolf/ai-resume-screening/backend/internal/handlers"
	"github.com/scalingwolf/ai-resume-screening/backend/internal/router"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("startup: %v", err)
	}

	ai := gemini.NewClient(cfg.GeminiAPIKey, cfg.GeminiModel)
	h := handlers.New(db, ai)
	r := router.New(h, cfg.AllowedOrigin)

	log.Printf("startup: listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("startup: server exited: %v", err)
	}
}
