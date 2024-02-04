package neighbourhood

import (
	"context"

	"github.com/boson-research/patterns/internal/models"
	"github.com/boson-research/patterns/internal/telemetry/logger"
	"go.opentelemetry.io/otel"
)

func ExtractNeighbourhoods(ctx context.Context, a models.Alphabet, c []models.Pattern) []models.Neighbourhood {
	ctx, span := otel.Tracer("neighbourhood").Start(ctx, "ExtractNeighbourhoods")
	defer span.End()
	l := logger.Logger(ctx)

	l.Debug("extracting neighbourhoods")

	neighbourhoods := make([]models.Neighbourhood, 0, len(c))
	for _, center := range c {
		var elements []models.Pattern
		for _, symbol := range a {
			elements = append(elements, models.Pattern{center[0], symbol, center[2]})
		}

		neighbourhoods = append(neighbourhoods, models.Neighbourhood{Center: center, Elements: elements})
	}
	return neighbourhoods
}
