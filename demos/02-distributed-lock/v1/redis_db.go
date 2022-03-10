package main

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

// getIPhoneStock 获取商品库存
func getIPhoneStock(c context.Context) (int, error) {
	// rdb.Get()
	key := iphoneStockKey()

	n, err := rdb.Get(c, key).Int()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return -1, nil
		}

		return 0, err
	}

	return n, nil
}

// iphoneStockKey 返回促销 key
func iphoneStockKey() string {
	return "promote:iphone:stock"
}

// addLuckyGuys 保存促销中奖用户到 set 中, 返回添加人数
func addLuckyGuys(c context.Context, uid string) int64 {
	key := luckyGuysKey()
	return rdb.SAdd(c, key, uid).Val()
}

// isLuckyGuy 判断是否为中奖用户
func isLuckyGuy(c context.Context, uid string) bool {
	return rdb.SIsMember(c, luckyGuysKey(), uid).Val()
}

// luckyGuysKey 返回幸运用户的集合
func luckyGuysKey() string {
	return "promote:iphone:luckyguys"
}

// decrIPHoneStock 减少库存
func decrIPhoneStock(c context.Context) {
	key := iphoneStockKey()
	_ = rdb.Decr(c, key).Val()
}
