package document

import (
	"fmt"

	"github.com/heizeisaburou/templatex/pkg/rnge"
)

// Region representa un rango semiabierto de posiciones de un documento.
type Region struct {
	rng rnge.Range[Position]
}

// NewRegion crea la región semiabierta [start, end).
func NewRegion(start, end Position) Region {
	return Region{
		rng: rnge.New(start, end),
	}
}

// Start devuelve el índice inicial del rango.
func (r Region) Start() Position {
	return r.rng.Start()
}

// End devuelve el índice final del rango.
func (r Region) End() Position {
	return r.rng.End()
}

// SetStart establece el índice inicial del rango.
//
// Entra en pánico si start es mayor que el índice final.
func (r *Region) SetStart(start Position) {
	r.rng.SetStart(start)
}

// SetEnd establece el índice final del rango.
//
// Entra en pánico si end es menor que el índice inicial.
func (r *Region) SetEnd(end Position) {
	r.rng.SetEnd(end)
}

// Set establece los índices inicial y final del rango.
//
// Entra en pánico si end es menor que start.
// Para asignar otro rango puede utilizarse [Range.SetRange].
func (r *Region) Set(start, end Position) {
	r.rng.Set(start, end)
}

// SetRange asigna a r los índices de other.
//
// Para establecer los índices por separado puede utilizarse [Range.Set].
func (r *Region) SetRange(other Region) {
	r.rng.SetRange(other.rng)
}

// Empty indica si el rango está vacío.
func (r Region) Empty() bool {
	return r.rng.Empty()
}

// Bounds devuelve el rango como una "tupla" de dos elementos individuales.
func (r Region) Bounds() (start, end Position) {
	return r.rng.Bounds()
}

// Array devuelve el rango como un array de dos elementos.
func (r Region) Array() [2]Position {
	return [2]Position{r.Start(), r.End()}
}

// String devuelve el rango con el formato "[start, end)".
// Se utiliza al formatear el rango mediante %s o %v.
func (r Region) String() string {
	return fmt.Sprintf("[%v, %v)", r.Start(), r.End())
}

// GoString devuelve una representación explícita del tipo y sus campos.
// Se utiliza al formatear el rango mediante %#v.
func (r Region) GoString() string {
	return fmt.Sprintf("%T{start: %v, end: %v}", r, r.Start(), r.End())
}
