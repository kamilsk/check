> # check [![Tweet](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/intent/tweet?text=Check%20Something%20as%20a%20Service&url=https://github.com/kamilsk/check&via=ikamilsk&hashtags=go,tool,website-audit)
> [![Analytics](https://ga-beacon.appspot.com/UA-109817251-19/check/readme?pixel)](https://github.com/kamilsk/check)
> Check Something as a Service.

[![Patreon](https://img.shields.io/badge/patreon-donate-orange.svg)](https://www.patreon.com/octolab)
[![Build Status](https://travis-ci.org/kamilsk/check.svg?branch=master)](https://travis-ci.org/kamilsk/check)
[![Code Coverage](https://scrutinizer-ci.com/g/kamilsk/check/badges/coverage.png?b=master)](https://scrutinizer-ci.com/g/kamilsk/check/?branch=master)
[![Code Quality](https://scrutinizer-ci.com/g/kamilsk/check/badges/quality-score.png?b=master)](https://scrutinizer-ci.com/g/kamilsk/check/?branch=master)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Usage

### Quick start

#### check urls

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

## Specification

### CLI

```bash
$ check --help
Usage:
  check [command]

Available Commands:
  completion  Print Bash or Zsh completion
  help        Help about any command
  urls        Check all internal URLs on availability
  version     Show application version

Flags:
  -h, --help   help for check

Use "check [command] --help" for more information about a command.
```

## Installation

### Brew

```bash
$ brew install kamilsk/tap/check
```

### Binary

```bash
$ export VER=1.0.0      # all available versions are on https://github.com/kamilsk/check/releases
$ export REQ_OS=Linux   # macOS and Windows are also available
$ export REQ_ARCH=64bit # 32bit is also available
$ wget -q -O check.tar.gz \
       https://github.com/kamilsk/check/releases/download/"${VER}/check_${VER}_${REQ_OS}-${REQ_ARCH}".tar.gz
$ tar xf check.tar.gz -C "${GOPATH}"/bin/ && rm check.tar.gz
```

### From source code

```bash
$ egg github.com/kamilsk/check@^1.0.0 -- make test install
```

#### Mirror

```bash
$ egg bitbucket.org/kamilsk/check@^1.0.0 -- make test install
```

> [egg](https://github.com/kamilsk/egg) is an `extended go get`.

### Bash and Zsh completions

You can find completion files [here](https://github.com/kamilsk/shared/tree/dotfiles/bash_completion.d) or
build your own using these commands

```bash
$ check completion bash > /path/to/bash_completion.d/check.sh
$ check completion zsh  > /path/to/zsh-completions/_check.zsh
```

## Notes

- brief roadmap
  - [x] v1: MVP
  - [ ] v2: check redirects
  - [ ] v3: check repository
  - [ ] v4: check package
  - [ ] v5: distributed run
  - [ ] integrate with Status, SaaS
- [research](../../tree/research)
- tested on Go 1.8, 1.9 and 1.10

---

[![@kamilsk](https://img.shields.io/badge/author-%40kamilsk-blue.svg)](https://twitter.com/ikamilsk)
[![@octolab](https://img.shields.io/badge/sponsor-%40octolab-blue.svg)](https://twitter.com/octolab_inc)

made with ❤️ by [OctoLab](https://www.octolab.org/)
