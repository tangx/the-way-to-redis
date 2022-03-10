# Redis 键操作

## 数据库操作

```bash
# 切换数据库 ， 0 ~ 15
127.0.0.1:6379> SELECT 1
OK

# 查看当前数据库中所有的 key 的数量
127.0.0.1:6379[1]> DBSIZE
(integer) 1

# 清空所有数据库
127.0.0.1:6379[1]> FLUSHALL
OK

# 清空当前数据库
127.0.0.1:6379[1]> FLUSHDB
OK
```

## key 的基本操作


### key 的查询

```bash
127.0.0.1:6379[1]> set k1 value
OK

# 查看所有 key
127.0.0.1:6379[1]> KEYS *
1) "k1"

# 匹配关键字 key
127.0.0.1:6379[1]> KEYS k*
1) "k2"
2) "k1"

# 检查 key 是否存在
127.0.0.1:6379[1]> EXISTS k1
(integer) 1

# 查看 key 类型
127.0.0.1:6379[1]> TYPE k1
string
```

### key 的删除


```bash
# 删除 key
127.0.0.1:6379[1]> DEL k1
(integer) 1

# 仅将keys从keyspace元数据中删除，真正的删除会在后续异步操作
127.0.0.1:6379[1]> UNLINK k2
(integer) 1
```


### key 的过期时间管理

```bash
# 设置 key 过期时间， 单位 秒(s)
127.0.0.1:6379[1]> EXPIRE k1 20
(integer) 1

# 查看 key 剩余过期时间, -1 表示不过期
127.0.0.1:6379[1]> TTL k1
(integer) 16

127.0.0.1:6379[1]> TTL key100
(integer) -1
```

> 注意: 如果对一个已有过期时间的 key 进行时间设置， 将覆盖原有过期时间而非累加。

```bash
# 第一次设置过期时间
127.0.0.1:6379[1]> EXPIRE key100 20
(integer) 1
127.0.0.1:6379[1]> TTL key100
(integer) 17

# 重置过期时间
127.0.0.1:6379[1]> EXPIRE key100 100
(integer) 1
127.0.0.1:6379[1]> TTL key100
(integer) 97
```