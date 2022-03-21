package generics

import "golang.org/x/exp/constraints"

type Slice[T any] []T

func (s *Slice[T]) Map(f func(T) T) *Slice[T] {
	for k, v := range *s {
		(*s)[k] = f(v)
	}
	return s
}

func (s *Slice[T]) Reduce(r T, f func(a, e T) T) T {
	for _, v := range *s {
		r = f(r, v)
	}
	return r
}

func Sum[T constraints.Ordered](a, e T) T {
	return a + e
}

func Double[T constraints.Ordered](v T) T {
	return v + v
}
