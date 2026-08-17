// Package clamp contiene utilidades para limitar índices y posiciones
// a los límites válidos de una secuencia.
package clamp

// Index representa un tipo entero que puede utilizarse como índice o posición.
//
// La restricción ~int permite utilizar tanto int como tipos definidos cuyo tipo
// subyacente sea int.
type Index interface {
	~int
}

// ClampIndex limita index al rango de índices válidos de una secuencia.
//
// length representa el número de elementos de la secuencia. Por tanto, el resultado
// pertenece al intervalo:
//
//	[0, length)
//
// Reglas:
//   - Si index es negativo, devuelve 0. 
//   - Si index es igual o superior a length, devuelve length - 1.
//
// La función provoca panic cuando length es menor o igual que cero, ya que una
// secuencia vacía no contiene ningún índice válido.
func ClampIndex[T Index](length, index T) T {
	if length <= 0 {
		panic("clamp: cannot clamp index of empty sequence")
	}

	if index < 0 {
		return 0
	}

	if index >= length {
		return length - 1
	}

	return index
}

// IsIndexValid indica si index identifica un elemento existente.
//
// Un índice es válido cuando pertenece al intervalo:
//
//	[0, length)
func IsIndexValid[T Index](length, index T) bool {
	return index >= 0 && index < length
}

// ClampPosition limita position al rango de posiciones válidas de un cursor.
//
// A diferencia de un índice, una posición puede ser igual a length. Esa
// posición representa el punto situado inmediatamente después del último
// elemento, normalmente utilizado como EOF.
//
// El resultado pertenece al intervalo:
//
//	[0, length]
//
// Si position es negativa, devuelve 0. Si es superior a length, devuelve
// length.
//
// La función provoca panic cuando length es negativo.
func ClampPosition[T Index](length, position T) T {
	if length < 0 {
		panic("clamp: negative sequence length")
	}

	if position < 0 {
		return 0
	}

	if position > length {
		return length
	}

	return position
}

// IsPositionValid indica si position es una posición válida de cursor.
//
// Una posición es válida cuando pertenece al intervalo:
//
//	[0, length]
//
// El valor length es válido y representa la posición posterior al último
// elemento.
func IsPositionValid[T Index](length, position T) bool {
	return position >= 0 && position <= length
}
