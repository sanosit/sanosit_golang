package main

import (
	"fmt"
)

const size = 10

func main() {
	var a [size]int = [size]int{10, 20, 30, 40, 50, 11, 21, 31, 41, 51}
	Even, Uneven := EvenUnevenSort(a)
	fmt.Printf("%v", a)
	fmt.Println()
	fmt.Printf("%v", Even)
	fmt.Println()
	fmt.Printf("%v", Uneven)
}
func EvenUnevenSort(a [size]int) (b [size / 2]int, c [size / 2]int) {
	var countEven, countUneven int
	for i := 0; i < size; i++ {
		if a[i]%2 == 0 {
			b[countEven] = a[i]
			countEven++
		} else {
			c[countUneven] = a[i]
			countUneven++
		}
	}
	return
}
