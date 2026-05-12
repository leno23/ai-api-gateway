package store

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func OpenPostgres(dsn string, log *zap.Logger) (*gorm.DB, error) {
	gl := gormlogger.Default.LogMode(gormlogger.Warn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gl,
	})
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)
	return db, nil
}
