package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Job struct {
	Index int
	Input int
	Output int
}

// add auto-scaling worker function to this code

func worker(id int, jobCh <- chan Job, resCh chan <- Job, failCh chan <- Job, wk *sync.WaitGroup){
	defer wk.Done()
	for job := range jobCh{
		
		attempts := 0
		RETRY:
			for attempts < 3{
				jobCtx, cancel := context.WithTimeout(context.Background(), 300 * time.Millisecond)
				done := make(chan bool, 1)
				go func (){
					time.Sleep(time.Duration(rand.Intn(430)) * time.Millisecond)
					done <- true
				}()

				select{
					case <- jobCtx.Done():
						fmt.Printf("worker %d job %d failed (attempts = %d)\n", id, job.Input, attempts)
						attempts ++
						cancel()
					case <- done:
						fmt.Printf("worker %d job %d success\n", id, job.Input)
						cancel()
						attempts = 0
						job.Output = job.Input * job.Input
						resCh <- job
						break RETRY
				}

				if attempts == 3{
					fmt.Printf("worker %d failed job %d three times", id, job.Input)
					attempts = 999
					job.Output = -1
					failCh <- job
				}
			}
	}
}

func main(){
	var wg sync.WaitGroup

	jobCh := make(chan Job)
	resCh := make(chan Job)
	failCh := make(chan Job)

	for i := 1; i <= 2; i ++ {
		wg.Add(1)
		go worker(i, jobCh, resCh, failCh, &wg)
	}

	go func(){
		for i := 1; i <= 50; i ++ {
			jobCh <- Job{
				Index: i - 1,
				Input: i,
				Output: 0,
			}
		}

		close(jobCh)
	}()

	go func(){
		wg.Wait()
		close(resCh)
		close(failCh)
	}()
	
	result := make([]int, 50)
	fail := make([]int, 0, 25)
	close := 0
	for close < 2 {
		select{
		case res, ok := <- resCh:
			if !ok {
				close ++
				continue
			}
			result[res.Index] = res.Output
		case f, ok := <- failCh:
			if !ok {
				close ++
				continue
			}
			fail = append(fail, f.Index)
		}
	}

	fmt.Println("final result: ", result)
	fmt.Println("failed indexes: ", fail)
}