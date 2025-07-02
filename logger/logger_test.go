package logger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// TestNewLogger tests the NewLogger function.
//
// The test does the following:
//
// 1. Unsets the APP_ENV environment variable.
// 2. Creates a new logger instance using the NewLogger function.
// 3. Asserts that the logger instance is not nil and that there is no error.
// 4. Logs a message using the Info, Debug, and Error methods.
func TestNewLogger(t *testing.T) {
	os.Unsetenv("APP_ENV")

	l, err := NewLogger("testservice")
	assert.NoError(t, err)
	assert.NotNil(t, l)

	l.Info("info message")
	l.Debug("debug message")
	l.Error("error message")
}

// TestLoggerMethods tests the Info, Debug, and Error methods of LoggerInterface.
//
// The test does the following:
//
// 1. Creates a new logger instance using the zap.New function.
// 2. Logs a message using the Info, Debug, and Error methods.
// 3. Asserts that the messages are logged with the correct level.
// 4. Asserts that the context map contains the correct key-value pairs.
func TestLoggerMethods(t *testing.T) {
	core, recorded := observer.New(zapcore.DebugLevel)
	zapLogger := zap.New(core)
	l := &Logger{Log: zapLogger}

	l.Info("info log", zap.String("key", "value"))
	l.Debug("debug log")
	l.Error("error log")

	logs := recorded.All()
	assert.Len(t, logs, 3)

	assert.Equal(t, "info log", logs[0].Message)
	assert.Equal(t, "debug log", logs[1].Message)
	assert.Equal(t, "error log", logs[2].Message)

	assert.Equal(t, zapcore.InfoLevel, logs[0].Level)
	assert.Equal(t, zapcore.DebugLevel, logs[1].Level)
	assert.Equal(t, zapcore.ErrorLevel, logs[2].Level)

	assert.Equal(t, "value", logs[0].ContextMap()["key"])
}
