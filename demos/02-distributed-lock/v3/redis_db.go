package main

import (
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// iphoneStockKey 返回促销 key
func iphoneStockKey() string {
	return "promote:iphone:stock"
}

// luckyGuysKey 返回幸运用户的集合
func luckyGuysKey() string {
	return "promote:iphone:luckyguys"
}

func luckyGuysPipeWithLua(c *gin.Context) error {

	rdb.SAdd(c, luckyGuysKey(), "123123")

	b, _ := os.ReadFile("promote.lua")
	lua := redis.NewScript(string(b))

	uid := userID()
	n, err := lua.Run(c, rdb, []string{iphoneStockKey(), luckyGuysKey()}, uid).Int()

	if err != nil {
		return err
	}

	switch n {
	case 0:
		return errors.New("活动已经结束")
	case -1:
		return errors.New("活动还没开始")
	case -2:
		return errors.New("已经是中奖用户")
	default:
		return nil
	}

}
