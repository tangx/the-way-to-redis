package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tangx/the-way-to-redis/demos/global"
)

func main() {

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	v, err := global.RdbClient.Set(ctx, key(), "10000", 0).Result()
	fmt.Printf("set val=%s, err=%v\n", v, err)

	for i := 0; i < 9; i++ {
		vv, err := global.RdbClient.Get(ctx, key()).Result()
		fmt.Printf("%d : get val=%s, err=%v\n", i, vv, err)
		time.Sleep(1 * time.Second)
	}
	// go run .
	// set val=OK, err=<nil>
	// 0 : get val=10000, err=<nil>
	// 1 : get val=10000, err=<nil>
	// 2 : get val=10000, err=<nil>
	// 3 : get val=10000, err=<nil>
	// 4 : get val=10000, err=<nil>
	// 5 : get val=, err=context deadline exceeded
	// 6 : get val=, err=context deadline exceeded
	// 7 : get val=, err=context deadline exceeded
	// 8 : get val=, err=context deadline exceeded
}

func key() string {
	return "go-redis-and-context"
}
