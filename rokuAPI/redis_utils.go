package rokuAPI

import (
	"github.com/go-redis/redis"
	"time"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)


type RedisClient struct {
	client *redis.Client
	addr, password string
	db int32
}

func CheckRedisRunning() (running bool) {
	running = true
	cmd := "redis-cli"
	args := []string{"ping"}
	Cmd := exec.Command(cmd, args...)
	Out, err := Cmd.Output()
	if err != nil {
		return false
	}

	out := string(Out)
	if !strings.Contains(out, "PONG") {
		running = false
	}

	return running
}

func StartRedisServer(wg *sync.WaitGroup) {
	defer wg.Done()

	// determine OS
	cmd := "uname"
	args := []string{}
	Cmd := exec.Command(cmd, args...)
	Out, err := Cmd.Output()
	if err != nil { // OS is Mac
		fmt.Println("Unable to run command; Error: ", err)
		panic(err)
	}

	formatted_out := strings.TrimSpace(strings.ToLower(string(Out)))
	if formatted_out == "darwin" {
		cmd = "redis-server"
	} else if formatted_out == "linux"{ // OS is linux
		cmd = "service"
		args = []string{"redis-server", "start"}
	} else {
		fmt.Println("Unknown OS; exiting...")
		panic("Unknown OS")
	}

	// start redis
	Cmd = exec.Command(cmd, args...)
	Out, err = Cmd.Output()
	if err != nil {
		fmt.Println("Unable to run command; Error: ", err)
		panic(err)
	}
}

func Init() *RedisClient {

	var client *redis.Client
	client = redis.NewClient(&redis.Options{
		Addr:         REDIS_ADDR,
		Password:     REDIS_PASSWORD,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
		DB:           REDIS_DB,
	})
	client.FlushDB()

	return &RedisClient{client: client, addr: REDIS_ADDR, db: REDIS_DB}
}

func NewRedisClient() *RedisClient {

	var client *redis.Client
	client = redis.NewClient(&redis.Options{
		Addr:     REDIS_ADDR,
		Password: REDIS_PASSWORD,
		DB:       REDIS_DB,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>

	return &RedisClient{client: client, addr: REDIS_ADDR, db: REDIS_DB}

}

func (redisClient *RedisClient) Quit() {
	err := redisClient.client.Quit().Err()
	if err != nil {
		panic(err)
	}
}

func (redisClient *RedisClient) Set(key string, value string) {
	err := redisClient.client.Set(key, value, 0).Err()
	if err != nil {
		panic(err)
	}

	// fmt.Println(fmt.Sprintf("Set {%s: %s}", key, value))
}

func (redisClient *RedisClient) Get(key string) (value string) {

	value, err := redisClient.client.Get(key).Result()
	if err == redis.Nil {
		fmt.Println(fmt.Sprintf("'%s' does not exists", key))
		value = ""
	} else if err != nil {
		panic(err)
	}else {
		// fmt.Println(fmt.Sprintf("Retrieved {%s: %s}", key, value))
	}
	return value
}
