package main

import (
	"fmt"
	"sync"
)

func stage1Worker(id int, in <- chan int, out chan <- int, wk *sync.WaitGroup){
	defer wk.Done()
	for i := range in{
		fmt.Printf("stage 1: worker %d processing value %d -> %d\n" , id, i, i * i)
		out <- i * i
	}
}

func stage2Worker(id int, in <- chan int, out chan <- int, wk *sync.WaitGroup){
	defer wk.Done()

	for i := range in{
		fmt.Printf("stage 2: worker %d processing value %d -> %d\n", id, i, i * 2)
		out <- i * 2
	}
}

func main(){
	
	numbersCh := make(chan int)
	squareCh := make(chan int)
	resultCh := make(chan int)

	var stage1Wg sync.WaitGroup
	var stage2Wg sync.WaitGroup
	
	stage1Wg.Add(2)
	stage2Wg.Add(2)
	
	go func(){
		for i := 1; i <= 20; i++{
			numbersCh <- i
		}
		close(numbersCh)
	}()
	
	go stage1Worker(1, numbersCh, squareCh, &stage1Wg)
	go stage1Worker(2, numbersCh, squareCh, &stage1Wg)
	go stage2Worker(1, squareCh, resultCh, &stage2Wg)
	go stage2Worker(2, squareCh, resultCh, &stage2Wg)
	
	go func(){
		stage1Wg.Wait()
		close(squareCh)
	}()

	go func(){
		stage2Wg.Wait()
		close(resultCh)
	}()

	for result := range resultCh{
		fmt.Printf("final result %d\n", result)
	} 

	
	
	
}