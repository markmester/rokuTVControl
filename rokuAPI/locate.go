package rokuAPI

import (
	"net"
	"fmt"
	"time"
	"regexp"
	"sync"
	"strings"
)

const (
	ssdpRequest = "M-SEARCH * HTTP/1.1\r\n" +
	"HOST: 239.255.255.250:1900\r\n" +
	"Man: \"ssdp:discover\"\r\n" +
	"MX: 5\r\n" +
	"ST: roku:ecp\r\n\r\n"

	host = "239.255.255.250"
	port = 1900
	protocol = "udp"
)


func CheckError(err error) {
	// CheckError: a generic function for handling socket connection errors
	if err  != nil {
		fmt.Println("Error: " , err)
	}
}


func Locate() (ip string){
	// Locate: attempts to discover the IP address of a Roku device in the area. Uses the standard SSDP multicast address
	// and port (239.255.255.250:1900) used for local area network communication.
	// Returns the Roku device IP address if found.
	inBuf := make([]byte, 1024)
	timeoutDuration := 10 * time.Second
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(host),
	}
	var readLen int
	var fromAddr *net.UDPAddr

	//Connect udp
	conn, err := net.ListenUDP(protocol, &addr)
	CheckError(err)
	conn.SetDeadline(time.Now().Add(timeoutDuration))
	defer conn.Close()


	//Write discovery
	sendLen, err := conn.WriteToUDP([]byte(ssdpRequest), &addr)
	CheckError(err)
	fmt.Println(fmt.Sprintf("Sent %d bytes from %s:%d", sendLen, host, port))

	//Read response
	readLen, fromAddr, err = conn.ReadFromUDP(inBuf)
	fmt.Println("Read", readLen, "bytes from", fromAddr)

	data := string(inBuf[:readLen])

	if data != "" {
		ip = ParseIP(data)
	}

	return ip

}

func ParseIP(data string) (ip string) {
	//tbd for parsing out IP fom resp
	r := regexp.MustCompile(`LOCATION:(.*)`)
	match := r.FindStringSubmatch(data)
	if len(match) > 1 {
		r := strings.NewReplacer("http://", "",)
		ip = r.Replace(strings.TrimSuffix(strings.TrimSpace(match[1]),"/"))
	}

	return ip
}

func LocateLoop(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(">>> Starting Roku Device Location Loop...")

	redis_ctx := *NewRedisClient()
	for {
		//fmt.Println(">>> Attempting to locate Roku device...")
		ip := Locate()

		if ip != "" {
			redis_ctx.Set("roku_address", ip)
		}
		time.Sleep(5 * time.Second)
	}
}