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


## 其他命令

更多 list 相关命令

> https://redis.io/commands/#list
