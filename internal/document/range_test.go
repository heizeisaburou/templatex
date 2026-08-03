package document

import "testing"

func TestRangeLen(t *testing.T) {
	r := NewRange(2, 5)

	if got, want := r.Len(), ByteOffset(3); got != want {
		t.Errorf("NewRange(2, 5).Len() = %d; want %d", got, want)
	}
}

func TestRangeEmpty(t *testing.T) {
	tests := []struct {
		name string
		rng  Range
		want bool
	}{
		{name: "empty", rng: NewRange(2, 2), want: true},
		{name: "not empty", rng: NewRange(2, 5), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rng.Empty(); got != tt.want {
				t.Errorf("%v.Empty() = %t; want %t", tt.rng, got, tt.want)
			}
		})
	}
}

func TestRegionEmpty(t *testing.T) {
	position := NewPosition(1, 2)
	region := NewRegion(position, position)

	if !region.Empty() {
		t.Errorf("%v.Empty() = false; want true", region)
	}
}
