# Redis Hook for [Logrus](https://github.com/Sirupsen/logrus) <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>

## Why?

Useful for centralized logging, using a RELK stack (Redis, Elasticsearch, Logstash and Kibana). When the hook is installed, all log messages are sent to a Redis server, in Logstash message V0 or V1 format, ready to be parsed/processed by Logstash.

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
	hook, err := logredis.NewHook("localhost", 6379, "my_redis_key", "my_app_name")
	if err == nil {
		log.AddHook(hook)
	}
}

func main() {
	// when hook is injected succesfully, logs will be send to redis server
	log.Info("just some logging...")
}
```

## Message types sent to redis

#### LogstashMessageV0
```
type LogstashMessageV0 struct {
	Type       string `json:"@type,omitempty"`
	Timestamp  string `json:"@timestamp"`
	Sourcehost string `json:"@source_host"`
	Message    string `json:"@message"`
	Fields     struct {
		Application string `json:"application"`
		File        string `json:"file"`
		Level       string `json:"level"`
	} `json:"@fields"`
}
```

#### LogstashMessageV1
```
type LogstashMessageV1 struct {
	Type        string `json:"@type,omitempty"`
	Timestamp   string `json:"@timestamp"`
	Sourcehost  string `json:"host"`
	Message     string `json:"message"`
	Application string `json:"application"`
	File        string `json:"file"`
	Level       string `json:"level"`
}
```

## Testing
Please see the `docker-compose` directory for information about how to test. There is a readme inside.

## In case of hook: disable writing to stdout
See this: https://github.com/Sirupsen/logrus/issues/328#issuecomment-210758435
