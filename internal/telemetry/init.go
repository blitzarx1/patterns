package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/boson-research/patterns/internal/telemetry/logger"
	"github.com/boson-research/patterns/internal/telemetry/tracing"
)

type Config struct {
	Name               string
	Version            string
	JaegerOTLPEndpoint string
}

func Init(ctx context.Context, cfg Config) (context.Context, func(ctx context.Context) error, error) {
	closer, err := tracing.InitTracerProvider(ctx, cfg.Name, cfg.Version, cfg.JaegerOTLPEndpoint, time.Second)
	if err != nil {
		return nil, nil, fmt.Errorf("initialize tracing: %w", err)
	}
	return logger.InitLogger(ctx), closer, nil
}
