package document

// Utiliza New para crear un Document; no crees un Document{} directamente.
type Document struct {
	lm  lineMap
	src []byte
}

func New(src []byte) Document {
	return Document{
		lm:  newLineMap(src),
		src: src,
	}
}

// ToPosition convierte un ByteOffset en una posición dentro del documento.
//
// Los ByteOffset situados fuera de los límites del documento se ajustan al rango válido.
func (d Document) ToPosition(offset ByteOffset) Position {
	return d.lm.toPosition(d.src, offset)
}

// ToRegion convierte un rango de bytes en un rango de posiciones dentro del
// documento.
//
// Los offsets situados fuera de los límites del documento se ajustan al rango válido.
func (d Document) ToRegion(r Range) Region {
	return d.lm.toRegion(d.src, r)
}
