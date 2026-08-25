package list

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// :TODO: Saburou: ya no queda ningún panic, falta terminar el resto de tests

var (
	ErrOutOfBounds = errors.New("out of bounds")
)

// Index define el tipo permitido para los índices de la lista.
// Acepta cualquier tipo subyacente que sea un int (int, int32, int64, etc.).
type Index interface {
	~int
}

// List representa una estructura de datos de lista dinámica y genérica.
// Utiliza genéricos con una plantilla I (para el índice) y T (para cualquier tipo de dato).
// Mantiene su estado interno encapsulado para evitar mutaciones externas indeseadas.
type List[I Index, T any] struct {
	items []T
}

// New inicializa y devuelve una nueva instancia de List.
// Toma un número variádico de elementos y realiza una copia segura en memoria
// para garantizar que referencias externas no puedan alterar la lista inicial.
func New[I Index, T any](items ...T) List[I, T] {
	return List[I, T]{
		items: append([]T(nil), items...),
	}
}

// Push añade un nuevo elemento al final de la lista.
func (l *List[I, T]) Push(item T) {
	l.items = append(l.items, item)
}

// Set asigna un elemento a un índice específico existente.
//
// Retorna un error que envuelve ErrOutOfBounds si el índice es menor a 0
// o mayor/igual al tamaño actual de la lista. En ese caso la lista no se modifica.
func (l *List[I, T]) Set(id I, item T) error {
	length := l.Len()
	if id < 0 || id >= length {
		return fmt.Errorf("%w: [len=%d, id=%d]", ErrOutOfBounds, length, id)
	}

	l.items[id] = item

	return nil
}

// Truncate reduce la longitud de la lista hasta end, descartando los
// elementos a partir de esa posición.
//
// Antes de reducir el slice, pone a valor cero los elementos descartados.
// Esto evita que el array subyacente siga reteniendo referencias a objetos
// que podrían ser reclamados por el garbage collector.
//
// Retorna un error que envuelve ErrOutOfBounds si end es menor a 0 o mayor
// que la longitud actual. En ese caso la lista no se modifica.
func (l *List[I, T]) Truncate(end I) error {
	if end < 0 || end > l.Len() {
		return fmt.Errorf("%w: [len=%d, end=%d]", ErrOutOfBounds, l.Len(), end)
	}

	n := int(end)

	clear(l.items[n:])
	l.items = l.items[:n]

	return nil
}

// ReplaceRange elimina los elementos del rango [start, end) y los sustituye
// por los elementos proporcionados en 'repl'. Maneja dinámicamente el
// redimensionamiento del slice subyacente.
//
// Retorna un error que envuelve ErrOutOfBounds si el rango es inválido
// (start < 0, end < start o end > Len). En ese caso la lista no se modifica.
func (l *List[I, T]) ReplaceRange(start, end I, repl ...T) error {
	if start < 0 || end < start || end > l.Len() {
		return fmt.Errorf(
			"%w: [len=%d, start=%d, end=%d]",
			ErrOutOfBounds,
			l.Len(),
			start,
			end,
		)
	}

	removed := int(end - start)
	inserted := len(repl)

	// Camino optimizado: si quitamos e insertamos la misma cantidad,
	// evitamos reasignar memoria y solo copiamos.
	if removed == inserted {
		copy(l.items[start:end], repl)
		return nil
	}

	// Camino dinámico: Creamos un nuevo slice con la capacidad exacta necesaria
	// para evitar múltiples reasignaciones durante los appends.
	next := make([]T, 0, len(l.items)-removed+inserted)
	next = append(next, l.items[:start]...)
	next = append(next, repl...)
	next = append(next, l.items[end:]...)
	l.items = next

	return nil
}

// At retorna el elemento ubicado en el índice especificado.
//
// El segundo valor sigue el convenio ok de Go: es true si el índice existe.
// Si está fuera de los límites de la lista retorna el valor cero de T y false.
func (l List[I, T]) At(id I) (T, bool) {
	if id < 0 || id >= l.Len() {
		var zero T
		return zero, false
	}

	return l.items[id], true
}

// Range devuelve una copia de los elementos dentro del rango [start, end).
// Retorna un nuevo slice, asegurando que las modificaciones externas no
// afecten el estado interno de la lista.
//
// Retorna un error que envuelve ErrOutOfBounds si el rango es inválido
// (start < 0, end < start o end > Len). En ese caso el slice devuelto es nil.
func (l List[I, T]) Range(start, end I) ([]T, error) {
	if start < 0 || end < start || end > l.Len() {
		return nil, fmt.Errorf(
			"%w: [len=%d, start=%d, end=%d]",
			ErrOutOfBounds,
			l.Len(),
			start,
			end,
		)
	}

	dest := make([]T, end-start)
	copy(dest, l.items[start:end])

	return dest, nil
}

// Slice retorna una copia completa de todos los elementos de la lista.
// Se utiliza append sobre un slice nil para garantizar la aislación de memoria:
// la copia es superficial, así que los punteros almacenados siguen compartidos.
func (l List[I, T]) Slice() []T {
	return append([]T(nil), l.items...)
}

// Len retorna la cantidad actual de elementos contenidos en la lista.
func (l List[I, T]) Len() I {
	return I(len(l.items))
}

// Empty evalúa si la lista carece de elementos, retornando true si el tamaño es 0.
func (l List[I, T]) Empty() bool {
	return len(l.items) == 0
}

// String implementa la interfaz fmt.Stringer para representar la lista
// en un formato de texto legible (ej: "{1, 2, 3}").
func (l List[I, T]) String() string {
	var b strings.Builder

	b.WriteString("{")

	for i, item := range l.items {
		if i > 0 {
			b.WriteString(", ")
		}
		// fmt.Fprint escribe directamente en el buffer, es seguro para genéricos
		fmt.Fprint(&b, item)
	}

	b.WriteString("}")

	return b.String()
}

// GoString implementa la interfaz fmt.GoStringer.
// Retorna la representación de la lista incluyendo el nombre del tipo de la estructura,
// reutilizando la lógica de formateo de la función String().
func (l List[I, T]) GoString() string {
	structName := reflect.TypeFor[List[I, T]]().String()

	// Concatenamos el nombre del struct con el resultado de String()
	// cumpliendo el principio DRY (Don't Repeat Yourself).
	return structName + l.String()
}
