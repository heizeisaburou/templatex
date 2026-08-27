package document

import "github.com/heizeisaburou/templatex/pkg/cursor"

type Cursor struct {
	cur *cursor.Cursor[byte, ByteOffset]
}

func NewCursor(doc Document) *Cursor {
	return &Cursor{
		cur: cursor.New[byte, ByteOffset](doc.src),
	}
}

// NewCursorAt crea un cursor situado en position.
//
// Devuelve ErrOutOfBounds cuando position no pertenece al intervalo
// [0, doc.Len()].
func NewCursorAt(doc Document, position ByteOffset) (*Cursor, error) {
	cur, err := cursor.NewAt(doc.src, position)
	if err != nil {
		return nil, outOfBounds(err)
	}

	return &Cursor{cur: cur}, nil
}

// Len devuelve la longitud del documento que recorre el cursor
func (c Cursor) Len() ByteOffset {
	return c.cur.Len()
}

// Pos devuelve la posición actual del cursor
func (c Cursor) Pos() ByteOffset {
	return c.cur.Pos()
}

// EOF devuelve true si se ha alcanzado el final del documento
func (c Cursor) EOF() bool {
	return c.cur.EOF()
}

// Peek devuelve el elemento situado bajo el cursor.
//
// El segundo resultado es false en EOF.
func (c Cursor) Peek() (byte, bool) {
	return c.cur.Peek()
}

// Next devuelve el elemento situado bajo el cursor y avanza una posición.
//
// El segundo resultado es false en EOF. En EOF no modifica el cursor.
func (c *Cursor) Next() (byte, bool) {
	return c.cur.Next()
}

// Seek mueve el cursor a una posición absoluta.
//
// Si position no pertenece al intervalo [0, Len()], devuelve
// ErrOutOfBounds y no modifica el cursor.
func (c *Cursor) Seek(position ByteOffset) error {
	return outOfBounds(c.cur.Seek(position))
}

// Move desplaza el cursor delta posiciones desde su posición actual.
//
// Si el destino no pertenece al intervalo [0, Len()], devuelve
// ErrOutOfBounds y no modifica el cursor. La validación se realiza antes
// de sumar para evitar un posible desbordamiento de enteros.
func (c *Cursor) Move(delta ByteOffset) error {
	return outOfBounds(c.cur.Move(delta))
}
