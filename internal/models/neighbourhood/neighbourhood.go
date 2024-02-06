package neighbourhood

import (
	"fmt"

	"github.com/boson-research/patterns/internal/models"
)

type Neighbourhood struct {
	Center   *models.Pattern
	Elements []*models.Pattern
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

func (n Neighbourhood) String() string {
	return fmt.Sprintf("Neighbourhood{center: %s, elements: %s}", n.Center, n.Elements)
}
