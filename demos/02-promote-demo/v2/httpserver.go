package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

	err := luckyGuysPipeWithWatch(c)
	if err != nil {
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
