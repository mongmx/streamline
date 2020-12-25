package redis

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Config redis.
type Config struct {
	Host string
	Port string
}

// LoadEnv - load configuration from env.
func LoadEnv() Config {
	return Config{
		Host: os.Getenv("REDIS_HOST"),
		Port: os.Getenv("REDIS_PORT"),
	}
}

// NewRedis creates new connection to redis and return the connection
func NewRedis(cfg Config) (*redis.Pool, error) {
	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	pool := &redis.Pool{
		MaxIdle:     3,                 // adjust to your needs
		IdleTimeout: 240 * time.Second, // adjust to your needs
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			pong, err := conn.Do("PING"); 
			if err != nil {
				conn.Close()
				return nil, err
			}
			_, err = redis.String(pong, err)
			if err != nil {
				conn.Close()
				return nil, err
			}
			return conn, nil
		},
	}
	return pool, nil
}

func db() int {
	d, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return 1
	}
	return d
}
