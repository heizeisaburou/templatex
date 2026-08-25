package clamp

import "testing"

func TestClampIndex(t *testing.T) {
	tests := []struct {
		name      string
		length    int
		index     int
		want      int
		wantPanic bool
	}{
		{
			name:   "negative index",
			length: 10,
			index:  -5,
			want:   0,
		},
		{
			name:   "index inside the range",
			length: 10,
			index:  5,
			want:   5,
		},
		{
			name:   "last valid index",
			length: 10,
			index:  9,
			want:   9,
		},
		{
			name:   "index equals to lenght",
			length: 10,
			index:  10,
			want:   9,
		},
		{
			name:   "index greater than lenght",
			length: 10,
			index:  11,
			want:   9,
		},
		{
			name:      "length zero",
			length:    0,
			index:     0,
			wantPanic: true,
		},
		{
			name:      "negative length",
			length:    -1,
			index:     0,
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if recover() == nil {
						t.Error("ClampIndex must be panic when lenght is equal or less to zero")
					}
				}()
			}

			got := ClampIndex(tt.length, tt.index)

			if tt.wantPanic {
				return
			}

			if got != tt.want {
				t.Errorf("ClampIndex(%d, %d) = %d, want %d", tt.length, tt.index, got, tt.want)

			}
		})
	}
}

func TestIsIndexValid(t *testing.T) {
	tests := []struct {
		name   string
		length int
		index  int
		want   bool
	}{
		{
			name:   "negative index",
			length: 5,
			index:  -1,
			want:   false,
		},
		{
			name:   "zero",
			length: 5,
			index:  0,
			want:   true,
		},
		{
			name:   "inside",
			length: 5,
			index:  2,
			want:   true,
		},
		{
			name:   "last valid",
			length: 5,
			index:  4,
			want:   true,
		},
		{
			name:   "equals length",
			length: 5,
			index:  5,
			want:   false,
		},
		{
			name:   "overflow",
			length: 5,
			index:  10,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIndexValid(tt.length, tt.index); got != tt.want {
				t.Errorf("IsIndexValid(%d, %d) = %t, want %t", tt.length, tt.index, got, tt.want)
			}
		})
	}
}

func TestClampPosition(t *testing.T) {
	tests := []struct {
		name      string
		length    int
		position  int
		want      int
		wantPanic bool
	}{
		{
			name:      "negative length panics",
			length:    -1,
			position:  0,
			wantPanic: true,
		},
		{
			name:     "negative position -> 0",
			length:   5,
			position: -1,
			want:     0,
		},
		{
			name:     "inside -> same",
			length:   5,
			position: 2,
			want:     2,
		},
		{
			name:     "equals length -> length",
			length:   5,
			position: 5,
			want:     5,
		},
		{
			name:     "overflow -> length",
			length:   5,
			position: 10,
			want:     5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if recover() == nil {
						t.Error("ClampPosition must panic when length is negative")
					}
				}()
			}

			got := ClampPosition(tt.length, tt.position)

			if tt.wantPanic {
				return
			}

			if got != tt.want {
				t.Errorf("ClampPosition(%d, %d) = %d, want %d", tt.length, tt.position, got, tt.want)
			}
		})
	}
}

func TestIsPositionValid(t *testing.T) {
	tests := []struct {
		name     string
		length   int
		position int
		want     bool
	}{
		{
			name:     "negative position",
			length:   5,
			position: -1,
			want:     false,
		},
		{
			name:     "zero",
			length:   5,
			position: 0,
			want:     true,
		},
		{
			name:     "inside",
			length:   5,
			position: 2,
			want:     true,
		},
		{
			name:     "equals length",
			length:   5,
			position: 5,
			want:     true,
		},
		{
			name:     "overflow",
			length:   5,
			position: 10,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPositionValid(tt.length, tt.position); got != tt.want {
				t.Errorf("IsPositionValid(%d, %d) = %t, want %t", tt.length, tt.position, got, tt.want)
			}
		})
	}
}
