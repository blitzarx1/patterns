package context

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Key string

const ctxKeyLogger Key = "logger"

type Context interface {
	context.Context
	WithLogger(l *logrus.Logger) Context
	Logger() *logrus.Entry
	StartSpan(name string) (Context, trace.Span)
}

type PatternContext struct {
	context.Context
}

func New(ctx context.Context) *PatternContext {
	return &PatternContext{ctx}
}

func WithValue(ctx context.Context, k Key, v interface{}) Context {
	return &PatternContext{context.WithValue(ctx, k, v)}
}

func (c *PatternContext) WithLogger(l *logrus.Logger) Context {
	return WithValue(c, ctxKeyLogger, l)
}

func (c *PatternContext) Logger() *logrus.Entry {
	l, ok := c.Value(ctxKeyLogger).(*logrus.Logger)
	if !ok {
		return logrus.New().WithContext(c)
	}

	return l.WithContext(c)
}

func (c *PatternContext) StartSpan(name string) (Context, trace.Span) {
	ctx, span := otel.Tracer("").Start(c.Context, name)
	return New(ctx), span
}
