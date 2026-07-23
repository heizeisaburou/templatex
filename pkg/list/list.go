package list

import "fmt"

// TODO: Crear String()
// TODO: Asegurarse de que se accede siempre a copias de los datos
// TODO: Documentar funciones

type Index interface {
	~int
}

type List[I Index, T any] struct {
	items []T
}

func New[I Index, T any](items ...T) List[I, T] {
	return List[I, T]{
		items: append([]T(nil), items...),
	}
}

func (l *List[I, T]) Push(item T) I {
	nextID := I(len(l.items))
	l.items = append(l.items, item)
	return nextID
}

func (l *List[I, T]) Set(id I, item T) {
	if id < 0 || int(id) >= len(l.items) {
		panic(fmt.Sprintf("list index out of range: %d", id))
	}

	l.items[id] = item
}

func (l *List[I, T]) Truncate(end I) {
	if end < 0 || int(end) > len(l.items) {
		panic(fmt.Sprintf("invalid list truncate: %d", end))
	}

	l.items = l.items[:end]
}

func (l *List[I, T]) ReplaceRange(start, end I, repl ...T) {
	if start < 0 || end < start || int(end) > len(l.items) {
		panic(fmt.Sprintf("invalid list range: %d:%d", start, end))
	}

	removed := int(end - start)
	inserted := len(repl)

	if removed == inserted {
		copy(l.items[start:end], repl)
		return
	}

	next := make([]T, 0, len(l.items)-removed+inserted)
	next = append(next, l.items[:start]...)
	next = append(next, repl...)
	next = append(next, l.items[end:]...)
	l.items = next
}

func (l List[I, T]) At(id I) (T, bool) {
	var zero T
	if id < 0 || int(id) >= len(l.items) {
		return zero, false
	}

	return l.items[id], true
}

func (l List[I, T]) MustAt(id I) T {
	item, ok := l.At(id)
	if !ok {
		panic(fmt.Sprintf("list index out of range: %d", id))
	}

	return item
}

func (l List[I, T]) Range(start, end I) []T {
	if start < 0 || end < start || int(end) > len(l.items) {
		panic(fmt.Sprintf("invalid list range: %d:%d", start, end))
	}

	return l.items[start:end]
}

func (l List[I, T]) Slice() []T {
	return append([]T(nil), l.items...)
}

func (l List[I, T]) Len() I {
	return I(len(l.items))
}

func (l List[I, T]) Empty() bool {
	return len(l.items) == 0
}

func (l List[I, T]) NextID() I {
	return I(len(l.items))
}

