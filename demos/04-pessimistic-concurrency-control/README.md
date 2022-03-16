# 悲观锁

![20220316221324](https://assets.tangx.in/blog/README/20220316221324.png)



## trouble shoot

### context 过期导致无法获取 redis 值

```go
// deleteLock 删除锁, 删除锁时判断锁是否为自己的。 这里的 CAS 是非原子操作。
func deleteLock(ctx context.Context, uuid string) {

	// 这里遇到一个问题， 当 context 已经被取消的时候， 获取到的值为空， 即使 key 存在
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
```