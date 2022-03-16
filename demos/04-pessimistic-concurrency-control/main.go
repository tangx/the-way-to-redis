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

		// go func() {
		// 	c, cancel := context.WithTimeout(ctx, 4*time.Second)
		// 	defer cancel()
		// 	// err := worker(c)
		// 	// if err != nil {
		// 	// 	fmt.Println(err.Error())
		// 	// }
		// 	worker(c)

		// }()
		go worker(ctx)

		time.Sleep(500 * time.Millisecond)
	}

	time.Sleep(20 * time.Second)

}

func worker(ctx context.Context) error {
	// 1. 设置大锁， 主工作区耗时不能超过 4s
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	id := uuid.New().String()

	// 2. 抢占锁
	if !setLock(ctx, id) {
		msg := fmt.Sprintf("%s xxxx> 没有抢到锁", id)
		return errors.New(msg)
	}
	// 2.1 锁抢占成功， 程序退出时主动释放锁
	defer deleteLock(ctx, id)

	// 3. 业务逻辑stmt
	func() {
		fmt.Println(id, "抢到了锁")
		n := rand.Intn(9)
		fmt.Println(id, "执行时间: ", n, "s")
		time.Sleep(time.Duration(n) * time.Second)
		fmt.Println(id, "yyyyy====> 执行完成。")
	}()

	return nil
}

// setLock 设置锁
func setLock(ctx context.Context, uuid string) bool {
	ok := rdb.SetNX(ctx, lockerKey(), uuid, 3*time.Second).Val()

	// 2.2 开启协程， 为锁自动续期
	if ok {
		go refreshLock(ctx, uuid)
	}
	return ok
}

// refreshLock 1. 负责刷新锁时间， 避免应用逻辑时间过长而丢失锁。
//             2. 当 context 收到结束信号时删除锁， 一直刷新抢锁而造成饥饿问题。
func refreshLock(ctx context.Context, uuid string) {
	for {
		select {
		case <-ctx.Done():
			// 2. 当 ctx 超时过期时， 释放锁

			// go-redis 使用被取消后的 ctx 后， 将无法获得数据。
			//          即使 redis 中的数据真实存在。
			//          get val=, err=context deadline exceeded
			ctx := context.Background()
			deleteLock(ctx, uuid)
			fmt.Println(uuid, "超时 ====> 释放了锁")

			return
		default:
			// 1. 存续期内， 自动刷新锁
			id := rdb.Get(ctx, lockerKey()).Val()
			if id == uuid {
				rdb.Expire(ctx, lockerKey(), 3*time.Second)
				fmt.Println(uuid, "====> 刷新了锁")
				time.Sleep(1 * time.Second)
			}
		}
	}
}

// deleteLock 删除锁, 删除锁时判断锁是否为自己的。 这里的 CAS 是非原子操作。
func deleteLock(ctx context.Context, uuid string) {

	// 这里遇到一个问题， 当 context 已经取消的时候，
	// .   获取到的值为空， 即使 key 存在
	// .   get val=, err=context deadline exceeded
	//     因此如果需要， 可以在这里因此这里强行覆盖 context
	//     但最好在上层控制传入的 context
	// ctx = context.Background()
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
