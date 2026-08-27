package document

import (
	"fmt"
	"unicode/utf8"
)

type lineMap struct {
	lineStarts []int
	src        []byte
}

// newLineMap construye el mapa de posiciones de un documento fuente inmutable.
func newLineMap(src []byte) lineMap {
	lm := lineMap{
		lineStarts: []int{0},
		src:        src,
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
// Si src[] esta vacio la funcion retornara (EOF) Position{1, 1}
func (m lineMap) toPosition(offset ByteOffset) (Position, error) {
	ln := ByteOffset(len(m.src))
	if offset < 0 || offset > ln {
		return Position{}, fmt.Errorf("%w [len=%d, offset=%d]", ErrOutOfBounds, ln, offset)
	}

	if offset < ln {
		// No fallará siempre y cuando src sea una secuencia utf8 válida.
		if !utf8.RuneStart(m.src[offset]) {
			return Position{}, ErrByteOffsetNotAtRuneBoundary
		}
	}

	line := 0

	lo := 0
	hi := len(m.lineStarts) - 1

	for lo <= hi {
		mid := lo + (hi-lo)/2

		if ByteOffset(m.lineStarts[mid]) <= offset {
			line = mid
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}

	// :PERF: Esto recuenta runas desde el inicio de la línea en cada consulta.
	//  Si el formateo de regiones se calienta, guardar índices byte-runa por
	//  línea o una caché pequeña.

	col := 0
	for pos := ByteOffset(m.lineStarts[line]); pos < offset; col++ {
		// Note: Ya que src es una secuencia utf8 válida, y que hemos validado que
		// offset es el comienzo de un caracter es imposible que obtengamos RuneError.
		_, sz := utf8.DecodeRune(m.src[pos:offset])
		pos += ByteOffset(sz)
	}

	return NewPosition(line+1, col+1), nil
}

// toRegion convierte un rango de bytes en un rango de posiciones dentro del
// documento.
//
// Los offsets situados fuera de los límites del documento se ajustan al rango válido.
func (m lineMap) toRegion(r Range) (Region, error) {

	start, err := m.toPosition(r.Start())
	if err != nil {
		return Region{}, err
	}

	end, err := m.toPosition(r.End())
	if err != nil {
		return Region{}, err
	}

	return NewRegion(start, end)
}
