package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log levels
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

// Logger represents a logger
type Logger struct {
	zap *zap.Logger
}

// Config represents logger configuration
type Config struct {
	Level      string
	AppEnv     string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

// NewLogger creates a new logger instance
func NewLogger(config Config) (*Logger, error) {
	// Set default level to info if not specified
	if config.Level == "" {
		config.Level = InfoLevel
	}

	// Parse log level
	level := zapcore.InfoLevel
	switch strings.ToLower(config.Level) {
	case DebugLevel:
		level = zapcore.DebugLevel
	case InfoLevel:
		level = zapcore.InfoLevel
	case WarnLevel:
		level = zapcore.WarnLevel
	case ErrorLevel:
		level = zapcore.ErrorLevel
	case FatalLevel:
		level = zapcore.FatalLevel
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create core with console output
	var core zapcore.Core

	// For development, use console encoder
	if config.AppEnv == "development" {
		// Create console encoder
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		// Create console writer
		consoleOutput := zapcore.AddSync(os.Stdout)
		// Create console core
		core = zapcore.NewCore(consoleEncoder, consoleOutput, level)
	} else {
		// For production, use JSON encoder
		jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
		// Create writer for file if specified
		var output zapcore.WriteSyncer
		if config.Filename != "" {
			// Create file writer
			file, err := os.OpenFile(config.Filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return nil, fmt.Errorf("failed to open log file: %w", err)
			}
			// Use multi-writer for both file and console
			multiWriter := io.MultiWriter(file, os.Stdout)
			output = zapcore.AddSync(multiWriter)
		} else {
			// Use only console writer
			output = zapcore.AddSync(os.Stdout)
		}
		// Create json core
		core = zapcore.NewCore(jsonEncoder, output, level)
	}

	// Create zap logger
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	// Return logger
	return &Logger{
		zap: zapLogger,
	}, nil
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...zapcore.Field) {
	l.zap.Debug(msg, fields...)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...zapcore.Field) {
	l.zap.Info(msg, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...zapcore.Field) {
	l.zap.Warn(msg, fields...)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...zapcore.Field) {
	l.zap.Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, fields ...zapcore.Field) {
	l.zap.Fatal(msg, fields...)
}

// Sync flushes the logger
func (l *Logger) Sync() error {
	return l.zap.Sync()
}

// Field creates a field
func Field(key string, value interface{}) zapcore.Field {
	return zap.Any(key, value)
}

// String creates a string field
func String(key, value string) zapcore.Field {
	return zap.String(key, value)
}

// Int creates an int field
func Int(key string, value int) zapcore.Field {
	return zap.Int(key, value)
}

// Error creates an error field
func Error(err error) zapcore.Field {
	return zap.Error(err)
}

// FiberLogger returns a fiber middleware logger
func FiberLogger(l *Logger) fiber.Handler {
	// Create custom format for fiber logger
	format := "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path}\n"

	// Return fiber logger middleware
	return logger.New(logger.Config{
		Format:     format,
		TimeFormat: time.RFC3339,
		TimeZone:   "Local",
		Output: zapWriter{
			logger: l,
		},
	})
}

// zapWriter implements io.Writer interface for zap logger
type zapWriter struct {
	logger *Logger
}

// Write implements io.Writer interface
func (w zapWriter) Write(p []byte) (n int, err error) {
	w.logger.Info(string(p))
	return len(p), nil
}
