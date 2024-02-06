package models

type Pattern struct {
	value []byte
}

func NewPattern(v []byte) *Pattern {
	return &Pattern{value: v}
}

func (p *Pattern) String() string {
	return string(p.value)
}

func (p *Pattern) Value() []byte {
	return p.value
}