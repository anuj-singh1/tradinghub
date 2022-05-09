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

func GetLogger() *logrus.Logger {
	return logger
}
