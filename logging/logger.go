package logging

import (
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
)

func init() {
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.JSONFormatter{})
}

// Trace loglevel - trace
func Trace(v interface{}) {
	log.Trace(v)
}

// Debug loglevel - debug
func Debug(v interface{}) {
	log.Debug(v)
}

// Info loglevel - info
func Info(v interface{}) {
	log.Info(v)
}

// Warn loglevel - warning
func Warn(v interface{}) {
	log.Warn(v)
}

// Error loglevel - error
func Error(v interface{}) {
	log.Error(v)
}

// Fatal loglevel - fatal
func Fatal(v interface{}) {
	log.Fatal(v)
}

// Panic loglevel - panic
func Panic(v interface{}) {
	log.Panic(v)
}
