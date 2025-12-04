---
id: 48
database_id: 319257956
node_id: MDU6SXNzdWUzMTkyNTc5NTY=
status: open
title: "improve code quality, second iteration"
labels: ["type:refactoring"]
url: https://github.com/kamilsk/check/issues/48
created_at: 2018-05-01T17:01:04Z
updated_at: 2018-05-01T17:01:04Z
---

# improve code quality, second iteration

- [ ] up code coverage to 100% https://scrutinizer-ci.com/g/kamilsk/check/code-structure/master/code-coverage
- [ ] refactor some worst functions https://scrutinizer-ci.com/g/kamilsk/check/code-structure/master
- [ ] solve gometalinter issues
```
http/availability/crawler.go:152:21:warning: error return value not checked (el.Request.Visit(href)) (errcheck)
http/availability/printer.go:73:22:warning: error return value not checked (critical().Fprintf(w, "report %q has error %q\n", site.Name, site.Error)) (errcheck)
http/availability/printer.go:75:23:warning: error return value not checked (critical().Fprintf(ioutil.Discard, "stack trace: %#+v\n", stack) // for future) (errcheck)
http/availability/printer.go:82:37:warning: error return value not checked (colorize(page.StatusCode).Fprintf(w, head, page.StatusCode, page.Location)) (errcheck)
http/availability/printer.go:86:39:warning: error return value not checked (colorize(link.StatusCode).Fprintf(w, foot, link.StatusCode, link.FullLocation(sep))) (errcheck)
http/availability/printer.go:89:38:warning: error return value not checked (colorize(link.StatusCode).Fprintf(w, body, link.StatusCode, link.FullLocation(sep))) (errcheck)
http/availability/printer.go:93:22:warning: error return value not checked (critical().Fprintf(w, "found problems on the site %q\n", site.Name)) (errcheck)
http/availability/printer.go:95:23:warning: error return value not checked (critical().Fprintf(w, "- [%d] %s `%+v`\n", i, problem.Message, problem.Context)) (errcheck)
http/availability/report.go:127::warning: cyclomatic complexity 17 of function (*Site).listen() is high (> 10) (gocyclo)
http/availability/report.go:177::warning: declaration of "found" shadows declaration at report.go:173 (vetshadow)
http/availability/report.go:226:36:warning: exported func NewReadableEventBus returns unexported type chan availability.event, which can be annoying to use (golint)
```
