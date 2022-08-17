package main

import (
	"fmt"
	"strings"
)

func main() {
	sentences := [2]string{"hello", "world"}
	chars := [3]rune{'h', 'e', 'l'}
	fmt.Printf("%v\n", sentences)
	fmt.Printf("%q\n", chars)
	fmt.Println()
	for i := 0; i < len(sentences); i++ {
		for j := 0; j < len(chars); j++ {
			tempIndex := strings.Index(sentences[i], string(chars[j]))
			fmt.Printf("Word: %v", sentences[i])
			fmt.Printf(". Letter: %q", chars[j])
			fmt.Printf(". Index: %v\n", tempIndex)
		}
	}
}
