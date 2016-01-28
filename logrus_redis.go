package logredis

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
)

// RedisHook to sends logs to Redis server
type RedisHook struct {
	RedisPool *redis.Pool
	RedisHost string
	RedisKey  string
	RedisPort int
}

// LogstashMessageV0 represents v0 format
type LogstashMessageV0 struct {
	Type       string `json:"@type,omitempty"`
	Timestamp  string `json:"@timestamp"`
	Sourcehost string `json:"@source_host"`
	Message    string `json:"@message"`
	Level      string `json:"@level"`
	Fields     struct {
		Exception struct {
			ExceptionClass   string `json:"exception_class"`
			ExceptionMessage string `json:"exception_message"`
			Stacktrace       string `json:"stacktrace"`
		} `json:"exception"`
		File      string `json:"file"`
		Level     string `json:"level"`
		Timestamp string `json:"timestamp"`
	} `json:"@fields"`
}

// NewHook creates a hook to be added to an instance of logger
func NewHook(host string, port int, key string) (*RedisHook, error) {
	pool := newRedisConnectionPool(host, port)

	// test if connection with REDIS can be established
	conn := pool.Get()
	defer conn.Close()

	// check connection
	_, err := conn.Do("PING")
	if err != nil {
		return nil, fmt.Errorf("unable to connect to REDIS: %s", err)
	}

	return &RedisHook{
		RedisHost: host,
		RedisPool: pool,
		RedisKey:  key,
	}, nil
}

// Fire is called when a log event is fired.
func (hook *RedisHook) Fire(entry *logrus.Entry) error {
	msg := createMessage(entry)

	js, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error creating message for REDIS: %s", err)
	}

	conn := hook.RedisPool.Get()
	defer conn.Close()

	_, err = conn.Do("RPUSH", hook.RedisKey, js)
	if err != nil {
		return fmt.Errorf("error sending message to REDIS: %s", err)
	}
	return nil
}

// Levels returns the available logging levels.
func (hook *RedisHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

func createMessage(entry *logrus.Entry) LogstashMessageV0 {
	m := LogstashMessageV0{}
	m.Timestamp = entry.Time.UTC().Format(time.RFC3339Nano)
	m.Sourcehost = reportHostname()
	m.Message = entry.Message
	m.Fields.Level = entry.Level.String()
	return m
}

func newRedisConnectionPool(server string, port int) *redis.Pool {
	hostPort := fmt.Sprintf("%s:%d", server, port)
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", hostPort)
			if err != nil {
				return nil, err
			}

			// if password != "" {
			// 	if _, err := c.Do("AUTH", password); err != nil {
			// 		c.Close()
			// 		return nil, err
			// 	}
			// }

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func reportHostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return h
}
