package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/tangx/the-way-to-redis/demos/global"
)

func Test_gen(t *testing.T) {
	u := &User{
		Name: "tangxin",
		Age:  20,
	}

	s := u.Marshal()

	key := userRdbKey("10001")
	global.RdbClient.Set(context.Background(), key, s, 0).Val()

}

func Test_GetUserById(t *testing.T) {
	u, err := GetUserById("10001")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", u)

	u, err = GetUserById("10002")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", u)
}
