package alphabet

import "github.com/samber/lo"

type Alphabet []byte

func New(raw []byte, symbolSize int) []Symbol {
	return lo.Map(lo.Chunk(raw, symbolSize), func(chunk []byte, _ int) Symbol {
		return Symbol(chunk)
	})
}

func (a Alphabet) String() string {
	return string(a)
}
