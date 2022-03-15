package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
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

	time.Sleep(time.Microsecond * time.Duration(rand.Intn(100)))

	// 1. 获取当前 locker 版本
	ver := rdb.Get(ctx, lockerKey()).Val()

	// 2. 获取账户余额
	_, err := rdb.Get(ctx, accountKey()).Int()
	if err != nil {
		return err
	}

	// 4. 事务执行
	fn := func(tx *redis.Tx) error {
		pipe := tx.Pipeline()
		// 4.1 获取 locker 版本， locker 版本于本地版本一致判断
		ver2 := tx.Get(ctx, lockerKey()).Val()
		if ver2 != ver {
			return errors.New("version 变更， 退出")
		}

		// 4.3 将账户余额存进 redis
		pipe.IncrBy(ctx, accountKey(), 100)
		pipe.Incr(ctx, lockerKey())
		_, err := pipe.Exec(ctx)
		return err
	}

	// 3. watch locker 版本
	err = rdb.Watch(
		ctx,
		fn,
		lockerKey(),
	)

	return err
}
