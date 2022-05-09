package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var logger *logrus.Logger
var loggerSetup sync.Once

func SetupLogger() {
	loggerSetup.Do(func() {
		logger = logrus.New()
	})

	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)
	var err error
	logger.Level, err = logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logger.Level = logrus.InfoLevel
	}
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func GetLogger() *logrus.Logger {
	return logger
}