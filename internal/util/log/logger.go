package log

import (
	"errors"
	"strings"
	"sync"

	config "github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the wrapper around zap's Sugared Logger.
type Logger struct {
	*zap.SugaredLogger
}

// logger is a centralized instance for logging. This is because many parts of
// adeia-api, that are not part of the methods of APIServer, need access to the
// logger.
var (
	logger  *Logger
	initLog = new(sync.Once)
)

// levels is a map of supported log levels.
var levels = map[string]zapcore.Level{
	"debug": zap.DebugLevel,
	"info":  zap.InfoLevel,
	"warn":  zap.WarnLevel,
	"error": zap.ErrorLevel,
	"panic": zap.PanicLevel,
	"fatal": zap.FatalLevel,
}

// Init initializes a new logger instance based on config.
func Init() error {
	err := errors.New("logger already initialized")

	initLog.Do(func() {
		err = nil

		// parse log level
		level, e := parseLevel(config.GetString("logger.level"))
		if e != nil {
			err = e
			return
		}

		// TODO: switch to custom config
		cfg := zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(level)
		// TODO: setup log rotation
		cfg.OutputPaths = config.GetStringSlice("logger.paths")

		// build logger from config
		l, e := cfg.Build(zap.AddCallerSkip(1))
		if e != nil {
			err = e
			return
		}
		logger = &Logger{l.Sugar()}
	})

	return err
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

// Set sets the logger.
func Set(l *zap.SugaredLogger) {
	logger.SugaredLogger = l
}

// ==========
// Wrapper methods
// ==========

// We use wrapper methods so that we need not have verbose func calls like logger.log.info(...).
// Doing this, we can do so with just logger.info(...).

// Sync wraps SugaredLogger's Sync.
func Sync() error {
	return logger.Sync()
}

// Debugf wraps SugaredLogger's Debugf.
func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

// Debug wraps SugaredLogger's Debug.
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Infof wraps SugaredLogger's Infof.
func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

// Info wraps SugaredLogger's Info.
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Warnf wraps SugaredLogger's Warnf.
func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

// Warn wraps SugaredLogger's Warn.
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Errorf wraps SugaredLogger's Errorf.
func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

// Error wraps SugaredLogger's Error.
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Panicf wraps SugaredLogger's Panicf.
func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args...)
}

// Panic wraps SugaredLogger's Panic.
func Panic(args ...interface{}) {
	logger.Panic(args...)
}

// Fatalf wraps SugaredLogger's Fatalf.
func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}

// Fatal wraps SugaredLogger's Fatal.
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}