package db

import (
	"fmt"
	"os"
	"strings"

	"Jougan-0/distributed-task-scheduler/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.DSN()

	gormLogLevelEnv := os.Getenv("GORM_LOG_LEVEL")
	if gormLogLevelEnv == "" {
		gormLogLevelEnv = "silent"
	}

	logLevel := parseGormLogLevel(gormLogLevelEnv)

	newLogger := logger.Default.LogMode(logLevel)

	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("gorm open error: %w", err)
	}
	return dbConn, nil
}

func parseGormLogLevel(levelStr string) logger.LogLevel {
	switch strings.ToLower(levelStr) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Silent
	}
}
