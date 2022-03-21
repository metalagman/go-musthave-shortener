package generics

import "fmt"

func isEven(e int) bool {
	return e%2 == 0
}

func ExampleSlice_Filter() {
	var si = Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	res := si.Filter(isEven)
	fmt.Println(res)
	// Output: &[2 4 6 8 10]
}
