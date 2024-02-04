package neighbourhood

import (
	"context"

	"github.com/boson-research/patterns/internal/models"
	"github.com/boson-research/patterns/internal/telemetry/logger"
	"go.opentelemetry.io/otel"
)

func ExtractCenters(ctx context.Context, a models.Alphabet) []models.Pattern {
	ctx, span := otel.Tracer("neighbourhood").Start(ctx, "ExtractCenters")
	defer span.End()
	l := logger.Logger(ctx)

	l.Debug("extracting centers")

	centers := make([]models.Pattern, 0, len(a))
	for _, s1 := range a {
		for _, s2 := range a {
			centers = append(centers, models.Pattern{s2, s1, s1})

		}
	}
	return centers
}
