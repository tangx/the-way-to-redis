# Redis 数据类型 Hash 映射

1. Redis hash 是一个 string 类型的 field 和 value 的映射表， hash 特别适合用于存储对象。
2. 一个 hash key 下可以保存多个键值对。

![20220311172131](https://assets.tangx.in/blog/redis-type-hash/20220311172131.png)


## 常用命令

> https://redis.io/commands/#hash

hash 的命令基本已 `H` 开头

### 设置与获取键值对: `HSET / HGET`

1. `HSET` 可以一次设置多组 filed-value 
2. `HGET` 一次只能获取 **一组** filed-value


```bash
HSET key field value [field value ...]
HGET key field
```


```bash
127.0.0.1:6379> HSET myhash k1 v1 k2 v2 k3 v3
(integer) 3

127.0.0.1:6379> HGET myhash k1
"v1"

127.0.0.1:6379> HGET myhash k1 k2  # 尝试使用 hget 获取多组 filed-value
(error) ERR wrong number of arguments for 'hget' command
```


### 设置与获取 `多组` 键值对: `HMSET / HMGET`

通常 `HGET / HSET` 一次只操作一组 filed-value。 而通过 `HMGET / HMSET` 操作多组 filed-value。

```
HMSET key field value [field value ...]
HMGET key field [field ...]
```

```bash
127.0.0.1:6379> HMSET myhash2 k3 v3 k4 v4
OK
127.0.0.1:6379> HMGET myhash2 k3 k4
1) "v3"
2) "v4"
```

### 获取所有键值对: `HGETALL`

通过 `HGetAll` 可以获取所有的 filed-value

```
HGETALL key
```

```bash
127.0.0.1:6379> HGETALL myhash2
1) "k3"
2) "v3"
3) "k4"
4) "v4"
```

### 删除键值对: `HDel`

通过 `HDel` 删除 **一个或多个** filed-value

```
HDEL key field [field ...]
```

```bash
127.0.0.1:6379> HDEL myhash k3 k5
(integer) 1
```

### 查询键值对是否存在: `HExists`

通过 `HExists` 查询某个 filed-value 是否存在。

```
HEXISTS key field
```

```bash
127.0.0.1:6379> HEXISTS myhash2 k3
(integer) 1
127.0.0.1:6379> HEXISTS myhash3 k5
(integer) 0
```


### 查询所有 filed 或 value: `HKeys` / `HVals`

1. `HKeys` 返回 hash 桶中的所有 filed
2. `HVals` 返回 hash 桶中的所有 value

```
HKEYS key
HVals key
```

```bash
127.0.0.1:6379> HKEYS myhash2
1) "k3"
2) "k4"
127.0.0.1:6379> HVALS myhash2
1) "v3"
2) "v4"
```


### 值增加: `HIncrBy`

`HINCRBY` 为 filed 对应的 value 增加相应的值。 

```
HINCRBY key field increment
```

`increment` 的值可以是负数，就表示了 **减** 操作。

> 注意： 没有 `HDecrBy` 命令进行减少操作。

```bash
127.0.0.1:6379> HSET myhash2 age 20
(integer) 1
127.0.0.1:6379> HINCRBY myhash2 age 10
(integer) 30
127.0.0.1:6379> HGET myhash2 age
"30"

127.0.0.1:6379> HINCRBY myhash2 age -15
(integer) 15
```

### 统计字段长度: `HStrLen`

使用 `HStrLen` 可以统计指定字段值的长度。

```
HSTRLEN key field
```

```bash
127.0.0.1:6379> HSET myhash f1 v1 f2 v2 f3 v3
(integer) 3
127.0.0.1:6379> HSTRLEN myhash f1
(integer) 2
```


### 获取随机字段: `HRandFiled`

使用 `HRandFiled` 可以随机返回 **一个或多个** **字段名称**， 及 **其值** （如指定）。

```
HRANDFIELD key [count [WITHVALUES]]
```

1. `count` 返回 N 个字段。 默认值为 **1**
2. `count WithValues` 返回字段名称及其值。

```bash
127.0.0.1:6379> HRANDFIELD myhash
"f3"

127.0.0.1:6379> HRANDFIELD myhash 1
1) "f3"

127.0.0.1:6379> HRANDFIELD myhash 1 withvalues
1) "f1"
2) "v1"

## WithValues 必须与 count 联用。
127.0.0.1:6379> HRANDFIELD myhash withvalues
(error) ERR value is not an integer or out of range
```

### 只设置 `不存在的字段` 值: `HSetNX`

`HSetNX` 只对 **不存在的字段** 生效。 如果字段已存在， 则略过。

```
HSETNX key field value
```

```bash
127.0.0.1:6379> HKEYS myhash
1) "f1"
2) "f2"
3) "f3"

127.0.0.1:6379> HSETNX myhash f1 VVV1
(integer) 0

127.0.0.1:6379> HSETNX myhash f99 vvv99
(integer) 1
```