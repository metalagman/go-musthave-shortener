package generics

import "fmt"

func ExampleListIter() {
	// конструируем List
	li := new(List[int])
	// наполняем
	for i := 1; i < 7; i++ {
		nl := new(List[int])
		nl.value = i
		nl.next = li
		li = nl
	}
	// конструируем итератор для списка
	var iter = Iter[int]{NewListIter(li)}
	// смотрим, что получилось
	for nxt, ok := iter.Next(); ok; nxt, ok = iter.Next() {
		fmt.Println(nxt)
	}
	// применяем метод Map()
	iter.Map(&lit, Double[int])
	// смотрим, что получилось
	fmt.Println("After Mapping")
	for nxt, ok := iter.Next(); ok; nxt, ok = iter.Next() {
		fmt.Println(nxt)
	}
}
