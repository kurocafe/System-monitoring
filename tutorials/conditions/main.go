package main

import "fmt"

func main(){
	var n int
	fmt.Println("enter a random number:")
	fmt.Scan(&n)

	isEven := IsEven(n)
	if isEven {
		fmt.Println("Even")
	}else{
		fmt.Println("Odd")
	}

}