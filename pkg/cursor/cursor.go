package cursor

import (
	"errors"
	"fmt"
)

type Index interface {
	~int
}

// Cursor mantiene una posición dentro de una secuencia.
//
// Sus posiciones válidas pertenecen al intervalo [0, Len()]. La posición
// Len() representa EOF y no identifica ningún elemento de la secuencia.
type Cursor[T any, I Index] struct {
	src    []T
	offset I
}

var ErrPositionOutOfRange = errors.New("cursor: position out of range")

// New crea un cursor situado al principio de src.
func New[T any, I Index](src []T) *Cursor[T, I] {
	return &Cursor[T, I]{src: src}
}

// NewAt crea un cursor situado en position.
//
// Devuelve ErrPositionOutOfRange cuando position no pertenece al intervalo
// [0, len(src)].
func NewAt[T any, I Index](src []T, position I) (*Cursor[T, I], error) {
	c := New[T, I](src)
	if err := c.Seek(position); err != nil {
		return &Cursor[T, I]{}, err
	}

	return c, nil
}

func (c Cursor[T, I]) Len() I {
	return I(len(c.src))
}

func (c Cursor[T, I]) Pos() I {
	return c.offset
}

func (c Cursor[T, I]) EOF() bool {
	return c.offset >= c.Len()
}

// Peek devuelve el elemento situado bajo el cursor.
//
// El segundo resultado es false en EOF.
func (c Cursor[T, I]) Peek() (T, bool) {
	if c.EOF() {
		var zero T
		return zero, false
	}

	return c.src[int(c.offset)], true
}

// Next devuelve el elemento situado bajo el cursor y avanza una posición.
//
// El segundo resultado es false en EOF. En EOF no modifica el cursor.
func (c *Cursor[T, I]) Next() (T, bool) {
	value, ok := c.Peek()
	if !ok {
		var zero T
		return zero, false
	}

	c.offset++
	return value, true
}

// Seek mueve el cursor a una posición absoluta.
//
// Si position no pertenece al intervalo [0, Len()], devuelve
// ErrPositionOutOfRange y no modifica el cursor.
func (c *Cursor[T, I]) Seek(position I) error {
	if position < 0 || position > c.Len() {
		return fmt.Errorf(
			"%w: %d not in [0,%d]",
			ErrPositionOutOfRange,
			position,
			c.Len(),
		)
	}

	c.offset = position
	return nil
}

// Move desplaza el cursor delta posiciones desde su posición actual.
//
// Si el destino no pertenece al intervalo [0, Len()], devuelve
// ErrPositionOutOfRange y no modifica el cursor. La validación se realiza antes
// de sumar para evitar un posible desbordamiento de enteros.
func (c *Cursor[T, I]) Move(delta I) error {
	minDelta := -c.offset
	maxDelta := c.Len() - c.offset

	if delta < minDelta || delta > maxDelta {
		return fmt.Errorf(
			"%w: cannot move %d from position %d; valid delta [%d,%d]",
			ErrPositionOutOfRange,
			delta,
			c.offset,
			minDelta,
			maxDelta,
		)
	}

	c.offset += delta
	return nil
}
