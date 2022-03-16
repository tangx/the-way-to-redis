package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func main() {

	for i := 0; i < 99; i++ {

		ctx := context.Background()

		go func() {
			c, cancel := context.WithTimeout(ctx, 4*time.Second)
			defer cancel()
			// err := worker(c)
			// if err != nil {
			// 	fmt.Println(err.Error())
			// }
			worker(c)

		}()

		time.Sleep(500 * time.Millisecond)
	}

	time.Sleep(20 * time.Second)

}

func worker(ctx context.Context) error {
	id := uuid.New().String()

	// setLock(id)
	if !setLock(ctx, id) {
		msg := fmt.Sprintf("%s xxxx> 没有抢到锁", id)
		return errors.New(msg)
	}
	defer deleteLock(ctx, id)
	go refreshLock(ctx, id)

	fmt.Println(id, "抢到了锁")
	// stmt
	n := rand.Intn(9)
	fmt.Println(id, "执行时间: ", n, "s")
	time.Sleep(time.Duration(n) * time.Second)
	fmt.Println(id, "yyyyy====> 执行完成。")
	return nil
}

// setLock 设置锁
func setLock(ctx context.Context, uuid string) bool {
	return rdb.SetNX(ctx, lockerKey(), uuid, 3*time.Second).Val()
}

// refreshLock 1. 负责刷新锁时间， 避免应用逻辑时间过长而丢失锁。
//             2. 当 context 收到结束信号时删除锁， 一直刷新抢锁而造成饥饿问题。
func refreshLock(ctx context.Context, uuid string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(uuid, "超时 ====> 释放了锁")
			deleteLock(context.Background(), uuid)
			return
		default:
			id := rdb.Get(ctx, lockerKey()).Val()
			if id == uuid {
				// fmt.Println("@@@ id=>", id, "  uuid=>", uuid)
				rdb.Expire(ctx, lockerKey(), 3*time.Second)
				fmt.Println(uuid, "====> 刷新了锁")
				time.Sleep(1 * time.Second)
			}
		}
	}
}

// deleteLock 删除锁, 删除锁时判断锁是否为自己的。 这里的 CAS 是非原子操作。
func deleteLock(ctx context.Context, uuid string) {

	// 这里遇到一个问题， 当 context 已经取消的时候， 获取到的值为空， 即使 key 存在
	// 因此这里强行覆盖 context
	ctx = context.Background()
	id := rdb.Get(ctx, lockerKey()).Val()
	// fmt.Println("id=>", id, "  uuid=>", uuid)
	if id == uuid {
		rdb.Del(ctx, lockerKey())
		fmt.Println(uuid, "释放了锁")
		return
	}

	fmt.Println(">>> ", uuid, "锁已经没有了")
}

func lockerKey() string {
	return "pcc:locker"
}
