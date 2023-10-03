package database

import (
	"context"

	"github.com/alhazenlabs/procure-hub/backend/internal/logger"
	"gorm.io/gorm"
)

func BeginTransaction(context context.Context) (*gorm.DB, func() error) {
	tx := GetDBWithCtx(context).Begin()
	return tx, func() error {
		if r := recover(); r != nil {
			tx.Rollback()
			logger.Error(r.(error).Error())
			return r.(error)
		}
		if err := tx.Error; err != nil {
			tx.Rollback()
			logger.Error(err.Error())
			return err
		}
		return tx.Commit().Error
	}
}

/*
// DeferFunc is a generic function to handle defer operations.
func DeferFunc(f func() error) {
	if err := f(); err != nil {
		logger.Fatal(err.Error())
	}
}
*/
