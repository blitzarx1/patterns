package processor

import (
	// "fmt"
	// "strings"

	"context"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/boson-research/patterns/internal/alphabet"
	"github.com/boson-research/patterns/internal/neighbourhood"
	"github.com/boson-research/patterns/internal/telemetry/logger"
	"go.opentelemetry.io/otel"
	// "github.com/samber/lo"
)

type Processor struct {
	neighbourhoods []*neighbourhood.Neighbourhood
}

func New(ctx context.Context) *Processor {
	return &Processor{}
}

func (p *Processor) AnalyzeAlphabet(ctx context.Context, a alphabet.Alphabet) {
	ctx, span := otel.Tracer("").Start(ctx, "AnalyzeAlphabet")
	defer span.End()

	logger.MustFromContext(ctx).Debug("analyzing alphabet")

	centers := p.extractCenters(ctx, a)
	p.neighbourhoods = p.extractNeighbourhoods(ctx, a, centers)

	logger.MustFromContext(ctx).Debugf("alphabet analyzed\n%s", p.neighbourhoods)
}

func (p *Processor) AnalyzeText(ctx context.Context, text []byte) {
	ctx, span := otel.Tracer("").Start(ctx, "AnalyzeText")
	defer span.End()

	logger.MustFromContext(ctx).Debug("analyzing text")

	p.findTextEntries(ctx, text)
	p.exportNeighbourhoods()
	// p.clusterize(ctx)

	// logger.MustFromContext(ctx).Info("text analyzed")

	// // print resulting neighbourhoods
	// fmt.Print(
	// 	strings.Join(
	// 		lo.Map(
	// 			lo.Filter(
	// 				p.neighbourhoods,
	// 				func(n *neighbourhood.Neighbourhood, _ int) bool { return n.TextEntries != nil },
	// 			),
	// 			func(n *neighbourhood.Neighbourhood, _ int) string { return n.String() },
	// 		),
	// 		"\n",
	// 	),
	// )
}

func (p *Processor) exportNeighbourhoods() {
	for _, n := range p.neighbourhoods {
		if len(n.TextEntries.Locations()) == 0 {
			continue
		}

		// write csv to file
		file, err := os.Create(fmt.Sprintf("output/%s.csv", n.Center.String()))
		if err != nil {
			panic(err)
		}

		w := csv.NewWriter(file)
		defer w.Flush()
		for i, loc := range n.TextEntries.Locations() {
			w.Write([]string{fmt.Sprintf("%d", loc), n.TextEntries.Patterns()[i].String()})
		}
	}
}

func (p *Processor) clusterize(ctx context.Context) {
	ctx, span := otel.Tracer("").Start(ctx, "clusterize")
	defer span.End()

	logger.MustFromContext(ctx).Debug("clusterizing")

	for _, n := range p.neighbourhoods {
		if len(n.TextEntries.Locations()) == 0 {
			continue
		}

		n.Clusterize(ctx)
	}
}

func (p *Processor) findTextEntries(ctx context.Context, text []byte) {
	ctx, span := otel.Tracer("").Start(ctx, "findTextEntries")
	defer span.End()

	logger.MustFromContext(ctx).Debug("finding text entries")

	for _, n := range p.neighbourhoods {
		n.FindTextEntries(ctx, text)
	}
}

func (p *Processor) extractCenters(ctx context.Context, a alphabet.Alphabet) []*alphabet.Pattern {
	ctx, span := otel.Tracer("").Start(ctx, "extractCenters")
	defer span.End()

	logger.MustFromContext(ctx).Debug("extracting centers")

	centers := make([]*alphabet.Pattern, 0, len(a))
	for _, s1 := range a {
		for _, s2 := range a {
			centers = append(centers, alphabet.NewPattern([]byte{s2, s1, s1}))
		}
	}

	logger.MustFromContext(ctx).Debugf("extracted centers: %v", centers)

	return centers
}

func (p *Processor) extractNeighbourhoods(ctx context.Context, a alphabet.Alphabet, c []*alphabet.Pattern) []*neighbourhood.Neighbourhood {
	ctx, span := otel.Tracer("").Start(ctx, "extractNeighbourhoods")
	defer span.End()

	logger.MustFromContext(ctx).Debug("extracting neighbourhoods")

	neighbourhoods := make([]*neighbourhood.Neighbourhood, 0, len(c))
	for _, center := range c {
		var elements []*alphabet.Pattern
		for _, symbol := range a {
			elements = append(elements, alphabet.NewPattern([]byte{center.Value()[0], symbol, center.Value()[2]}))
		}

		neighbourhoods = append(neighbourhoods, neighbourhood.New(center).WithElements(elements))
	}

	logger.MustFromContext(ctx).Debugf("extracted neighbourhoods: %v", neighbourhoods)

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
