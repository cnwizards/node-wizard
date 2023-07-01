package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func SetupLogger() {
	logLevel, _ := os.LookupEnv("LOG_LEVEL")
	switch logLevel {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
