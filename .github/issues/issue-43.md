---
id: 43
database_id: 318817524
node_id: MDU6SXNzdWUzMTg4MTc1MjQ=
status: open
title: "check redirects"
labels: ["type:feature"]
url: https://github.com/kamilsk/check/issues/43
created_at: 2018-04-30T08:29:56Z
updated_at: 2018-08-14T11:29:51Z
---

# check redirects

```bash
$ check urls https://some.site/ | tee sitemap.txt | check redirects > redirect.txt
```
