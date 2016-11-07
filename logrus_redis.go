package logredis

import (
	"encoding/json"
	"fmt"
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
	Hostname       string
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
func NewHook(redisHost string, port int, key string, format string, appname string, hostname string) (*RedisHook, error) {
	pool := newRedisConnectionPool(redisHost, port)

	// test if connection with REDIS can be established
	conn := pool.Get()
	defer conn.Close()

	// check connection
	_, err := conn.Do("PING")
	if err != nil {
		return nil, fmt.Errorf("unable to connect to REDIS: %s", err)
	}

	// by default, use V0 format
	if strings.ToLower(format) != "v0" && strings.ToLower(format) != "v1" && strings.ToLower(format) != "custom" {
		format = "v0"
	}

	return &RedisHook{
		RedisHost:      redisHost,
		RedisPool:      pool,
		RedisKey:       key,
		LogstashFormat: format,
		AppName:        appname,
		Hostname:       hostname,
	}, nil
}

// Fire is called when a log event is fired.
func (hook *RedisHook) Fire(entry *logrus.Entry) error {
	var msg interface{}

	switch hook.LogstashFormat {
	case "v0":
		msg = createV0Message(entry, hook.AppName, hook.Hostname)
	case "v1":
		msg = createV1Message(entry, hook.AppName, hook.Hostname)
	case "custom":
		msg = createCustomMessage(entry, hook.AppName, hook.Hostname)
	default:
		fmt.Println("Invalid LogstashFormat")
	}

	js, err := json.Marshal(msg)
	fmt.Println(string(js))
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

func createV0Message(entry *logrus.Entry, appName, hostname string) LogstashMessageV0 {
	m := LogstashMessageV0{}
	m.Timestamp = entry.Time.UTC().Format(time.RFC3339Nano)
	m.Sourcehost = hostname
	m.Message = entry.Message
	m.Fields.Level = entry.Level.String()
	m.Fields.Application = appName
	return m
}

func createV1Message(entry *logrus.Entry, appName, hostname string) LogstashMessageV1 {
	m := LogstashMessageV1{}
	m.Timestamp = entry.Time.UTC().Format(time.RFC3339Nano)
	m.Sourcehost = hostname
	m.Message = entry.Message
	m.Level = entry.Level.String()
	m.Application = appName
	return m
}

func createCustomMessage(entry *logrus.Entry, appName, hostname string) map[string]interface{} {
	m := make(map[string]interface{})
	m["@timestamp"] = entry.Time.UTC().Format(time.RFC3339Nano)
	m["host"] = hostname
	m["message"] = entry.Message
	m["level"] = entry.Level.String()
	m["application"] = appName
	for k, v := range entry.Data {
		m[k] = v
	}
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

			// In case redis needs authentication
			// https://github.com/rogierlommers/logrus-redis-hook/issues/2
			//
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
