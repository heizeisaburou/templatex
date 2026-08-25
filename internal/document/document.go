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


// New este el el constructor para documentos
// New construye un documento a partir de un []byte
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

// Len devuelve la longitud del documento
func (d Document) Len() ByteOffset {
	return ByteOffset(len(d.src))
}

// Range retorna una copia independiente de la región solicitada.
//
// Convenio de nombres del proyecto (ver pkg/list): Range devuelve una copia de
// un subrango, mientras que Slice devuelve una copia del contenido completo.
// Este método se llamaba Slice; renombrado para no chocar con ese convenio.
//
// :TODO: Saburou (2026-08-25): NO revertir a Slice. Si hace falta una copia del
// documento entero, añadir un Slice() []byte aparte, sin parámetros.
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
