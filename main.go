package main

import (
	"fmt"
	"time"
)

const (
	modelTime      = 10000
	entrancingTime = 130
	servingTime    = 150
	flowsCount     = 3
	channels       = 2
)

func main() {

	fmt.Println("======================================")

	modelingTime := time.Duration(modelTime) * time.Millisecond
	explTIme := time.Duration(entrancingTime) * time.Millisecond

	nower := time.Now()
	serv := NewService(servingTime, 0, channels, nower)

	breaked := false

	after := time.After(modelingTime)
	//ticker := time.NewTicker(explTIme)
	records := 0

	pusher := make(chan int)

	go func() {
		for range time.Tick(explTIme) {
			for range make([]int, flowsCount) {
				pusher <- 0
			}
		}
	}()

	go func() {
		for {
			fmt.Printf("/ loading\r")
			<-time.After(150 * time.Millisecond)
			fmt.Printf("- loading\r")
			<-time.After(150 * time.Millisecond)
			fmt.Printf("\\ loading\r")
			<-time.After(150 * time.Millisecond)
			fmt.Printf("| loading\r")
			<-time.After(150 * time.Millisecond)
		}
	}()

	for {
		if breaked {
			break
		}
		select {
		case <-after:
			fmt.Println("Process finished:")
			//ticker.Stop()
			breaked = true
		case <-pusher:
			records++
			serv.Process()
		}
	}

	serv.Processed(records, float64(channels*entrancingTime)/float64(flowsCount*servingTime))
	fmt.Println("======================================")
}
