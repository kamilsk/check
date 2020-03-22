> # 🔬 check
>
> Tool to check something.

[![Build][build.icon]][build.page]
[![Template][template.icon]][template.page]

## 💡 Idea

...

Full description of the idea is available [here][design.page].

## 🏆 Motivation

...

## 🤼‍♂️ How to

### check urls

Fast website link checker.

```bash
$ check urls https://kamil.samigullin.info/
# [200] https://kamil.samigullin.info/
#     ├───[200] https://howilive.ru/en/
#     ├───[200] https://kamil.samigullin.info/ru/
#     ├───[301] https://m.do.co/c/b2a387de5da4 -> (Moved Permanently) -> https://...
#     ├───[302] https://kamil.samigullin.info/goto?url=https://github.com/kamilsk -> (Found) -> https://...
#     ├───[302] https://kamil.samigullin.info/goto?url=https://twitter.com/ikamilsk -> (Found) -> https://...
#     └───[302] https://kamil.samigullin.info/goto?url=https://www.linkedin.com/in/kamilsk -> (Found) -> https://...
# [200] https://kamil.samigullin.info/ru/
#     ├───[200] https://howilive.ru/
#     ├───[200] https://kamil.samigullin.info/
#     ├───[301] https://m.do.co/c/b2a387de5da4 -> (Moved Permanently) -> https://...
#     ├───[302] https://kamil.samigullin.info/goto?url=https://github.com/kamilsk -> (Found) -> https://...
#     ├───[302] https://kamil.samigullin.info/goto?url=https://twitter.com/ikamilsk -> (Found) -> https://...
#     └───[302] https://kamil.samigullin.info/goto?url=https://www.linkedin.com/in/kamilsk -> (Found) -> https://...
$ check urls https://www.octolab.org/ | grep '\[3[0-9][0-9]\]'
#     ├───[301] https://m.do.co/c/b2a387de5da4 -> (Moved Permanently) -> https://...
#     ├───...
```

## 🧩 Installation

### Homebrew

```bash
$ brew install kamilsk/tap/check
```

### Binary

```bash
$ curl -sSL https://bit.ly/install-check | sh
# or
$ wget -qO- https://bit.ly/install-check | sh
```

### Source

```bash
# use standard go tools
$ go get -u github.com/kamilsk/check
# or use egg tool
$ egg tools add github.com/kamilsk/check
```

> [egg][]<sup id="anchor-egg">[1](#egg)</sup> is an `extended go get`.

### Bash and Zsh completions

```bash
$ check completion bash > /path/to/bash_completion.d/check.sh
$ check completion zsh  > /path/to/zsh-completions/_check.zsh
```

<sup id="egg">1</sup> The project is still in prototyping. [↩](#anchor-egg)

---

made with ❤️ for everyone

[build.page]:       https://travis-ci.org/kamilsk/check
[build.icon]:       https://travis-ci.org/kamilsk/check.svg?branch=master
[design.page]:      https://www.notion.so/33715348cc114ea79dd350a25d16e0b0?r=0b753cbf767346f5a6fd51194829a2f3
[promo.page]:       https://github.com/kamilsk/check
[template.page]:    https://github.com/octomation/go-tool
[template.icon]:    https://img.shields.io/badge/template-go--tool-blue

[egg]:              https://github.com/kamilsk/egg
