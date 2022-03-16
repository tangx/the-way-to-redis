# the way to redis

## 目录

[目录 SUMMARY.md](./docs/SUMMARY.md)



## Demo

1. [短信验证码 demo](./demos/01-sms-code-verify/README.md)
2. [秒杀 demo](./demos/02-promote-demo/README.md)
  + [v1 - 购买基本功能实现, 存在超卖问题](./demos/02-promote-demo/v1/README.md)
  + [v2 - 使用 redis 事务，依然存在超卖](./demos/02-promote-demo/v2/README.md)
  + [v3 - 使用 lua 实现原子操作](./demos/02-promote-demo/v3/README.md)
3. [分布式锁-乐观锁 - watch + 事务](./demos/03-optimistic-concurrency-control/README.md)
4. [分布式锁-悲观锁 - 大锁 + 小锁](./demos/04-pessimistic-concurrency-control/README.md)
5. [context 过期对 go-redis 的影响](./demos/05-go-redis-and-context/README.md)
