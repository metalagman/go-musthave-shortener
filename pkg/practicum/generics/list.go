package generics

type Iterator[V any] interface {
	Next() (V, bool)
	Set(V)
}

func Map[I Iterator[V], V any](it I, f func(V) V) {
	for nxt, ok := it.Next(); ok; nxt, ok = it.Next() {
		it.Set(f(nxt))
	}
}

type Iter[I any] struct {
	Iterator[I]
}

func (iter *Iter[V]) Map(f func(V) V) *Iter[V] {
	for nxt, ok := iter.Next(); ok; nxt, ok = iter.Next() {
		iter.Set(f(nxt))
	}
	return iter
}

type List[T any] struct {
	next  *List[T]
	value T
}

type ListIter[T any] struct {
	lst *List[T]
	cur *List[T]
}

func NewListIter[V any](l List[V]) *ListIter[V] {
	var i ListIter[V]
	i.lst = *l
	i.cur = *l
	return *i
}

func (l *ListIter[T]) Next() (T, bool) {
	if l.cur.next != nil {
		l.cur = l.cur.next
		return l.cur.value, true
	}
	// перезагружаем итератор
	// для повторного использования
	l.cur = l.lst
	return l.cur.value, false
}

func (l *ListIter[T]) Set(v T) {
	l.cur.value = v
}
