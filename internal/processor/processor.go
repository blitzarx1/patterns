package processor

import (
	"context"

	"github.com/boson-research/patterns/internal/models"
	"github.com/boson-research/patterns/internal/models/neighbourhood"
	"github.com/boson-research/patterns/internal/telemetry/logger"
	"go.opentelemetry.io/otel"
)

type Processor struct {
	neighbourhoods []*neighbourhood.Neighbourhood
	stats          map[*neighbourhood.Neighbourhood]*neighbourhood.Stat
}

func New(ctx context.Context) *Processor {
	return &Processor{}
}

func (p *Processor) AnalyzeAlphabet(ctx context.Context, a models.Alphabet) {
	ctx, span := otel.Tracer("processor").Start(ctx, "AnalyzeAlphabet")
	defer span.End()
	l := logger.Logger(ctx)

	l.Debug("analyzing alphabet")

	centers := p.extractCenters(ctx, a)
	p.neighbourhoods = p.extractNeighbourhoods(ctx, a, centers)

	l.Debug("alphabet analyzed")
}

func (p *Processor) AnalyzeText(ctx context.Context, text []byte) {
	ctx, span := otel.Tracer("processor").Start(ctx, "AnalyzeText")
	defer span.End()
	l := logger.Logger(ctx)

	l.Debug("analyzing text")

	p.stats = make(map[*neighbourhood.Neighbourhood]*neighbourhood.Stat, len(p.neighbourhoods)*2)
	for it := range text {
		for _, n := range p.neighbourhoods {
			for _, pat := range n.Elements {
				if checkPattern(pat, text, it) {
					if stat, ok := p.stats[n]; !ok {
						stat := neighbourhood.NewStat()
						stat.Add(it, pat)
						p.stats[n] = stat
					} else {
						stat.Add(it, pat)
					}
				}
			}
		}
	}

	l.Debugf("text analyzed: %v", p.stats)
}

func (p *Processor) extractCenters(ctx context.Context, a models.Alphabet) []*models.Pattern {
	ctx, span := otel.Tracer("processor").Start(ctx, "extractCenters")
	defer span.End()
	l := logger.Logger(ctx)

	l.Debug("extracting centers")

	centers := make([]*models.Pattern, 0, len(a))
	for _, s1 := range a {
		for _, s2 := range a {
			centers = append(centers, models.NewPattern([]byte{s2, s1, s1}))
		}
	}

	l.Debugf("extracted centers: %v", centers)

	return centers
}

func (p *Processor) extractNeighbourhoods(ctx context.Context, a models.Alphabet, c []*models.Pattern) []*neighbourhood.Neighbourhood {
	ctx, span := otel.Tracer("processor").Start(ctx, "extractNeighbourhoods")
	defer span.End()
	l := logger.Logger(ctx)

	l.Debug("extracting neighbourhoods")

	neighbourhoods := make([]*neighbourhood.Neighbourhood, 0, len(c))
	for _, center := range c {
		var elements []*models.Pattern
		for _, symbol := range a {
			elements = append(elements, models.NewPattern([]byte{center.Value()[0], symbol, center.Value()[2]}))
		}

		neighbourhoods = append(neighbourhoods, neighbourhood.New(center).WithElements(elements))
	}

	l.Debugf("extracted neighbourhoods: %v", neighbourhoods)

	return neighbourhoods
}

func mergeStatsNeighbourhoods(a *neighbourhood.Stat, b *neighbourhood.Stat) *neighbourhood.Stat {
	if a == nil && b == nil {
		return nil
	}

	if a == nil || len(a.Locations()) == 0 {
		return b
	}

	if b == nil || len(b.Locations()) == 0 {
		return a
	}

	res := neighbourhood.NewStatWithSize(len(a.Locations()) + len(b.Locations()))
	ia, ib := 0, 0
	for {
		if ia == len(a.Locations()) {
			res.AddMany(b.Locations()[ib:], b.Patterns()[ib:])
			return res
		}

		if ib == len(b.Locations()) {
			res.AddMany(a.Locations()[ia:], a.Patterns()[ia:])
			return res
		}

		if a.Locations()[ia] >= b.Locations()[ib] {
			res.Add(b.Locations()[ib], b.Patterns()[ib])
			ib++
		} else {
			res.Add(a.Locations()[ia], a.Patterns()[ia])
			ia++
		}
	}
}

func checkPattern(p *models.Pattern, text []byte, it int) bool {
	for ip := range p.Value() {
		if it+ip >= len(text) {
			return false
		}

		if p.Value()[ip] != text[it+ip] {
			return false
		}
	}

	return true
}
