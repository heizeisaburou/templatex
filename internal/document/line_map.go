package document

import (
	"unicode/utf8"
	"github.com/heizeisaburou/templatex/internal/clamp"
)

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
  ln := ByteOffset(len(src))
  offset = clamp.ClampPosition(ln, offset)

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
		_, sz := utf8.DecodeRune(src[pos:offset])
		pos += ByteOffset(sz)
	}

  return NewPosition(line, col)
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
