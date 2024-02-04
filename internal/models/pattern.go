package models

type (
	Pattern    []byte
	MapPattern string
)

func (p Pattern) String() string {
	return string(p)
}

func NewMapPattern(p Pattern) MapPattern {
	return MapPattern(p)
}

func (m MapPattern) Pattern() Pattern {
	return Pattern(m)
}
