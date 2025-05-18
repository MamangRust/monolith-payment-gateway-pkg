package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	Log *zap.Logger
}

//go:generate mockgen -source=logger.go -destination=mocks/logger.go
type LoggerInterface interface {
	Info(message string)
	Fatal(message string, fields ...zap.Field)
	Debug(message string, fields ...zap.Field)
	Error(message string, fields ...zap.Field)
}

func NewLogger() (LoggerInterface, error) {

	config := zap.NewDevelopmentConfig()

	encoderConfig := zap.NewDevelopmentEncoderConfig()

	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = ""
	config.EncoderConfig = encoderConfig

	zapLogger, err := config.Build(zap.AddCallerSkip(1))

	if err != nil {
		return nil, err
	}

	return &Logger{
		Log: zapLogger,
	}, nil
}

func (Logger *Logger) Info(message string) {

	Logger.Log.Info(message)
}

func (Logger *Logger) Fatal(message string, fields ...zap.Field) {
	Logger.Log.Fatal(message, fields...)
}

func (Logger *Logger) Debug(message string, fields ...zap.Field) {
	Logger.Log.Debug(message, fields...)

}

func (Logger *Logger) Error(message string, fields ...zap.Field) {
	Logger.Log.Error(message, fields...)
}
