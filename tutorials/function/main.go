package main

import "fmt"

func main(){
	example := []int{1,2,3,4}
	result := Filter(example, func(i int) bool {return i % 2 == 1})
	fmt.Println(result)
}