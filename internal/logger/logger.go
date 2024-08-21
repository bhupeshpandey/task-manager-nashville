package logger

import "github.com/bhupeshpandey/task-manager-nashville/internal/models"

func NewLogger(config *models.LoggingConfig) models.Logger {
	switch config.Type {
	case "zap":
		return newZapLogger(config)
	}
	return nil
}
