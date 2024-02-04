package models

type Alphabet []byte

func (a Alphabet) String() string {
	return string(a)
}
