package padawan

import (
	"errors"
	"testing"
)

func TestDividirNumero(t *testing.T) {
	tests := []struct {
		name      string
		dividend  int
		divisor   int
		want      int
		wantError error
	}{
		{
			name:     "simple division",
			dividend: 8,
			divisor:  4,
			want:     2,
		},
		{
			name:      "division by zero",
			dividend:  8,
			divisor:   0,
			wantError: DivisionByZeroError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := Divide(tt.dividend, tt.divisor)

			if !errors.Is(err, tt.wantError) {
				t.Fatalf(
					"Divide(%d, %d) %d, %v; want err = %v",
					tt.dividend,
					tt.divisor,
					got,
					err,
					tt.wantError,
				)
			}

			if tt.wantError == nil && got != tt.want {
				t.Errorf(
					"Divide(%d, %d) = %d; want %d",
					tt.dividend,
					tt.divisor,
					got,
					tt.want,
				)
			}
		})
	}
}
