> # check
> > Check Something as a Service.

[![Patreon](https://img.shields.io/badge/patreon-donate-orange.svg)](https://www.patreon.com/octolab)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Usage

### Quick start

```bash
$ check urls https://kamil.samigullin.info/
# [200] https://kamil.samigullin.info/
# ├─── [200] https://howilive.ru/en/
# ├─── [200] https://github.com/kamilsk
# ├─── [200] https://twitter.com/ikamilsk
# ├─── [999] https://www.linkedin.com/in/kamilsk/en
# └─── [200] https://kamil.samigullin.info/ru/
# [200] https://kamil.samigullin.info/ru/
# ├─── [200] https://howilive.ru/
# ├─── [200] https://github.com/kamilsk
# ├─── [200] https://twitter.com/ikamilsk
# ├─── [999] https://www.linkedin.com/in/kamilsk
# └─── [200] https://kamil.samigullin.info/
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

## Notes

- tested on Go 1.5, 1.6, 1.7, 1.8, 1.9 and 1.10

---

[![@kamilsk](https://img.shields.io/badge/author-%40kamilsk-blue.svg)](https://twitter.com/ikamilsk)
[![@octolab](https://img.shields.io/badge/sponsor-%40octolab-blue.svg)](https://twitter.com/octolab_inc)

made with ❤️ by [OctoLab](https://www.octolab.org/)
