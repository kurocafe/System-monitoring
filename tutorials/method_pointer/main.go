package main

import "fmt"

type Counter struct {
	value int
}

func (counter *Counter) Increment() {
	counter.value += 1
}

func (counter Counter) Value() int{
	return counter.value
}

func main(){
	counter := Counter{0}
	fmt.Printf("value: %d \n",counter.Value())
	counter.Increment()

	fmt.Printf("value: %d \n",counter.Value())
}

// pointer is needed for Increment()because it is modifying the number inside value. While Value() is just referencing the actual value inside the variable