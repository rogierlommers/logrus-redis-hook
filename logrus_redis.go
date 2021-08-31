package logredis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
)

// HookConfig stores configuration needed to setup the hook
type HookConfig struct {
	Key      string
	Format   string
	App      string
	Host     string
	Password string
	Hostname string
	Port     int
	DB       int
	TTL      int
}

// RedisHook to sends logs to Redis server
type RedisHook struct {
	RedisPool      *redis.Pool
	RedisHost      string
	RedisKey       string
	LogstashFormat string
	AppName        string
	Hostname       string
	RedisPort      int
	TTL            int
	DialOptions    []redis.DialOption
}

// NewHook creates a hook to be added to an instance of logger
func NewHook(config HookConfig, options ...redis.DialOption) (*RedisHook, error) {
	pool := newRedisConnectionPool(config.Host, config.Password, config.Port, config.DB, options...)

	if config.Format != "v0" && config.Format != "v1" && config.Format != "access" {
		return nil, fmt.Errorf("unknown message format")
	}

	// test if connection with REDIS can be established
	conn := pool.Get()
	defer conn.Close()

	// check connection
	_, err := conn.Do("PING")
	if err != nil {
		return nil, fmt.Errorf("unable to connect to REDIS: %s", err)
	}

	return &RedisHook{
		RedisHost:      config.Host,
		RedisPool:      pool,
		RedisKey:       config.Key,
		LogstashFormat: config.Format,
		AppName:        config.App,
		Hostname:       config.Hostname,
		TTL:            config.TTL,
		DialOptions:    options,
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
	case "access":
		msg = createAccessLogMessage(entry, hook.AppName, hook.Hostname)
	default:
		fmt.Println("Invalid LogstashFormat")
	}

	// Marshal into json message
	js, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error creating message for REDIS: %s", err)
	}

	// get connection from pool
	conn := hook.RedisPool.Get()
	defer conn.Close()

	// send message
	_, err = conn.Do("RPUSH", hook.RedisKey, js)
	if err != nil {
		return fmt.Errorf("error sending message to REDIS: %s", err)
	}

	if hook.TTL != 0 {
		_, err = conn.Do("EXPIRE", hook.RedisKey, hook.TTL)
		if err != nil {
			return fmt.Errorf("error setting TTL to key: %s, %s", hook.RedisKey, err)
		}
	}

	return nil
}

// Levels returns the available logging levels.
func (hook *RedisHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.TraceLevel,
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

func createV0Message(entry *logrus.Entry, appName, hostname string) map[string]interface{} {
	m := make(map[string]interface{})
	m["@timestamp"] = entry.Time.UTC().Format(time.RFC3339Nano)
	m["@source_host"] = hostname
	m["@message"] = entry.Message

	fields := make(map[string]interface{})
	fields["level"] = entry.Level.String()
	fields["application"] = appName

	for k, v := range entry.Data {
		fields[k] = v
	}
	m["@fields"] = fields

	return m
}

func createV1Message(entry *logrus.Entry, appName, hostname string) map[string]interface{} {
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

func createAccessLogMessage(entry *logrus.Entry, appName, hostname string) map[string]interface{} {
	m := make(map[string]interface{})
	m["message"] = entry.Message
	m["@source_host"] = hostname

	fields := make(map[string]interface{})
	fields["application"] = appName

	for k, v := range entry.Data {
		fields[k] = v
	}
	m["@fields"] = fields

	return m
}

func newRedisConnectionPool(server, password string, port int, db int, options ...redis.DialOption) *redis.Pool {
	hostPort := fmt.Sprintf("%s:%d", server, port)
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", hostPort, redis.DialDatabase(db), redis.DialPassword(password), options...)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
