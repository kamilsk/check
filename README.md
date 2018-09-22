> # 🔬 check [![Tweet][icon_twitter]][twitter_publish]
> [![Analytics][analytics_pixel]][page_promo]
> Check Something as a Service.

[![Patreon][icon_patreon]](https://www.patreon.com/octolab)
[![Build Status][icon_build]][page_build]
[![Code Coverage][icon_coverage]][page_quality]
[![Code Quality][icon_quality]][page_quality]
[![Research][icon_research]](../../tree/research)
[![License][icon_license]](LICENSE)

## Roadmap

- [x] v1: MVP
- [ ] v2: check redirects
- [ ] v3: check repositories
- [ ] v4: check packages

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

#### Bash and Zsh completions

You can find completion files [here](https://github.com/kamilsk/shared/tree/dotfiles/bash_completion.d) or
build your own using these commands

```bash
$ check completion bash > /path/to/bash_completion.d/check.sh
$ check completion zsh  > /path/to/zsh-completions/_check.zsh
```

## Installation

### Brew

```bash
$ brew install kamilsk/tap/check
```

### Binary

```bash
$ export REQ_VER=1.0.0  # all available versions are on https://github.com/kamilsk/check/releases
$ export REQ_OS=Linux   # macOS and Windows are also available
$ export REQ_ARCH=64bit # 32bit is also available
$ # wget -q -O check.tar.gz
$ curl -sL -o check.tar.gz \
       https://github.com/kamilsk/check/releases/download/"${REQ_VER}/check_${REQ_VER}_${REQ_OS}-${REQ_ARCH}".tar.gz
$ tar xf check.tar.gz -C "${GOPATH}"/bin/ && rm check.tar.gz
```

### From source code

```bash
$ egg github.com/kamilsk/check@^1.0.0 -- make test install
$ # or use mirror
$ egg bitbucket.org/kamilsk/check@^1.0.0 -- make test install
```

> [egg](https://github.com/kamilsk/egg)<sup id="anchor-egg">[1](#egg)</sup> is an `extended go get`.

<sup id="egg">1</sup> The project is still in prototyping. [↩](#anchor-egg)

---

[![@kamilsk][icon_tw_author]](https://twitter.com/ikamilsk)
[![@octolab][icon_tw_sponsor]](https://twitter.com/octolab_inc)

made with ❤️ by [OctoLab](https://www.octolab.org/)

[analytics_pixel]: https://ga-beacon.appspot.com/UA-109817251-19/check/readme?pixel

[icon_build]:      https://travis-ci.org/kamilsk/check.svg?branch=master
[icon_coverage]:   https://scrutinizer-ci.com/g/kamilsk/check/badges/coverage.png?b=master
[icon_gitter]:     https://badges.gitter.im/Join%20Chat.svg
[icon_license]:    https://img.shields.io/badge/license-MIT-blue.svg
[icon_patreon]:    https://img.shields.io/badge/patreon-donate-orange.svg
[icon_quality]:    https://scrutinizer-ci.com/g/kamilsk/check/badges/quality-score.png?b=master
[icon_research]:   https://img.shields.io/badge/research-in%20progress-yellow.svg
[icon_tw_author]:  https://img.shields.io/badge/author-%40kamilsk-blue.svg
[icon_tw_sponsor]: https://img.shields.io/badge/sponsor-%40octolab-blue.svg
[icon_twitter]:    https://img.shields.io/twitter/url/http/shields.io.svg?style=social

[page_build]:      https://travis-ci.org/kamilsk/check
[page_promo]:      https://github.com/kamilsk/check
[page_quality]:    https://scrutinizer-ci.com/g/kamilsk/check/?branch=master

[twitter_publish]: https://twitter.com/intent/tweet?text=Check%20Something%20as%20a%20Service&url=https://github.com/kamilsk/check&via=ikamilsk&hashtags=go,tool,website-audit
