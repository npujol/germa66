package utils

import (
	log "github.com/sirupsen/logrus"
)

// SetLogger sets up the logger with JSON formatter
func SetLogger() {
	log.SetFormatter(&log.JSONFormatter{})
}

// LogFatalf logs fatal messages and exits
func LogFatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// LogInfo logs informational messages
func LogInfo(args ...interface{}) {
	log.Println(args...)
}
