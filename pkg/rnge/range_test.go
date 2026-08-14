package rnge

import "testing"

func TestNewOrdered(t *testing.T) {
	r := New(2, 5)

	if start, end := r.Bounds(); start != 2 || end != 5 {
		t.Fatalf("New(2, 5).Bounds() = (%d, %d); want (2, 5)", start, end)
	}
}

func TestNewCompare(t *testing.T) {
	type position struct {
		line int
		col  int
	}

	compare := func(a, b position) int {
		switch {
		case a.line < b.line || a.line == b.line && a.col < b.col:
			return -1
		case a == b:
			return 0
		default:
			return 1
		}
	}

	start := position{line: 1, col: 8}
	end := position{line: 2, col: 3}
	r := NewCompare(start, end, compare)

	if got := r.Start(); got != start {
		t.Fatalf("Start() = %v; want %v", got, start)
	}
	if got := r.End(); got != end {
		t.Fatalf("End() = %v; want %v", got, end)
	}
}
