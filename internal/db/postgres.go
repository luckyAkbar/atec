// Package db contains the function to initialize and utilize connection to databases
package db

import (
	"log"
	"os"
	"time"

	"github.com/luckyAkbar/atec/internal/config"
	"github.com/sirupsen/logrus"

	libdb "github.com/sweet-go/stdlib/db"
	"gorm.io/gorm"

	"gorm.io/gorm/logger"
)

var (
	// PostgresDB is the global postgres connection
	PostgresDB *gorm.DB
)

// InitializePostgresConn initializes the postgres connection
func InitializePostgresConn() {
	retryLimit := 10
	retryDelaySec := 3

	success := false

	for range retryLimit {
		conn, err := libdb.NewPostgresDB(config.PostgresDSN())
		if err != nil {
			logrus.WithError(err).Error("failed to initialize postgres connection")

			time.Sleep(time.Duration(retryDelaySec) * time.Second)
		} else {
			success = true
			PostgresDB = conn

			break
		}
	}

	sqlDB, err := PostgresDB.DB()
	if err != nil {
		logrus.WithError(err).Error("failed to get postgres connection")
		os.Exit(1)
	}

	sqlDB.SetMaxIdleConns(config.DBMaxIdleConn())
	sqlDB.SetMaxOpenConns(config.DBMaxOpenConn())
	sqlDB.SetConnMaxLifetime(config.DBConnMaxLifetime())

	if !success {
		logrus.Error("failed to initialize postgres connection")
		os.Exit(1)
	}

	var logLevel logger.LogLevel

	switch config.LogLevel() {
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	default:
		logLevel = logger.Info
	}

	const defaultSlowThreshold = 500 * time.Millisecond

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             defaultSlowThreshold, // Slow SQL threshold
			LogLevel:                  logLevel,             // Log level
			IgnoreRecordNotFoundError: true,                 // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,                 // Don't include params in the SQL log
			Colorful:                  true,                 // Disable color
		},
	)

	PostgresDB.Logger = gormLogger

	logrus.Info("Connected to Postgres Database")
}

// TxController create a transaction controller
func TxController() *gorm.DB {
	return PostgresDB.Begin()
}
