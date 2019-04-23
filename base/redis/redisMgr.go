package redis

import (
	"fmt"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

var RedisPool sync.Map

type RedisInstance struct {
	pool *redis.Pool
	cfg  *RedisConfig
}

func NewRedisInstance(name string, cfg *RedisConfig) {
	in := &RedisInstance{
		cfg: cfg,
	}
	RedisPool.Store(name, in)
}

func GetInstance(name string) *RedisInstance {
	i, ok := RedisPool.Load(name)
	if !ok {
		return nil
	}
	return i.(*RedisInstance)
}

// GetRedisConn 获取一个redis连接
func (r *RedisInstance) GetConn() redis.Conn {
	if r.pool == nil {
		r.initPool()
	}
	c := r.pool.Get()
	if c.Err() != nil {
		c.Close()
	}
	return c
}

// initPool 初始化, 创建redis连接池
func (r *RedisInstance) initPool() {
	if r.pool == nil {
		rawURL := fmt.Sprintf("redis://%s/%d", cfg.Addr, cfg.Index)
		r.pool = &redis.Pool{
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
