# Redis 数据类型 Set 集合

Redis  `SET` 对外提供的功能与 list 类似是一个列表的功能，特殊之处在于 set 是可以 **自动排重的** ，当你需要存储一个列表数据，又不希望出现重复数据时， set 是一个很好的选择，并且set提供了**判断某个成员是否在一个 set 集合内的重要接口**， 这个也是 list 所不能提供的。

Redis 的 Set 是 string 类型的 **无序集合** 。 它 **底层其实是一个 value 为 null 的 hash 表** ，所以 **添加，删除，查找的复杂度都是 `O(1)`** 。 `O(1)` 是一个算法，随着数据的增加，执行时间的长短，如果是 O(1) ，数据增加、查找数据的时间不变。


## 常用命令

集合的命令通常都是以 `S` 开头。 入 `SADD / SREM`

> https://redis.io/commands/#set


### 添加元素: `SADD`

向集合中添加一个或多个成员 (元素)。 

这里用 **成员 member** 进行描述， 能更好的应对命令中的参数。

```
SADD key member [member ...]
```

在添加成员的过程中， 如果遇到相同成员， 只会保留一个。

```bash
127.0.0.1:6379> SADD myset v1 v2 v3 v1
(integer) 3
127.0.0.1:6379> SMEMBERS myset
1) "v2"
2) "v3"
3) "v1"
```

### 查看所有成员: `SMEMBERS`

返回集合中的所有成员


```bash
SMEMBERS key
```

```bash
127.0.0.1:6379> SMEMBERS myset
1) "v2"
2) "v3"
3) "v1"
```

### 成员是否存在于集合: `SISMEMBER`

检查成员是否存在于集合

```
SISMEMBER key member
```

```bash
127.0.0.1:6379> SISMEMBER myset v3
(integer) 1
127.0.0.1:6379> SISMEMBER myset v5
(integer) 0
```

### 查询集合中的成员数量: `SCARD`

返回一个集合总的所有成员数量总和。

```
SCARD key
```

```bash
127.0.0.1:6379> SMEMBERS myset
1) "v2"
2) "v3"
3) "v1"
127.0.0.1:6379> SCARD myset
(integer) 3
```

### 删除集合中的元素: `SREM`

从集合中删除一个或多个成员

```
SREM key member [member ...]
```

```bash
127.0.0.1:6379> SMEMBERS myset
1) "v2"
2) "v3"
3) "v1"
127.0.0.1:6379> SREM myset v3 v2
(integer) 2
127.0.0.1:6379> SMEMBERS myset
1) "v1"
```

### 从集合中抛出成员: `SPOP`

随机从集合中 **抛出 pop** **指定数量 (默认为 1) ** 的 **随机** 成员， 这些成员将在集合中被删除。

```
SPOP key [count]
```

```bash
127.0.0.1:6379> SMEMBERS myset
1) "v4"
2) "v2"
3) "v3"
4) "v1"
5) "v5"
127.0.0.1:6379> SPOP myset 3
1) "v2"
2) "v5"
3) "v3"
127.0.0.1:6379> SMEMBERS myset
1) "v1"
2) "v4"
```


### 从集合中随机选择成员: `SRandMember`

从集合中 **随机** 选择 **指定数量（默认为 1）** 的成员。 但这些成员 **不会** 被删除。

```
SRANDMEMBER key [count]
```


```
127.0.0.1:6379> SMEMBERS myset
1) "v1"
2) "v4"
127.0.0.1:6379> SRANDMEMBER myset
"v4"
127.0.0.1:6379> SMEMBERS myset
1) "v1"
2) "v4"
```


### 成员在集合之间转移: `SMove`

将成员从一个集合移动到另外一个集合

```
SMOVE source destination member
```


```bash
127.0.0.1:6379> SMEMBERS myset
1) "v1"
2) "v4"

127.0.0.1:6379> SMOVE myset otherset v4
(integer) 1
127.0.0.1:6379> SMEMBERS myset
1) "v1"
127.0.0.1:6379>
127.0.0.1:6379> SMEMBERS otherset
1) "v4"
```


### 集合的交、并、差集: `SInter / SUnion / SDiff`

对多个集合求值

1. `SInter` 对多个集合求交集。
2. `SUnion` 对多个集合求并集。
3. `SDiff` 对多个集合求差集。 即在 **源集合** 中出现但未在 **目标集合** 中出现的成员。

```
SINTER key [key ...]
SUNION key [key ...]
SDIFF  key [key ...]
```

```bash
127.0.0.1:6379> sadd set_a v1 v2 v3 v4 v5
(integer) 5
127.0.0.1:6379> sadd set_b v1 v3 v5 v7 v9
(integer) 5

## 求交集
127.0.0.1:6379> SINTER set_a set_b
1) "v3"
2) "v1"
3) "v5"

## 求并集
127.0.0.1:6379> SUNION set_a set_b
1) "v2"
2) "v5"
3) "v7"
4) "v1"
5) "v4"
6) "v3"
7) "v9"

## 求差集
127.0.0.1:6379> SDIFF set_a set_b  # a 中出现， b 中没出现。
1) "v2"
2) "v4"
127.0.0.1:6379> SDIFF set_b set_a  # b 中出现， a 中没出现。
1) "v7"
2) "v9"
```

### 集合交、并、差集另存为: `SInterStore / SUnionStore / SDiffStore`

将集合的 **交集、 并集、 差集** 结果保存到 **目标集合** 中。

```bash
SInterSTORE destination key [key ...]
SUnionSTORE destination key [key ...]
SDiffSTORE  destination key [key ...]
```

注意： 命令后 **先** 跟目标集合， 再跟 **源集合**

```bash
127.0.0.1:6379> SUNIONSTORE s_union set_a set_b
(integer) 7
127.0.0.1:6379> SINTERSTORE s_inter set_a set_b
(integer) 3
127.0.0.1:6379> SDIFFSTORE s_diff_ab set_a set_b
(integer) 2
```