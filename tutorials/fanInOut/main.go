package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func worker(id int, ctx context.Context, jobCh <- chan int, resCh chan <- int, wk *sync.WaitGroup){
	
	defer wk.Done()
	for {
		select{
			case <- ctx.Done():
				fmt.Printf("worker %d cancelled\n", id)
				return
			case v, ok := <- jobCh:
				time.Sleep(200 * time.Millisecond)
				if !ok {
					fmt.Println("channel closed")
					return
				}

				fmt.Printf("worker %d processing job %d\n", id, v)
				resCh <- v * v
		}
	}
}

func main(){
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(1 * time.Second, cancel)

	jobCh := make(chan int, 10)
	resCh := make(chan int, 10)

	var wg sync.WaitGroup
	for i := 1; i <= 4; i ++{
		wg.Add(1)
		go worker(i, ctx, jobCh, resCh, &wg)
	}
	
	go func(){
		for i := 1; i <= 50; i ++{
			select{
			case <- ctx.Done():
				close(jobCh)
				return
			default:
				jobCh <- i
			}
		}

		close(jobCh)
	}()

	go func(){
		wg.Wait()
		close(resCh)
	}()
	res := []int{}
	for result := range resCh {
		res = append(res, result)
	}

	fmt.Println("Colected result: ", res)
	fmt.Println("Total precessed: ", len(res))
}