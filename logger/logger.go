package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//go:generate mockgen -source=logger.go -destination=mocks/logger.go
type LoggerInterface interface {
	Info(message string, fields ...zap.Field)
	Fatal(message string, fields ...zap.Field)
	Debug(message string, fields ...zap.Field)
	Error(message string, fields ...zap.Field)
}

type Logger struct {
	Log *zap.Logger
}

var once sync.Once
var instance LoggerInterface

// NewLogger returns a singleton instance of LoggerInterface with a logger that writes
// JSON-formatted log messages to both stdout and a log file. The log file name is
// determined by the service parameter. If the service parameter is empty, the log file
// name will be "app.log". The logger will be configured to write logs in the
// following locations:
//
//   - stdout, if the APP_ENV environment variable is not set
//   - /var/log/app/<service>.log, if the APP_ENV environment variable is set to
//     "docker", "production", or "kubernetes"
//   - ./logs/<service>.log, otherwise
//
// The logger will fallback to stdout only if it fails to create the log directory or
// open the log file.
//
// The logger will use the DebugLevel by default.
func NewLogger(service string) (LoggerInterface, error) {
	var setupErr error

	once.Do(func() {
		env := os.Getenv("APP_ENV")
		if env == "" {
			env = "development"
		}

		logDir := "./logs"
		if env == "docker" || env == "production" || env == "kubernetes" {
			logDir = "/var/log/app"
		}

		if err := os.MkdirAll(logDir, 0755); err != nil {
			setupErr = fmt.Errorf("failed to create log directory '%s': %w", logDir, err)
			log.Println("[WARN] Fallback to stdout only:", setupErr)
		}

		logPath := filepath.Join(logDir, fmt.Sprintf("%s.log", service))

		var logFile *os.File
		if setupErr == nil {
			var err error
			logFile, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				setupErr = fmt.Errorf("failed to open log file '%s': %w", logPath, err)
				log.Println("[WARN] Fallback to stdout only:", setupErr)
			}
		}

		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		cores := []zapcore.Core{
			zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig),
				zapcore.AddSync(os.Stdout),
				zapcore.DebugLevel,
			),
		}

		if logFile != nil {
			cores = append(cores, zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig),
				zapcore.AddSync(logFile),
				zapcore.DebugLevel,
			))
		}

		logger := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1))
		instance = &Logger{Log: logger}
	})

	return instance, setupErr
}

// Info logs a message at the Info level. It accepts a message string
// and optional zap fields to include additional context in the log entry.
func (l *Logger) Info(message string, fields ...zap.Field) {
	l.Log.Info(message, fields...)
}

// Fatal logs a message at the Fatal level. It accepts a message string
// and optional zap fields to include additional context in the log entry.
// The function will log the message and then terminate the application.
func (l *Logger) Fatal(message string, fields ...zap.Field) {
	l.Log.Fatal(message, fields...)
}

// Debug logs a message at the Debug level. It accepts a message string
// and optional zap fields to include additional context in the log entry.
func (l *Logger) Debug(message string, fields ...zap.Field) {
	l.Log.Debug(message, fields...)
}

// Error logs a message at the Error level. It accepts a message string
// and optional zap fields to include additional context in the log entry.
func (l *Logger) Error(message string, fields ...zap.Field) {
	l.Log.Error(message, fields...)
}
