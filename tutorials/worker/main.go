package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wk *sync.WaitGroup) {
	defer wk.Done()
	fmt.Printf("Worker %d Starting \n", id)
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("Worker %d Done\n", id)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	go worker(1, &wg)
	go worker(2, &wg)
	go worker(3, &wg)

	wg.Wait()

}
