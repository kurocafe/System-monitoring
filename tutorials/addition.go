package main

func add(x int, y int) int{
    return x + y
}

func add_ez(x, y int) int{
	return x + y
}

func naked_return(x, y int) (a, b int){
	a = x + y
	b = x - y
	return 
}