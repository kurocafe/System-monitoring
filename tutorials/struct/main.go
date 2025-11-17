package main

import "fmt"

type Student struct {
    Name string
    Age  int
    Grade float64
}

func (student *Student) IsAdult() bool {
	
	is_adult := false
	if student.Age >= 18 {
		is_adult = true
	}

	return is_adult
}

func main(){
	child := Student{Name: "nachi", Age: 13, Grade: 5.0}

	fmt.Println(child.IsAdult())

	adult := Student{Name: "gentleman", Age: 30, Grade: 6.0}
	fmt.Println(adult.IsAdult())
}