package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	ctx := context.Background()
	for i := 0; i < 99; i++ {
		i := i
		go func() {
			err := worker(ctx)
			if err != nil {
				fmt.Printf("go( %d ) err: %v \n", i, err.Error())
				return
			}
			fmt.Printf("go( %d ): success\n", i)
		}()

	}
	time.Sleep(1 * time.Second)

	result(ctx)

}

func result(ctx context.Context) {
	fmt.Println("locker version =>", rdb.Get(ctx, lockerKey()).Val())
	fmt.Println("account money  =>", rdb.Get(ctx, accountKey()).Val())
}

func lockerKey() string {
	return "occ:locker:version"
}
func accountKey() string {
	return "occ:account:money"
}

func worker(ctx context.Context) error {

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))

	/* 数据获取
	1. 获取当前 locker 版本
	2. 获取账户余额
	// ver := rdb.Get(ctx, lockerKey()).Val()
	// _, err := rdb.Get(ctx, accountKey()).Int()
	// if err != nil {
	// 	return err
	// }
	*/
	vals := rdb.MGet(ctx, lockerKey(), accountKey()).Val()
	ver := vals[0]
	account, err := atoi(vals[1])
	if err != nil {
		return err
	}
	newAccount := account + 100
	fmt.Println("newAccount:=>", newAccount)

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(200)))

	// 4. 事务执行
	fn := func(tx *redis.Tx) error {

		// 4.1 获取 locker 版本， locker 版本于本地版本一致判断
		ver2 := tx.Get(ctx, lockerKey()).Val()
		if ver2 != ver {
			return errors.New("version 变更， 退出")
		}

		pipe := tx.Pipeline()
		// 4.3 将账户余额存进 redis
		// pipe.IncrBy(ctx, accountKey(), 100)
		pipe.Set(ctx, accountKey(), newAccount, 0)
		pipe.Incr(ctx, lockerKey())
		_, err := pipe.Exec(ctx)
		return err
	}

	// 3. watch locker 版本
	err = rdb.Watch(
		ctx,
		fn,
		lockerKey(),
		accountKey(),
	)

	return err
}

// atoi 将 redis 结果转化为 int 类型
func atoi(v interface{}) (int, error) {
	switch vv := v.(type) {
	case string:
		return strconv.Atoi(vv)
	case int, int8, int16, int32, int64:
		return vv.(int), nil
	default:
		return 0, errors.New("uknown type")
	}
}
