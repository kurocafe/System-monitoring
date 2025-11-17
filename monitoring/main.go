package main

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/cpu"
)

func main(){
	fmt.Println("System monitoring")
	fmt.Println("-----------------")

	// CPU percentage (per second = false)
  cpuPercent, _ := cpu.Percent(0, false)
  fmt.Println("CPU:", cpuPercent[0])
}