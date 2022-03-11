# Redis 数据类型 String

1. String 是 Redis 最基本的类型，**一个key对应一个value** 。
2. String 类型是 **二进制安全** 的。 意味着 Redis 的 string 可以包含任何数据。 比如jpg图片或者序列化的对象； 也可以 **数字** 类型进行累加操作 `INCR`
3. String 类型是 Redis 最基本的数据类型，一个 Redis 中字符串 value 最多可以是 **512M**


## 基本 `SET` 命令

最简单的设置方式是 

```bash
SET key value
```

但完整语法是

```bash
SET key value [EX seconds|PX milliseconds|EXAT timestamp|PXAT milliseconds-timestamp|KEEPTTL] [NX|XX] [GET]
```

1. 过期时间设置
  + `EX seconds` : 设置 key 的时候同时设置 n 秒过期。与 `SETEX` 命令行为一致。 `SET key100 value100 EX 20`。
  + `PX milliseconds`: 设置 n 毫秒过期， 与 EX 互斥。 `SET key1 value1 PX 100`
  + `KEEPTTL`: 保留剩余的过期时间。 不使用此参数则剩余过期时间将被覆盖。
2. Key 存在条件
  + `NX` : 当 Key 不存在时生效. 与 `SETNX` 行为一致。
  + `XX` : 当 Key 存在时生效。
3. `GET`: 存数据的同时取出原来的值。 如果 key 不存在返回 `nil`； 如果 `TYPE key` 不是 string 的报错。


## 扩展 SET 命令


1. `SETEX key value seconds`: 保存 key 并设置过期时间
2. `SETNX key value`: 当 key 不存在的时候生效。
3. `MSET key value [key value ...]`: 同时设置多个 key-value 对


###  `SET ... keepttl` 参数 demo

```bash
127.0.0.1:6379[1]> SET k100 value EX 50
OK
127.0.0.1:6379[1]> TTL k100
(integer) 47
# 使用 KeepTTL 保留剩余过期时间
127.0.0.1:6379[1]> SET k100 value2 keepttl
OK
127.0.0.1:6379[1]> ttl k100
(integer) 28

127.0.0.1:6379[1]> set k200 value2 ex 50
OK
127.0.0.1:6379[1]> ttl k200
(integer) 41

# 没有使用 KeepTTL, 过期时间被覆盖，
127.0.0.1:6379[1]> set k200 value4
OK
127.0.0.1:6379[1]> ttl k200
(integer) -1
```


### `SET ... GET` demo

```bash
127.0.0.1:6379[1]> SET key 100
OK
127.0.0.1:6379[1]> SET key 200 GET
"100"

# key 不存在返回 nil
127.0.0.1:6379[1]> SET key99 value GET
(nil)
```


## 获取值 `GET` and `MGET`

1. `GET`: 获取单个 key 的值
2. `MGET`: 同时获取多个 key 的值

```bash
# get
127.0.0.1:6379[1]> GET key
"200"

# mget
127.0.0.1:6379[1]> MGET key key99 key9 key200
1) "200"
2) "value"
3) "GET"
4) (nil)
```

## 字符串操作


### 字符串拼接: `APPEND`

1. `APPEND key value` 在原有字符串后面拼接新的字段。 如果原 key 不存在则行为类似 `set key value`

```bash
127.0.0.1:6379> APPEND k9 value
(integer) 5
127.0.0.1:6379> get k9
"value"
127.0.0.1:6379> APPEND k9 1234
(integer) 9
127.0.0.1:6379> get k9
"value1234"
```

### 长度统计: `STRLEN`

1. `STRLEN key`: 统计 key 的值的长度。

```bash
127.0.0.1:6379> get k9
"value1234"
127.0.0.1:6379> STRLEN k9
(integer) 9
```



### 片段设置 `SETRANAGE` and `GETRANAGE`

1. `SETRANGE key offset value`: 从指定的 `offset` 位置操作片段， 替换原有片段。 如果 offset 位置比原 key 长度大， 将用 `zero` 占位。

2. `GETRANGE key start end`: 获取指定 **收尾** 的字符串。 **左右包含**


```bash
# SetRange
127.0.0.1:6379> SET key1 "Hello World"
OK
127.0.0.1:6379>
127.0.0.1:6379> SETRANGE key1 6 "Redis"
(integer) 11
127.0.0.1:6379> GET key1
"Hello Redis"


# SetRange: 用 zero 占位
127.0.0.1:6379> SETRANGE key2 6 "Redis"
(integer) 11
127.0.0.1:6379>
127.0.0.1:6379> GET key2
"\x00\x00\x00\x00\x00\x00Redis"


# GetRange
127.0.0.1:6379> set n 0123456
OK
127.0.0.1:6379> GETRANGE n 0 3
"0123"
127.0.0.1:6379> GETRANGE n 1 3
"123"
127.0.0.1:6379> GETRANGE n 1 10
"123456"


redis> SET mykey "This is a string"
"OK"
redis> GETRANGE mykey 0 3
"This"
redis> GETRANGE mykey -3 -1
"ing"
redis> GETRANGE mykey 0 -1
"This is a string"
redis> GETRANGE mykey 10 100
"string"
redis>
```

## 算数操作

string 类型对纯数字类型的值可以进行 **加减操作**

1. `INCR key`: 累加 1
2. `INCRBY key increment`: 指定累加值。 `ex: INCRBY key 10`
3. `DECR key`: 累减 1
4. `DECRBY key increment`: 指定减少值。 `ex: DECRBY key 10`

