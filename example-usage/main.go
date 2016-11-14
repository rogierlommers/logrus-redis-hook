package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/rogierlommers/logrus-redis-hook"
)

func init() {
	hook, err := logredis.NewHook("localhost",
		"my_redis_key", // key to use
		"v0",           // logstash format (v0, v1 or custom)
		"my_app_name",  // your application name
		"my_hostname",  // your hostname
		"",             // password for redis authentication, leave empty for no authentication
		6379,           // redis port
	)
	if err == nil {
		log.AddHook(hook)
	} else {
		log.Error(err)
	}
}

func main() {
	// when hook is injected succesfully, logs will be send to redis server
	log.Info("just some info logging...")

	// we also support log.WithFields()
	log.WithFields(log.Fields{"animal": "walrus",
		"foo":  "bar",
		"this": "that"}).
		Info("A walrus appears")
}
