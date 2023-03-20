package main

import (
	"fmt"
	"math"
	"os"
)

func main() {
	var n, m, k float64
	fmt.Println("Введите количество джунов: ")
	fmt.Fscan(os.Stdin, &n) //Чтение с консоли
	fmt.Println("Введите количество сеньоров: ")
	fmt.Fscan(os.Stdin, &m)
	fmt.Println("Введите количество необходимых проверок: ")
	fmt.Fscan(os.Stdin, &k)
	time := math.Ceil(n * k / m)
	fmt.Println(time)
}
