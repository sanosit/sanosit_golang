package main

import (
	"fmt"
)

const sizeOne = 4
const sizeTwo = 5
const sizeMerged = sizeOne + sizeTwo

func combine(a [sizeOne]int, b [sizeTwo]int) [sizeMerged]int {
	var merged [sizeMerged]int
	j, k := 0, 0
	for i := 0; i < sizeMerged; i++ {
		if j >= len(a) {
			merged[i] = b[k]
			k++
			continue
		} else if k >= len(b) {
			merged[i] = a[j]
			j++
			continue
		}
		if a[j] > b[k] {
			merged[i] = b[k]
			k++
		} else {
			merged[i] = a[j]
			j++
		}
	}
	return merged
}
func main() {

	a := [sizeOne]int{1, 3, 5, 7}
	b := [sizeTwo]int{2, 4, 6, 8, 9}

	fmt.Println(a)
	fmt.Println(b)

	c := combine(a, b)
	fmt.Println(c)
}
