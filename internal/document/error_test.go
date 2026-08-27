package document

import (
	"errors"
	"testing"
)

// Los tests de este archivo fijan el contrato de errores del paquete: toda
// operación que por dentro se apoye en pkg/rnge o pkg/cursor debe fallar con un
// centinela de este paquete, para que quien use document no tenga que conocer
// los errores de los paquetes de apoyo.

func TestNewRangeErrors(t *testing.T) {
	got, err := NewRange(5, 2)

	if !errors.Is(err, ErrInvalidRange) {
		t.Errorf("NewRange(5, 2) err = %v; want = %v", err, ErrInvalidRange)
	}

	if got != (Range{}) {
		t.Errorf("NewRange(5, 2) = %v; want = %v", got, Range{})
	}
}

func TestRangeSettersErrors(t *testing.T) {
	tests := []struct {
		name string
		set  func(r *Range) error
	}{
		{
			name: "SetStart por encima del final",
			set:  func(r *Range) error { return r.SetStart(9) },
		},
		{
			name: "SetEnd por debajo del inicio",
			set:  func(r *Range) error { return r.SetEnd(0) },
		},
		{
			name: "Set con end menor que start",
			set:  func(r *Range) error { return r.Set(9, 1) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mustRange(t, 2, 4)

			if err := tt.set(&r); !errors.Is(err, ErrInvalidRange) {
				t.Errorf("err = %v; want = %v", err, ErrInvalidRange)
			}

			// Una asignación fallida no debe modificar el rango.
			if got := r.Array(); got != [2]ByteOffset{2, 4} {
				t.Errorf("Array() = %v; want = %v", got, [2]ByteOffset{2, 4})
			}
		})
	}
}

// mustRegion construye una Region o aborta la prueba.
func mustRegion(t *testing.T, start, end Position) Region {
	t.Helper()

	r, err := NewRegion(start, end)
	if err != nil {
		t.Fatalf("NewRegion(%v, %v) = %v; want nil", start, end, err)
	}

	return r
}

func TestNewRegionErrors(t *testing.T) {
	start := NewPosition(3, 1)
	end := NewPosition(1, 1)

	got, err := NewRegion(start, end)

	if !errors.Is(err, ErrInvalidRange) {
		t.Errorf("NewRegion(%v, %v) err = %v; want = %v", start, end, err, ErrInvalidRange)
	}

	if got != (Region{}) {
		t.Errorf("NewRegion(%v, %v) = %v; want = %v", start, end, got, Region{})
	}
}

func TestRegionSettersErrors(t *testing.T) {
	tests := []struct {
		name string
		set  func(r *Region) error
	}{
		{
			name: "SetStart por encima del final",
			set:  func(r *Region) error { return r.SetStart(NewPosition(9, 1)) },
		},
		{
			name: "SetEnd por debajo del inicio",
			set:  func(r *Region) error { return r.SetEnd(NewPosition(1, 1)) },
		},
		{
			name: "Set con end menor que start",
			set:  func(r *Region) error { return r.Set(NewPosition(9, 1), NewPosition(1, 1)) },
		},
	}

	want := [2]Position{NewPosition(2, 1), NewPosition(4, 1)}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mustRegion(t, want[0], want[1])

			if err := tt.set(&r); !errors.Is(err, ErrInvalidRange) {
				t.Errorf("err = %v; want = %v", err, ErrInvalidRange)
			}

			// Una asignación fallida no debe modificar la región.
			if got := r.Array(); got != want {
				t.Errorf("Array() = %v; want = %v", got, want)
			}
		})
	}
}

func TestRegionSetters(t *testing.T) {
	t.Run("SetStart, SetEnd y Set aceptan valores válidos", func(t *testing.T) {
		r := mustRegion(t, NewPosition(2, 1), NewPosition(4, 1))

		if err := r.SetStart(NewPosition(1, 1)); err != nil {
			t.Fatalf("SetStart() = %v; want nil", err)
		}

		if err := r.SetEnd(NewPosition(6, 1)); err != nil {
			t.Fatalf("SetEnd() = %v; want nil", err)
		}

		if err := r.Set(NewPosition(2, 2), NewPosition(3, 3)); err != nil {
			t.Fatalf("Set() = %v; want nil", err)
		}

		want := [2]Position{NewPosition(2, 2), NewPosition(3, 3)}
		if got := r.Array(); got != want {
			t.Errorf("Array() = %v; want = %v", got, want)
		}
	})
}
