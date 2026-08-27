package document

import (
	"errors"
	"testing"
)

// mustCursorAt construye un Cursor situado en position o aborta la prueba.
func mustCursorAt(t *testing.T, d Document, position ByteOffset) *Cursor {
	t.Helper()

	c, err := NewCursorAt(d, position)
	if err != nil {
		t.Fatalf("NewCursorAt(%v, %d) = %v; want nil", d, position, err)
	}

	return c
}

func TestNewCursor(t *testing.T) {
	tests := []struct {
		name string
		src  []byte
	}{
		{name: "documento vacío", src: []byte("")},
		{name: "documento de un byte", src: []byte("a")},
		{name: "documento normal", src: []byte("hola")},
		{name: "documento multibyte", src: []byte("a界b")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := mustNew(t, tt.src)
			c := NewCursor(d)

			if got := c.Pos(); got != 0 {
				t.Errorf("Pos() = %d; want = 0", got)
			}

			if got := c.Len(); got != d.Len() {
				t.Errorf("Len() = %d; want = %d", got, d.Len())
			}
		})
	}
}

func TestNewCursorAt(t *testing.T) {
	tests := []struct {
		name     string
		src      []byte
		position ByteOffset
	}{
		{name: "inicio del documento", src: []byte("hola"), position: 0},
		{name: "posición interior", src: []byte("hola"), position: 2},
		{name: "EOF", src: []byte("hola"), position: 4},
		{name: "documento vacío", src: []byte(""), position: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := mustNew(t, tt.src)
			c := mustCursorAt(t, d, tt.position)

			if got := c.Pos(); got != tt.position {
				t.Errorf("Pos() = %d; want = %d", got, tt.position)
			}
		})
	}
}

func TestNewCursorAtErrors(t *testing.T) {
	tests := []struct {
		name     string
		src      []byte
		position ByteOffset
	}{
		{name: "posición negativa", src: []byte("hola"), position: -1},
		{name: "posición más allá de EOF", src: []byte("hola"), position: 5},
		{name: "documento vacío con posición no nula", src: []byte(""), position: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := mustNew(t, tt.src)

			got, err := NewCursorAt(d, tt.position)

			if !errors.Is(err, ErrOutOfBounds) {
				t.Errorf(
					"NewCursorAt(%q, %d) err = %v; want = %v",
					tt.src,
					tt.position,
					err,
					ErrOutOfBounds,
				)
			}

			if got != nil {
				t.Errorf(
					"NewCursorAt(%q, %d) = %v; want = nil",
					tt.src,
					tt.position,
					got,
				)
			}
		})
	}
}

func TestCursorEOF(t *testing.T) {
	tests := []struct {
		name     string
		src      []byte
		position ByteOffset
		want     bool
	}{
		{name: "documento vacío", src: []byte(""), position: 0, want: true},
		{name: "inicio del documento", src: []byte("hola"), position: 0, want: false},
		{name: "último byte", src: []byte("hola"), position: 3, want: false},
		{name: "final del documento", src: []byte("hola"), position: 4, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := mustNew(t, tt.src)
			c := mustCursorAt(t, d, tt.position)

			if got := c.EOF(); got != tt.want {
				t.Errorf("EOF() = %t; want = %t", got, tt.want)
			}
		})
	}
}

func TestCursorPeek(t *testing.T) {
	tests := []struct {
		name     string
		src      []byte
		position ByteOffset
		want     byte
		wantOK   bool
	}{
		{name: "primer byte", src: []byte("hola"), position: 0, want: 'h', wantOK: true},
		{name: "byte interior", src: []byte("hola"), position: 2, want: 'l', wantOK: true},
		{name: "último byte", src: []byte("hola"), position: 3, want: 'a', wantOK: true},
		{name: "EOF", src: []byte("hola"), position: 4, want: 0, wantOK: false},
		{name: "documento vacío", src: []byte(""), position: 0, want: 0, wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := mustNew(t, tt.src)
			c := mustCursorAt(t, d, tt.position)

			got, ok := c.Peek()

			if got != tt.want || ok != tt.wantOK {
				t.Errorf(
					"Peek() = (%q, %t); want = (%q, %t)",
					got,
					ok,
					tt.want,
					tt.wantOK,
				)
			}

			// Peek no debe mover el cursor.
			if pos := c.Pos(); pos != tt.position {
				t.Errorf("Pos() = %d; want = %d", pos, tt.position)
			}
		})
	}
}

func TestCursorNext(t *testing.T) {
	t.Run("recorre el documento byte a byte", func(t *testing.T) {
		src := []byte("hola")
		d := mustNew(t, src)
		c := NewCursor(d)

		for i, want := range src {
			got, ok := c.Next()

			if !ok {
				t.Fatalf("Next() ok = false en la posición %d; want = true", i)
			}

			if got != want {
				t.Errorf("Next() = %q; want = %q", got, want)
			}

			if pos := c.Pos(); pos != ByteOffset(i+1) {
				t.Errorf("Pos() = %d; want = %d", pos, i+1)
			}
		}

		if !c.EOF() {
			t.Errorf("EOF() = false; want = true tras consumir %q", src)
		}
	})

	t.Run("un documento multibyte se recorre byte a byte", func(t *testing.T) {
		src := []byte("a界b")
		d := mustNew(t, src)
		c := NewCursor(d)

		var got []byte
		for {
			b, ok := c.Next()
			if !ok {
				break
			}
			got = append(got, b)
		}

		if len(got) != len(src) {
			t.Errorf("Next() recorrió %d bytes; want = %d", len(got), len(src))
		}
	})

	t.Run("en EOF no avanza", func(t *testing.T) {
		d := mustNew(t, []byte("hola"))
		c := mustCursorAt(t, d, 4)

		got, ok := c.Next()

		if ok {
			t.Errorf("Next() = (%q, true); want ok = false", got)
		}

		if pos := c.Pos(); pos != 4 {
			t.Errorf("Pos() = %d; want = 4", pos)
		}
	})
}

func TestCursorSeek(t *testing.T) {
	tests := []struct {
		name     string
		src      []byte
		position ByteOffset
	}{
		{name: "al inicio", src: []byte("hola"), position: 0},
		{name: "a una posición interior", src: []byte("hola"), position: 2},
		{name: "a EOF", src: []byte("hola"), position: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := mustNew(t, tt.src)
			c := NewCursor(d)

			if err := c.Seek(tt.position); err != nil {
				t.Fatalf("Seek(%d) = %v; want nil", tt.position, err)
			}

			if got := c.Pos(); got != tt.position {
				t.Errorf("Pos() = %d; want = %d", got, tt.position)
			}
		})
	}
}

func TestCursorSeekErrors(t *testing.T) {
	tests := []struct {
		name     string
		src      []byte
		position ByteOffset
	}{
		{name: "posición negativa", src: []byte("hola"), position: -1},
		{name: "posición más allá de EOF", src: []byte("hola"), position: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := mustNew(t, tt.src)
			c := mustCursorAt(t, d, 2)

			err := c.Seek(tt.position)

			if !errors.Is(err, ErrOutOfBounds) {
				t.Errorf(
					"Seek(%d) err = %v; want = %v",
					tt.position,
					err,
					ErrOutOfBounds,
				)
			}

			// Un Seek fallido no debe modificar el cursor.
			if got := c.Pos(); got != 2 {
				t.Errorf("Pos() = %d; want = 2", got)
			}
		})
	}
}

func TestCursorMove(t *testing.T) {
	tests := []struct {
		name  string
		start ByteOffset
		delta ByteOffset
		want  ByteOffset
	}{
		{name: "hacia adelante", start: 0, delta: 2, want: 2},
		{name: "hacia atrás", start: 3, delta: -2, want: 1},
		{name: "delta nulo", start: 2, delta: 0, want: 2},
		{name: "hasta el inicio", start: 3, delta: -3, want: 0},
		{name: "hasta EOF", start: 1, delta: 3, want: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := mustNew(t, []byte("hola"))
			c := mustCursorAt(t, d, tt.start)

			if err := c.Move(tt.delta); err != nil {
				t.Fatalf("Move(%d) = %v; want nil", tt.delta, err)
			}

			if got := c.Pos(); got != tt.want {
				t.Errorf("Pos() = %d; want = %d", got, tt.want)
			}
		})
	}
}

func TestCursorMoveErrors(t *testing.T) {
	tests := []struct {
		name  string
		start ByteOffset
		delta ByteOffset
	}{
		{name: "por debajo del inicio", start: 1, delta: -2},
		{name: "más allá de EOF", start: 2, delta: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := mustNew(t, []byte("hola"))
			c := mustCursorAt(t, d, tt.start)

			err := c.Move(tt.delta)

			if !errors.Is(err, ErrOutOfBounds) {
				t.Errorf(
					"Move(%d) err = %v; want = %v",
					tt.delta,
					err,
					ErrOutOfBounds,
				)
			}

			// Un Move fallido no debe modificar el cursor.
			if got := c.Pos(); got != tt.start {
				t.Errorf("Pos() = %d; want = %d", got, tt.start)
			}
		})
	}
}
