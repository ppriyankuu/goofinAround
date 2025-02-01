package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	// Set the log output to stdout
	log.Out = os.Stdout

	// Set the log level to Info by default
	log.Level = logrus.InfoLevel

	// Set the log format to JSON
	log.Formatter = &logrus.JSONFormatter{}
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return log.WithFields(fields)
}
