package main

import (
	"fmt"
	"math/rand"
	"time"
)

const n = 10

func CountAfter(a [n]int, value int) (count int) {
	count = 0
	for i := 0; i < n; i++ {
		if a[i] == value {
			count++
			fmt.Print("Искомое число: ")
		}
		if count != 0 {
			count++
		}
		fmt.Println(a[i])
	}
	return count - 2
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var a [n]int
	for i := 0; i < n; i++ {
		a[i] = rand.Intn(10 * n)
	}
	value := a[5] // Задаем искомое число, для простоты проверки, оно задано по массиву.
	count := CountAfter(a, value)
	fmt.Printf("Количество чисел после искомого: %v", count)

}
