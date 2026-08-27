package document

import (
	"errors"
	"fmt"
)

// Los centinelas de este paquete son los únicos que sus usuarios necesitan
// conocer.
//
// Las operaciones que se apoyan en pkg/rnge o pkg/cursor envuelven el error de
// origen con el centinela equivalente de aquí, de modo que
// errors.Is(err, ErrOutOfBounds) funciona sin importar en qué paquete interno
// se haya originado el fallo. El error original permanece en la cadena para
// conservar su mensaje, pero no forma parte del contrato: no dependas de él.
var (
	ErrInvalidUTF8                 = errors.New("invalid utf8 sequence")
	ErrOutOfBounds                 = errors.New("out of bounds")
	ErrInvalidRange                = errors.New("invalid range")
	ErrByteOffsetNotAtRuneBoundary = errors.New("el offset de bytes no apunta al inicio de un carácter codificado en UTF-8")
)

// outOfBounds envuelve err con ErrOutOfBounds conservándolo en la cadena.
//
// Devuelve nil si err es nil, para poder usarse directamente sobre el
// resultado de la operación envuelta.
func outOfBounds(err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: %w", ErrOutOfBounds, err)
}

// invalidRange envuelve err con ErrInvalidRange conservándolo en la cadena.
//
// Devuelve nil si err es nil, para poder usarse directamente sobre el
// resultado de la operación envuelta.
func invalidRange(err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: %w", ErrInvalidRange, err)
}
