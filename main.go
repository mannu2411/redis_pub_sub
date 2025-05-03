package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	RedisPub redis.Conn
	psc      redis.PubSubConn
)

func NewRedisProvider() {
	var err error
	RedisPub, err = redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Printf("Unable to connect to redis client")
	}
	RedisSub, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Printf("Unable to connect to redis client")
	}
	psc = redis.PubSubConn{Conn: RedisSub}
	err = psc.Subscribe("order")
	if err != nil {
		fmt.Printf("Unable to subscribe to order channel")
	}
}
func Publish(key string, val interface{}) error {
	_, err := RedisPub.Do("PUBLISH", key, val)
	if err != nil {
		fmt.Printf("Unable to publish to order channel")
		return err
	}
	return nil
}
func Consume() error {
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("Recieved: %v", string(v.Data))
		case error:
			return v
		}
	}
}

func main() {
	NewRedisProvider()
	time.Sleep(time.Second * 5)
	go Publish("order", "Hello")
	go Consume()
	time.Sleep(time.Minute * 2)
}
