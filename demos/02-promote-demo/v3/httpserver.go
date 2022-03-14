package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

func serve() {
	httpserver.POST("promote/iphone/:no", handlerPromoteSET)
	httpserver.GET("/promote/iphone", handlerPromote)
	_ = httpserver.Run(":8081")
}

func handlerPromoteSET(c *gin.Context) {
	no := c.Param("no")
	rdb.Set(c, iphoneStockKey(), no, -1)
	rdb.Del(c, luckyGuysKey())
}

func handlerPromote(c *gin.Context) {

	err := luckyGuysPipeWithLua(c)
	if err != nil && !errors.Is(err, redis.Nil) {
		c.JSON(http.StatusForbidden, map[string]interface{}{
			"msg": "forbiden",
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"msg": "恭喜中奖",
	})
	return

}

func userID() string {
	return uuid.New().String()
}
