package logger

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

// NewLogger creates a new logr.Logger instance with production configuration
func NewLogger() (logr.Logger, error) {
	zapLog, err := zap.NewProduction()
	if err != nil {
		return logr.Logger{}, err
	}

	return zapr.NewLogger(zapLog), nil
}

// NewDevelopmentLogger creates a new logr.Logger instance with development configuration
func NewDevelopmentLogger() (logr.Logger, error) {
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		return logr.Logger{}, err
	}

	return zapr.NewLogger(zapLog), nil
}
