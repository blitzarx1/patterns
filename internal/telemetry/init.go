package telemetry

import (
	pkgContext "context"
	"fmt"
	"time"

	"github.com/boson-research/patterns/internal/context"
	"github.com/boson-research/patterns/internal/telemetry/logger"
	"github.com/boson-research/patterns/internal/telemetry/tracing"
)

type Config struct {
	Name               string
	Version            string
	JaegerOTLPEndpoint string
}

func Init(ctx context.Context, cfg Config) (context.Context, func(ctx pkgContext.Context) error, error) {
	var closer func(pkgContext.Context) error
	var err error
	if cfg.JaegerOTLPEndpoint == "" {
		closer = func(_ pkgContext.Context) error {
			return nil
		}
	} else {
		closer, err = tracing.InitTracerProvider(ctx, cfg.Name, cfg.Version, cfg.JaegerOTLPEndpoint, time.Second)
		if err != nil {
			return nil, nil, fmt.Errorf("initialize tracing: %w", err)
		}
	}

	return logger.MustCreate(ctx), closer, nil
}
