---
id: 30
database_id: 318345555
node_id: MDU6SXNzdWUzMTgzNDU1NTU=
status: closed
title: "bug with empty url"
labels: ["type:bug"]
url: https://github.com/kamilsk/check/issues/30
created_at: 2018-04-27T09:47:34Z
updated_at: 2018-04-27T13:33:41Z
---

# bug with empty url

```
panic: not consistent fetch result. link "" not found

goroutine 20 [running]:
github.com/kamilsk/check/http/availability.(*Site).listen(0xc4200a7e80, 0xc4200b88a0)
        /Users/kamilsk/Development/go/src/github.com/kamilsk/check/http/availability/report.go:158 +0x196b
github.com/kamilsk/check/http/availability.(*Site).Fetch.func1(0xc420216340, 0xc4200a7e80, 0xc4200b88a0)
        /Users/kamilsk/Development/go/src/github.com/kamilsk/check/http/availability/report.go:103 +0x6d
created by github.com/kamilsk/check/http/availability.(*Site).Fetch
        /Users/kamilsk/Development/go/src/github.com/kamilsk/check/http/availability/report.go:101 +0x10b
exit status 2
```
