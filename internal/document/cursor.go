package document

import "github.com/heizeisaburou/templatex/pkg/cursor"

type Cursor struct {
	*cursor.Cursor[byte, ByteOffset]
}

func NewCursor(doc Document) *Cursor {
	return &Cursor{
		Cursor: cursor.New[byte, ByteOffset](doc.src),
	}
}

func NewCursorAt(doc Document, position ByteOffset) (*Cursor, error) {
	cur, err := cursor.NewAt(doc.src, position)
	if err != nil {
		return nil, err
	}

	return &Cursor{Cursor: cur}, nil
}
