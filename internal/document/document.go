package document

type Document struct {
	lm  lineMap
	src []byte
}

// ToPosition convierte un offset de bytes en una posición dentro del documento.
//
// Los offsets situados fuera de los límites del documento se ajustan al rango válido.
func (d Document) ToPosition(offset ByteOffset) Position {
	return d.lm.toPosition(d.src, offset)
}

// toRegion convierte un rango de bytes en un rango de posiciones dentro del
// documento.
//
// Los offsets situados fuera de los límites del documento se ajustan al rango válido.
func (d Document) ToRegion(r Range) Region {
	return d.lm.toRegion(d.src, r)
}
