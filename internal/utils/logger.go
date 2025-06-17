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
	log.Info(args...)
}

// LogInfof logs formatted informational messages
func LogInfof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// LogError logs error messages
func LogError(args ...interface{}) {
	log.Error(args...)
}

// LogErrorf logs formatted error messages
func LogErrorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// LogDebug logs debug messages
func LogDebug(args ...interface{}) {
	log.Debug(args...)
}

// LogDebugf logs formatted debug messages
func LogDebugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// LogWarn logs warning messages
func LogWarn(args ...interface{}) {
	log.Warn(args...)
}

// LogWarnf logs formatted warning messages
func LogWarnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}
