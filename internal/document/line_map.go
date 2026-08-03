package document

type lineMap struct {
	lineStarts []int
}

// newLineMap construye el mapa de posiciones de un documento fuente inmutable.
func newLineMap(src []byte) lineMap {
	lm := lineMap{
		lineStarts: []int{0},
	}

	for i, b := range src {
		if b == '\n' {
			lm.lineStarts = append(lm.lineStarts, i+1)
		}
	}

	return lm
}

// toPosition convierte un ByteOffset en una posición dentro del documento.
//
// Los ByteOffset situados fuera de los límites del documento se ajustan al rango válido.
func (m lineMap) toPosition(src []byte, offset ByteOffset) Position {
	panic("termina esta función")
}

// toRegion convierte un rango de bytes en un rango de posiciones dentro del
// documento.
//
// Los offsets situados fuera de los límites del documento se ajustan al rango válido.
func (m lineMap) toRegion(src []byte, r Range) Region {
	return NewRegion(
		m.toPosition(src, r.Start()),
		m.toPosition(src, r.End()),
	)
}
