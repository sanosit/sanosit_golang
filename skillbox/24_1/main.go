package main

import "fmt"

const size = 10

func sortIns(input [size]int) [size]int {
	for i := 0; i < size; i++ {
		j := i
		for j > 0 {
			if input[j-1] > input[j] {
				input[j-1], input[j] = input[j], input[j-1]
			}
			j = j - 1
		}
	}
	return input
}
func main() {
	a := [size]int{10, 30, 20, 5, 8, 6, 1, 9, 0, 15}
	fmt.Println("До: ", a)
	b := sortIns(a)
	fmt.Println("После: ", b)
}
