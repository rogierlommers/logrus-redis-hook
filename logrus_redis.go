package logredis

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
)

// RedisHook to sends logs to Redis server
type RedisHook struct {
	RedisPool      *redis.Pool
	RedisHost      string
	RedisKey       string
	LogstashFormat string
	AppName        string
	RedisPort      int
}

// LogstashMessageV0 represents v0 format
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

// LogstashMessageV1 represents v1 format
type LogstashMessageV1 struct {
	Type        string `json:"@type,omitempty"`
	Timestamp   string `json:"@timestamp"`
	Sourcehost  string `json:"host"`
	Message     string `json:"message"`
	Application string `json:"application"`
	File        string `json:"file"`
	Level       string `json:"level"`
}

// NewHook creates a hook to be added to an instance of logger
func NewHook(host string, port int, key string, format string, appname string) (*RedisHook, error) {
	pool := newRedisConnectionPool(host, port)

	// test if connection with REDIS can be established
	conn := pool.Get()
	defer conn.Close()

	// check connection
	_, err := conn.Do("PING")
	if err != nil {
		return nil, fmt.Errorf("unable to connect to REDIS: %s", err)
	}

	// by default, use V0 format
	if strings.ToLower(format) != "v0" && strings.ToLower(format) != "v1" {
		format = "v0"
	}

	return &RedisHook{
		RedisHost:      host,
		RedisPool:      pool,
		RedisKey:       key,
		LogstashFormat: format,
		AppName:        appname,
	}, nil
}

// Fire is called when a log event is fired.
func (hook *RedisHook) Fire(entry *logrus.Entry) error {
	var msg interface{}

	switch hook.LogstashFormat {
	case "v0":
		msg = createV0Message(entry, hook.AppName)
	case "v1":
		msg = createV1Message(entry, hook.AppName)
	}

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

func createV0Message(entry *logrus.Entry, appName string) LogstashMessageV0 {
	m := LogstashMessageV0{}
	m.Timestamp = entry.Time.UTC().Format(time.RFC3339Nano)
	m.Sourcehost = reportHostname()
	m.Message = entry.Message
	m.Fields.Level = entry.Level.String()
	m.Fields.Application = appName
	return m
}

func createV1Message(entry *logrus.Entry, appName string) LogstashMessageV1 {
	m := LogstashMessageV1{}
	m.Timestamp = entry.Time.UTC().Format(time.RFC3339Nano)
	m.Sourcehost = reportHostname()
	m.Message = entry.Message
	m.Level = entry.Level.String()
	m.Application = appName
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
