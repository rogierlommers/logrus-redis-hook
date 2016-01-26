package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/rogierlommers/logrus-redis-hook"
)

var log = logrus.New()

func init() {
	redis_host := os.Getenv("REDIS_HOST")
	if redis_host == "" {
		redis_host = "localhost"
	}
	redis_key := os.Getenv("REDIS_KEY")
	if redis_key == "" {
		redis_key = "mykey"
	}

	fmt.Printf("Connecting to redis://%s. Pushing to key '%s'\n", redis_host, redis_key)

	hook, err := logredis.NewHook(redis_host, 6379, redis_key)
	if err == nil {
		log.Hooks.Add(hook)
	}
}

func main() {
	// send 1000 records to redis
	for i := 0; i < 1000; i++ {
		log.Infof("logrule, number: %d", i)
		time.Sleep(1 * time.Second)
	}
}
