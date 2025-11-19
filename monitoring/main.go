package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type Metrics struct{
	CPU float64
	Memory float64
	Disk float64
}

func cpuCollector(ch chan float64){
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		v, _ := cpu.Percent(0, false)
		ch <- v[0]
	}
}

func memoryCollector(ch chan float64){
	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		m, _ := mem.VirtualMemoryWithContext(
			context.Background(),
		)
		ch <- m.UsedPercent 
	}
}

func DiskCollector(ch chan float64){
	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		d, _ := disk.Usage("/")
		ch <- d.UsedPercent
	}
}

func main(){

	procPath := os.Getenv("GOMONITOR_PROC")
	rootPath := os.Getenv("GOMONITOR_ROOT")

	if procPath == "" {
		procPath = "/proc"
	}
	if rootPath == "" {
		rootPath = "/"
	}

	fmt.Println("System monitoring")
	fmt.Println("-----------------")

	cpuCh := make(chan float64)
	memCh := make(chan float64)
	diskCh := make(chan float64)

	go cpuCollector(cpuCh)
	go memoryCollector(memCh)
	go DiskCollector(diskCh)

	metrics := Metrics{}

	for {
		select{
		case cpuVal := <- cpuCh:
			metrics.CPU = cpuVal
		case memVal := <- memCh:
			metrics.Memory = memVal
		case diskVal := <- diskCh:
			metrics.Disk = diskVal
		}

		fmt.Print("\033[H\033[2J") // clear terminal
		fmt.Println("System Monitor (updates every 1s)")
		fmt.Println("--------------------------------")
		fmt.Printf("CPU:    %.1f %%\n", metrics.CPU)
		fmt.Printf("Memory: %.1f %%\n", metrics.Memory)
		fmt.Printf("Disk:   %.1f %%\n", metrics.Disk)
	}


}
