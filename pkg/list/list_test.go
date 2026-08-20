package list

import (
	"slices"
	"testing"
)

// :TODO: Saburou: voy a hacer los tests de aquí

func TestNew(t *testing.T) {
	t.Run("new stores items", func(t *testing.T) {
		got := New[int](1, 2, 3)
		want := []int{1, 2, 3}
		if !slices.Equal(got.items, want) {
			t.Errorf("New(%v) = %v; want %v", []int{1, 2, 3}, got, want)
		}
	})

	t.Run("new perfoms shallow copy", func(t *testing.T) {
		value := 1
		items := []*int{&value}

		got := New[int](items...)

		*items[0] = 69
		if *got.items[0] != 69 {
			t.Errorf("New(%v) = %v; want shared value %v", items, got.items, 69)
		}
	})
}

func TestPush(t *testing.T) {
	tests := []struct {
		name  string
		items []int
		push  int
		want  []int
	}{
		{name: "empty", push: 69, want: []int{69}},
		{name: "one", items: []int{68}, push: 69, want: []int{68, 69}},
		{name: "two", items: []int{68, 69}, push: 70, want: []int{68, 69, 70}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New[int](tt.items...)

			got.Push(tt.push)

			if got, want := got.items, tt.want; !slices.Equal(got, want) {
				t.Errorf("New(%v) = %v; want %v", tt.items, got, want)
			}
		})
	}
}

// Camino válido asigna al elemento correcto
// Camino inválido lanza error y no modifica la lista
func TestSet(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		id        int
		item      int
		want      int
		wantError error
	}{}
}
