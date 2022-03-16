# go-redis and context


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