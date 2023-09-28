package database

import (
	"github.com/alhazenlabs/procure-hub/backend/internal/logger"
	"github.com/alhazenlabs/procure-hub/backend/internal/models"
	"gorm.io/gorm"
)

func RunMigrations() {
	db := GetDB()
	if !checkEnumExists(db, "user_type") {
		if err := db.Exec(`CREATE TYPE user_type AS ENUM ('primary_owner', 'client')`).Error; err != nil {
			logger.Fatal(err.Error())
		}
	}
	db.AutoMigrate(&models.User{})
}

func checkEnumExists(db *gorm.DB, name string) bool {
	var exists bool
	if err := db.Raw(`SELECT 1 FROM pg_type WHERE typname = ?`, name).Scan(&exists).Error; err != nil {
		logger.Fatal(err.Error())
	}
	return exists
}
