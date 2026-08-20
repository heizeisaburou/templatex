package document

import (
	"bytes"
	"errors"
	"fmt"
	"unicode/utf8"
)

type Document struct {
	lm  lineMap
	src []byte
}

var (
	ErrInvalidUTF8 = errors.New("La secuencia utf8 no es valida")
	ErrOutOfBounds = errors.New("Indice fuera de los limites del documento")
)

// New este el el constructor para documentos
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

// Slice retorna una copia independiente de la región solicitada.
func (d Document) Slice(rng Range) ([]byte, error) {
	start := rng.Start()
	end := rng.End()

	if start < 0 || end > d.Len() {
		return nil, fmt.Errorf(
			"%w [range = %v, document.Len() = %d]",
			ErrOutOfBounds,
			rng,
			d.Len(),
		)
	}

	subSlice := d.src[start:end]
	return bytes.Clone(subSlice), nil
}
