package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Println("Press ctrl+c to stop")
	go func() {
		var i int
		for {
			select {
			case <-done:
				fmt.Println("Exit")
				return
			default:
				i = i + 1
				fmt.Println(i * i)
				time.Sleep(time.Second)
			}
		}
	}()
	wg.Wait()
}
