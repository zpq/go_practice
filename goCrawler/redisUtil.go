package main

import (
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
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

func (r *redisPool) getOneURL(tag string) string {
	url, err := r.pool.Get().Do("lpop", tag)
	if err != nil || url == nil {
		// log.Println(err.Error())
		return ""
	}
	return string(url.([]byte))
}

func (r *redisPool) addOneURL(tag, url string) {
	_, err := r.pool.Get().Do("rpush", tag, url)
	if err != nil {
		log.Println(err.Error())
	}
}
