package main

import (
	"fmt"
	"sync"
)

func worker(id int, req chan int, res chan int, wk *sync.WaitGroup) {
	defer wk.Done()
	val := <-req
	fmt.Printf("Worker %d processing Job %d \n", id, val)
	res <- val * val
}

func main() {
	jobCh := make(chan int)
	resCh := make(chan int)
	var wg sync.WaitGroup
	wg.Add(10)

	go worker(1, jobCh, resCh, &wg)
	go worker(2, jobCh, resCh, &wg)
	go worker(3, jobCh, resCh, &wg)

	go func() {
		for i := 1; i <= 10; i++ {
			jobCh <- i
		}
	}()

	for i := 0; i < 10; i++ {
		fmt.Printf("Result: %d\n", <-resCh)
	}

	wg.Wait()
}
