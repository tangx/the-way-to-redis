# 分布式锁

应用场景：

商场有一个促销活动， 推出 10 台 iphone13 的 0 元购活动。

要求：
1. 使用 redis 存储剩余数量， 并完成售卖功能。
2. 使用 redis 完成分布式锁解决超卖的问题。
3. 如果数量为 500 时， 锁可能导致购买失败，产品剩余。
4. redis 连接池大小问题

```go
Error distribution:
  [15]	Get "http://127.0.0.1:8081/promote/iphone": dial tcp 127.0.0.1:8081: connect: connection reset by peer
  [293]	Get "http://127.0.0.1:8081/promote/iphone": dial tcp 127.0.0.1:8081: socket: too many open files
  [1]	Get "http://127.0.0.1:8081/promote/iphone": read tcp 127.0.0.1:59348->127.0.0.1:8081: read: connection reset by peer
  [1]	Get "http://127.0.0.1:8081/promote/iphone": read tcp 127.0.0.1:59349->127.0.0.1:8081: read: connection reset by peer
  [1]	Get "http://127.0.0.1:8081/promote/iphone": read tcp 127.0.0.1:59350->127.0.0.1:8081: read: connection reset by peer
  ```

  ## v1 

  1. 完成基础功能
  2. 所有 redis db 操作单独封装函数， 有超卖现象。

  