package main

import "fmt"


func main(){
	numbers := []int{3, 1, 4, 1, 5, 9}
	result := []int{}

	for _, number := range numbers {
		if number == 1{
			continue
		}
		fmt.Println(number)
		result = append(result, number)
	}

	fmt.Println(result)
}