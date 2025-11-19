package main

import "fmt"

type Worker interface {
	Work() string
}

type Programmer struct {
	Language string
}

type Designer struct {
	Tool string
}

func (p Programmer) Work() string {
	
	work := "this guy is working on " + p.Language + " project"
	
	return work
}

func (d Designer) Work() string {
	work := "this gay is working with " + d.Tool
	return work
}

func main(){

	workers := []Worker{
		Programmer{
			Language: "Go lang",
		},
		 Designer {
			Tool: "gomi",
		},
	}

	for _, worker := range workers {
		fmt.Println(worker.Work())
	}
}