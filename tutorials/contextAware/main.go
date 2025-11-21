package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func worker(ctx context.Context, id int, jobCh <-chan int, resCh chan<- int, wg *sync.WaitGroup){
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d STOPPED\n", id)
			return

		case v, ok := <- jobCh:

			if !ok {
				return
			}
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("Worker %d processing job %d \n", id, v)
			resCh <- v * v
		}
	}
}

func main(){
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(500 * time.Millisecond, cancel)

	jobCh := make(chan int)
	resCh := make(chan int)
	var wg sync.WaitGroup
	wg.Add(3)

	go worker(ctx, 1, jobCh, resCh, &wg)
	go worker(ctx, 2, jobCh, resCh, &wg)
	go worker(ctx, 3, jobCh, resCh, &wg)

	go func() {
		jobs := []int{1,2,3,4,5,6,7,8,9,10}
		for _, v := range jobs{
			select {
				case <- ctx.Done():
					close(jobCh)
					return
				default:
					jobCh <- v
			}
		}

		close(jobCh)
	}()

	go func() {
		wg.Wait()
		close(resCh)
	}()

	for r := range resCh {
		fmt.Printf("Result: %d\n", r)
	}

}