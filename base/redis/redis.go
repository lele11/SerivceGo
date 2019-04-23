package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	// redis连接池
	pool *redis.Pool
	cfg  RedisConfig
)

type RedisConfig struct {
	Addr        string
	Index       uint32
	MaxActive   int
	MaxIdle     int
	IdleTimeout uint32
	Password    string
}

func SetConfig(addr string, password string) {
	cfg.Addr = addr
	cfg.Password = password
	cfg.Index = 0
	cfg.IdleTimeout = 10
	cfg.MaxActive = 100
	cfg.MaxIdle = 10
}

// GetRedisConn 获取一个redis连接
func GetRedisConn() redis.Conn {
	if pool == nil {
		initPool()
	}
	c := pool.Get()
	if c.Err() != nil {
		fmt.Println(c.Err())
		c.Close()
	}
	return c
}

// initPool 初始化, 创建redis连接池
func initPool() {
	if pool == nil {
		rawURL := fmt.Sprintf("redis://%s/%d", cfg.Addr, cfg.Index)
		pool = &redis.Pool{
			MaxIdle:     cfg.MaxIdle,
			MaxActive:   cfg.MaxActive,
			Wait:        false,
			IdleTimeout: time.Duration(cfg.IdleTimeout) * time.Second,
			Dial: func() (redis.Conn, error) {
				if cfg.Password != "" {
					return redis.DialURL(rawURL, redis.DialPassword(cfg.Password))
				}
				return redis.DialURL(rawURL)
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
		}
	}
}
