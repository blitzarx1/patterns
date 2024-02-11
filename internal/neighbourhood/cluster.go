package neighbourhood

import (
	"fmt"
	"strings"
)

type Cluster struct {
	center  float64
	entries []*TextEntry
}

func (c *Cluster) String() string {
	b := strings.Builder{}
	for _, e := range c.entries {
		b.WriteString(fmt.Sprintf("%.2f-{%d - %s}\n", c.center, e.Loc(), e.Pattern()))
	}

	return b.String()
}
