package main

import (
	"fmt"
	"os"
)

func main() {
	var a, b int
	fmt.Println("Введите число а: ")
	fmt.Fscan(os.Stdin, &a)      //Чтение с консоли
	if a < -32000 || a > 32000 { //проверка значения вводимого числа
		fmt.Println("Число а выходит за рамки условий")
		os.Exit(0)
	}
	fmt.Println("Введите число b: ")
	fmt.Fscan(os.Stdin, &b)
	if b < -32000 || b > 32000 {
		fmt.Println("Число b выходит за рамки условий")
		os.Exit(0)
	}
	c := a + b
	fmt.Print("Сумма а и b равна: ", c)

}
