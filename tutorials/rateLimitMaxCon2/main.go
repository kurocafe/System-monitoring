package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(job int, wk *sync.WaitGroup){
	defer wk.Done()
	fmt.Printf("result: %d\n", job * 2)
}

func main(){
	sem := make(chan int, 3)
	var wg sync.WaitGroup
	wg.Add(30)

	go func(){
		ticker := time.NewTicker(1 * time.Second)
		cnt := 0
		jobs := make([]int, 30)
		for i := 0; i < 30; i ++{
			jobs[i] = i + 1
		}
		
		for i := 0; i <= 30; {
			
			select {
				case <- ticker.C:
					cnt = 0
				case v, ok := <-sem:
					if !ok{
						return
					}
					go worker(v, &wg)
				default:
					if cnt >= 5{
						continue
					}
					sem <- jobs[i]
					cnt += 1
					i ++
			}
		}
		defer close(sem)
	}()

	wg.Wait()

}