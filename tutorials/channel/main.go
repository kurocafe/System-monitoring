package main

import (
	"fmt"
)

func square(n int, ch chan int) {
	ch <- n * n
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go square(3, ch1)
	go square(4, ch2)
	go square(5, ch3)

	fmt.Printf("%d, %d, %d \n", <-ch1, <-ch2, <-ch3)

}
