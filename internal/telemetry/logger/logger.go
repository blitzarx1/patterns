package logger

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const logLevelEnv = "LOG_LEVEL"

func InitLogger(ctx context.Context) context.Context {
	logLevelRaw := os.Getenv(logLevelEnv)

	logLevel, err := logrus.ParseLevel(logLevelRaw)
	if err != nil {
		logLevel = logrus.InfoLevel
	}

	return InjectLogLevel(ctx, logLevel)
}

func Logger(ctx context.Context) *logrus.Entry {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.00000",
	})

	l.SetLevel(LogLevel(ctx))

	l.AddHook(logrusTraceHook{})

	return l.WithContext(ctx)
}

// logrusTraceHook is a hook that;
// (a) adds TraceIds & spanIds to logs of all LogLevels
// (b) adds logs to the active span as events.
type logrusTraceHook struct{}

func (t logrusTraceHook) Levels() []logrus.Level { return logrus.AllLevels }

func (t logrusTraceHook) Fire(entry *logrus.Entry) error {
	ctx := entry.Context
	if ctx == nil {
		return nil
	}
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return nil
	}

	{ // (a) adds TraceIds & spanIds to logs.
		sCtx := span.SpanContext()
		if sCtx.HasTraceID() {
			entry.Data["traceId"] = sCtx.TraceID().String()
		}
		if sCtx.HasSpanID() {
			entry.Data["spanId"] = sCtx.SpanID().String()
		}
	}

	{ // (b) adds logs to the active span as events.

		// code from: https://github.com/uptrace/opentelemetry-go-extra/tree/main/otellogrus
		// whose license(BSD 2-Clause) can be found at: https://github.com/uptrace/opentelemetry-go-extra/blob/v0.1.18/LICENSE
		attrs := make([]attribute.KeyValue, 0)
		logSeverityKey := attribute.Key("log.severity")
		logMessageKey := attribute.Key("log.message")
		attrs = append(attrs, logSeverityKey.String(entry.Level.String()))
		attrs = append(attrs, logMessageKey.String(entry.Message))

		span.AddEvent("log", trace.WithAttributes(attrs...))
		if entry.Level <= logrus.ErrorLevel {
			span.SetStatus(codes.Error, entry.Message)
		}
	}

	return nil
}
