package logger

import (
	"adeia-api/internal/config"
	"errors"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logger is a centralized instance for logging. This is because many parts of
// adeia-api, that are not part of the methods of APIServer, need access to the
// logger.
var logger *zap.SugaredLogger

// levels is a map of supported log levels.
var levels = map[string]zapcore.Level{
	"debug": zap.DebugLevel,
	"info":  zap.InfoLevel,
	"warn":  zap.WarnLevel,
	"error": zap.ErrorLevel,
	"panic": zap.PanicLevel,
	"fatal": zap.FatalLevel,
}

// Init initializes a new logger instance based on passed-in config.
func Init(conf *config.LoggerConfig) error {
	// parse log level
	level, err := parseLevel(conf.Level)
	if err != nil {
		return err
	}

	// TODO: switch to custom config
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(level)
	// TODO: set log output to a file
	cfg.OutputPaths = []string{"stdout"}

	// build logger from config
	l, err := cfg.Build()
	if err != nil {
		return err
	}
	logger = l.Sugar()

	return nil
}

// parseLevel returns the appropriate zapcore.Level for the passed-in string.
func parseLevel(s string) (zapcore.Level, error) {
	l, ok := levels[strings.ToLower(s)]
	if !ok {
		return 0, errors.New("specified log level not one of " +
			"['debug', 'info', 'warn', 'error', 'panic', 'fatal']")
	}
	return l, nil
}

// Get returns the logger instance
func Get() *zap.SugaredLogger {
	return logger
}
