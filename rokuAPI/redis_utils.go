package rokuAPI

import (
	"github.com/go-redis/redis"
	"time"
	"fmt"
)

const (
	redisAddr = "localhost:6379"
	redisPassword = "" // no password set
	redisDB = 0 // using default DB
)

type RedisClient struct {
	client *redis.Client
	addr, password string
	db int32
}


func Init() *RedisClient {
	var client *redis.Client
	client = redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		Password: 	 redisPassword,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
		DB: 			 redisDB,
	})
	client.FlushDB()

	return &RedisClient{client: client, addr: redisAddr, db: redisDB}
}

func NewRedisClient() *RedisClient {
	var client *redis.Client
	client = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>

	return &RedisClient{client: client, addr: redisAddr, db: redisDB}

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

	fmt.Println(fmt.Sprintf("Set {%s: %s}", key, value))
}

func (redisClient *RedisClient) Get(key string) (value string) {

	value, err := redisClient.client.Get(key).Result()
	if err == redis.Nil {
		fmt.Println(fmt.Sprintf("'%s' does not exists", key))
		value = ""
	} else if err != nil {
		panic(err)
	}else {
		fmt.Println(fmt.Sprintf("Retrieved {%s: %s}", key, value))
	}
	return value
}

//func main() {
//	//var c RedisClient = *Init()
//	c := *NewRedisClient()
//	c.Set("key", "testing")
//	c.Get("key")
//	fmt.Println(c)
//
//}