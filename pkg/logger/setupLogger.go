package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func SetupLogger() {
	logLevel, _ := os.LookupEnv("LOG_LEVEL")

	logrus := log.New()
	logrus.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	switch logLevel {
	case "trace":
		logrus.SetLevel(log.TraceLevel)
	case "debug":
		logrus.SetLevel(log.DebugLevel)
	case "warn":
		logrus.SetLevel(log.WarnLevel)
	case "error":
		logrus.SetLevel(log.ErrorLevel)
	default:
		logrus.SetLevel(log.InfoLevel)
	}
}
