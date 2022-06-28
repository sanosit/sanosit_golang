package main

import (
	"fmt"
)

const aSize = 6

func bubbleSort(a [aSize]int) [aSize]int {
	for i := 0; i < len(a); i++ {
		for j := i; j < len(a); j++ {
			if a[i] > a[j] {
				tmp := a[i]
				a[i] = a[j]
				a[j] = tmp
			}
		}
	}
	return a
}
func main() {
	a := [aSize]int{0, 2, -1, 1, 3, 4}
	fmt.Println(a)
	b := bubbleSort(a)
	fmt.Println(b)
}
