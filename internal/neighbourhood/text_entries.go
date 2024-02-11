package neighbourhood

import (
	"fmt"
	"strings"

	"github.com/boson-research/patterns/internal/alphabet"
	"github.com/samber/lo"
)

type TextEntry struct {
	loc     int
	pattern *alphabet.Pattern
}

func (te *TextEntry) Loc() int {
	return te.loc
}

func (te *TextEntry) Pattern() *alphabet.Pattern {
	return te.pattern
}

func (te *TextEntry) String() string {
	return fmt.Sprintf("{%d - %s}", te.loc, te.pattern)
}

type TextEntries struct {
	locations []int
	patterns  []*alphabet.Pattern
}

func NewTextEntries() *TextEntries {
	return &TextEntries{}
}

func NewTextEntriesWithSize(size int) *TextEntries {
	return &TextEntries{
		locations: make([]int, 0, size),
		patterns:  make([]*alphabet.Pattern, 0, size),
	}
}

func (te *TextEntries) Add(loc int, pat *alphabet.Pattern) {
	te.locations = append(te.locations, loc)
	te.patterns = append(te.patterns, pat)
}

func (te *TextEntries) AddMany(locs []int, pats []*alphabet.Pattern) {
	te.locations = append(te.locations, locs...)
	te.patterns = append(te.patterns, pats...)
}

func (te *TextEntries) Locations() []int {
	if te == nil {
		return nil
	}

	return te.locations
}

func (te *TextEntries) Patterns() []*alphabet.Pattern {
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
