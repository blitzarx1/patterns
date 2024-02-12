package neighbourhood

import (
	"fmt"
	"sort"
	"strings"

	"github.com/boson-research/patterns/internal/alphabet"
	"github.com/boson-research/patterns/internal/cluster/kmeans"
	"github.com/boson-research/patterns/internal/context"
	"github.com/samber/lo"
)

type Neighbourhood struct {
	Center      *alphabet.Pattern
	Elements    []*alphabet.Pattern
	TextEntries *TextEntries
	Clusters    []*Cluster
}

func New(c *alphabet.Pattern) *Neighbourhood {
	return &Neighbourhood{
		Center: c,
	}
}

func (n *Neighbourhood) WithElements(elements []*alphabet.Pattern) *Neighbourhood {
	n.Elements = elements
	return n
}

func (n *Neighbourhood) FindTextEntries(ctx context.Context, text []byte) {
	ctx, span := ctx.StartSpan("FindTextEntries")
	defer span.End()

	ctx.Logger().Debugf("finding entries in text for %s", n)

	for it := range text {
		for _, pat := range n.Elements {
			if checkPattern(pat, text, it) {
				ctx.Logger().Tracef("adding entry %s at index %d", pat, it)

				if n.TextEntries == nil {
					n.TextEntries = NewTextEntries()
				}

				n.TextEntries.Add(it, pat)
			}
		}
	}
}

func (n *Neighbourhood) Clusterize(ctx context.Context) {
	ctx, span := ctx.StartSpan("Clusterize")
	defer span.End()

	ctx.Logger().Debugf("clusterizing %s", n)

	entryByLoc := make(map[int]*TextEntry, len(n.TextEntries.Locations()))
	for i, loc := range n.TextEntries.Locations() {
		entryByLoc[loc] = &TextEntry{loc: loc, pattern: n.TextEntries.Patterns()[i]}
	}

	clusterInput := lo.Map(n.TextEntries.Locations(), func(loc int, _ int) float64 {
		return float64(loc)
	})

	fmt.Printf("computing kmeans for neighbourhood with center: %s\n", n.Center)
	centroids, labels := kmeans.KMeans(ctx, clusterInput)
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

	// sort clusters by center
	sort.Slice(n.Clusters, func(i, j int) bool {
		return n.Clusters[i].center < n.Clusters[j].center
	})
}

func checkPattern(p *alphabet.Pattern, text []byte, it int) bool {
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

	b.WriteString(fmt.Sprintf("neighbourhood{%s}", n.Center))

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
