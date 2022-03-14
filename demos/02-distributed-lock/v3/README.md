# v1 版本总结

1. 实现了需求的基础能力， 能通过 redis 进行库存管理和增减。 
2. 但是由于没有做并发控制， 因此在并发高的时候会出现超卖的情况。
  + 可以使用分布式锁式控制

3. redis 的每一步操作都封装成函数了， 有点分散。



## v2 版本

v2 版本使用了 watch 和 pipeline ， 但是并不能解决超卖的问题。 只能降低超卖的数量。

在 pipeline 事务中， 只有在 `pipe.EXEC()` 的时候才会被执行。
因此在代码逻辑中使用 `pipe.Get(ctx, iphoneStock)` 获取库存量的时候， 结果永远是 0。 

**总结**:  只要 **库存数量判断** 和 **库存修改操作** 不在一个事务中， 就会超卖。


### v3 版本

error log

```json
{
  "err": "ERR Error running script (call to f_963c5ad20bfbb28aa69e9b0222df3c80413f3ac8): @enable_strict_lua:15: user_script:14: Script attempted to access nonexistent global variable 'toboolean'",
  "msg": "forbiden"
}
```

错误位置: `user_script:14` 在 lua 的 14 行。
