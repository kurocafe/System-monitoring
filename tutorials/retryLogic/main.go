package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func worker(id int, jobCh <- chan int, resCh chan <- int, failCh chan <- int, wk *sync.WaitGroup){
	defer wk.Done()
	
	for job := range jobCh {
		attempt := 0
		RETRY:
			for attempt < 3{
				ctx, cancel := context.WithTimeout(context.Background(), 300 * time.Millisecond)
				done := make(chan bool)
				
				go func(){
					time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

					done <- true
				}()

				select{
				case <- ctx.Done():
					attempt ++
					fmt.Printf("worker %d job %d failed\n", id, job)
					cancel()
				case <- done:
					attempt = 0
					fmt.Printf("Worker %d job %d success\n", id, job)
					cancel()
					resCh <- job * job
					break RETRY
				}

				if attempt == 3{
					fmt.Printf("worker %d job %d completely failed\n", id, job)
					failCh <- job
					attempt = 999
				}
			}
	}
}

func main(){
	jobCh := make(chan int)
	resCh := make(chan int)
	failCh := make(chan int)

	var wg sync.WaitGroup
	for i := 1; i <= 3; i ++ {
		wg.Add(1)
		go worker(i, jobCh, resCh, failCh, &wg)
	}

	go func() {
		for i := 1; i <= 20; i ++{
			jobCh <- i
		}

		close(jobCh)
	}()

	go func() {
		wg.Wait()
		close(resCh)
		close(failCh)
	}()

	res := []int {}

	fail := []int{}
	closed := 0

	for closed < 2 {
		select {
		case r, ok := <-resCh:
			if !ok {
				closed++
				continue
			}
			res = append(res, r)

		case f, ok := <-failCh:
			if !ok {
				closed++
				continue
			}
			fail = append(fail, f)
		}
	}

	

	fmt.Println("Processed results: ", res)
	fmt.Println("failed job: ", fail)
	fmt.Println("Total processed: ", len(res))
	fmt.Println("total failed: ", len(fail))
}	