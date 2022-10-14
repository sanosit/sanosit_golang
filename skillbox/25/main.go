package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	var str string
	var substr string

	flag.StringVar(&str, "str", "string", "start string")
	flag.StringVar(&substr, "substr", "substring", "end string")

	flag.Parse()

	fmt.Println(str, substr)
	cont(&str, &substr)
}

func cont(a, b *string) {
	contain := strings.Contains(*a, *b)
	fmt.Print(contain)
}
