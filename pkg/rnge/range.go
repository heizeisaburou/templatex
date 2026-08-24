package rnge

import (
	"errors"
	"fmt"
)

// ErrInvalidRange indica que el rango es inválido porque start es mayor que end.
var ErrInvalidRange = errors.New("invalid range")

// Ordered representa un valor que puede compararse con otro del mismo tipo.
type Ordered[T any] interface {
	Compare(T) int
}

// Range representa el rango semiabierto [start, end).
//
// Reglas de start y end:
//   - start debe ser menor o igual que end.
//
// Notar:
//   - Si start y end son iguales, el rango está vacío.
//   - Cuando el rango está vacío, start y end conservan su valor porque indican
//     su posición y pueden representar un punto de inserción o la posición de
//     un cursor.
type Range[I Ordered[I]] struct {
	start I
	end   I
}

// New crea un nuevo [Range] delimitado por start y end.
//
// Devuelve error si end es menor que start.
func New[I Ordered[I]](start, end I) (Range[I], error) {
	if start.Compare(end) > 0 {
		return Range[I]{}, fmt.Errorf("%w: [%v, %v)", ErrInvalidRange, start, end)
	}

	return Range[I]{
		start: start,
		end:   end,
	}, nil
}

// Start devuelve el índice inicial del rango.
func (r Range[I]) Start() I {
	return r.start
}

// End devuelve el índice final del rango.
func (r Range[I]) End() I {
	return r.end
}

// SetStart establece el índice inicial del rango.
//
// Devuelve error si start es mayor que el índice final.
func (r *Range[I]) SetStart(start I) error {
	if start.Compare(r.end) > 0 {
		return fmt.Errorf("%w: [%v, %v)", ErrInvalidRange, start, r.end)
	}

	r.start = start
	return nil
}

// SetEnd establece el índice final del rango.
//
// Devuelve error si end es menor que el índice inicial.
func (r *Range[I]) SetEnd(end I) error {
	if r.start.Compare(end) > 0 {
		return fmt.Errorf("%w: [%v, %v)", ErrInvalidRange, r.start, end)
	}
	r.end = end
	return nil
}

// Set establece los índices inicial y final del rango.
//
// Devuelve error si end es menor que start.
// Para asignar otro rango puede utilizarse [Range.SetRange].
func (r *Range[I]) Set(start, end I) error {
	if start.Compare(end) > 0 {
		return fmt.Errorf("%w: [%v, %v)", ErrInvalidRange, start, end)
	}

	r.start = start
	r.end = end
	return nil
}

// SetRange asigna a r los índices de other.
//
// Para establecer los índices por separado puede utilizarse [Range.Set].
func (r *Range[I]) SetRange(other Range[I]) {
	r.start = other.start
	r.end = other.end
}

// Empty indica si el rango está vacío.
func (r Range[I]) Empty() bool {
	return r.start.Compare(r.end) == 0
}

// Bounds devuelve el rango como una "tupla" de dos elementos individuales.
func (r Range[I]) Bounds() (start, end I) {
	return r.start, r.end
}

// Array devuelve el rango como un array de dos elementos.
func (r Range[I]) Array() [2]I {
	return [2]I{r.start, r.end}
}

// String devuelve el rango con el formato "[start, end)".
// Se utiliza al formatear el rango mediante %s o %v.
func (r Range[I]) String() string {
	return fmt.Sprintf("[%v, %v)", r.start, r.end)
}

// GoString devuelve una representación explícita del tipo y sus campos.
// Se utiliza al formatear el rango mediante %#v.
func (r Range[I]) GoString() string {
	return fmt.Sprintf("%T{start: %v, end: %v}", r, r.start, r.end)
}
