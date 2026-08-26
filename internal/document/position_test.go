package document

import "testing"

func TestNewPosition(t *testing.T) {
	t.Run("Validating constructor", func(t *testing.T) {
		p := NewPosition(5, 10)
		if got := p.line; got != 5 {
			t.Errorf("Retun %d, want 5", got)
		}
		if got := p.col; got != 10 {
			t.Errorf("Retun %d, want 5", got)
		}

		t.Run("Validating panic wheen line or col are smaller than zero", func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Error("NewPosition() must panic when line or col are smaller than 0")
				}
			}()

			NewPosition(-5, -10)
		})
	})
}

func TestCompare(t *testing.T) {
	tests := []struct {
		name  string
		p     Position
		other Position
		want  int
	}{
		{
			name:  "equal positions",
			p:     NewPosition(5, 10),
			other: NewPosition(5, 10),
			want:  0,
		},
		{
			name:  "p line is smaller",
			p:     NewPosition(3, 10),
			other: NewPosition(5, 10),
			want:  -1,
		},
		{
			name:  "p line is bigger",
			p:     NewPosition(5, 10),
			other: NewPosition(3, 10),
			want:  1,
		},
		{
			name:  "same line but p column is smaller",
			p:     NewPosition(5, 6),
			other: NewPosition(5, 10),
			want:  -1,
		},
		{
			name:  "same line but p column is bigger",
			p:     NewPosition(5, 10),
			other: NewPosition(5, 6),
			want:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Compare(tt.other); got != tt.want {
				t.Errorf("Compare() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestLine(t *testing.T) {
	p := NewPosition(5, 10)

	if got := p.Line(); got != 5 {
		t.Errorf("p Line() returns %d, want 5", got)
	}
}

func TestCol(t *testing.T) {
	p := NewPosition(5, 10)

	if got := p.Col(); got != 10 {
		t.Errorf("p Col() returns %d, want 10", got)
	}
}

func TestSetLine(t *testing.T) {
	t.Run("SetLine updates the line", func(t *testing.T) {
		p := NewPosition(5, 10)
		p.SetLine(6)

		if got := p.Line(); got != 6 {
			t.Errorf("SetLine(6).Line() = %d, want 6", got)
		}
	})

	t.Run("SetLine panics with negative value", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Error("SetLine() must panic when value is smaller than 0")
			}
		}()
		p := NewPosition(5, 10)
		p.SetLine(-1)
	})
}

func TestSetCol(t *testing.T) {
	t.Run("SetCol updates the column", func(t *testing.T) {
		p := NewPosition(5, 10)
		p.SetCol(6)

		if got := p.Col(); got != 6 {
			t.Errorf("SetCol(6).Col() = %d, want 6", got)
		}
	})

	t.Run("SetCol panics with negative value", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Error("SetCol() must panic when value is smaller than 0")
			}
		}()
		p := NewPosition(5, 10)
		p.SetCol(-1)
	})
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		position Position
		want     string
	}{
		{
			name:     "normal position",
			position: NewPosition(5, 10),
			want:     "5:10",
		},
		{
			name:     "zero position",
			position: NewPosition(0, 0),
			want:     "0:0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.position.String(); got != tt.want {
				t.Errorf("Compare() = %s, want %s", got, tt.want)
			}
		})
	}
}
