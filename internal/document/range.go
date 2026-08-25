package document

import (
	"fmt"

	"github.com/heizeisaburou/templatex/pkg/rnge"
)

// Range representa un rango semiabierto de ByteOffset.
type Range struct {
	rng rnge.Range[ByteOffset]
}

// NewRange crea el rango semiabierto [start, end).
func NewRange(start, end ByteOffset) (Range, error) {
	return Range{
		rng: rnge.New(start, end),
	}, nil
}

// Len devuelve la longitud del rango.
func (r Range) Len() ByteOffset {
	return r.End() - r.Start()
}

// Start devuelve el índice inicial del rango.
func (r Range) Start() ByteOffset {
	return r.rng.Start()
}

// End devuelve el índice final del rango.
func (r Range) End() ByteOffset {
	return r.rng.End()
}

// SetStart establece el índice inicial del rango.
//
// Entra en pánico si start es mayor que el índice final.
func (r *Range) SetStart(start ByteOffset) {
	r.rng.SetStart(start)
}

// SetEnd establece el índice final del rango.
//
// Entra en pánico si end es menor que el índice inicial.
func (r *Range) SetEnd(end ByteOffset) {
	r.rng.SetEnd(end)
}

// Set establece los índices inicial y final del rango.
//
// Entra en pánico si end es menor que start.
// Para asignar otro rango puede utilizarse [Range.SetRange].
func (r *Range) Set(start, end ByteOffset) {
	r.rng.Set(start, end)
}

// SetRange asigna a r los índices de other.
//
// Para establecer los índices por separado puede utilizarse [Range.Set].
func (r *Range) SetRange(other Range) {
	r.rng.SetRange(other.rng)
}

// Empty indica si el rango está vacío.
func (r Range) Empty() bool {
	return r.rng.Empty()
}

// Bounds devuelve el rango como una "tupla" de dos elementos individuales.
func (r Range) Bounds() (start, end ByteOffset) {
	return r.rng.Bounds()
}

// Array devuelve el rango como un array de dos elementos.
func (r Range) Array() [2]ByteOffset {
	return [2]ByteOffset{r.Start(), r.End()}
}

// String devuelve el rango con el formato "[start, end)".
// Se utiliza al formatear el rango mediante %s o %v.
func (r Range) String() string {
	return fmt.Sprintf("[%v, %v)", r.Start(), r.End())
}

// GoString devuelve una representación explícita del tipo y sus campos.
// Se utiliza al formatear el rango mediante %#v.
func (r Range) GoString() string {
	return fmt.Sprintf("%T{start: %v, end: %v}", r, r.Start(), r.End())
}
