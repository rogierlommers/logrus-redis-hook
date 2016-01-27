# Redis Hook for [Logrus](https://github.com/Sirupsen/logrus) <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>

## Why?

Useful for centralized logging, using a RELK stack (Redis, Elasticsearch, Logstash and Kibana).

## Install

```shell
$ go get github.com/rogierlommers/logrus-redis-hook
```

![Colored](http://i.imgur.com/3sWfI4s.jpg)

## Usage

```go
package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/rogierlommers/logrus-redis-hook"
)

func init() {
	hook, err := logredis.NewHook("localhost", 6379, "my_redis_key")
	if err == nil {
		log.AddHook(hook)
	}
}

func main() {
	// when hook is injected succesfully, logs will be send to redis server
	log.Info("just some logging...")
}
```

## License
*MIT*
