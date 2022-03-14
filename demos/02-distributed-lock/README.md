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

# v1 版本总结

1. 实现了需求的基础能力， 能通过 redis 进行库存管理和增减。 
2. 但是由于没有做并发控制， 因此在并发高的时候会出现超卖的情况。
  + 可以使用分布式锁式控制

3. redis 的每一步操作都封装成函数了， 有点分散。



## v2 版本总结

v2 版本使用了 watch 和 pipeline ， 但是并不能解决超卖的问题。 只能降低超卖的数量。

在 pipeline 事务中， 只有在 `pipe.EXEC()` 的时候才会被执行。
因此在代码逻辑中使用 `pipe.Get(ctx, iphoneStock)` 获取库存量的时候， 结果永远是 0。 

**总结**:  只要 **库存检查** 和 **库存修改操作** 不在一个事务中， 就会超卖。



## v3 版本， 使用 lua 执行原子操作。

为了解决 **库存检查与库存修改** 的原子性问题， 可能只有 lua 是相对好的方案了。

> https://help.aliyun.com/document_detail/63920.html


