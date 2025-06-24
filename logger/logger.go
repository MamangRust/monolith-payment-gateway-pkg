package logger

import (
	"fmt"
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

func NewLogger(service string) (LoggerInterface, error) {
	var setupErr error

	once.Do(func() {
		env := os.Getenv("APP_ENV")

		if env == "" {
			env = "development"
		}

		logDir := "./logs"

		if env == "docker" || env == "production" {
			logDir = "/var/log/app"
		}

		if err := os.MkdirAll(logDir, 0755); err != nil {
			setupErr = fmt.Errorf("failed to create log directory: %w", err)
			return
		}

		logPath := filepath.Join(logDir, fmt.Sprintf("%s.log", service))

		logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			setupErr = fmt.Errorf("failed to open log file: %w", err)
			return
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

		core := zapcore.NewTee(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig),
				zapcore.AddSync(logFile),
				zapcore.DebugLevel,
			),
			zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig),
				zapcore.AddSync(os.Stdout),
				zapcore.DebugLevel,
			),
		)

		logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		instance = &Logger{Log: logger}
	})

	return instance, setupErr
}

func (l *Logger) Info(message string, fields ...zap.Field) {
	l.Log.Info(message, fields...)
}

func (l *Logger) Fatal(message string, fields ...zap.Field) {
	l.Log.Fatal(message, fields...)
}

func (l *Logger) Debug(message string, fields ...zap.Field) {
	l.Log.Debug(message, fields...)
}

func (l *Logger) Error(message string, fields ...zap.Field) {
	l.Log.Error(message, fields...)
}
