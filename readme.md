# Redis Hook for [Logrus](https://github.com/Sirupsen/logrus) <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>

## Install

```shell
$ go get github.com/rogierlommers/logrus-redis-hook
```

![Colored](http://i.imgur.com/3sWfI4s.jpg)

## Usage

```go
package main

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/rogierlommers/logrus-redis-hook"
)

var log = logrus.New()

func init() {
	hook, err := logredis.NewHook("localhost", 6379, "mykey")
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
```

## License
*MIT*
