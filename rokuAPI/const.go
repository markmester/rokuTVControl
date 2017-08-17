package rokuAPI

const (
	// lirc server settings
	LIRC_SERVER_ADDR = "172.24.1.106:5000"

	// roku ssdp settings
	SSDP_REQUEST = "M-SEARCH * HTTP/1.1\r\n" +
		"HOST: 239.255.255.250:1900\r\n" +
		"Man: \"ssdp:discover\"\r\n" +
		"MX: 5\r\n" +
		"ST: roku:ecp\r\n\r\n"

	SSDP_HOST     = "239.255.255.250"
	SSDP_PORT     = 1900
	SSDP_PROTOCOL = "udp"

	// redis connection settings
	REDIS_ADDR     = "localhost:6379"
	REDIS_PASSWORD = "" // no password set
	REDIS_DB       = 0  // using default DB

)