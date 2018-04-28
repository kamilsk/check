> # check
> [![Analytics](https://ga-beacon.appspot.com/UA-109817251-19/check/readme?pixel)](https://github.com/kamilsk/check)
> Check Something as a Service.

[![Patreon](https://img.shields.io/badge/patreon-donate-orange.svg)](https://www.patreon.com/octolab)
[![Build Status](https://travis-ci.org/kamilsk/check.svg?branch=master)](https://travis-ci.org/kamilsk/check)
[![Code Coverage](https://scrutinizer-ci.com/g/kamilsk/check/badges/coverage.png?b=master)](https://scrutinizer-ci.com/g/kamilsk/check/?branch=master)
[![Code Quality](https://scrutinizer-ci.com/g/kamilsk/check/badges/quality-score.png?b=master)](https://scrutinizer-ci.com/g/kamilsk/check/?branch=master)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Usage

### Quick start

```bash
$ check urls https://kamil.samigullin.info/
# [200] https://kamil.samigullin.info/
#     ├─── [200] https://howilive.ru/en/
#     ├─── [200] https://github.com/kamilsk
#     ├─── [200] https://twitter.com/ikamilsk
#     ├─── [200] https://kamil.samigullin.info/ru/
#     └─── [999] https://www.linkedin.com/in/kamilsk
# [200] https://kamil.samigullin.info/ru/
#     ├─── [200] https://howilive.ru/
#     ├─── [200] https://github.com/kamilsk
#     ├─── [200] https://twitter.com/ikamilsk
#     ├─── [200] https://kamil.samigullin.info/
#     └─── [999] https://www.linkedin.com/in/kamilsk
$ check urls https://www.octolab.org/ | grep '\[3[0-9][0-9]\]'
#     ├─── [301] https://m.do.co/c/b2a387de5da4 -> https://www.digitalocean.com...
```

## Installation

```bash
$ egg github.com/kamilsk/check@^1.0.0 -- make test install
```

### Mirror

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

- tested on Go 1.8, 1.9 and 1.10

---

[![@kamilsk](https://img.shields.io/badge/author-%40kamilsk-blue.svg)](https://twitter.com/ikamilsk)
[![@octolab](https://img.shields.io/badge/sponsor-%40octolab-blue.svg)](https://twitter.com/octolab_inc)

made with ❤️ by [OctoLab](https://www.octolab.org/)
