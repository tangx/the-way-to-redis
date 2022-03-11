# Redis 数据类型 ZSet 有序集合

1. Redis有序集合 zset 与普通集合 set 非常相似， 是一个 **没有重复元素** 的字符串集合。
2. 不同之处是有序集合的每个成员都关联了一个 **评分（score）**, 这个评分（score）被用来按照从最低分到最高分的方式排序集合中的成员。 集合的成员是 **唯一** 的，但是评分可以是重复了 。
3. 因为元素是有序的, 所以你也可以很快的根据评分（score）或者次序（position）来获取一个范围的元素。
4. 访问有序集合的中间元素也是非常快的, 因此你能够使用有序集合作为一个没有重复成员的智能列表。


## 常用命令

> https://redis.io/commands/#sorted-set

与集合相关类似的命令就不写了。 参考 [集合 SET](./redis-type-set.md)。

**有序集合** 的命令通常以 Z 开头。

### 添加集合元素: `ZAdd`

向集合中添加成员和其评分权重。

```bash
ZADD key [NX|XX] [GT|LT] [CH] [INCR] score member [score member ...]
```

1. `socre member` : 权重 成员名
2. **存在条件** 语句
  + `NX`: 成员不存在时执行
  + `XX`: 仅成员存在时生效
3. **权重条件** 语句
  + `GT`: 仅当元素存在， 且 **新权重大于旧权重** 时才会更新权重。 不会产生新元素。
  + `LT`: 仅当元素存在， 且 **新权重小于旧权重** 时才会更新权重。 不会产生新元素。

```bash
127.0.0.1:6379> ZADD myzset 200 user1 100 user2 300 user3
(integer) 3
```

### 遍历所有字段: `ZRange`

便利返回字段名称。 可以根据条件返回。

```bash
ZRANGE key min max [BYSCORE|BYLEX] [REV] [LIMIT offset count] [WITHSCORES]
```

1. `min` / `max`: 指定起止位置。 第一个成员为 0 ， 最后一个成员为 -1。
2. 返回顺序
  + `ByScore` : 根据 **评分顺序** 返回， **默认**。
  + `ByLex `: 根据 **字母顺序** 返回
3. `Rev` 逆序返回， 可以与 `ByScore` 和 `ByLex` 联合使用。
4. `WithSocres` 同时返回评分。

```bash
# 默认
127.0.0.1:6379> ZRANGE myzset 0 -1
1) "user2"
2) "user1"
3) "user3"

# 逆序返回
127.0.0.1:6379> ZRANGE myzset 0 -1 rev
1) "user3"
2) "user1"
3) "user2"

# 同时返回权重
127.0.0.1:6379>  ZRANGE myzset 0 -1 withscores
1) "user2"
2) "100"
3) "user1"
4) "200"
5) "user3"
6) "300"

## 字母顺序不支持与权重共存
127.0.0.1:6379> ZRANGE myzset 0 -1 bylex withscores
(error) ERR syntax error, WITHSCORES not supported in combination with BYLEX
```


### 根据分数区间统计成员个数: `ZCount`

`ZCount` 根据分数区间统计成员个数。

```
ZCOUNT key min max
```

1. `min / max` 为 **最小/最大** 分数

```bash
ZCOUNT myzset 0 1000
(integer) 3
```

### 删除成员: `ZRem`

`ZRem` 根据名称删除一个或多个成员。

```
ZREM key member [member ...]
```

```bash
127.0.0.1:6379> ZRANGE myzset 0 -1
1) "user2"
2) "user1"
3) "user3"
127.0.0.1:6379> ZREM myzset user1
(integer) 1
127.0.0.1:6379> ZRANGE myzset 0 -1
1) "user2"
2) "user3"
```


### 成员排名: `ZRank`

`ZRank` 返回成员在集合中的排名位置， **根据权重排序**

```
ZRANK key member
```

1. 成员排名从 0 开始计算。
2. 如果成员不存在， 返回 nil 。

```bash
127.0.0.1:6379> ZADD myzset 1 one
(integer) 1
127.0.0.1:6379> ZADD myzset 2 two
(integer) 1
127.0.0.1:6379> ZADD myzset 3 three
(integer) 1

# 排名从 0 开始。
127.0.0.1:6379> ZRANK myzset one
(integer) 0
127.0.0.1:6379> ZRANK myzset three
(integer) 2

## 不存在返回 nil
127.0.0.1:6379> ZRANK myzset four
(nil)
```

## 数据结构

SortedSet(zset) 是 Redis 提供的一个非常特别的数据结构， 一方面它等价于Java的数据结构`Map<String, Double>`，可以给每一个元素value赋予一个权重score，另一方面它又类似于TreeSet，内部的元素会按照权重score进行排序，可以得到每个元素的名次，还可以通过score的范围来获取元素的列表。
zset底层使用了两个数据结构

1. hash，hash 的作用就是关联元素 value 和权重 score ，保障元素 value 的唯一性，可以通过元素 value 找到相应的 score 值。
2. 跳跃表，跳跃表的目的在于给元素 value 排序，根据 score 的范围获取元素列表。



## 小案例

根据访问量权重排序热门文章 top3

```bash
# 添加 6 个文章
127.0.0.1:6379> ZADD news 339 t1 123 t2 543 t3 443 t4 223 t5 111 t6
(integer) 6

# 查看分值排序
127.0.0.1:6379> ZRANGE news 0 -1
1) "t6"
2) "t2"
3) "t5"
4) "t1"
5) "t4"
6) "t3"

# 返回前三名及分值。 从 0 开始， 左右闭合区间。
127.0.0.1:6379> ZRANGE news 0 2 WithScores
1) "t6"
2) "111"
3) "t2"
4) "123"
5) "t5"
6) "223"
```

