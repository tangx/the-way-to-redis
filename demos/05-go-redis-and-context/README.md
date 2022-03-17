# go-redis and context


这里遇到一个问题， 当 context 已经取消的时候，
.   获取到的值为空， 即使 key 存在
.   get val=, err=context deadline exceeded
    因此如果需要， 可以在这里因此这里强行覆盖 context
    但最好在上层控制传入的 context

```bash
→ go run .
set val=OK, err=<nil>
0 : get val=10000, err=<nil>
1 : get val=10000, err=<nil>
2 : get val=10000, err=<nil>
3 : get val=10000, err=<nil>
4 : get val=10000, err=<nil>
5 : get val=, err=context deadline exceeded
6 : get val=, err=context deadline exceeded
7 : get val=, err=context deadline exceeded
8 : get val=, err=context deadline exceeded
```