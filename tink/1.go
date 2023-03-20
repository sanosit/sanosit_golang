package main

import (
	"fmt"
	"os"
)

func main() {
	var h1, h2, h3, h4 int
	fmt.Println("Введите рост первого человека: ")
	fmt.Fscan(os.Stdin, &h1) //Чтение с консоли
	fmt.Println("Введите рост второго человека: ")
	fmt.Fscan(os.Stdin, &h2)
	fmt.Println("Введите рост третьего человека: ")
	fmt.Fscan(os.Stdin, &h3)
	fmt.Println("Введите рост четвертого человека: ")
	fmt.Fscan(os.Stdin, &h4)
	if h1 >= h4 {
		if h1 >= h2 {
			if h2 >= h3 {
				if h3 >= h4 {
					fmt.Print("YES")
					os.Exit(0)
				}
			}
		}
	} else if h1 <= h2 {
		if h2 <= h3 {
			if h3 <= h4 {
				fmt.Print("YES")
				os.Exit(0)
			}
		}
	}
	fmt.Println("NO")
}
