package rokuAPI

import (
	"net"
	"fmt"
	"io/ioutil"
	"strings"
	"bufio"
)


func RokuRequest(host string, route string, method string) (data string) {
	request := fmt.Sprintf("%s %s HTTP/1.1\r\nHOST: %s\r\n\r\n", strings.ToUpper(method), route, host)

	conn, err := net.Dial("tcp", host)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	//Write request
	sendLen, err := conn.Write([]byte(request))
	fmt.Sprintf("Sent %d bytes from %s", sendLen, host)

	//Read response
	raw_data, _ := ioutil.ReadAll(conn)
	fmt.Sprintln("Read", len(raw_data), "bytes from", host)

	//Decode response
	data = string(raw_data[:])

	return data
}

func LircServerRequest(host string, route string, method string) (resp string) {
	request := fmt.Sprintf("%s %s HTTP/1.1\r\nHOST: %s\r\n\r\n", strings.ToUpper(method),route, host)
	conn, err := net.Dial("tcp", host)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	//Write request
	sendLen, err := conn.Write([]byte(request))
	fmt.Println(fmt.Sprintf("Sent %d bytes", sendLen))

	//Read response
	raw_data, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Read", len(raw_data), "bytes from", host)

	//Decode response
	resp = string(raw_data[:])

	return resp
}


