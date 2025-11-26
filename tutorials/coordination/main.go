package main

import (
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

func worker(id int, jobCh <- chan Job, resCh chan <- Job, wk *sync.WaitGroup){
	defer wk.Done()

	for job := range jobCh {
		time.Sleep(time.Duration(200 + rand.Intn(300)) * time.Millisecond)
		job.Output = job.Input * job.Input
		fmt.Printf("worker %d processing job %d -> %d\n", id, job.Input, job.Output)

		resCh <- job
	}
}

func main(){
	var wg sync.WaitGroup

	jobCh := make(chan Job)
	resCh := make(chan Job)

	for i := 1; i <= 3; i ++ {
		wg.Add(1)
		go worker(i, jobCh, resCh, &wg)
	}

	go func() {
		for i := 1; i <= 10; i ++ {
			job := Job{
				Index: i - 1,
				Input: i,
				Output: 0,
			}

			jobCh <- job
		}
		close(jobCh)
	}()

	go func(){
		wg.Wait()
		close(resCh)
	}()
	res := make([]int, 10)
	for result := range resCh {
		res[result.Index] = result.Output 
	}

	fmt.Println("final output result: ", res)
	
}