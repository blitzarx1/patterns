package processor

import (
	"context"

	"github.com/boson-research/patterns/internal/models"
	"github.com/boson-research/patterns/internal/neighbourhood"
	"github.com/boson-research/patterns/internal/telemetry/logger"
	"go.opentelemetry.io/otel"
)

type Processor struct {
	alphabet       models.Alphabet
	neighbourhoods []models.Neighbourhood
}

func NewProcessor(ctx context.Context, alphabet models.Alphabet) *Processor {
	ctx, span := otel.Tracer("processor").Start(ctx, "NewProcessor")
	defer span.End()
	l := logger.Logger(ctx)

	l.Debug("creating new processor")

	centers := neighbourhood.ExtractCenters(ctx, alphabet)

	l.Debugf("found centers: %v", centers)

	neighbourhoods := neighbourhood.ExtractNeighbourhoods(ctx, alphabet, centers)

	l.Debugf("found neighbourhoods: %v", neighbourhoods)

	return &Processor{
		alphabet:       alphabet,
		neighbourhoods: neighbourhoods,
	}
}

func (p *Processor) PatternsLocations(ctx context.Context, text []byte) map[models.MapPattern][]int {
	ctx, span := otel.Tracer("processor").Start(ctx, "PatternsLocations")
	defer span.End()
	l := logger.Logger(ctx)

	l.Debug("finding patterns locations")

	res := make(map[models.MapPattern][]int, len(p.neighbourhoods)*2)
	for it := range text {
		for _, n := range p.neighbourhoods {
			for _, pg := range n.Elements {
				if checkPattern(ctx, pg, text, it) {
					res[models.NewMapPattern(pg)] = append(res[models.NewMapPattern(pg)], it)
				}
			}
		}
	}

	return res
}

func checkPattern(ctx context.Context, pattern models.Pattern, text []byte, it int) bool {
	l := logger.Logger(ctx)

	l.Tracef("checking pattern: %s", pattern)

	for ip := range pattern {
		if it+ip >= len(text) {
			return false
		}

		if pattern[ip] != text[it+ip] {
			return false
		}
	}

	return true
}
