package processor

import (
	"context"
	"fmt"

	"github.com/boson-research/patterns/internal/alphabet"
	"github.com/boson-research/patterns/internal/neighbourhood"
	"github.com/boson-research/patterns/internal/telemetry/logger"
	"go.opentelemetry.io/otel"
)

type Processor struct {
	neighbourhoods []*neighbourhood.Neighbourhood
}

func New(ctx context.Context) *Processor {
	return &Processor{}
}

func (p *Processor) AnalyzeAlphabet(ctx context.Context, a alphabet.Alphabet) {
	ctx, span := otel.Tracer("processor").Start(ctx, "AnalyzeAlphabet")
	defer span.End()
	l := logger.Logger(ctx)

	l.Debug("analyzing alphabet")

	centers := p.extractCenters(ctx, a)
	p.neighbourhoods = p.extractNeighbourhoods(ctx, a, centers)

	l.Debugf("alphabet analyzed\n%s", p.neighbourhoods)
}

func (p *Processor) AnalyzeText(ctx context.Context, text []byte) {
	ctx, span := otel.Tracer("processor").Start(ctx, "AnalyzeText")
	defer span.End()
	l := logger.Logger(ctx)

	l.Debug("analyzing text")

	p.findTextEntries(ctx, text)
	p.clusterize(ctx)

	l.Info("text analyzed")
	for _, n := range p.neighbourhoods {
		if n.TextEntries == nil {
			continue
		}
		fmt.Printf("%s", n)
	}
}

func (p *Processor) clusterize(ctx context.Context) {
	ctx, span := otel.Tracer("processor").Start(ctx, "clusterize")
	defer span.End()
	l := logger.Logger(ctx)

	l.Debug("clusterizing")

	for _, n := range p.neighbourhoods {
		if len(n.TextEntries.Locations()) == 0 {
			continue
		}

		n.Clusterize(ctx)
	}
}

func (p *Processor) findTextEntries(ctx context.Context, text []byte) {
	ctx, span := otel.Tracer("processor").Start(ctx, "findTextEntries")
	defer span.End()
	l := logger.Logger(ctx)

	l.Debug("finding text entries")

	for _, n := range p.neighbourhoods {
		n.FindTextEntries(ctx, text)
	}
}

func (p *Processor) extractCenters(ctx context.Context, a alphabet.Alphabet) []*alphabet.Pattern {
	ctx, span := otel.Tracer("processor").Start(ctx, "extractCenters")
	defer span.End()
	l := logger.Logger(ctx)

	l.Debug("extracting centers")

	centers := make([]*alphabet.Pattern, 0, len(a))
	for _, s1 := range a {
		for _, s2 := range a {
			centers = append(centers, alphabet.NewPattern([]byte{s2, s1, s1}))
		}
	}

	l.Debugf("extracted centers: %v", centers)

	return centers
}

func (p *Processor) extractNeighbourhoods(ctx context.Context, a alphabet.Alphabet, c []*alphabet.Pattern) []*neighbourhood.Neighbourhood {
	ctx, span := otel.Tracer("processor").Start(ctx, "extractNeighbourhoods")
	defer span.End()
	l := logger.Logger(ctx)

	l.Debug("extracting neighbourhoods")

	neighbourhoods := make([]*neighbourhood.Neighbourhood, 0, len(c))
	for _, center := range c {
		var elements []*alphabet.Pattern
		for _, symbol := range a {
			elements = append(elements, alphabet.NewPattern([]byte{center.Value()[0], symbol, center.Value()[2]}))
		}

		neighbourhoods = append(neighbourhoods, neighbourhood.New(center).WithElements(elements))
	}

	l.Debugf("extracted neighbourhoods: %v", neighbourhoods)

	return neighbourhoods
}

func mergeStatsNeighbourhoods(a *neighbourhood.TextEntries, b *neighbourhood.TextEntries) *neighbourhood.TextEntries {
	if a == nil && b == nil {
		return nil
	}

	if a == nil || len(a.Locations()) == 0 {
		return b
	}

	if b == nil || len(b.Locations()) == 0 {
		return a
	}

	res := neighbourhood.NewTextEntriesWithSize(len(a.Locations()) + len(b.Locations()))
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
