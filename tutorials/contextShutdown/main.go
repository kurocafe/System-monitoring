package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func dummy_worker(id int, ctx context.Context, wk *sync.WaitGroup){
	timer := time.NewTicker(100 * time.Millisecond)
	
	defer wk.Done()
	for {
		value := rand.Intn(100)
		select{
			case <- ctx.Done():
				fmt.Printf("worker %d cancelled\n", id)
				return
			case <- timer.C:
				fmt.Printf("worker %d processing job %d\n", id, value)
		}

	}
	
	
}

func main(){
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(1500 * time.Millisecond, cancel)
	var wg sync.WaitGroup
	wg.Add(5)

	go dummy_worker(1, ctx, &wg)
	go dummy_worker(2, ctx, &wg)
	go dummy_worker(3, ctx, &wg)
	go dummy_worker(4, ctx, &wg)
	go dummy_worker(5, ctx, &wg)

	wg.Wait()
}