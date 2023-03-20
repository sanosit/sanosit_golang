package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	var s string
	var size, countA, countB, countC, countD int
	fmt.Println("Введите размер строки: ")
	fmt.Fscan(os.Stdin, &size)
	fmt.Println("Введите строку: ")
	fmt.Fscan(os.Stdin, &s)
	ss := strings.Split(s, "")
	fmt.Println(ss)
	for i := 0; i < size; i++ {
		if ss[i] == "a" {
			countA = i + 1
		}
		if ss[i] == "b" {
			countB = i + 1
		}
		if ss[i] == "c" {
			countC = i + 1
		}
		if ss[i] == "d" {
			countD = i + 1
		}
	}
	fmt.Println(countA)
	fmt.Println(countB)
	fmt.Println(countC)
	fmt.Println(countD)

	if countA != 0 && countB != 0 && countC != 0 && countD != 0 {
		max := max4(countA, countB, countC, countD)
		min := min4(countA, countB, countC, countD)
		fmt.Println(max, min)
		len := max - min + 1
		fmt.Println(len)
	} else {
		fmt.Println(-1)
	}
}

func max4(a, b, c, d int) (max int) {
	max = a
	if max < b {
		max = b
	}
	if max < c {
		max = c
	}
	if max < d {
		max = d
	}
	return max
}
func min4(a, b, c, d int) (min int) {
	min = a
	if min > b {
		min = b
	}
	if min > c {
		min = c
	}
	if min > d {
		min = d
	}
	return min
}
