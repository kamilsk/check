---
id: 22
database_id: 317601406
node_id: MDU6SXNzdWUzMTc2MDE0MDY=
status: closed
title: "configure user agent"
labels: ["type:refactoring"]
url: https://github.com/kamilsk/check/issues/22
created_at: 2018-04-25T12:13:26Z
updated_at: 2018-04-26T15:28:41Z
---

# configure user agent

```go
func UserAgent() func(*colly.Collector) {
	return colly.UserAgent("check")
}
```
