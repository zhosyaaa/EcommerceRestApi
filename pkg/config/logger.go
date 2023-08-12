package config

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func init() {
	// Read log level
	logLevel := GetEnvVar("LOG_LEVEL")

	// Replace with switch or add more levels as per need
	if logLevel == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Set time format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Set Global fields
	host, err := os.Hostname()
	if err != nil {
		log.Logger = log.With().Str("host", "unknown").Logger()
	} else {
		log.Logger = log.With().Str("host", host).Logger()
	}

	log.Logger = log.With().Str("service", "gin-boilerplate").Logger()

	log.Logger = log.With().Caller().Logger()
}
