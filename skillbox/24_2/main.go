package main

import "fmt"

const size = 10

func sortBubble(input [size]int) [size]int {
	for i := 0; i < size-1; i++ {
		if input[i] > input[i+1] {
			input[i], input[i+1] = input[i+1], input[i]
		}
	}
	return input
}
func reverse(input [size]int) [size]int {
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		input[i], input[j] = input[j], input[i]
	}
	return input
}
func merge(input [size]int) [size]int {
	input = sortBubble(input)
	input = reverse(input)
	return input
}
func main() {
	a := [size]int{10, 30, 20, 5, 8, 6, 1, 9, 0, 15}
	fmt.Println("До: ", a)
	b := merge(a)
	fmt.Println("После: ", b)
}
