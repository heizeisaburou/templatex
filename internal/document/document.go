package document

// :TODO: as un motodo len, que retorne el tamaño (retornarlo como un ByteOffset)
// :TODO: hacer un test para comprobar que todo en linemap a partir del toPosistion funcione correctamente.

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

// ToPosition convierte un offset de bytes en una posición dentro del documento.
//
// Los offsets situados fuera de los límites del documento se ajustan al rango válido.
func (d Document) ToPosition(offset ByteOffset) Position {
	return d.lm.toPosition(offset)
}

// ToRegion convierte un rango de bytes en un rango de posiciones dentro del
// documento.
//
// Los offsets situados fuera de los límites del documento se ajustan al rango válido.
func (d Document) ToRegion(r Range) Region {
	return d.lm.toRegion(r)
}

func (d Document) Len() ByteOffset {
	return ByteOffset(len(d.src))
}
