package neighbourhood

import (
	"fmt"

	"github.com/boson-research/patterns/internal/models"
)

type Stat struct {
	locations []int
	patterns  []*models.Pattern
}

func NewStat() *Stat {
	return &Stat{}
}

func NewStatWithSize(size int) *Stat {
	return &Stat{
		locations: make([]int, 0, size),
		patterns:  make([]*models.Pattern, 0, size),
	}
}

func (s *Stat) Add(loc int, pat *models.Pattern) {
	s.locations = append(s.locations, loc)
	s.patterns = append(s.patterns, pat)
}

func (s *Stat) AddMany(locs []int, pats []*models.Pattern) {
	s.locations = append(s.locations, locs...)
	s.patterns = append(s.patterns, pats...)
}

func (s *Stat) Locations() []int {
	return s.locations
}

func (s *Stat) Patterns() []*models.Pattern {
	return s.patterns
}

func (s *Stat) String() string {
	return fmt.Sprintf("Stat{locations: %v, patterns: %v}", s.locations, s.patterns)
}
