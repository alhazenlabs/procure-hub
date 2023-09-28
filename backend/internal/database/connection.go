package database

import (
	"context"
	"fmt"

	"sync"

	"github.com/alhazenlabs/procure-hub/backend/internal/config"
	"github.com/alhazenlabs/procure-hub/backend/internal/logger"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	DB   *gorm.DB
	once sync.Once
)

func getPostgresDsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", config.PG_HOST,
		config.PG_USER, config.PG_PASSWORD, config.PG_DB, config.PG_PORT, config.PG_SSLMODE)
}

func InitializeDB() {
	db, err := gorm.Open(postgres.Open(getPostgresDsn()), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info)})
	if err != nil {
		logger.Fatal("failed to connect database")
	}
	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		logger.Fatal("unable to add telemetery to gorm")
	}
	DB = db
}

func GetDB() *gorm.DB {
	if DB == nil {
		logger.Info("database not connected, connecting")
		once.Do(InitializeDB)
	}
	return DB
}

func GetDBWithCtx(context context.Context) *gorm.DB {
	return GetDB().WithContext(context)
}
