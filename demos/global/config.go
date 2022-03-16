package global

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var RdbClient *redis.Client

func init() {
	RdbClient = redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "redis123",
			DB:       0,
		},
	)

	c := context.Background()
	if err := RdbClient.Ping(c).Err(); err != nil {
		panic(err)
	}

}
