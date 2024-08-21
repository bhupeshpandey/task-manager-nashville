package logger

import (
	"fmt"
	"github.com/bhupeshpandey/task-manager-nashville/internal/models"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.Logger
}

func newZapLogger(cfg *models.LoggingConfig) models.Logger {

	var zapAtomicLevel = zap.NewAtomicLevelAt(zap.InfoLevel)

	switch cfg.LogLevel {
	case "debug":
		zapAtomicLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "error":
		zapAtomicLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}

	var env = true

	if cfg.Environment == "prod" {
		env = false
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:            zapAtomicLevel,
		Development:      env,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout", "logfile"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	return &zapLogger{
		logger: logger,
	}
}

func (zapLG *zapLogger) Log(logType models.LogLevel, message string, logs ...interface{}) {
	var fields []zapcore.Field
	for id, lg := range logs {
		fields = append(fields, zap.Any(fmt.Sprintf("%d", id), lg))
	}
	var logLevel = zapcore.InfoLevel
	switch logType {
	case models.ErrorLevel:
		logLevel = zapcore.ErrorLevel
	case models.DebugLevel:
		logLevel = zapcore.DebugLevel
	case models.WarnLevel:
		logLevel = zapcore.WarnLevel
	}
	zapLG.logger.Log(logLevel, message, fields...)
}
