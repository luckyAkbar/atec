// Package console contains all functionality accessible by cmd / terminal
package console

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/sweet-go/stdlib/cmd"
)

// RootCMD is the root command of the application
var rootCMD = cmd.CobraInitializer()

// Execute will be the entry point of all the registered command
func Execute() {
	if err := rootCMD.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
