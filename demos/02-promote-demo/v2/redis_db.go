package main

import (
	"errors"

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

func luckyGuysPipeWithWatch(c *gin.Context) error {

	key := iphoneStockKey()
	fn := func(tx *redis.Tx) error {
		n, err := tx.Get(c, key).Int()
		if err != nil && errors.Is(err, redis.Nil) {
			return errors.New("活动还没开始")
		}

		if n == 0 {
			return errors.New("活动已结束")
		}
		if n < 0 {
			return errors.New("活动已超卖")
		}

		uid := userID()
		if tx.SIsMember(c, luckyGuysKey(), uid).Val() {
			return errors.New("已经中奖")
		}

		pipe := rdb.TxPipeline()
		// n2, _ := pipe.Get(c, key).Int()
		// fmt.Println("stock n==>", n2)
		pipe.Decr(c, iphoneStockKey())
		pipe.SAdd(c, luckyGuysKey(), uid)
		_, err = pipe.Exec(c)

		return err
	}

	err := rdb.Watch(
		c,
		fn,
		key,
	)
	return err
}
