package neighbourhood

import (
	"fmt"
	"strings"

	"github.com/boson-research/patterns/internal/models"
	"github.com/samber/lo"
)

type TextEntry struct {
	loc     int
	pattern *models.Pattern
}

func (te *TextEntry) Loc() int {
	return te.loc
}

func (te *TextEntry) Pattern() *models.Pattern {
	return te.pattern
}

func (te *TextEntry) String() string {
	return fmt.Sprintf("{%d - %s}", te.loc, te.pattern)
}

type TextEntries struct {
	locations []int
	patterns  []*models.Pattern
}

func NewTextEntries() *TextEntries {
	return &TextEntries{}
}

func NewTextEntriesWithSize(size int) *TextEntries {
	return &TextEntries{
		locations: make([]int, 0, size),
		patterns:  make([]*models.Pattern, 0, size),
	}
}

func (te *TextEntries) Add(loc int, pat *models.Pattern) {
	te.locations = append(te.locations, loc)
	te.patterns = append(te.patterns, pat)
}

func (te *TextEntries) AddMany(locs []int, pats []*models.Pattern) {
	te.locations = append(te.locations, locs...)
	te.patterns = append(te.patterns, pats...)
}

func (te *TextEntries) Locations() []int {
	if te == nil {
		return nil
	}

	return te.locations
}

func (te *TextEntries) Patterns() []*models.Pattern {
	if te == nil {
		return nil
	}

	return te.patterns
}

func (te *TextEntries) String() string {
	b := strings.Builder{}
	b.WriteString(strings.Join(lo.Map(te.locations, func(_ int, i int) string {
		return fmt.Sprintf("{%d - %s}", te.locations[i], te.patterns[i])
	}), "\n"))

	return b.String()
}
