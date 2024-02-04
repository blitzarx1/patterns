package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

type logLevelKeyType string

const logLevelKey logLevelKeyType = "log_level"

func InjectLogLevel(ctx context.Context, level logrus.Level) context.Context {
	return context.WithValue(ctx, logLevelKey, level)
}

func LogLevel(ctx context.Context) logrus.Level {
	if level, ok := ctx.Value(logLevelKey).(logrus.Level); ok {
		return level
	}

	return logrus.InfoLevel
}
