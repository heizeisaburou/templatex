package document

import (
	"errors"
	"unicode/utf8"
)

// :TODO: hacer un test para comprobar que todo en linemap a partir del toPosistion funcione correctamente.

type Document struct {
	lm  lineMap
	src []byte
}

var (
	ErrInvalidUTF8 = errors.New("La secuencia utf8 no es valida")
)

func New(src []byte) (Document, error) {
	if !utf8.Valid(src) {
		return Document{}, ErrInvalidUTF8
	}

	return Document{lm: newLineMap(src), src: src}, nil
}

// ToPosition convierte un offset de bytes en una posición dentro del documento.
//
// Los offsets situados fuera de los límites del documento se ajustan al rango válido.
func (d Document) ToPosition(offset ByteOffset) (Position, error) {
	return d.lm.toPosition(offset)
}

// ToRegion convierte un rango de bytes en un rango de posiciones dentro del
// documento.
//
// Los offsets situados fuera de los límites del documento se ajustan al rango válido.
func (d Document) ToRegion(r Range) (Region, error) {
	return d.lm.toRegion(r)
}

func (d Document) Len() ByteOffset {
	return ByteOffset(len(d.src))
}
