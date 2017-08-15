package main

import (
	"net"
	"fmt"
	"io/ioutil"
)

const (
	Request = "GET /query/apps HTTP/1.1\r\n" +
		"HOST: 172.24.1.99:8060\r\n\r\n"
	myhost = "172.24.1.99:8060"
)

func main() {
	conn, err := net.Dial("tcp", myhost)
	if err != nil {
		// handle error
	}


	//Write discovery
	sendLen, err := conn.Write([]byte(Request))
	fmt.Println(fmt.Sprintf("Sent %d bytes from %s", sendLen, myhost))

	//Read response
	raw_data, _ := ioutil.ReadAll(conn)
	fmt.Println("Read", len(raw_data), "bytes from", myhost)
	defer conn.Close()

	//Decode response
	data := string(raw_data[:])

	fmt.Println(data)
}



