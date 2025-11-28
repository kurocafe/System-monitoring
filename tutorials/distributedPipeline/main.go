package main

//maybe finished???
import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
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

func abs(value int) int{
	result := value
	if value < 0{
		result = - result
	}
	return result
}

func validationWorker(id int, jobCh <- chan Job, resCh chan <- Job, failCh chan <- Job, stopCh <- chan struct{}, statusCh chan <- Status, wk *sync.WaitGroup,  completedJobs *int64){
	defer wk.Done()
	for {
		select{
		case job, ok := <- jobCh:
			
			if !ok{
				fmt.Println("job channel closed!!")
				return
			}
			
			statusCh <- Status{
				ID: id,
				Busy: true,
			}

			fmt.Printf("Worker %d: processing job %d\n", id, job.Index)
			attempts := 0

			for attempts < 3{
				time.Sleep(time.Duration(100 + rand.Intn(200)) * time.Millisecond)
				
				if rand.Intn(100) < 50 {
					attempts++
					fmt.Printf("Worker %d: failed job %d (attempts = %d)\n", id, job.Index, attempts)
					continue
				}
				fmt.Printf("Worker %d: completed job %d\n", id, job.Index)
				resCh <- Job{
					Index: job.Index,
					Input: job.Input,
					Output:abs(job.Input),
				}
				atomic.AddInt64(completedJobs, 1)
				
				statusCh <- Status{
					ID: id,
					Busy: false,
				}

				break
			}

			if attempts == 3{
				attempts = 999
				fmt.Printf("worker %d job %d failed completely....\n", id, job.Index)
				failCh <- Job{
					Index: job.Index,
					Input: job.Input,
					Output: job.Output,
				}
				atomic.AddInt64(completedJobs, 1)
				statusCh <- Status{
					ID: id,
					Busy: false,
				}
				continue
			}
		
		case <- stopCh:
			fmt.Printf("worker %d shutting down", id)
			return
		}
	}
}


//create squareWoker() that retrieves value from valiCh and send squared value to squareCh

func main(){
	var wg sync.WaitGroup
	var pendingJobs int64
	var completedJobs int64

	jobCh := make(chan Job, 200)
	stopCh := make(map[int]chan struct{})
	statusCh := make(chan Status)

	valiCh := make(chan Job, 50)
	failCh := make(chan Job, 50)

	st1Min := 2
	st1Max := 10
	st1Active := 0
	idleTimers := make([]time.Time, st1Max)
	var mu sync.Mutex

	monitorSt1Ctx, stopmonitor := context.WithCancel(context.Background())

	go func(){
		for {
			select{
			case <- monitorSt1Ctx.Done():
				fmt.Println("stop monitoring stage1...")
				return
			case s, ok := <- statusCh:
				if !ok {
					fmt.Println("status channel closed!!")
					return
				}

				if s.Busy{
					fmt.Printf("worker %d has just started working!\n", s.ID)
					continue
				}

				idleTimers[s.ID] = time.Now()
				
			default:
				time.Sleep(500 * time.Millisecond)
				backlog := atomic.LoadInt64(&pendingJobs) - atomic.LoadInt64(&completedJobs)
				mu.Lock()
				count := st1Active
				mu.Unlock()
				fmt.Println("len job: ", backlog)
				if backlog > 5 && count < st1Max - 1{
					
					mu.Lock()
					st1Active ++
					id := st1Active
					mu.Unlock()

					wg.Add(1)
					stopCh[id] = make(chan struct{})
					go validationWorker(id, jobCh, valiCh, failCh, stopCh[id], statusCh, &wg, &completedJobs)
					fmt.Printf("new stage 1 worker initialized. id = %d\n", id)
				}

				for i, t := range idleTimers {
					if time.Since(t) > 2 * time.Second && st1Active > st1Min{
						select {
							case stopCh[i] <- struct{}{}: // send stop signal
							default:
						}
						
						mu.Lock()
						id := st1Active
						st1Active--
						mu.Unlock()

						fmt.Printf("stage 1 worker terminated. id = %d\n", id) 
					}
				}
			}
		}
	}()

	for i := 0; i < 2; i ++ {
		wg.Add(1)
		stopCh[i] = make(chan struct{})
		go validationWorker(i, jobCh, valiCh, failCh, stopCh[i], statusCh, &wg, &completedJobs)
	}

	// send jobs
	go func() {
			for i := 1; i <= 200; i++ {
					jobCh <- Job{Index: i - 1, Input: i}
					atomic.AddInt64(&pendingJobs, 1)
			}
			close(jobCh)
	}()

	// wait for workers
	go func() {
			wg.Wait()
			close(valiCh)
			close(failCh)
			close(statusCh)
			for id, ch := range stopCh {
					close(ch) // safe to close now because no one is sending anymore
					fmt.Printf("[shutdown] closed stop channel for worker %d\n", id)
			}
			stopmonitor()
	}()

	res := make([]int, 0, 50)
	fail := make([]int, 0, 50)

	for v := range valiCh {
			res = append(res, v.Output)
	}
	for f := range failCh {
			fail = append(fail, f.Index)
	}

	fmt.Println("final result: ", res)
	fmt.Println("fail indexes: ", fail)
}