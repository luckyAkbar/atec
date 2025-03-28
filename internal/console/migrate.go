// Package console holds all the service functionality accessible from command prompt
package console

import (
	"strconv"

	"github.com/luckyAkbar/atec/internal/config"
	"github.com/luckyAkbar/atec/internal/db"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/sweet-go/stdlib/helper"
)

var migrateCMD = &cobra.Command{
	Use:  "migrate",
	Long: "run migrations",
	Run:  migrateFn,
}

//nolint:gochecknoinits
func init() {
	migrateCMD.PersistentFlags().Int("step", 0, "maximum migration steps")
	migrateCMD.PersistentFlags().String("direction", "up", "migration direction")
	rootCMD.AddCommand(migrateCMD)
}

func migrateFn(cmd *cobra.Command, _ []string) {
	direction := cmd.Flag("direction").Value.String()
	stepStr := cmd.Flag("step").Value.String()
	step, err := strconv.Atoi(stepStr)

	if err != nil {
		logrus.WithField("stepStr", stepStr).Fatal("Failed to parse step to int: ", err)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "./db/migration",
	}

	migrate.SetTable("schema_migrations")

	db.InitializePostgresConn()

	pgdb, err := db.PostgresDB.DB()
	if err != nil {
		logrus.WithField("DatabaseDSN", config.PostgresDSN()).Fatal("failed to run migration")
	}

	var n int
	if direction == "down" {
		n, err = migrate.ExecMax(pgdb, "postgres", migrations, migrate.Down, step)
	} else {
		n, err = migrate.ExecMax(pgdb, "postgres", migrations, migrate.Up, step)
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"migrations": helper.Dump(migrations),
			"direction":  direction}).
			Fatal("Failed to migrate database: ", err)
	}

	logrus.Infof("Applied %d migrations!\n", n)
}
