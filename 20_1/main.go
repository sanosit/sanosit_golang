package main

import "fmt"

const aSize = 3
const bSize = 3

func det(m [aSize][bSize]int64) int64 {
	x := m[1][1]*m[2][2] - m[2][1]*m[1][2]
	y := m[1][0]*m[2][2] - m[2][0]*m[1][2]
	z := m[1][0]*m[2][1] - m[2][0]*m[1][1]
	det := m[0][0]*x - m[0][1]*y + m[0][2]*z

	return det
}
func main() {
	matrix := [aSize][bSize]int64{
		{1, 0, 3},
		{4, 5, 6},
		{7, 0, 9},
	}
	fmt.Println(det(matrix))
}
