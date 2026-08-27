package document

import (
	"errors"
	"slices"
	"testing"
)

func TestNewLineMapConstructor(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
		want    []int
	}{
		{
			name:    "empty src",
			content: []byte{},
			want:    []int{0},
		},
		{
			name:    "just \\n",
			content: []byte{'\n'},
			want:    []int{0, 1},
		},
		{
			name:    "char and \\n",
			content: []byte{' ', '\n'},
			want:    []int{0, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lm := newLineMap(tt.content)
			got := lm.lineStarts

			if !slices.Equal(got, tt.want) {
				t.Errorf("newLineMap(%q).lineStarts = %v; want %v", tt.content, got, tt.want)
			}
		})
	}
}

func TestNewLineMapToPosition(t *testing.T) {
	tests := []struct {
		name   string
		offset ByteOffset
		src    []byte
		want   Position
	}{
		{
			name:   "Base: Inicio exacto del archivo",
			offset: ByteOffset(0),
			src:    []byte("xooo"),
			want:   NewPosition(1, 1),
		},
		{
			name:   "Frontera: Inicio exacto después de un salto de línea",
			offset: ByteOffset(2),
			src:    []byte("a\nb"),    // a=0, \n=1, b=2
			want:   NewPosition(2, 1), // Inicio de la segunda línea.
		},
		{
			name:   "EOF después de carácter UTF-8 de 3 bytes",
			offset: ByteOffset(3),
			src:    []byte("界"),
			want:   NewPosition(1, 2),
		},
		{
			name:   "Fuente vacía retorna posición inicial",
			offset: ByteOffset(0),
			src:    []byte(""),
			want:   NewPosition(1, 1),
		},
		{
			name:   "Salto de línea después de carácter multibyte",
			offset: ByteOffset(5),
			src:    []byte("ab界\ndef"),
			want:   NewPosition(1, 4),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lm := newLineMap(tt.src)
			got, err := lm.toPosition(tt.offset)

			if err != nil {
				t.Fatalf(
					"\nPrueba:  %s\nFuente:  %q (Longitud: %d bytes)\nOffset:  %d\nObtuvo:  %v\nEsperó:  %v error: %v", tt.name, tt.src, len(tt.src), tt.offset, got, tt.want, err,
				)
			}

			// Formateo estricto del Error para ayudar al debugging visual
			if got.Compare(tt.want) != 0 {
				t.Errorf(
					"\nPrueba:  %s\nFuente:  %q (Longitud: %d bytes)\nOffset:  %d\nObtuvo:  %v\nEsperó:  %v",
					tt.name, tt.src, len(tt.src), tt.offset, got, tt.want,
				)
			}
		})
	}
}

func TestNewLineMapToPositionErrors(t *testing.T) {
	tests := []struct {
		name    string
		offset  ByteOffset
		src     []byte
		wantErr error
	}{
		{
			name:    "Offset negativo",
			offset:  ByteOffset(-10),
			src:     []byte("xooo"),
			wantErr: ErrOutOfBounds,
		},
		{
			name:    "Offset justo después del final",
			offset:  ByteOffset(5),
			src:     []byte("xooo"),
			wantErr: ErrOutOfBounds,
		},
		{
			name:    "Offset muy por encima de la longitud",
			offset:  ByteOffset(999),
			src:     []byte("a\nb\nc"),
			wantErr: ErrOutOfBounds,
		},
		{
			name:    "Offset muy por encima de la longitud tras salto de línea final",
			offset:  ByteOffset(9999),
			src:     []byte("a\nb\nc\n"),
			wantErr: ErrOutOfBounds,
		},
		{
			name:    "Offset en el primer byte interior de rune UTF-8 de 3 bytes",
			offset:  ByteOffset(1),
			src:     []byte("界"),
			wantErr: ErrByteOffsetNotAtRuneBoundary,
		},
		{
			name:    "Offset en el segundo byte interior de rune UTF-8 de 3 bytes",
			offset:  ByteOffset(2),
			src:     []byte("界"),
			wantErr: ErrByteOffsetNotAtRuneBoundary,
		},
		{
			name:    "Offset interior de rune UTF-8 de 3 bytes con prefijo ASCII",
			offset:  ByteOffset(3),
			src:     []byte("a界b"),
			wantErr: ErrByteOffsetNotAtRuneBoundary,
		},
		{
			name:    "Offset en el primer byte interior de emoji UTF-8 de 4 bytes",
			offset:  ByteOffset(2),
			src:     []byte("x🌎y"),
			wantErr: ErrByteOffsetNotAtRuneBoundary,
		},
		{
			name:    "Offset en el segundo byte interior de emoji UTF-8 de 4 bytes",
			offset:  ByteOffset(3),
			src:     []byte("x🌎y"),
			wantErr: ErrByteOffsetNotAtRuneBoundary,
		},
		{
			name:    "Offset en el tercer byte interior de emoji UTF-8 de 4 bytes",
			offset:  ByteOffset(4),
			src:     []byte("x🌎y"),
			wantErr: ErrByteOffsetNotAtRuneBoundary,
		},
		{
			name:    "Offset interior de rune multibyte después de salto de línea",
			offset:  ByteOffset(3),
			src:     []byte("a\n界b"),
			wantErr: ErrByteOffsetNotAtRuneBoundary,
		},
		{
			name:    "Offset interior de rune multibyte en una línea posterior",
			offset:  ByteOffset(5),
			src:     []byte("a\nb\n界"),
			wantErr: ErrByteOffsetNotAtRuneBoundary,
		},
		{
			name:    "Offset interior de segundo rune multibyte consecutivo",
			offset:  ByteOffset(5),
			src:     []byte("界界"),
			wantErr: ErrByteOffsetNotAtRuneBoundary,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lm := newLineMap(tt.src)
			got, err := lm.toPosition(tt.offset)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("\nPrueba:  %s\nFuente:  %q (Longitud: %d bytes)\nOffset:  %d \nPosicion obtenida: %v \nError Obtenido:  %v\nError Esperado:  %v", tt.name, tt.src, len(tt.src), tt.offset, got, err, tt.wantErr)
			}

		})
	}

}
