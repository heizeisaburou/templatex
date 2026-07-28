package rnge

import (
	"fmt"
	"reflect"
)

// Index representa los tipos enteros que pueden utilizarse como índices.
type Index interface {
	~int
}

// Rnge representa un rango delimitado por los índices start y end.
//
// Reglas de start y end:
//   - start debe ser mayor o igual que 0.
//   - start debe ser menor o igual que end.
//
// Notar:
//   - Si start y end son iguales, el rango está vacío.
//   - Cuando el rango está vacío, start y end conservan su valor porque indican
//     su posición y pueden representar un punto de inserción o la posición de
//     un cursor.
type Rnge[I Index] struct {
	start I
	end   I
}

// New crea un nuevo [Rnge] delimitado por start y end.
//
// Entra en pánico si start es negativo o end es menor que start.
func New[I Index](start, end I) Rnge[I] {
	if start < 0 || end < start {
		panic(fmt.Sprintf("invalid range: %d:%d", start, end))
	}

	return Rnge[I]{
		start: start,
		end:   end,
	}
}

// Start devuelve el índice inicial del rango.
func (s Rnge[I]) Start() I {
	return s.start
}

// End devuelve el índice final del rango.
func (s Rnge[I]) End() I {
	return s.end
}

// SetStart establece el índice inicial del rango.
//
// Entra en pánico si start es negativo.
func (s *Rnge[I]) SetStart(start I) {
	if start < 0 {
		panic(fmt.Sprintf("invalid range: %d:%d", start, s.end))
	}

	s.start = start
}

// SetEnd establece el índice final del rango.
//
// Entra en pánico si end es menor que el índice inicial.
func (s *Rnge[I]) SetEnd(end I) {
	if end < s.start {
		panic(fmt.Sprintf("invalid range: %d:%d", s.start, end))
	}
	s.end = end
}

// Set establece los índices inicial y final del rango.
//
// Entra en pánico si start es negativo o end es menor que start.
// Para asignar otro rango puede utilizarse [Rnge.SetRnge].
func (s *Rnge[I]) Set(start, end I) {
	if start < 0 || end < start {
		panic(fmt.Sprintf("invalid range: %d:%d", start, end))
	}

	s.start = start
	s.end = end
}

// SetRnge asigna a s los índices de rng.
//
// Para establecer los índices por separado puede utilizarse [Rnge.Set].
func (s *Rnge[I]) SetRnge(rng Rnge[I]) {
	s.start = rng.start
	s.end = rng.end
}

// Len devuelve la longitud del rango.
func (s Rnge[I]) Len() I {
	return s.end - s.start
}

// Empty indica si el rango está vacío.
func (s Rnge[I]) Empty() bool {
	return (s.end - s.start) == 0
}

// Bounds devuelve el rango como una "tupla" de dos elementos individuales.
func (s Rnge[I]) Bounds() (start, end I) {
	return s.start, s.end
}

// Array devuelve el rango como un array de dos elementos.
func (s Rnge[I]) Array() [2]I {
	return [2]I{s.start, s.end}
}

// String devuelve el rango con el formato "start:end".
// Se utiliza al formatear el rango mediante %s o %v.
func (s Rnge[I]) String() string {
	return fmt.Sprintf("%d:%d", s.start, s.end)
}

// GoString devuelve una representación explícita del tipo y sus campos.
// Se utiliza al formatear el rango mediante %#v.
func (s Rnge[I]) GoString() string {
	structName := reflect.TypeFor[Rnge[I]]().String()
	return fmt.Sprintf("%s{start: %d, end: %d}", structName, s.start, s.end)
}
