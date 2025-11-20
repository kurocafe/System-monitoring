package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string, 2)

	go func() {
		ch <- "A"
		ch <- "B"
		ch <- "C"
	}()

	for i := 0; i < 3; i++ {
		time.Sleep(500 * time.Millisecond)
		fmt.Println(<-ch)
	}
}
