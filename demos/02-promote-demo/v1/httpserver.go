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

	// 检查库存
	n, err := getIPhoneStock(c)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"msg":   "内部错误",
			"error": err.Error(),
		})

		return
	}

	if n == -1 {
		c.JSON(http.StatusForbidden, map[string]interface{}{
			"msg": "促销没开始",
		})

		return
	}

	if n == 0 {
		c.JSON(http.StatusOK, map[string]interface{}{
			"msg": "活动已经结束了",
		})
		return
	}

	uid := userID()
	// 判断用户是否中奖
	if isLuckyGuy(c, uid) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"msg": "success",
			"err": "已经中奖",
		})

		return
	}

	// 添加心中奖用户，并减少库存。
	decrIPhoneStock(c)
	last := addLuckyGuys(c, uid)

	c.JSON(http.StatusOK, map[string]interface{}{
		"msg":  "恭喜中奖",
		"uid":  uid,
		"last": last,
	})
}

func userID() string {
	return uuid.New().String()
}
