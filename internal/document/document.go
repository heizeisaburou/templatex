package document

import (
	"bytes"
	"fmt"
	"unicode/utf8"
)

type Document struct {
	lm  lineMap
	src []byte
}

// New construye un documento a partir de src.
//
// Devuelve ErrInvalidUTF8 si src no es una secuencia UTF-8 válida.
func New(src []byte) (Document, error) {
	if !utf8.Valid(src) {
		return Document{}, ErrInvalidUTF8
	}

	return Document{lm: newLineMap(src), src: src}, nil
}

// ToPosition convierte un offset de bytes en una posición dentro del documento.
//
// Devuelve ErrOutOfBounds si offset no pertenece al intervalo [0, Len()] y
// ErrByteOffsetNotAtRuneBoundary si offset no apunta al inicio de un carácter.
func (d Document) ToPosition(offset ByteOffset) (Position, error) {
	return d.lm.toPosition(offset)
}

// ToRegion convierte un rango de bytes en un rango de posiciones dentro del
// documento.
//
// Propaga el error de [Document.ToPosition] si alguno de los extremos del rango
// no es un offset válido dentro del documento.
func (d Document) ToRegion(r Range) (Region, error) {
	return d.lm.toRegion(r)
}

// Len devuelve la longitud en bytes del documento.
func (d Document) Len() ByteOffset {
	return ByteOffset(len(d.src))
}

// Range devuelve una copia independiente de los bytes contenidos en rng.
//
// La copia es independiente: modificarla no afecta al documento.
//
// Devuelve ErrOutOfBounds si rng no está contenido en [0, Len()].
//
// El nombre sigue el convenio de pkg/list: Range copia un subrango y recibe los
// límites, mientras que Slice copia el contenido completo y no recibe
// parámetros. Por eso este método no se llama Slice; si en el futuro hace falta
// una copia del documento entero, corresponde a un Slice() []byte aparte.
func (d Document) Range(rng Range) ([]byte, error) {
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
