// Package db contains the function to initialize and utilize connection to databases
package db

import (
	"os"

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
	conn, err := libdb.NewPostgresDB(config.PostgresDSN())
	if err != nil {
		logrus.Error("failed to initialize postgres connection")
		os.Exit(1)
	}

	PostgresDB = conn

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
