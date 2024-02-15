package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/boson-research/patterns/internal/telemetry/logger"
	"github.com/boson-research/patterns/internal/telemetry/tracing"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Name               string
	Version            string
	JaegerOTLPEndpoint string
}

func Init(ctx context.Context, cfg Config) (*logrus.Logger, func(ctx context.Context) error, error) {
	if cfg.JaegerOTLPEndpoint == "" {
		return nil, func(_ context.Context) error {
			return nil
		}, nil
	}

	closer, err := tracing.InitTracerProvider(ctx, cfg.Name, cfg.Version, cfg.JaegerOTLPEndpoint, time.Second)
	if err != nil {
		return nil, nil, fmt.Errorf("initialize tracing: %w", err)
	}

	return logger.MustCreate(), closer, nil
}
