package main

import "fmt"

// S = 2 × x + y ^ 2 − 3/z,
//где x — int16, y — uint8, a z — float32.
func MathCalc(x int16, y uint8, z float32) float32 {
	return 2*float32(x) + float32(y*y) - 3/z
}
func main() {
	fmt.Println(MathCalc(10, 20, 30.5))
}
