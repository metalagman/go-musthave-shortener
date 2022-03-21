package generics

import "fmt"

func ExampleSlice_Map() {
	var si = Slice[int]{1, 2, 3, 4, 5}
	sum := si.Reduce(0, Sum[int])
	fmt.Println(sum) // 15
	// теперь цепочку
	res := si.Map(Double[int]).Reduce(0, Sum[int])
	fmt.Println(res) // 30
	// теперь для строк
	var ss = Slice[string]{"foo", "bar", "buzz"}
	res1 := ss.Map(Double[string]).Reduce("", Sum[string])
	fmt.Println(res1)
	// Output: 15
	// 30
	// foofoobarbarbuzzbuzz
}
