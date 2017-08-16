package main

import (
	"github.com/markmester/rokuTVControl/rokuAPI"
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	// ---------------------------- Locate Loop
	wg.Add(1)
	go rokuAPI.LocateLoop(&wg)

	// ---------------------------- Webserver
	wg.Add(2)
	go rokuAPI.MuxServer(&wg)


	fmt.Println("Main: Waiting for GoRoutines to finish")
	wg.Wait()

	//const (
	//	request = "POST /keypress/Power HTTP/1.1\r\n" +
	//		"HOST: 172.24.1.99:8060\r\n\r\n"
	//	myhost = "172.24.1.99:8060"
	//)
}
