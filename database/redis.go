package database

import "github.com/go-redis/redis"

var client *redis.Client

func InitRedis() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}