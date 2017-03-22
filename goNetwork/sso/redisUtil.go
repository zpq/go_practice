package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	redisIP   = "192.168.2.111"
	redisPORT = "6379"
)

type redisPool struct {
	pool *redis.Pool
}

func newRedisPool(password string) (*redisPool, error) {
	server := redisIP + ":" + redisPORT
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}

			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return &redisPool{pool}, nil
}
