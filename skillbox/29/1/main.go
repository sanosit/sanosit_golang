package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)
	fc := input(3, &wg)
	sc := square(fc, &wg)
	tc := multi(sc, &wg)
	fmt.Print("Произведение: ", <-tc)
}

func input(a int, wg *sync.WaitGroup) chan int {
	defer wg.Done()
	firstChan := make(chan int)
	go func() {
		firstChan <- a
		close(firstChan)
	}()
	fmt.Print("Ввод: ", a)
	fmt.Println()
	return firstChan
}
func square(firstChan chan int, wg *sync.WaitGroup) chan int {
	defer wg.Done()
	secondChan := make(chan int)
	go func() {
		for {
			a, ok := <-firstChan
			if !ok {
				break
			}
			b := a * a
			secondChan <- b
			fmt.Print("Квадрат: ", b)
			fmt.Println()
		}
		close(secondChan)
	}()
	return secondChan
}
func multi(secondChan chan int, wg *sync.WaitGroup) chan int {
	defer wg.Done()
	thirdChan := make(chan int)
	go func() {
		for {
			a, ok := <-secondChan
			if !ok {
				break
			}
			b := a * 2
			thirdChan <- b
		}
		close(thirdChan)
	}()
	return thirdChan
}
