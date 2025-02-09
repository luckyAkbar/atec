// Package db contains the function to initialize and utilize connection to databases
package db

import (
	"os"
	"time"

	"github.com/luckyAkbar/atec/internal/config"
	"github.com/sirupsen/logrus"

	libdb "github.com/sweet-go/stdlib/db"
	"gorm.io/gorm"

	gormLogger "gorm.io/gorm/logger"
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

	for i := 0; i < retryLimit; i++ {
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

	if !success {
		logrus.Error("failed to initialize postgres connection")
		os.Exit(1)
	}

	switch config.LogLevel() {
	case "error":
		PostgresDB.Logger = PostgresDB.Logger.LogMode(gormLogger.Error)
	case "warn":
		PostgresDB.Logger = PostgresDB.Logger.LogMode(gormLogger.Warn)
	default:
		PostgresDB.Logger = PostgresDB.Logger.LogMode(gormLogger.Info)
	}

	logrus.Info("Connected to Postgres Database")
}

// TxController create a transaction controller
func TxController() *gorm.DB {
	return PostgresDB.Begin()
}
