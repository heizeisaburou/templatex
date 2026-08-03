package document

import "github.com/heizeisaburou/templatex/pkg/rnge"

// Region representa un rango semiabierto de posiciones de un documento.
type Region struct {
	rnge.Range[Position]
}

// NewRegion crea la región semiabierta [start, end).
func NewRegion(start, end Position) Region {
	return Region{
		Range: rnge.New(start, end),
	}
}
