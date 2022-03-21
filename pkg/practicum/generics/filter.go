package generics

func (s *Slice[T]) Filter(allow func(e T) bool) *Slice[T] {
	for i := len(*s) - 1; i >= 0; i-- {
		if !allow((*s)[i]) {
			*s = (*s)[:i+copy((*s)[i:], (*s)[i+1:])]
		}
	}
	return s
}
