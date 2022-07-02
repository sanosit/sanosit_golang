package main

import "fmt"

const rowsOne = 3
const colsOne = 5
const rowsTwo = 5
const colsTwo = 4

func multiply(matOne [rowsOne][colsOne]int, matTwo [rowsTwo][colsTwo]int) [rowsOne][colsTwo]int {
	var m [rowsOne][colsTwo]int
	for i := 0; i < rowsOne; i++ {
		for j := 0; j < colsTwo; j++ {
			for k := 0; k < rowsTwo; k++ {
				m[i][j] += matOne[i][k] * matTwo[k][j]

			}
		}
	}
	return m
}
func main() {
	matrixOne := [rowsOne][colsOne]int{
		{1, 1, 1, 1, 1},
		{2, 2, 2, 2, 2},
		{3, 3, 3, 3, 3},
	}
	matrixTwo := [rowsTwo][colsTwo]int{
		{4, 4, 4, 4},
		{5, 5, 5, 5},
		{6, 6, 6, 6},
		{7, 7, 7, 7},
		{8, 8, 8, 8},
	}
	fmt.Println(multiply(matrixOne, matrixTwo))

}
