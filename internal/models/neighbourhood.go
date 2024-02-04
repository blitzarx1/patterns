package models

import "fmt"

type Neighbourhood struct {
	Center   Pattern
	Elements []Pattern
}

func (n Neighbourhood) String() string {
	return fmt.Sprintf("Neighbourhood{center: %s, elements: %s}", n.Center, n.Elements)
}
