package database

import (
	"fmt"
	"log"

	"github.com/scalingwolf/ai-resume-screening/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect opens the PostgreSQL connection and runs GORM auto-migration for all
// domain models. It returns a live *gorm.DB or an error the caller should treat
// as fatal.
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("database: open: %w", err)
	}

	// pgcrypto/uuid generation is done application-side (BeforeCreate), so we
	// only need the schema itself.
	if err := db.AutoMigrate(&models.Candidate{}, &models.Job{}, &models.Match{}); err != nil {
		return nil, fmt.Errorf("database: automigrate: %w", err)
	}

	log.Println("database: connected and migrated")
	return db, nil
}
