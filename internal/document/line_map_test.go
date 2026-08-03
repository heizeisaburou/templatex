package document

import (
	"fmt"
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

			assertIntSliceEqual(t, fmt.Sprintf("newLineMap(%q).lineStarts", tt.content), got, tt.want)
		})
	}
}

func assertIntSliceEqual(t *testing.T, name string, got, want []int) {
	t.Helper()

	if !slices.Equal(got, want) {
		t.Errorf("%s = %v; want %v", name, got, want)
	}
}
