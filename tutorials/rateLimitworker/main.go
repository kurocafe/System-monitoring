package main

import (
	"fmt"
	"sync"
	"time"
)

type Result struct {
	id int
	value int
	job int
}

func worker(id int, jobCh chan int, resCh chan Result, wk *sync.WaitGroup){
	defer wk.Done()
	ticker := time.NewTicker(200 * time.Millisecond)
	for job := range jobCh{
		<- ticker.C
		resCh <- Result{
			id: id,
			value: job * job,
			job: job,
		}
	}
}

func main(){
	jobCh := make(chan int)
	resCh := make(chan Result)

	var wg sync.WaitGroup
	wg.Add(3)

	go worker(1, jobCh, resCh, &wg)
	go worker(2, jobCh, resCh, &wg)
	go worker(3, jobCh, resCh, &wg)

	go func() {
		jobs := []int{1,2,3,4,5,6,7,8,9,10}
		for _, v := range jobs {
			jobCh <- v
		}
		close(jobCh)
	}()

	go func(){
		wg.Wait()
		close(resCh)
	}()

	for r := range resCh {
		fmt.Printf("Worker %d job %d done (200ms)\t result: %d \n", r.id, r.job, r.value)
	}
}