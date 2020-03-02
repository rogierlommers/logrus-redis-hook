# Redis Hook for [Logrus](https://github.com/Sirupsen/logrus) <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>

[![Go Report Card](https://goreportcard.com/badge/github.com/rogierlommers/logrus-redis-hook)](https://goreportcard.com/report/github.com/rogierlommers/logrus-redis-hook)

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
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/rogierlommers/logrus-redis-hook"
)

func init() {
	hookConfig := logredis.HookConfig{
		Host:     "localhost",
		Key:      "my_redis_key",
		Format:   "v0",
		App:      "my_app_name",
		Port:     6379,
		Hostname: "my_app_hostname", // will be sent to field @source_host
		DB:       0, // optional
		TTL:      3600,
	}

	hook, err := logredis.NewHook(hookConfig)
	if err == nil {
		logrus.AddHook(hook)
	} else {
		logrus.Errorf("logredis error: %q", err)
	}
}

func main() {
	// when hook is injected succesfully, logs will be sent to redis server
	logrus.Info("just some info logging...")

	// we also support log.WithFields()
	logrus.WithFields(logrus.Fields{
		"animal": "walrus",
		"foo":    "bar",
		"this":   "that"}).
		Info("additional fields are being logged as well")

	// If you want to disable writing to stdout, use setOutput
	logrus.SetOutput(ioutil.Discard)
	logrus.Info("This will only be sent to Redis")
}

```


## Testing
Please see the `docker-compose` directory for information about how to test. There is a readme inside.

## In case of hook: disable writing to stdout
See this: https://github.com/Sirupsen/logrus/issues/328#issuecomment-210758435
