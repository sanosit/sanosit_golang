package main

import "fmt"

func main() {
	defer fmt.Println("Рассчёт окончен")
	A := func(a, b int) int { return a + b }
	fmt.Println(A(1, 2))
	fmt.Println(A(3, 4))
	fmt.Println(A(6, 8))
}
