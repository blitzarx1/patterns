package neighbourhood

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/boson-research/patterns/internal/cluster/kmeans"
	"github.com/boson-research/patterns/internal/models"
	"github.com/boson-research/patterns/internal/telemetry/logger"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type Neighbourhood struct {
	Center      *models.Pattern
	Elements    []*models.Pattern
	TextEntries *TextEntries
	Clusters    []*Cluster
}

func New(c *models.Pattern) *Neighbourhood {
	return &Neighbourhood{
		Center: c,
	}
}

func (n *Neighbourhood) WithElements(elements []*models.Pattern) *Neighbourhood {
	n.Elements = elements
	return n
}

func (n *Neighbourhood) FindTextEntries(ctx context.Context, text []byte) {
	ctx, span := otel.Tracer("neighbourhood").Start(ctx, "FindTextEntries")
	defer span.End()
	l := logger.Logger(ctx)

	l.Debugf("finding entries in text for %s", n)

	for it := range text {
		for _, pat := range n.Elements {
			if checkPattern(pat, text, it) {
				l.Tracef("adding entry %s at index %d", pat, it)
				if n.TextEntries == nil {
					n.TextEntries = NewTextEntries()
				}

				n.TextEntries.Add(it, pat)
			}
		}
	}
}

func (n *Neighbourhood) Clusterize(ctx context.Context) {
	ctx, span := otel.Tracer("neighbourhood").Start(ctx, "Cluterize")
	defer span.End()
	l := logger.Logger(ctx)

	l.Debugf("clusterizing %s", n)

	entryByLoc := make(map[int]*TextEntry, len(n.TextEntries.Locations()))
	for i, loc := range n.TextEntries.Locations() {
		entryByLoc[loc] = &TextEntry{loc: loc, pattern: n.TextEntries.Patterns()[i]}
	}

	clusterInput := lo.Map(n.TextEntries.Locations(), func(loc int, _ int) float64 {
		return float64(loc)
	})
	// TODO: silhoutte
	centroids, labels := kmeans.KMeans(ctx, clusterInput, int(math.Min(float64(len(clusterInput)), 2)), 100)
	for label, centroid := range centroids {
		entries := make([]*TextEntry, 0, len(labels))

		for i, l := range labels {
			if l == label {
				entries = append(entries, entryByLoc[int(clusterInput[i])])
			}
		}

		n.Clusters = append(n.Clusters, &Cluster{
			center:  centroid,
			entries: entries,
		})
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

func (n *Neighbourhood) String() string {
	b := strings.Builder{}

	b.WriteString(fmt.Sprintf("\nneighbourhood{%s}", n.Center))

	if n.TextEntries != nil {
		b.WriteString(fmt.Sprintf("\ntext entries %d:\n%s", len(n.TextEntries.locations), n.TextEntries))
	} else {
		return b.String()
	}
	if n.Clusters != nil {
		b.WriteString(fmt.Sprintf("\nclusters %d:\n", len(n.Clusters)))
		for _, c := range n.Clusters {
			b.WriteString(c.String())
		}
	} else {
		return b.String()
	}

	return b.String()
}
