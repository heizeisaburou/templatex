package ranges

import "fmt"

// TODO: Documentar funciones
// TODO: Asegurarse de que se accede siempre a copias de los datos

type Index interface {
	~int
}

type Range[T Index] struct {
	start T
	end   T
}

func New[T Index](start, end T) Range[T] {
	if start < 0 || end < start {
		panic(fmt.Sprintf("invalid range: %d:%d", start, end))
	}

	return Range[T]{
		start: start,
		end:   end,
	}
}

func (s Range[T]) Start() T {
	return s.start
}

func (s Range[T]) End() T {
	return s.end
}

func (s *Range[T]) SetStart(start T) {
	if start < 0 {
		panic(fmt.Sprintf("invalid range: %d:%d", start, s.end))
	}

	s.start = start
}

func (s *Range[T]) SetEnd(end T) {
	if end < s.start {
		panic(fmt.Sprintf("invalid range: %d:%d", s.start, end))
	}
	s.end = end
}

func (s *Range[T]) Set(start, end T) {
	if start < 0 || end < start {
		panic(fmt.Sprintf("invalid range: %d:%d", start, end))
	}

	s.start = start
	s.end = end
}

func (s *Range[T]) SetRange(rng Range[T]) {
	s.start = rng.start
	s.end = rng.end
}

func (s Range[T]) Len() T {
	return s.end - s.start
}

func (s Range[T]) Empty() bool {
	return (s.end - s.start) == 0
}

func (s Range[T]) String() string {
	return fmt.Sprintf("%d:%d", s.start, s.end)
}
