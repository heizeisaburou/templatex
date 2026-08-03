package document

import "github.com/heizeisaburou/templatex/pkg/rnge"

// Range representa un rango semiabierto de ByteOffset.
type Range struct {
	rnge.Range[ByteOffset]
}

// NewRange crea el rango semiabierto [start, end).
func NewRange(start, end ByteOffset) Range {
	return Range{
		Range: rnge.New(start, end),
	}
}

// Len devuelve la longitud del rango.
func (r Range) Len() ByteOffset {
	return r.End() - r.Start()
}
