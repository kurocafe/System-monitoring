package main

//not working properly (idk why)
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

type Status struct {
	ID int
	Busy bool
}

// add auto-scaling worker function to this code

func worker(id int, jobCh <- chan Job, resCh chan <- Job, failCh chan <- Job, statusCh chan <- Status, stopCh <- chan struct{}, wk *sync.WaitGroup){
	defer wk.Done()
	for{
		select{
		case job, ok := <- jobCh:
			if !ok {
				fmt.Println("job channel closed...")
				return
			}
			attempts := 0
			RETRY:
				for attempts < 3{
					jobCtx, cancel := context.WithTimeout(context.Background(), 300 * time.Millisecond)
					done := make(chan bool, 1)
					
					statusCh <- Status{
						ID: job.Index,
						Busy: true,
					}

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
							statusCh <- Status{
								ID: id - 1,
								Busy: false,
							}
							break RETRY
						
						case <- stopCh:
							fmt.Printf("shutting down worker %d", id)
							cancel()
							return
						
					}

					if attempts == 3{
						fmt.Printf("worker %d failed job %d three times\n", id, job.Input)
						attempts = 999
						job.Output = -1
						failCh <- job
						statusCh <- Status{
							ID: id,
							Busy: false,
						}
					}
				}
			case <- stopCh:
				fmt.Printf("worker %d shutting down...", id)
				return
		}
		
	}
}

func main(){
	var wg sync.WaitGroup

	jobCh := make(chan Job)
	resCh := make(chan Job)
	failCh := make(chan Job)
	statusCh := make(chan Status)
	stopCh := make(map[int]chan struct{})

	maxworker := 10
	minworker := 2
	activeworker := 2
	idleTimers := make([]time.Time, maxworker)
	var mu sync.Mutex
	monitorCtx, stopMonitor := context.WithCancel(context.Background())
	
	go func(){
		for {
			time.Sleep(500 * time.Millisecond)

			select{
			case <- monitorCtx.Done():
				fmt.Println("monitor exiting....")
				return
			case status, ok := <- statusCh:
				if !ok {
					return
				}

				if !status.Busy {
					mu.Lock()
					idleTimers[status.ID] = time.Now()
					mu.Unlock()
				}

			default:
				backlog := len(jobCh)
				mu.Lock()
				aw := activeworker
				mu.Unlock()

				if backlog > 5 && aw < maxworker {
					//add new worker
					mu.Lock()
					activeworker ++
					id := activeworker
					stopCh[id - 1] = make(chan struct{})
					mu.Unlock()

					wg.Add(1)
					go worker(id, jobCh, resCh, failCh, statusCh, stopCh[id - 1], &wg)
					fmt.Println("Spawned new worker:", activeworker)
				}

				mu.Lock()
				currentIdle := idleTimers
				mu.Unlock()

				for id, last := range currentIdle {
					if time.Since(last) > 2 * time.Second && activeworker > minworker{
						stopCh[id] <- struct{}{}
						close(stopCh[id])
						delete(stopCh, id)

						mu.Lock()
						activeworker --
						mu.Unlock()
						fmt.Println("auto scale down...")
					}
				}
			}
		}
	}()

	for i := 1; i <= 2; i ++ {
		wg.Add(1)
		stopCh[i - 1] = make(chan struct{})
		go worker(i, jobCh, resCh, failCh, statusCh, stopCh[i - 1], &wg)
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
		stopMonitor()
		
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