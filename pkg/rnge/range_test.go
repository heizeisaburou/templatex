package rnge

import (
	"errors"
	"testing"
)

type Int int

func (i Int) Compare(other Int) int {
	switch {
	case i < other:
		return -1
	case i > other:
		return 1
	default:
		return 0
	}
}

// newRange construye un rango válido o aborta el test.
func newRange(t *testing.T, start, end Int) Range[Int] {
	t.Helper()

	r, err := New(start, end)
	if err != nil {
		t.Fatalf("New(%v, %v) unexpected_error: %v", start, end, err)
	}

	return r
}

func TestNew(t *testing.T) {
	t.Run("start menor que end", func(t *testing.T) {
		r := newRange(t, 5, 10)

		if got := r.Start(); got != 5 {
			t.Errorf("Start() = %d; want = 5", got)
		}

		if got := r.End(); got != 10 {
			t.Errorf("End() = %d; want = 10", got)
		}
	})

	t.Run("start igual que end", func(t *testing.T) {
		r := newRange(t, 7, 7)

		if got := r.Start(); got != 7 {
			t.Errorf("Start() = %d; want = 7", got)
		}

		if got := r.End(); got != 7 {
			t.Errorf("End() = %d; want = 7", got)
		}
	})

	t.Run("start mayor que end devuelve ErrInvalidRange", func(t *testing.T) {
		got, err := New[Int](10, 5)

		if !errors.Is(err, ErrInvalidRange) {
			t.Errorf("New(10, 5) err = %v; want = %v", err, ErrInvalidRange)
		}

		if got != (Range[Int]{}) {
			t.Errorf("New(10, 5) = %v; want = %v", got, Range[Int]{})
		}
	})
}

func TestSet(t *testing.T) {
	t.Run("Set establece ambos índices", func(t *testing.T) {
		r := newRange(t, 2, 4)

		if err := r.Set(5, 10); err != nil {
			t.Fatalf("Set(5, 10) unexpected_error: %v", err)
		}

		if got := r.Start(); got != 5 {
			t.Errorf("Start() = %d; want = 5", got)
		}

		if got := r.End(); got != 10 {
			t.Errorf("End() = %d; want = 10", got)
		}
	})

	t.Run("Set devuelve ErrInvalidRange si end es menor que start", func(t *testing.T) {
		r := newRange(t, 2, 4)

		if err := r.Set(10, 5); !errors.Is(err, ErrInvalidRange) {
			t.Errorf("Set(10, 5) err = %v; want = %v", err, ErrInvalidRange)
		}

		// El rango no debe modificarse cuando Set falla.
		if got := r.Array(); got != [2]Int{2, 4} {
			t.Errorf("Array() = %v; want = %v", got, [2]Int{2, 4})
		}
	})
}

func TestSetStart(t *testing.T) {
	t.Run("SetStart establece el índice inicial", func(t *testing.T) {
		r := newRange(t, 2, 10)

		if err := r.SetStart(5); err != nil {
			t.Fatalf("SetStart(5) unexpected_error: %v", err)
		}

		if got := r.Start(); got != 5 {
			t.Errorf("Start() = %d; want = 5", got)
		}
	})

	t.Run("SetStart devuelve ErrInvalidRange si start supera a end", func(t *testing.T) {
		r := newRange(t, 2, 10)

		if err := r.SetStart(11); !errors.Is(err, ErrInvalidRange) {
			t.Errorf("SetStart(11) err = %v; want = %v", err, ErrInvalidRange)
		}

		if got := r.Start(); got != 2 {
			t.Errorf("Start() = %d; want = 2", got)
		}
	})
}

func TestSetEnd(t *testing.T) {
	t.Run("SetEnd establece el índice final", func(t *testing.T) {
		r := newRange(t, 2, 10)

		if err := r.SetEnd(20); err != nil {
			t.Fatalf("SetEnd(20) unexpected_error: %v", err)
		}

		if got := r.End(); got != 20 {
			t.Errorf("End() = %d; want = 20", got)
		}
	})

	t.Run("SetEnd devuelve ErrInvalidRange si end es menor que start", func(t *testing.T) {
		r := newRange(t, 2, 10)

		if err := r.SetEnd(1); !errors.Is(err, ErrInvalidRange) {
			t.Errorf("SetEnd(1) err = %v; want = %v", err, ErrInvalidRange)
		}

		if got := r.End(); got != 10 {
			t.Errorf("End() = %d; want = 10", got)
		}
	})
}

func TestSetRange(t *testing.T) {
	t.Run("SetRange asigna a r los índices de other", func(t *testing.T) {
		r := newRange(t, 2, 4)
		other := newRange(t, 10, 20)

		r.SetRange(other)

		if got := r.Start(); got != other.Start() {
			t.Errorf("Start() = %d; want = %d", got, other.Start())
		}

		if got := r.End(); got != other.End() {
			t.Errorf("End() = %d; want = %d", got, other.End())
		}
	})
}

func TestEmpty(t *testing.T) {
	tests := []struct {
		name  string
		start Int
		end   Int
		want  bool
	}{
		{name: "start igual que end", start: 5, end: 5, want: true},
		{name: "rango nulo", start: 0, end: 0, want: true},
		{name: "start menor que end", start: 5, end: 10, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRange(t, tt.start, tt.end)

			if got := r.Empty(); got != tt.want {
				t.Errorf(
					"Range%v.Empty() = %t; want = %t",
					r,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestBounds(t *testing.T) {
	t.Run("Bounds devuelve ambos índices", func(t *testing.T) {
		r := newRange(t, 5, 10)

		start, end := r.Bounds()

		if start != 5 || end != 10 {
			t.Errorf("Bounds() = (%d, %d); want = (5, 10)", start, end)
		}
	})
}

func TestArray(t *testing.T) {
	t.Run("Array devuelve los índices como array", func(t *testing.T) {
		r := newRange(t, 5, 10)

		if got := r.Array(); got != [2]Int{5, 10} {
			t.Errorf("Array() = %v; want = %v", got, [2]Int{5, 10})
		}
	})
}

func TestString(t *testing.T) {
	t.Run("String usa el formato [start, end)", func(t *testing.T) {
		r := newRange(t, 5, 10)

		if got := r.String(); got != "[5, 10)" {
			t.Errorf("String() = %q; want = %q", got, "[5, 10)")
		}
	})
}
