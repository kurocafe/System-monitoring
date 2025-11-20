package main

import (
	"fmt"
	"time"
)

func printNumbers(id int) {
	for i := 1; i <= 5; i++ {
		fmt.Printf("Worker %d: %d \n", id, i)
		time.Sleep(200 * time.Millisecond)
	}
}

func main() {
	go printNumbers(1)
	go printNumbers(2)
	go printNumbers(3)
	time.Sleep(1 * time.Second)
}
