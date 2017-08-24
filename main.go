package main

import (
	//"github.com/markmester/rokuTVControl/rokuAPI"
	"github.com/markmester/rokuTVControl/alexa"
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	//// ---------------------------- Locate Loop
	//wg.Add(1)
	//go rokuAPI.LocateLoop(&wg)
	//
	//// ---------------------------- Webserver
	//wg.Add(2)
	//go rokuAPI.MuxServer(&wg)

	// ----------------------------- SQS polling
	wg.Add(3)
	go alexa.PollQueue(&wg, "roku-control",  5)

	fmt.Println("Main: Waiting for GoRoutines to finish")
	wg.Wait()
}
