package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, jobCh <-chan int, resCh chan<- int, wk *sync.WaitGroup){
	defer wk.Done()
	for {
		v, ok := <- jobCh
		if !ok {
			return
		}
		fmt.Printf("Worker %d processing Job %d \n", id, v)
		resCh <- v * 2
	}
}

func main(){
	jobCh := make(chan int)
	resCh := make(chan int, 5)
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		cnt := 0
		jobs := make([]int, 30)
		for i := 0; i < 30; i ++{
			jobs[i] = i + 1
		}
		
		for i := 0; i < 30; {
			select {
			case <- ticker.C:
				cnt = 0
			default:
				if cnt >= 5{
					continue
				}
				jobCh <- jobs[i]
				cnt += 1
				i ++
			}
		}
		defer close(jobCh)
	}()

	go worker(1, jobCh, resCh, &wg)
	go worker(2, jobCh, resCh, &wg)
	go worker(3, jobCh, resCh, &wg)
	
	go func() {
		wg.Wait()
		close(resCh)
	}()

	for r := range resCh{
		fmt.Printf("result: %d\n", r)
	}
}