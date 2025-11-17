package main

import "fmt"

func main(){
	for i := 1; i <= 100; i ++ {
		
		if i % 3 == 0{
			continue
		}

		fmt.Println(i)

		if i == 79 {
			break
		}
	}
}