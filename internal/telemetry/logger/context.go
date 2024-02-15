package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ContextKey string

const ContextKeyLogger ContextKey = "logger"

// InjectIntoContext injects the logger into the context. Supposed to be used in the main or other entrypoint functions.
func InjectIntoContext(ctx context.Context, logger *logrus.Logger) context.Context {
	return context.WithValue(ctx, ContextKeyLogger, logger)
}

// FromContext returns the logger from the context. If the logger is not found, the second return value is false.
func FromContext(ctx context.Context) (*logrus.Entry, bool) {
	logger, ok := ctx.Value(ContextKeyLogger).(*logrus.Logger)
	if !ok {
		return nil, false
	}

	return logger.WithContext(ctx), ok
}

func MustFromContext(ctx context.Context) *logrus.Entry {
	logger, ok := FromContext(ctx)
	if !ok {
		panic("logger not found in context")
	}

	return logger
}
