package processor

import (
	"fmt"
	"strings"

	"github.com/boson-research/patterns/internal/alphabet"
	"github.com/boson-research/patterns/internal/context"
	"github.com/boson-research/patterns/internal/neighbourhood"
	"github.com/samber/lo"
)

type Processor struct {
	neighbourhoods []*neighbourhood.Neighbourhood
}

func New(ctx context.Context) *Processor {
	return &Processor{}
}

func (p *Processor) AnalyzeAlphabet(ctx context.Context, a alphabet.Alphabet) {
	ctx, span := ctx.StartSpan("AnalyzeAlphabet")
	defer span.End()

	ctx.Logger().Debug("analyzing alphabet")

	centers := p.extractCenters(ctx, a)
	p.neighbourhoods = p.extractNeighbourhoods(ctx, a, centers)

	ctx.Logger().Debugf("alphabet analyzed\n%s", p.neighbourhoods)
}

func (p *Processor) AnalyzeText(ctx context.Context, text []byte) {
	ctx, span := ctx.StartSpan("AnalyzeText")
	defer span.End()

	ctx.Logger().Debug("analyzing text")

	p.findTextEntries(ctx, text)
	p.clusterize(ctx)

	ctx.Logger().Info("text analyzed")

	// print resulting neighbourhoods
	fmt.Print(
		strings.Join(
			lo.Map(
				lo.Filter(
					p.neighbourhoods,
					func(n *neighbourhood.Neighbourhood, _ int) bool { return n.TextEntries != nil },
				),
				func(n *neighbourhood.Neighbourhood, _ int) string { return n.String() },
			),
			"\n",
		),
	)
}

func (p *Processor) clusterize(ctx context.Context) {
	ctx, span := ctx.StartSpan("clusterize")
	defer span.End()

	ctx.Logger().Debug("clusterizing")

	for _, n := range p.neighbourhoods {
		if len(n.TextEntries.Locations()) == 0 {
			continue
		}

		n.Clusterize(ctx)
	}
}

func (p *Processor) findTextEntries(ctx context.Context, text []byte) {
	ctx, span := ctx.StartSpan("findTextEntries")
	defer span.End()

	ctx.Logger().Debug("finding text entries")

	for _, n := range p.neighbourhoods {
		n.FindTextEntries(ctx, text)
	}
}

func (p *Processor) extractCenters(ctx context.Context, a alphabet.Alphabet) []*alphabet.Pattern {
	ctx, span := ctx.StartSpan("extractCenters")
	defer span.End()

	ctx.Logger().Debug("extracting centers")

	centers := make([]*alphabet.Pattern, 0, len(a))
	for _, s1 := range a {
		for _, s2 := range a {
			centers = append(centers, alphabet.NewPattern([]byte{s2, s1, s1}))
		}
	}

	ctx.Logger().Debugf("extracted centers: %v", centers)

	return centers
}

func (p *Processor) extractNeighbourhoods(ctx context.Context, a alphabet.Alphabet, c []*alphabet.Pattern) []*neighbourhood.Neighbourhood {
	ctx, span := ctx.StartSpan("extractNeighbourhoods")
	defer span.End()

	ctx.Logger().Debug("extracting neighbourhoods")

	neighbourhoods := make([]*neighbourhood.Neighbourhood, 0, len(c))
	for _, center := range c {
		var elements []*alphabet.Pattern
		for _, symbol := range a {
			elements = append(elements, alphabet.NewPattern([]byte{center.Value()[0], symbol, center.Value()[2]}))
		}

		neighbourhoods = append(neighbourhoods, neighbourhood.New(center).WithElements(elements))
	}

	ctx.Logger().Debugf("extracted neighbourhoods: %v", neighbourhoods)

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
