---
id: 35
database_id: 318617483
node_id: MDU6SXNzdWUzMTg2MTc0ODM=
status: closed
title: "handle panics"
labels: ["status:blocked","type:critical"]
url: https://github.com/kamilsk/check/issues/35
created_at: 2018-04-28T09:15:06Z
updated_at: 2018-04-29T15:35:20Z
---

# handle panics

```
not handled panic in goroutines
	if err := func() (err error) {
		defer grace.Recover(&err)
		err = cmd.RootCmd.Execute()
		return
	}(); err != nil {
if goroutine has recover then data race is appeared

...
==================
Found 1 data race(s)
exit status 66
make: *** [cmd-urls] Error 1
```
