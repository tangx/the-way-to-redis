package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "redis123",
			DB:       0,
		},
	)

	c := context.Background()
	if err := rdb.Ping(c).Err(); err != nil {
		panic(err)
	}

	rdb.Set(c, lockerKey(), 0, 0)
	rdb.Set(c, accountKey(), 0, 0)

}
