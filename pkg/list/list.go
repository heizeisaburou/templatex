package list

import (
	"fmt"
	"reflect"
	"strings"
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
// Retorna el índice (ID) en el que fue insertado el elemento.
func (l *List[I, T]) Push(item T) I {
	nextID := I(len(l.items))
	l.items = append(l.items, item)
	return nextID
}

// Set asigna un elemento a un índice específico existente.
// Hace panic si el índice es menor a 0 o mayor/igual al tamaño actual de la lista.
func (l *List[I, T]) Set(id I, item T) {
	if id < 0 || int(id) >= len(l.items) {
		panic(fmt.Sprintf("list index out of range: %d", id))
	}

	l.items[id] = item
}

// Truncate reduce la longitud de la lista hasta end.
//
// Antes de reducir el slice, pone a valor cero los elementos descartados.
// Esto evita que el array subyacente siga reteniendo referencias a objetos
// que podrían ser reclamados por el garbage collector.
func (l *List[I, T]) Truncate(end I) {
	if end < 0 || int(end) > len(l.items) {
		panic(fmt.Sprintf("invalid list truncate: %d", end))
	}

	n := int(end)

	clear(l.items[n:])
	l.items = l.items[:n]
}

// ReplaceRange elimina los elementos en el rango [start:end] y los sustituye
// por los elementos proporcionados en 'repl'. Maneja dinámicamente el
// redimensionamiento del slice subyacente.
func (l *List[I, T]) ReplaceRange(start, end I, repl ...T) {
	if start < 0 || end < start || int(end) > len(l.items) {
		panic(fmt.Sprintf("invalid list range: %d:%d", start, end))
	}

	removed := int(end - start)
	inserted := len(repl)

	// Camino optimizado: si quitamos e insertamos la misma cantidad,
	// evitamos reasignar memoria y solo copiamos.
	if removed == inserted {
		copy(l.items[start:end], repl)
		return
	}

	// Camino dinámico: Creamos un nuevo slice con la capacidad exacta necesaria
	// para evitar múltiples reasignaciones durante los appends.
	next := make([]T, 0, len(l.items)-removed+inserted)
	next = append(next, l.items[:start]...)
	next = append(next, repl...)
	next = append(next, l.items[end:]...)
	l.items = next
}

// At retorna el elemento ubicado en el índice especificado.
// Hace panic si el índice está fuera de los límites de la lista.
func (l List[I, T]) At(id I) T {
	if id < 0 || int(id) >= len(l.items) {
		panic(fmt.Sprintf("invalid index: %d", id))
	}

	return l.items[id]
}

// Range devuelve una copia de los elementos dentro del rango [start:end].
// Retorna un nuevo slice, asegurando que las modificaciones externas no
// afecten el estado interno de la lista.
func (l List[I, T]) Range(start, end I) []T {
	if start < 0 || end < start || int(end) > len(l.items) {
		panic(fmt.Sprintf("invalid list range: %d:%d", start, end))
	}

	dest := make([]T, end-start)
	copy(dest, l.items[start:end])

	return dest
}

// Slice retorna una copia completa de todos los elementos de la lista.
// Se utiliza append sobre un slice nil para garantizar la aislación de memoria.
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
