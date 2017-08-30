package main

import (
	"github.com/markmester/rokuTVControl/rokuAPI"
	"github.com/markmester/rokuTVControl/alexa"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	// ------------------------------ Redis
	wg.Add(1)
	redis_proc := rokuAPI.CheckRedisRunning()
	if !redis_proc {
		go rokuAPI.StartRedisServer(&wg)
		time.Sleep(5 * time.Second) // give time for redis to start
	}

	// ---------------------------- Locate Loop
	wg.Add(2)
	go rokuAPI.LocateLoop(&wg)

	// ----------------------------- SQS polling
	wg.Add(3)
	go alexa.PollQueue(&wg, "roku-control",  5)

	fmt.Println("Main: Waiting for GoRoutines to finish")
	wg.Wait()
}
