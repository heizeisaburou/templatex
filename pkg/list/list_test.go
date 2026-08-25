package list

import (
	"errors"
	"slices"
	"testing"
)

// ptrs devuelve un puntero a cada elemento de values.
func ptrs(values []int) []*int {
	got := make([]*int, len(values))
	for i := range values {
		got[i] = &values[i]
	}

	return got
}

// wantSharedPointers comprueba que got sigue apuntando a los enteros de values,
// es decir, que la copia hecha por op fue superficial.
func wantSharedPointers(t *testing.T, op string, got []*int, values []int) {
	t.Helper()

	for i := range values {
		values[i] += 69
	}

	gotValues := make([]int, len(got))
	for i, p := range got {
		gotValues[i] = *p
	}

	if !slices.Equal(gotValues, values) {
		t.Errorf("%s hizo una copia profunda: got %v; want %v", op, gotValues, values)
	}
}

func TestNew(t *testing.T) {
	t.Run("new stores items", func(t *testing.T) {
		got := New[int](1, 2, 3)
		want := []int{1, 2, 3}
		if !slices.Equal(got.items, want) {
			t.Errorf("New(%v) = %v; want %v", []int{1, 2, 3}, got, want)
		}
	})

	t.Run("new performs shallow copy", func(t *testing.T) {
		values := []int{1}

		got := New[int](ptrs(values)...)

		wantSharedPointers(t, "New()", got.items, values)
	})

	t.Run("new does not alias the source slice", func(t *testing.T) {
		items := []int{1, 2, 3}

		got := New[int](items...)

		items[0] = 69
		if got.items[0] != 1 {
			t.Errorf("New(%v) aliased its argument: items = %v; want %v", items, got.items, []int{1, 2, 3})
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

func TestSet(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		id        int
		item      int
		want      []int
		wantError error
	}{
		{
			name:  "set first",
			items: []int{0},
			id:    0,
			item:  69,
			want:  []int{69},
		},
		{
			name:  "set last",
			items: []int{0, 1, 2},
			id:    2,
			item:  69,
			want:  []int{0, 1, 69},
		},
		{
			name:      "set negative",
			id:        -1,
			wantError: ErrOutOfBounds,
		},
		{
			name:      "set to EOF",
			id:        0,
			wantError: ErrOutOfBounds,
		},
		{
			name:      "set to post-EOF",
			items:     []int{1, 2, 3},
			id:        69,
			want:      []int{1, 2, 3},
			wantError: ErrOutOfBounds,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New[int](tt.items...)

			gotErr := got.Set(tt.id, tt.item)
			if !errors.Is(gotErr, tt.wantError) {
				t.Fatalf("List.Set(%d) = %v; want = %v", tt.id, gotErr, tt.wantError)
			}

			// En los casos de error tt.want es el contenido original,
			// así que esto comprueba también que la lista quedó intacta.
			if !slices.Equal(got.items, tt.want) {
				t.Errorf("List = %v; want = %v", got.items, tt.want)
			}
		})
	}
}

// Camino válido: la lista queda recortada a lo que hemos determinado
// Camino inválido: la lista queda inalterada
func TestTruncate(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		end       int
		want      []int
		wantError error
	}{
		{
			name: "truncate empty at end", // also end == len
			end:  0,
			want: []int{},
		},
		{
			name:  "truncate first",
			items: []int{1},
			end:   0,
			want:  []int{},
		},
		{
			name:  "truncate first of three",
			items: []int{1, 2, 3},
			end:   1,
			want:  []int{1},
		},
		{
			name:      "negative end",
			end:       -1,
			wantError: ErrOutOfBounds,
		},
		{
			name:      "post end",
			items:     []int{1, 2, 3},
			end:       69,
			want:      []int{1, 2, 3},
			wantError: ErrOutOfBounds,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New[int](tt.items...)

			gotErr := got.Truncate(tt.end)
			if !errors.Is(gotErr, tt.wantError) {
				t.Fatalf("List.Truncate(%d) = %v; want = %v", tt.end, gotErr, tt.wantError)
			}

			// En los casos de error tt.want es el contenido original,
			// así que esto comprueba también que la lista quedó intacta.
			if !slices.Equal(got.items, tt.want) {
				t.Errorf("List = %v; want = %v", got.items, tt.want)
			}
		})
	}
}

func TestReplaceRange(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		start     int
		end       int
		repl      []int
		want      []int
		wantError error
	}{
		{name: "do nothing"},
		{
			name:  "replace first",
			items: []int{0},
			start: 0,
			end:   1,
			repl:  []int{69},
			want:  []int{69},
		},
		{
			name:  "replace middle",
			items: []int{1, 2, 3},
			start: 1,
			end:   2,
			repl:  []int{69},
			want:  []int{1, 69, 3},
		},
		{
			name:  "append at start",
			items: []int{1, 2, 3},
			start: 0,
			end:   0,
			repl:  []int{69},
			want:  []int{69, 1, 2, 3},
		},
		{
			name:  "append at end", // start == end == len
			items: []int{1, 2, 3},
			start: 3,
			end:   3,
			repl:  []int{69},
			want:  []int{1, 2, 3, 69},
		},
		{
			name:  "replace and extend middle",
			items: []int{1, 2, 3},
			start: 1,
			end:   2,
			repl:  []int{68, 69, 70},
			want:  []int{1, 68, 69, 70, 3},
		},
		{
			name:  "replace all",
			items: []int{1, 2, 3},
			start: 0,
			end:   3,
			repl:  []int{69},
			want:  []int{69},
		},
		{
			name:  "replace and shrink middle",
			items: []int{1, 2, 3, 4, 5},
			start: 1,
			end:   4,
			repl:  []int{69},
			want:  []int{1, 69, 5},
		},
		{
			name:      "negative start",
			items:     []int{1, 2, 3},
			start:     -1,
			end:       2,
			repl:      []int{69},
			want:      []int{1, 2, 3},
			wantError: ErrOutOfBounds,
		},
		{
			name:      "end before start",
			items:     []int{1, 2, 3},
			start:     2,
			end:       1,
			repl:      []int{69},
			want:      []int{1, 2, 3},
			wantError: ErrOutOfBounds,
		},
		{
			name:      "end after list",
			items:     []int{1, 2, 3},
			start:     1,
			end:       69,
			repl:      []int{69},
			want:      []int{1, 2, 3},
			wantError: ErrOutOfBounds,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New[int](tt.items...)

			gotErr := got.ReplaceRange(tt.start, tt.end, tt.repl...)
			if !errors.Is(gotErr, tt.wantError) {
				t.Fatalf(
					"List.ReplaceRange(%d, %d, %v) = %v; want = %v",
					tt.start,
					tt.end,
					tt.repl,
					gotErr,
					tt.wantError,
				)
			}

			// En los casos de error tt.want es el contenido original,
			// así que esto comprueba también que la lista quedó intacta.
			if !slices.Equal(got.items, tt.want) {
				t.Errorf("List = %v; want = %v", got.items, tt.want)
			}
		})
	}
}

func TestReplaceRangePerformsShallowCopy(t *testing.T) {
	values := []int{1}

	got := New[int, *int]()
	if err := got.ReplaceRange(0, 0, ptrs(values)...); err != nil {
		t.Fatalf("List.ReplaceRange(0, 0, %v) = %v; want nil", values, err)
	}

	wantSharedPointers(t, "List.ReplaceRange()", got.items, values)
}

// si el indice existe devuelve elemento, true
// si el indice no existe devuelve zero, false
func TestAt(t *testing.T) {
	tests := []struct {
		name   string
		items  []int
		id     int
		want   int
		wantOk bool
	}{
		{
			name:   "first",
			items:  []int{69},
			id:     0,
			want:   69,
			wantOk: true,
		},
		{
			name: "empty",
		},
		{
			name:   "middle",
			items:  []int{1, 2, 3},
			id:     1,
			want:   2,
			wantOk: true,
		},
		{
			name:   "last",
			items:  []int{1, 2, 3},
			id:     2,
			want:   3,
			wantOk: true,
		},
		{
			name:  "negative",
			items: []int{1, 2, 3},
			id:    -1,
		},
		{
			name:  "post-EOF",
			items: []int{1, 2, 3},
			id:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lst := New[int](tt.items...)

			got, gotOk := lst.At(tt.id)
			if gotOk != tt.wantOk {
				t.Fatalf("List.At(%d) ok = %t; want %t", tt.id, gotOk, tt.wantOk)
			}

			// En los casos inválidos tt.want es el valor cero de T,
			// así que esto comprueba también el zero value documentado.
			if got != tt.want {
				t.Errorf("List.At(%d) = %d; want %d", tt.id, got, tt.want)
			}
		})
	}
}

// En caso de acertar devuelve []int
// En caso de fallar devuelve nil y un ErrOutOfBounds
func TestRange(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		start     int
		end       int
		want      []int
		wantError error
	}{
		{name: "void"},
		{
			name:  "empty",
			items: []int{1, 2, 3},
			start: 1,
			end:   1,
			want:  []int{},
		},
		{
			name:  "one",
			items: []int{69},
			start: 0,
			end:   1,
			want:  []int{69},
		},
		{
			name:  "middle",
			items: []int{1, 2, 3},
			start: 1,
			end:   2,
			want:  []int{2},
		},
		{
			name:  "middle to end",
			items: []int{1, 2, 3},
			start: 1,
			end:   3,
			want:  []int{2, 3},
		},
		{
			name:  "all the elements",
			items: []int{1, 2, 3},
			start: 0,
			end:   3,
			want:  []int{1, 2, 3},
		},
		{
			name:      "negative start",
			items:     []int{1, 2, 3},
			start:     -1,
			wantError: ErrOutOfBounds,
		},
		{
			name:      "start over end",
			items:     []int{1, 2, 3},
			start:     1,
			end:       0,
			wantError: ErrOutOfBounds,
		},
		{
			name:      "end post end",
			items:     []int{1, 2, 3},
			start:     0,
			end:       4,
			wantError: ErrOutOfBounds,
		},
		{
			name:      "both whatever",
			items:     []int{1, 2, 3},
			start:     -69,
			end:       69,
			wantError: ErrOutOfBounds,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lst := New[int](tt.items...)

			got, gotErr := lst.Range(tt.start, tt.end)

			if !errors.Is(gotErr, tt.wantError) {
				t.Fatalf(
					"List.Range(%d, %d) error = %v; want %v",
					tt.start,
					tt.end,
					gotErr,
					tt.wantError,
				)
			}

			if tt.wantError != nil {
				return
			}

			if !slices.Equal(got, tt.want) {
				t.Errorf("List.Range(%d, %d) = %v; want %v", tt.start, tt.end, got, tt.want)
			}
		})
	}
}

func TestRangePerformsShallowCopy(t *testing.T) {
	values := []int{1, 2, 3}
	lst := New[int](ptrs(values)...)

	got, err := lst.Range(0, 3)
	if err != nil {
		t.Fatalf("List.Range(0, 3) error = %v; want nil", err)
	}

	wantSharedPointers(t, "List.Range()", got, values)
}

// El slice devuelto es una copia: escribir en él no debe tocar la lista.
func TestRangeDoesNotAliasList(t *testing.T) {
	lst := New[int](1, 2, 3)

	got, err := lst.Range(0, 3)
	if err != nil {
		t.Fatalf("List.Range(0, 3) error = %v; want nil", err)
	}

	got[0] = 69
	if !slices.Equal(lst.items, []int{1, 2, 3}) {
		t.Errorf("List.Range() aliased the list: items = %v; want %v", lst.items, []int{1, 2, 3})
	}
}

func TestSlice(t *testing.T) {
	tests := []struct {
		name  string
		items []int
		want  []int
	}{
		{
			name: "empty",
		},
		{
			name:  "one",
			items: []int{69},
			want:  []int{69},
		},
		{
			name:  "three in a row",
			items: []int{1, 2, 3},
			want:  []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lst := New[int](tt.items...)
			got := lst.Slice()
			if !slices.Equal(got, tt.want) {
				t.Errorf("List.Slice() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestSlicePerformsShallowCopy(t *testing.T) {
	values := []int{1, 2, 3}
	lst := New[int](ptrs(values)...)

	wantSharedPointers(t, "List.Slice()", lst.Slice(), values)
}

// El slice devuelto es una copia: escribir en él no debe tocar la lista.
func TestSliceDoesNotAliasList(t *testing.T) {
	lst := New[int](1, 2, 3)

	got := lst.Slice()

	got[0] = 69
	if !slices.Equal(lst.items, []int{1, 2, 3}) {
		t.Errorf("List.Slice() aliased the list: items = %v; want %v", lst.items, []int{1, 2, 3})
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		name  string
		items []int
		want  int
	}{
		{name: "empty", want: 0},
		{name: "one", items: []int{69}, want: 1},
		{name: "three", items: []int{1, 2, 3}, want: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lst := New[int](tt.items...)
			if got, want := lst.Len(), tt.want; got != want {
				t.Errorf("List(%v).Len() = %d; want %d", tt.items, got, want)
			}
		})
	}
}

func TestEmpty(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		wantEmpty bool
	}{
		{name: "empty", wantEmpty: true},
		{name: "non-empty", items: []int{69}, wantEmpty: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lst := New[int](tt.items...)
			if got, want := lst.Empty(), tt.wantEmpty; got != want {
				t.Errorf("List(%v).Empty() = %t; want %t", tt.items, got, want)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name  string
		items []int
		want  string
	}{
		{name: "empty", want: "{}"},
		{name: "one", items: []int{69}, want: "{69}"},
		{name: "three", items: []int{1, 2, 3}, want: "{1, 2, 3}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lst := New[int](tt.items...)

			if got, want := lst.String(), tt.want; got != want {
				t.Errorf("List(%v).String() = %q; want %q", lst, got, want)
			}
		})
	}
}

func TestGoString(t *testing.T) {
	tests := []struct {
		name  string
		items []int
		want  string
	}{
		{name: "empty", want: "list.List[int,int]{}"},
		{name: "one", items: []int{69}, want: "list.List[int,int]{69}"},
		{name: "three", items: []int{1, 2, 3}, want: "list.List[int,int]{1, 2, 3}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lst := New[int](tt.items...)

			if got, want := lst.GoString(), tt.want; got != want {
				t.Errorf("List(%v).GoString() = %q; want %q", lst, got, want)
			}
		})
	}
}
