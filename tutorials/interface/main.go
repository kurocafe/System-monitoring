package main

import (
	"fmt"
	"math"
)

type Shape interface {
    Area() float64
}

type Circle struct {
	Radius float64
}

type Rectangle struct {
	W float64
	H float64
}

func (circle *Circle) Area() float64 {
	return circle.Radius * circle.Radius * math.Pi
}

func (rect *Rectangle) Area() float64 {
	return rect.H * rect.W
}

func main(){	
	shapes := []Shape{
    &Circle{Radius: 5},
    &Rectangle{W: 3, H: 4},
	}	

	for _, s := range shapes {
			fmt.Println(s.Area())
	}

}