# Redis 数据类型: List 列表

1. 单键多值
2. Redis 列表是简单的 **字符串列表** ，按照插入顺序排序。你可以添加一个元素到列表的 **头部（左边）** 或者 **尾部（右边）** 。

3. 它的底层实际是个 **双向链表** ，对两端的操作性能很高，通过索引下标的操作 *中间的节点性能会较差* 。

![20220311102948](https://assets.tangx.in/blog/redis-type-list/20220311102948.png)



## 常用命令

把双向链表想象成一根水管， 加入数据 (Push) 就是向水管里 **放入** 乒乓球。 弹出数据(Pop) 就是从水管理 **拿出** 乒乓球。

由于可以左右两边分别操作， 因此使用 `L / R` 表示方向 

![20220311102927](https://assets.tangx.in/blog/redis-type-list/20220311102927.png)

### 添加数据 `LPush` and `RPush`

`LPush` 从左侧加入数据， `RPush` 从右侧加入数据。 一次操作可以 **添加** 多个数据。

```bash
LPUSH key element [element ...]
RPUSH key element [element ...]
```

从下面案例可以得知， **无论从左右侧添加， 先添加的数据都更靠中间** 

```bash
127.0.0.1:6379[10]> lpush chan zero
(integer) 1
127.0.0.1:6379[10]> lrange chan 0 -1
1) "zero"
127.0.0.1:6379[10]> lrange chan 0 -1
1) "zero"
127.0.0.1:6379[10]> lpush chan l1 l2 l3
(integer) 4
127.0.0.1:6379[10]> rpush chan r1 r2 r3
(integer) 7
127.0.0.1:6379[10]> lrange chan 0 -1
1) "l3"
2) "l2"
3) "l1"
4) "zero"
5) "r1"
6) "r2"
7) "r3"
```


### 当数组存在时再添加: `LPushX` and `RPushX`

只有到目标数组 **存在时** ， `LPushX` 和 `RPushX` 才会将元素加入到数组中。

```bash
LPUSHX key element [element ...]
RPUSHX key element [element ...]
```

分别 **从左或从右** 向数组中添加 **至少一个** 元素。

```bash
127.0.0.1:6379> lpush list_x "world"
(integer) 1
127.0.0.1:6379> LPUSHX list_x "hello"
(integer) 2
127.0.0.1:6379> LRANGE list_x 0 -1
1) "hello"
2) "world"
127.0.0.1:6379>
127.0.0.1:6379> LPUSHX lsit_no "hello"
(integer) 0
127.0.0.1:6379> LRANGE list_no 0 -1
(empty array)
127.0.0.1:6379>
```

### 遍历元素: `LRANGE`

遍历数据只有命令 `LRange`， 即只能从左侧开始，不能从右侧开始。

```bash
LRANGE key start stop
```

`start/stop` 表示了数据的下标位置。  `0` 表示第一个元素， `-1` 表示最后一个元素。

```bash
127.0.0.1:6379[10]> LRANGE chan 0 -1
1) "l3"
2) "l2"
3) "l1"
4) "zero"
5) "r1"
6) "r2"
7) "r3"
127.0.0.1:6379[10]> LRANGE chan 3 -3
1) "zero"
2) "r1"
```

### 取出元素: `LPOP` 和 `RPOP`

与 push 相反， pop 是取出左右两边最外面的 **N个** 元素

```bash
LPOP key [count]

```

`count` 为可选参数， 默认值为 **1**。

```bash
127.0.0.1:6379[10]> LPOP chan 2
1) "l3"
2) "l2"

127.0.0.1:6379[10]> RPOP chan
"r3"

127.0.0.1:6379[10]> LRANGE chan 0 -1
1) "l1"
2) "zero"
3) "r1"
4) "r2"
```

### 根据索引获取列表元素: `LIndex`

根据列表中元素的索引获取元素。 

```
LINDEX key index
```

`index` 表示元素索引位置。 `0` 表示第一个元素， `-1` 表示最后一个元素。

```bash
127.0.0.1:6379[10]> LRANGE chan 0 -1
1) "l1"
2) "zero"
3) "r1"
4) "r2"
127.0.0.1:6379[10]> LINDEX chan 0
"l1"
127.0.0.1:6379[10]> LINDEX chan 2
"r1"
127.0.0.1:6379[10]> LINDEX chan -2
"r1"
```

### 列表长度: `LLEN`

返回列表元素个数

```
LLEN key
```

ex:

```bash
127.0.0.1:6379[10]> LLEN chan
(integer) 4
```

### 在列表中间插入元素: `LINSERT`

在 **列表** 所有两侧 **push/pop** 数据效率是最高的， 但 redis 也支持在列表中间添加元素。

```
LINSERT key BEFORE|AFTER pivot element
```

`pivot` 


```
127.0.0.1:6379[10]> LRANGE chan 0 -1
1) "l1"
2) "zero"
3) "r1"
4) "r2"
127.0.0.1:6379[10]> LINSERT chan BEFORE zero insL1
(integer) 5
127.0.0.1:6379[10]> LINSERT chan AFTER zero insR1
(integer) 6
127.0.0.1:6379[10]> LRANGE chan 0 -1
1) "l1"
2) "insL1"
3) "zero"
4) "insR1"
5) "r1"
6) "r2"
```

### 删除元素: `LREM`

通过 `LREM` 可以删除列表中 **指定数量** 的 **特定元素**

```
LREM key count element
```

1. `element` 需要倍删除的元素。
2. `count` 删除数量个数
  + `count > 0` **从左往右** 计数。
  + `count < 0` **从右往左** 计数。
  + `count = 0` **删除所有** 符合的元素。

```bash
127.0.0.1:6379> RPUSH mylist hello
(integer) 1
127.0.0.1:6379> RPUSH mylist hello
(integer) 2
127.0.0.1:6379> RPUSH mylist foo
(integer) 3
127.0.0.1:6379> RPUSH mylist hello
(integer) 4
127.0.0.1:6379> LRANGE mylist 0 -1
1) "hello"
2) "hello"  # 被删
3) "foo"
4) "hello"  # 被删
127.0.0.1:6379> LREM mylist -2 hello
(integer) 2
127.0.0.1:6379> LRANGE mylist 0 -1
1) "hello"
2) "foo"
```


### 替换元素: `LSET`

使用 `LSET` 根据 **下标** 替换数组中的对应元素。

```
LSET key index element
```

1. `index` 下标位置。 `0` 为第一元素； `-1` 为最后一个元素。
2. `element` 新元素的值

```bash
127.0.0.1:6379> LRANGE mylist 0 -1
1) "foo"
2) "foo"
3) "hello"
4) "foo"

127.0.0.1:6379> LSET mylist 0 FOO
OK
127.0.0.1:6379> LSET mylist -1 __FOO__
OK
127.0.0.1:6379> LRANGE mylist 0 -1
1) "FOO"
2) "foo"
3) "hello"
4) "__FOO__"
```


### 元素转移: 右出左进 `RpopLpush`

将一个元素从 **数组A** 的 **右侧** 取出， 并放入 **数组B** 的左侧

```
RPOPLPUSH source destination
```

1. `source` 源数组
2. `destination` 目标数组


```bash
127.0.0.1:6379> rpush list_A 1 2 3 4
(integer) 4
127.0.0.1:6379> rpush list_B a b c d
(integer) 4
127.0.0.1:6379> RPOPLPUSH list_A list_B
"4"
127.0.0.1:6379> LRANGE list_A 0 -1
1) "1"
2) "2"
3) "3"
127.0.0.1:6379> lrange list_B 0 -1
1) "4"
2) "a"
3) "b"
4) "c"
5) "d"
```

## 其他命令

更多 list 相关命令

> https://redis.io/commands/#list


## 数据结构

1. List 的数据结构为 **快速链表 `quickList`** 。
2. 首先在列表 **元素较少的情况下会使用一块连续的内存存储** ，这个结构是 ziplist，也即是 **压缩列** 。它将所有的元素紧挨着一起存储，分配的是一块连续的内存。

3. **当数据量比较多的时候才会改成 quicklist** 。 因为普通的链表需要的附加指针空间太大，会比较浪费空间。比如这个列表里存的只是 int 类型的数据，结构上还需要两个额外的指针prev和next。
 
4. Redis 将链表和 ziplist 结合起来组成了 quicklist 。也就是将多个 ziplist 使用双向指针串起来使用。 这样既满足了快速的插入删除性能， 又不会出现太大的空间冗余。

![20220311145251](https://assets.tangx.in/blog/redis-type-list/20220311145251.png)