package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {

}

func serve() {

	httpserver.POST("/promote/iphone")
	_ = httpserver.Run(":8081")
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

	uid := userID()

	if isLuckyGuy(c, uid) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"msg": "已经中奖",
		})

		return
	}

}

func userID() string {
	return uuid.New().String()
}
