package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger wraps zerolog.Logger
type Logger struct {
	zerolog.Logger
}

// New creates a new Logger instance
func New(level string, isDevelopment bool) *Logger {
	var output io.Writer = os.Stdout

	// Use console writer for development
	if isDevelopment {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	logLevel := parseLevel(level)

	logger := zerolog.New(output).
		Level(logLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	return &Logger{Logger: logger}
}

// parseLevel parses log level string to zerolog.Level
func parseLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
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

// WithField adds a field to the logger context
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{Logger: l.Logger.With().Interface(key, value).Logger()}
}

// WithFields adds multiple fields to the logger context
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	ctx := l.Logger.With()
	for key, value := range fields {
		ctx = ctx.Interface(key, value)
	}
	return &Logger{Logger: ctx.Logger()}
}

// WithError adds an error to the logger context
func (l *Logger) WithError(err error) *Logger {
	return &Logger{Logger: l.Logger.With().Err(err).Logger()}
}
