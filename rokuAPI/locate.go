package rokuAPI

import (
	"net"
	"fmt"
	"time"
	"regexp"
	"sync"
	"strings"
)

func CheckError(err error) {
	// CheckError: a generic function for handling socket connection errors
	if err  != nil {
		fmt.Println("Error: " , err)
	}
}


func Locate() (ip string){
	// Locate: attempts to discover the IP address of a Roku device in the area. Uses the standard SSDP multicast address
	// and SSDP_PORT (239.255.255.250:1900) used for local area network communication.
	// Returns the Roku device IP address if found.
	inBuf := make([]byte, 1024)
	timeoutDuration := 10 * time.Second
	addr := net.UDPAddr{
		Port: SSDP_PORT,
		IP:   net.ParseIP(SSDP_HOST),
	}
	var readLen int
	var fromAddr *net.UDPAddr

	//Connect udp
	conn, err := net.ListenUDP(SSDP_PROTOCOL, &addr)
	CheckError(err)
	conn.SetDeadline(time.Now().Add(timeoutDuration))
	defer conn.Close()


	//Write discovery
	sendLen, err := conn.WriteToUDP([]byte(SSDP_REQUEST), &addr)
	CheckError(err)
	fmt.Sprintf("Sent %d bytes from %s:%d", sendLen, SSDP_HOST, SSDP_PORT)

	//Read response
	readLen, fromAddr, err = conn.ReadFromUDP(inBuf)
	fmt.Sprintln("Read", readLen, "bytes from", fromAddr)

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

func WakeRoku() {
	fmt.Println(">>> Waking Roku...")
	_ = LircServerRequest(LIRC_SERVER_ADDR, "/power", "GET")
	time.Sleep(10 * time.Second)
	_ = LircServerRequest(LIRC_SERVER_ADDR, "/power", "GET")
}

func LocateLoop(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(">>> Starting Roku Device Location Loop...")

	attempts := 0
	redis_ctx := *NewRedisClient()
	for {
		attempts += 1
		roku_addr := redis_ctx.Get("roku_address")

		if attempts % 3 == 0 && roku_addr == ""{
			WakeRoku()
			attempts = 0
		}

		ip := Locate()
		if ip != "" {
			redis_ctx.Set("roku_address", ip)
			attempts = 0
		}
		time.Sleep(5 * time.Second)
	}
}