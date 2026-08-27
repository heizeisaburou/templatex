package document

import (
	"errors"
	"slices"
	"testing"
)

// mustRange construye un Range o aborta la prueba.
func mustRange(t *testing.T, start, end ByteOffset) Range {
	t.Helper()

	r, err := NewRange(start, end)
	if err != nil {
		t.Fatalf("NewRange(%d, %d) = %v; want nil", start, end, err)
	}

	return r
}

// mustNew construye un Document o aborta la prueba.
func mustNew(t *testing.T, src []byte) Document {
	t.Helper()

	d, err := New(src)
	if err != nil {
		t.Fatalf("New(%q) = %v; want nil", src, err)
	}

	return d
}

func TestNew(t *testing.T) {
	t.Run("secuencias utf8 válidas", func(t *testing.T) {
		tests := []struct {
			name string
			src  []byte
		}{
			{name: "fuente vacía", src: []byte{}},
			{name: "fuente nil", src: nil},
			{name: "solo ascii", src: []byte("hola")},
			{name: "multibyte", src: []byte("ab界\n🌎")},
			{name: "solo saltos de línea", src: []byte("\n\n\n")},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				d, err := New(tt.src)
				if err != nil {
					t.Fatalf("New(%q) = %v; want nil", tt.src, err)
				}

				if got, want := d.Len(), ByteOffset(len(tt.src)); got != want {
					t.Errorf("New(%q).Len() = %d; want %d", tt.src, got, want)
				}
			})
		}
	})

	t.Run("secuencias utf8 inválidas", func(t *testing.T) {
		tests := []struct {
			name string
			src  []byte
		}{
			{name: "byte de continuación suelto", src: []byte{0x80}},
			{name: "rune de 3 bytes truncada", src: []byte{0xE7, 0x95}},
			{name: "byte inválido tras ascii", src: []byte{'a', 0xFF, 'b'}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				d, err := New(tt.src)
				if !errors.Is(err, ErrInvalidUTF8) {
					t.Errorf("New(%v) error = %v; want %v", tt.src, err, ErrInvalidUTF8)
				}

				if got := d.Len(); got != 0 {
					t.Errorf("New(%v) = documento con Len() = %d; want documento cero", tt.src, got)
				}
			})
		}
	})
}

func TestDocumentLen(t *testing.T) {
	tests := []struct {
		name string
		src  []byte
		want ByteOffset
	}{
		{name: "fuente vacía", src: []byte(""), want: 0},
		{name: "ascii", src: []byte("hola"), want: 4},
		{name: "rune de 3 bytes cuenta bytes, no runas", src: []byte("界"), want: 3},
		{name: "emoji de 4 bytes cuenta bytes, no runas", src: []byte("🌎"), want: 4},
		{name: "varias líneas", src: []byte("a\nb\nc"), want: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := mustNew(t, tt.src)

			if got := d.Len(); got != tt.want {
				t.Errorf("New(%q).Len() = %d; want %d", tt.src, got, tt.want)
			}
		})
	}
}

func TestDocumentRange(t *testing.T) {
	tests := []struct {
		name  string
		src   []byte
		start ByteOffset
		end   ByteOffset
		want  []byte
	}{
		{
			name:  "primer byte",
			src:   []byte("hola"),
			start: 0,
			end:   1,
			want:  []byte("h"),
		},
		{
			name:  "documento completo",
			src:   []byte("hola"),
			start: 0,
			end:   4,
			want:  []byte("hola"),
		},
		{
			name:  "subrango interior",
			src:   []byte("hola"),
			start: 1,
			end:   3,
			want:  []byte("ol"),
		},
		{
			name:  "rango vacío en el interior",
			src:   []byte("hola"),
			start: 2,
			end:   2,
			want:  []byte(""),
		},
		{
			name:  "rango vacío en EOF",
			src:   []byte("hola"),
			start: 4,
			end:   4,
			want:  []byte(""),
		},
		{
			name:  "rango vacío en documento vacío",
			src:   []byte(""),
			start: 0,
			end:   0,
			want:  []byte(""),
		},
		{
			name:  "rune multibyte completa",
			src:   []byte("a界b"),
			start: 1,
			end:   4,
			want:  []byte("界"),
		},
		{
			name:  "bytes interiores de una rune multibyte",
			src:   []byte("界"),
			start: 1,
			end:   2,
			want:  []byte{0x95},
		},
		{
			name:  "rango que cruza un salto de línea",
			src:   []byte("ab\ncd"),
			start: 1,
			end:   4,
			want:  []byte("b\nc"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := mustNew(t, tt.src)
			rng := mustRange(t, tt.start, tt.end)

			got, err := d.Range(rng)
			if err != nil {
				t.Fatalf("Document.Range(%v) error = %v; want nil", rng, err)
			}

			if !slices.Equal(got, tt.want) {
				t.Errorf("Document.Range(%v) = %q; want %q", rng, got, tt.want)
			}
		})
	}
}

func TestDocumentRangeErrors(t *testing.T) {
	tests := []struct {
		name  string
		src   []byte
		start ByteOffset
		end   ByteOffset
	}{
		{name: "end mayor que Len()", src: []byte("hola"), start: 0, end: 5},
		{name: "rango entero fuera del documento", src: []byte("hola"), start: 10, end: 20},
		{name: "start negativo", src: []byte("hola"), start: -1, end: 2},
		{name: "documento vacío con rango no vacío", src: []byte(""), start: 0, end: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := mustNew(t, tt.src)
			rng := mustRange(t, tt.start, tt.end)

			got, err := d.Range(rng)
			if !errors.Is(err, ErrOutOfBounds) {
				t.Errorf("Document.Range(%v) error = %v; want %v", rng, err, ErrOutOfBounds)
			}

			if got != nil {
				t.Errorf("Document.Range(%v) = %q; want nil", rng, got)
			}
		})
	}
}

// TestDocumentRangeIsACopy comprueba el contrato de Range: el llamante recibe
// una copia independiente y no puede modificar el documento a través de ella.
func TestDocumentRangeIsACopy(t *testing.T) {
	src := []byte("hola")
	d := mustNew(t, src)
	rng := mustRange(t, 0, 4)

	got, err := d.Range(rng)
	if err != nil {
		t.Fatalf("Document.Range(%v) error = %v; want nil", rng, err)
	}

	got[0] = 'X'

	again, err := d.Range(rng)
	if err != nil {
		t.Fatalf("Document.Range(%v) error = %v; want nil", rng, err)
	}

	if !slices.Equal(again, []byte("hola")) {
		t.Errorf("tras mutar la copia, Document.Range(%v) = %q; want %q", rng, again, "hola")
	}
}

func TestDocumentToPosition(t *testing.T) {
	tests := []struct {
		name   string
		src    []byte
		offset ByteOffset
		want   Position
	}{
		{
			name:   "inicio del documento",
			src:    []byte("ab\ncd"),
			offset: 0,
			want:   NewPosition(1, 1),
		},
		{
			name:   "segunda columna de la primera línea",
			src:    []byte("ab\ncd"),
			offset: 1,
			want:   NewPosition(1, 2),
		},
		{
			name:   "el salto de línea pertenece a su línea",
			src:    []byte("ab\ncd"),
			offset: 2,
			want:   NewPosition(1, 3),
		},
		{
			name:   "inicio de la segunda línea",
			src:    []byte("ab\ncd"),
			offset: 3,
			want:   NewPosition(2, 1),
		},
		{
			name:   "EOF",
			src:    []byte("ab\ncd"),
			offset: 5,
			want:   NewPosition(2, 3),
		},
		{
			name:   "documento vacío",
			src:    []byte(""),
			offset: 0,
			want:   NewPosition(1, 1),
		},
		{
			name:   "las columnas cuentan runas, no bytes",
			src:    []byte("界界a"),
			offset: 6,
			want:   NewPosition(1, 3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := mustNew(t, tt.src)

			got, err := d.ToPosition(tt.offset)
			if err != nil {
				t.Fatalf("Document.ToPosition(%d) error = %v; want nil", tt.offset, err)
			}

			if got.Compare(tt.want) != 0 {
				t.Errorf("Document.ToPosition(%d) = %v; want %v", tt.offset, got, tt.want)
			}
		})
	}
}

func TestDocumentToPositionErrors(t *testing.T) {
	tests := []struct {
		name    string
		src     []byte
		offset  ByteOffset
		wantErr error
	}{
		{
			name:    "offset negativo",
			src:     []byte("hola"),
			offset:  -1,
			wantErr: ErrOutOfBounds,
		},
		{
			name:    "offset mayor que Len()",
			src:     []byte("hola"),
			offset:  5,
			wantErr: ErrOutOfBounds,
		},
		{
			name:    "documento vacío con offset positivo",
			src:     []byte(""),
			offset:  1,
			wantErr: ErrOutOfBounds,
		},
		{
			name:    "offset en el interior de una rune multibyte",
			src:     []byte("界"),
			offset:  1,
			wantErr: ErrByteOffsetNotAtRuneBoundary,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := mustNew(t, tt.src)

			got, err := d.ToPosition(tt.offset)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf(
					"Document.ToPosition(%d) error = %v; want %v (posición devuelta: %v)",
					tt.offset, err, tt.wantErr, got,
				)
			}
		})
	}
}

func TestDocumentToRegion(t *testing.T) {
	tests := []struct {
		name      string
		src       []byte
		start     ByteOffset
		end       ByteOffset
		wantStart Position
		wantEnd   Position
	}{
		{
			name:      "rango dentro de una línea",
			src:       []byte("ab\ncd"),
			start:     0,
			end:       2,
			wantStart: NewPosition(1, 1),
			wantEnd:   NewPosition(1, 3),
		},
		{
			name:      "rango que cruza un salto de línea",
			src:       []byte("ab\ncd"),
			start:     1,
			end:       4,
			wantStart: NewPosition(1, 2),
			wantEnd:   NewPosition(2, 2),
		},
		{
			name:      "rango vacío conserva la posición",
			src:       []byte("ab\ncd"),
			start:     3,
			end:       3,
			wantStart: NewPosition(2, 1),
			wantEnd:   NewPosition(2, 1),
		},
		{
			name:      "documento completo",
			src:       []byte("ab\ncd"),
			start:     0,
			end:       5,
			wantStart: NewPosition(1, 1),
			wantEnd:   NewPosition(2, 3),
		},
		{
			name:      "documento vacío",
			src:       []byte(""),
			start:     0,
			end:       0,
			wantStart: NewPosition(1, 1),
			wantEnd:   NewPosition(1, 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := mustNew(t, tt.src)
			rng := mustRange(t, tt.start, tt.end)

			got, err := d.ToRegion(rng)
			if err != nil {
				t.Fatalf("Document.ToRegion(%v) error = %v; want nil", rng, err)
			}

			if got.Start().Compare(tt.wantStart) != 0 || got.End().Compare(tt.wantEnd) != 0 {
				t.Errorf(
					"Document.ToRegion(%v) = %v; want [%v, %v)",
					rng, got, tt.wantStart, tt.wantEnd,
				)
			}
		})
	}
}

func TestDocumentToRegionErrors(t *testing.T) {
	tests := []struct {
		name    string
		src     []byte
		start   ByteOffset
		end     ByteOffset
		wantErr error
	}{
		{
			name:    "start negativo",
			src:     []byte("ab\ncd"),
			start:   -1,
			end:     2,
			wantErr: ErrOutOfBounds,
		},
		{
			name:    "end mayor que Len()",
			src:     []byte("ab\ncd"),
			start:   0,
			end:     6,
			wantErr: ErrOutOfBounds,
		},
		{
			name:    "start en el interior de una rune multibyte",
			src:     []byte("界a"),
			start:   1,
			end:     3,
			wantErr: ErrByteOffsetNotAtRuneBoundary,
		},
		{
			name:    "end en el interior de una rune multibyte",
			src:     []byte("a界"),
			start:   0,
			end:     2,
			wantErr: ErrByteOffsetNotAtRuneBoundary,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := mustNew(t, tt.src)
			rng := mustRange(t, tt.start, tt.end)

			got, err := d.ToRegion(rng)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf(
					"Document.ToRegion(%v) error = %v; want %v (región devuelta: %v)",
					rng, err, tt.wantErr, got,
				)
			}
		})
	}
}
