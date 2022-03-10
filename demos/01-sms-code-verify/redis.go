package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

func init() {

	// fmt.Println("init")
	rdb = redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "redis123",
			DB:       0,
		},
	)

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}

// getUserSmsCode 获取短信验证码
func getUserSmsCode(phone string) string {
	key := userSmsCodeKey(phone)

	code := rdb.Get(ctx, key).Val()

	// fmt.Println("code = ", code)
	return code
}

// delUserSmsCode 删除用户验证码
func delUserSmsCode(phone string) error {
	key := userSmsCodeKey(phone)

	err := rdb.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		return nil
	}

	return err
}

// setUserSmsCode 保存用户的短信验证码
func setUserSmsCode(phone string, code string) error {

	// 需求1:
	key := userSmsCodeKey(phone)

	// 1.2. 验证码保存到 redis 中，`userSmsCodeValue:13912341234` ， 过期时间 2 分钟。
	err := rdb.SetEX(ctx, key, code, 2*time.Minute).Err()
	if err != nil {
		return err
	}

	// 3.3. 在 redis 中验证码发送次数 `userSmsCodeTimes:13912341234` 累加 1
	if err := incrUserSmsCodeTimes(phone); err != nil {
		return err
	}
	return nil
}

// incrUserSmsCodeTimes 增加用户发送验证码次数
func incrUserSmsCodeTimes(phone string) error {
	key := userSmsCodeTimes(phone)

	/* todo: 优化点 事务 */
	_, err := rdb.Get(ctx, key).Int()

	// 3.3. 判断是否存在验证码发送计数 key (`userSmsCodeTimes:13912341234`)
	// 如果不存在， 则创建并设置过期时间。
	if err != nil && errors.Is(redis.Nil, err) {
		ttl := ttlSeconds()
		rdb.SetEX(ctx, key, 0, time.Duration(ttl)*time.Second)
	}

	// 在 redis 中验证码发送计数 累加 1。
	err = rdb.Incr(ctx, key).Err()
	return err
}

// getUserSmsCodeTimes 获取用户发送验证码次数
func getUserSmsCodeTimes(phone string) int {
	key := userSmsCodeTimes(phone)

	times, err := rdb.Get(ctx, key).Int()
	if err != nil {

		/* 这种写法很丑， 依赖了 error 的字面值 */
		// if err.Error() == "redis: nil" {
		// 	return 0
		// }

		if errors.Is(err, redis.Nil) {
			return 0
		}

		panic(err)
	}

	return times
}

// userSmsCodeKey 用户当前验证码
func userSmsCodeKey(phone string) (key string) {
	return fmt.Sprintf("userSmsCodeValue:%s", phone)
}

// userSmsCodeTimes 用户发送验证码次数 key
func userSmsCodeTimes(phone string) (key string) {
	return fmt.Sprintf("userSmsCodeTimes:%s", phone)
}

// ttlSeconds 返回到 23:59:59 的剩余秒数
func ttlSeconds() int {

	now := time.Now().Local()

	// todayLast := fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day()) + " 23:59:59"
	todayLast := fmt.Sprintf("%s %s", now.Format("2006-01-02"), "23:59:59")
	today, _ := time.ParseInLocation("2006-01-02 15:04:05", todayLast, now.Location())

	sec := today.Unix() - now.Unix()

	return int(sec)
}
