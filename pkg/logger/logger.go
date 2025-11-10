package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Level       string
	Environment string
	Output      io.Writer
}

func Init(cfg Config) {
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}

	level := parseLogLevel(cfg.Level)
	zerolog.SetGlobalLevel(level)

	if cfg.Environment == "development" {
		output := zerolog.ConsoleWriter{
			Out:        cfg.Output,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		}
		log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()
		return
	}

	log.Logger = zerolog.New(cfg.Output).With().Timestamp().Caller().Logger()
}

func parseLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}
