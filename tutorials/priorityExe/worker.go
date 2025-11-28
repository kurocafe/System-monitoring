package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func worker(id int, ctx context.Context, jobCh <- chan Job, resCh chan <- Job, logCh chan <- Log, wk *sync.WaitGroup, completeJobs *int64){
	defer wk.Done()
	var mu sync.Mutex

	for {
		select{
		case <- ctx.Done():
			fmt.Println("shutting down worker....")
			return
		case job, ok := <- jobCh:
			if !ok {
				fmt.Println("job channel closed!!!")
				return
			}
			
			attempts := 0
			done := make(chan bool, 1)
			var jobWg sync.WaitGroup
			logCh <- Log{
				ID: id,
				Time: time.Now(),
				Busy: true,
			}

			for attempts < 3{
				jobWg.Add(1)
				go func(){
					defer jobWg.Done()
					time.Sleep(time.Duration(200 + rand.Intn(300)) * time.Millisecond)
					r := rand.Intn(100)

					if r < 30{
						mu.Lock()
						attempts ++
						mu.Unlock()	
					}

					done <- !(r < 30)
				}()

				fin := <- done

				if fin {
					fmt.Printf("worker %d finished processing job %d (priority = %d)\n", id, job.ID, job.Priority)
					job.Output = job.Input * job.Input
					resCh <- job
					atomic.AddInt64(completeJobs, 1)
					jobWg.Wait()
					break
				}else{
					time.Sleep(100 * time.Millisecond)
					continue
				}
			}

			logCh <- Log{
				ID: id,
				Time: time.Now(),
				Busy: false,
			}

			if attempts == 3{
				fmt.Printf("worker %d failed processing job %d (priority = %d) \n", id, job.ID, job.Priority)
				jobWg.Wait()
				continue
			}
		default:
			logCh <- Log{
				ID: id,
				Time: time.Now(),
				Busy: false,
			}
			continue
		}
	}
}