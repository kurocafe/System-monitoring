package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)



func main(){
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(1000 * time.Millisecond, cancel)
	var wg sync.WaitGroup
	wg.Add(1)

	go func(){
		defer wg.Done()

		timer := time.NewTicker(100 * time.Millisecond)
		for i := 1; i <= 100; {
			select{
			case <- ctx.Done():
				fmt.Println("timer stopped!!")
				return
			case <- timer.C:
				fmt.Printf("%d \n", i)
				i ++
				continue
			default:
				continue
			}
		}
	}()
	
	wg.Wait()

}