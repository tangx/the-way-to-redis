package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/tangx/the-way-to-redis/demos/global"
)

func main() {

	u, err := GetUserById("10001")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", u)
}

func GetUserById(id string) (*User, error) {
	key := userRdbKey(id)

	r, err := global.RdbClient.Get(context.Background(), key).Result()
	if err == nil {
		user := &User{}
		user.Unmarshal(r)
		return user, nil
	}

	if errors.Is(err, redis.Nil) {
		// read from db
		// return nil, errors.New("read from db failed")

		u := &User{
			Name: "from db",
			Age:  999,
		}

		cacheUser(id, u)

		fmt.Println("从数据库中获取数据")
		return u, nil

	}

	return nil, err
}

func userRdbKey(id string) string {
	k := fmt.Sprintf("appname:userid:%s", id)
	// fmt.Println(k)
	return k
}

type User struct {
	Name string
	Age  int
}

func (user *User) Marshal() string {
	b, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (user *User) Unmarshal(s string) {
	b := []byte(s)

	err := json.Unmarshal(b, user)
	if err != nil {
		panic(err)
	}
	return
}

func cacheUser(id string, u *User) {
	key := userRdbKey(id)

	global.RdbClient.Set(context.Background(), key, u.Marshal(), 0)
}
