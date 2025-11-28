package main

import (
	"container/heap"
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Priority int
const (
	Low Priority = iota
	Medium
	High
)

type Job struct{
	ID int
	Priority Priority
	Input int
	Output int
}

type Log struct {
	ID int
	Time time.Time
	Busy bool
}

func main(){
	jobCh := make(chan Job, 500)
	resCh := make(chan Job, 500)
	logCh := make(chan Log)
	jobNum := 200
	var wg sync.WaitGroup
	var mu sync.Mutex
	var pendingJobs int64
	var completeJobs int64

	pq := &JobPQ{} 
	heap.Init(pq)

	active := 1
	maxWorker := 10
	workerCtxs := []context.Context{}
	cancels := []context.CancelFunc{}
	ctx, cancel := context.WithCancel(context.Background())
	idleTimes := make([]*Log, maxWorker)

	workerCtxs = append(workerCtxs, ctx)
	cancels = append(cancels, cancel)

	wg.Add(1)
	go worker(0, ctx, jobCh, resCh, logCh, &wg, &completeJobs)

	go func(){

		for i := 0; i < jobNum; i ++ {
			heap.Push(pq, Job{ID: i, Priority: Priority(rand.Intn(3))})
		}
		for i := 0; i < jobNum; i ++ {
			job := heap.Pop(pq).(Job)
			jobCh <- job

			atomic.AddInt64(&pendingJobs, 1)
		}

		close(jobCh)
	}()
	monitorCtx, stopMonitor := context.WithCancel(context.Background())
	logCtx, stopLog := context.WithCancel(context.Background())
	
	minWorker := 1
	go func() {
		for {
			select{
			case <- monitorCtx.Done():
				fmt.Println("stop monitoring workers...")
				return
			default:
				time.Sleep(1 * time.Second)
				backlogs := pendingJobs - completeJobs
				fmt.Printf("backlogs: %d\n", backlogs)

				if backlogs > 5 && active < maxWorker - 1	{
					mu.Lock()
					active ++
					id := active 
					ctx, cancel := context.WithCancel(context.Background())
					workerCtxs = append(workerCtxs, ctx)
					cancels = append(cancels, cancel)
					mu.Unlock()
					fmt.Println("create new worker !!!")
					wg.Add(1)
					go worker(id, ctx, jobCh, resCh, logCh, &wg, &completeJobs)
				}
				
				for i, l := range idleTimes{
					if l == nil{
						continue
					}

					fmt.Println("ID: ", idleTimes[i].ID)
					fmt.Println("Time: ", idleTimes[i].Time)
					fmt.Println("Busy: ", idleTimes[i].Busy)

					if time.Since(idleTimes[i].Time) > 2 * time.Second && !idleTimes[i].Busy && active > minWorker{
						fmt.Println("worker cancelled")

						mu.Lock()
						active --
						mu.Unlock()

						cancels[i]()
					} 
				}
			}
		}
	}()
	
	go func(){
		for {
			select{
			case <- logCtx.Done():
				fmt.Println("log cancelled!!")
				return
			case log, ok := <- logCh:
				if !ok{
					fmt.Println("log channel closed!!")
					return
				}

				idleTimes[log.ID] = &log
			}
		}
	}()

	go func ()  {
		wg.Wait()
		for i := range cancels{
			cancels[i]()
		}
		stopMonitor()
		stopLog()
		close(resCh)
	}()
	res := make([]int, jobNum)
	for r := range resCh{
		res[r.ID] = r.Output
	}
}