package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var httpserver *gin.Engine
var rdb *redis.Client

func init() {
	httpserver = gin.Default()

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "redis123",
		DB:       1,
		PoolSize: 500,
	})
}
