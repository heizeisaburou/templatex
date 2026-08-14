package cursor

import (
	"errors"
	"slices"
	"testing"
)

func TestNew(t *testing.T) {
	src := []byte{'a', 'b', 'c'}

	cur := New[byte, int](src)

	if got, want := cur.Pos(), 0; got != want {
		t.Fatalf("cur.Pos() = %d, want %d", got, want)
	}

	if !slices.Equal(src, cur.src) {
		t.Fatalf("cur.src = %#v, want %#v", cur.src, src)
	}

	src[0] = 'x'

	got, ok := cur.Peek()
	if !ok {
		t.Fatal("cur.Peek() ok = false, want true")
	}

	if want := byte('x'); got != want {
		t.Fatalf(
			"cur.Peek() = %q, want %q; src appears to have been copied",
			got,
			want,
		)
	}
}

func TestNewAt(t *testing.T) {
	src := []byte{'a', 'b', 'c'}

	cur, err := NewAt(src, len(src))
	if err != nil {
		t.Fatalf("NewAt(%#v, %d) unexpected error: %v", src, len(src), err)
	}

	if got, want := cur.Pos(), 3; got != want {
		t.Fatalf("cur.Pos() = %d, want %d", got, want)
	}

	if !slices.Equal(src, cur.src) {
		t.Fatalf("cur.src = %#v, want %#v", cur.src, src)
	}

	src[0] = 'x'

	if err := cur.Seek(0); err != nil {
		t.Fatalf(
			"cur.Seek(%d) with len(%d) unexpected error: %v",
			0,
			len(src),
			err,
		)
	}

	got, ok := cur.Peek()
	if !ok {
		t.Fatal("cur.Peek() ok = false, want true")
	}

	if want := byte('x'); got != want {
		t.Fatalf(
			"cur.Peek() = %q, want %q; src appears to have been copied",
			got,
			want,
		)
	}
}

func TestNewAtOutOfRange(t *testing.T) {
	src := []byte{'a', 'b', 'c'}

	_, err := NewAt(src, len(src)+1)
	if !errors.Is(err, ErrPositionOutOfRange) {
		t.Fatalf(
			"NewAt(%q, %d) error = %v, want %v",
			src,
			len(src)+1,
			err,
			ErrPositionOutOfRange,
		)
	}
}

func TestCursorLen(t *testing.T) {
	tests := []struct {
		name string
		src  []byte
		want int
	}{
		{name: "0 elements", src: []byte{}, want: 0},
		{name: "three elements", src: []byte{'a', 'b', 'c'}, want: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cur := New[byte, int](tt.src)

			if got := cur.Len(); got != tt.want {
				t.Errorf("cur.Len() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCursorPos(t *testing.T) {
	tests := []struct {
		name string
		src  []byte
		seek int
	}{
		{name: "empty (end)"},
		{name: "first", src: []byte{'a', 'b'}},
		{name: "last", src: []byte{'a', 'b'}, seek: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cur := New[byte, int](tt.src)

			if err := cur.Seek(tt.seek); err != nil {
				t.Fatalf("cur.Seek(%d) error, want nil: %v", tt.seek, err)
			}

			if got, want := cur.Pos(), tt.seek; got != want {
				t.Errorf("cur.Pos() = %d, want %d", got, want)
			}
		})
	}
}

func TestCursorEOF(t *testing.T) {
	tests := []struct {
		name    string
		src     []byte
		pos     int
		wantEOF bool
	}{
		{name: "empty", wantEOF: true},
		{name: "one, at first", src: []byte{'a'}},
		{name: "one, at EOF", src: []byte{'a'}, pos: 1, wantEOF: true},
		{name: "two, at last", src: []byte{'a', 'b'}, pos: 1},
		{name: "two, at EOF", src: []byte{'a', 'b'}, pos: 2, wantEOF: true},
		{name: "three, at middle", src: []byte{'a', 'b', 'c'}, pos: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cur := New[byte, int](tt.src)

			if err := cur.Seek(tt.pos); err != nil {
				t.Fatalf("cur.Seek(%d) unexpected error: %v", tt.pos, err)
			}

			if got := cur.EOF(); got != tt.wantEOF {
				t.Fatalf(
					"cur.EOF() = %t, want %t [cur.Pos() = %d, src = %q]",
					got,
					tt.wantEOF,
					cur.Pos(),
					tt.src,
				)
			}
		})
	}
}

func TestCursorPeek(t *testing.T) {
	tests := []struct {
		name      string
		src       []byte
		pos       int
		wantValue byte
		wantOk    bool
	}{
		{
			name:      "available element",
			src:       []byte{'a'},
			wantValue: 'a',
			wantOk:    true,
		},
		{name: "empty sequence"},
		{name: "EOF", src: []byte{'a'}, pos: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cur, err := NewAt(tt.src, tt.pos)
			if err != nil {
				t.Fatalf("NewAt(%q, %d): %v", tt.src, tt.pos, err)
			}

			gotValue, gotOk := cur.Peek()

			if gotOk != tt.wantOk {
				t.Fatalf("cur.Peek() ok = %t, want %t", gotOk, tt.wantOk)
			}

			if gotOk && gotValue != tt.wantValue {
				t.Errorf(
					"cur.Peek() value = %q, want %q",
					gotValue,
					tt.wantValue,
				)
			}
		})
	}
}

func TestCursorNext(t *testing.T) {
	tests := []struct {
		name      string
		src       []byte
		pos       int
		wantValue byte
		wantOk    bool
	}{
		{
			name:      "available element",
			src:       []byte{'a'},
			wantValue: 'a',
			wantOk:    true,
		},
		{
			name:   "empty sequence",
			wantOk: false,
		},
		{
			name:   "EOF",
			src:    []byte{'a'},
			pos:    1,
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cur, err := NewAt(tt.src, tt.pos)
			if err != nil {
				t.Fatalf("NewAt(%q, %d): %v", tt.src, tt.pos, err)
			}

			gotValue, gotOk := cur.Next()

			if gotOk != tt.wantOk {
				t.Fatalf("cur.Peek() ok = %t, want %t", gotOk, tt.wantOk)
			}

			if gotOk && gotValue != tt.wantValue {
				t.Errorf(
					"cur.Peek() value = %q, want %q",
					gotValue,
					tt.wantValue,
				)
			}
		})
	}
}

func TestCursorSeek(t *testing.T) {
	tests := []struct {
		name    string
		src     []byte
		seek    int
		wantErr bool
	}{
		{name: "empty sequence, to EOF"},
		{name: "available element", src: []byte{'a'}},
		{name: "to EOF", src: []byte{'a', 'b', 'c'}, seek: 3},
		{
			name:    "negative position",
			src:     []byte{'a', 'b', 'c'},
			seek:    -1,
			wantErr: true,
		},
		{
			name:    "beyond EOF",
			src:     []byte{'a', 'b', 'c'},
			seek:    4,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cur := New[byte, int](tt.src)

			err := cur.Seek(tt.seek)

			if tt.wantErr {
				if !errors.Is(err, ErrPositionOutOfRange) {
					t.Fatalf(
						"cur.Seek(%d) error = %v, want %v",
						tt.seek,
						err,
						ErrPositionOutOfRange,
					)
				}

				if got, want := cur.Pos(), 0; got != want {
					t.Fatalf(
						"cur.Pos() = %d, want %d; cursor should stay unchanged after error",
						got,
						want,
					)
				}

				return
			}

			if err != nil {
				t.Fatalf("cur.Seek(%d) unexpected error: %v", tt.seek, err)
			}

			if got, want := cur.Pos(), tt.seek; got != want {
				t.Fatalf("cur.Pos() = %d, want %d", got, want)
			}
		})
	}
}

func TestCursorMove(t *testing.T) {
	abc := []byte{'a', 'b', 'c'}

	tests := []struct {
		name        string
		src         []byte
		originalPos int
		delta       int
		wantPos     int
		wantErr     bool
	}{
		{
			name:        "to EOF",
			src:         abc,
			originalPos: 0,
			delta:       3,
			wantPos:     3,
		},
		{
			name:        "to start",
			src:         abc,
			originalPos: 3,
			delta:       -3,
			wantPos:     0,
		},
		{
			name:    "negative delta",
			src:     abc,
			delta:   -1,
			wantErr: true,
		},
		{
			name:    "over EOF delta",
			src:     abc,
			delta:   4,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cur, err := NewAt(tt.src, tt.originalPos)
			if err != nil {
				t.Fatalf("NewAt(%q, %d): %v", tt.src, tt.originalPos, err)
			}

			if gotErr := cur.Move(tt.delta); (gotErr != nil) != tt.wantErr {
				t.Fatalf("cur.Move(%d) error = %v; want %v", tt.delta, gotErr, tt.wantErr)
			}

			if tt.wantErr {
				// offset no cambia ante un error
				if got, want := cur.Pos(), tt.originalPos; got != want {
					t.Errorf(
						"cur.Pos() = %d want %d; cursor should stay unchanged after error",
						got,
						want,
					)
				}
				return
			}

			if gotPos := cur.Pos(); gotPos != tt.wantPos {
				t.Fatalf("cur.Pos() = %d; want %d", gotPos, tt.wantPos)
			}
		})
	}
}
