package main

import (
	"fmt"
)

const n = 12

func main() {
	var a [n]int = [n]int{1, 2, 2, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println(a)
	value := 2
	index := find(a, value)
	fmt.Printf("Искомое число: %v\n", value)
	fmt.Printf("Индекс: %v\n", index)
}

func find(a [n]int, value int) (index int) {
	index = -1
	min := 0
	max := n - 1
	for max >= min {
		middle := (max + min) / 2
		if a[middle] == value {
			index = middle
			for i := 0; i < middle; i++ {
				if a[middle-i] == value {
					index = middle - i
				}
			}
			break
		} else if a[middle] > value {
			max = middle - 1
		} else {
			min = middle + 1
		}
	}
	return
}
